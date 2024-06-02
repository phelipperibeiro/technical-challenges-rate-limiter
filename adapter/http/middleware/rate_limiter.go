package middleware

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	db "github.com/phelipperibeiro/technical-challenges-rate-limiter/adapter/db/redis"
)

// RateLimiter define a estrutura do limitador de taxa com parâmetros de configuração
type RateLimiter struct {
	maxIPRequests    int      // Máximo de requisições permitidas por IP
	maxTokenRequests int      // Máximo de requisições permitidas por token
	blockDuration    int      // Duração do bloqueio em segundos após exceder o limite
	cache            db.Cache // Interface do cache (possivelmente implementado com Redis)
}

// Limit é o método que implementa o middleware de limite de taxa
func (rateLimiter *RateLimiter) Limit(next http.Handler) http.Handler {

	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {

		// Extrai o IP do cliente do request
		id, err := extractClientIP(request)
		if err != nil {
			log.Println("Failed to get remote address: ", err)
			http.Error(responseWriter, "Internal Server Error", http.StatusBadRequest)
			return
		}

		kind := "addr"                           // Define o tipo padrão como endereço IP
		if request.Header.Get("API_KEY") != "" { // Verifica se o header API_KEY está presente
			id = request.Header.Get("API_KEY")
			kind = "token" // Se presente, define o tipo como token
		}

		// Formata a chave para o cache usando o tipo e o identificador (IP ou token)
		key := fmt.Sprintf("%s:%s", kind, id)
		log.Println("Key: ", key)

		// Obtém o valor correspondente à chave do cache
		val, err := rateLimiter.cache.Get(key)
		if err != nil {
			log.Println("Failed to get value from cache: ", err)
			http.Error(responseWriter, "Internal Server Error", http.StatusBadGateway)
			return
		}

		if val == "" { // Se o valor não existe no cache
			val = "1"                                                                                    // Inicializa o valor como "1"
			err := rateLimiter.cache.Set(key, val, time.Duration(rateLimiter.blockDuration)*time.Second) // Define o valor no cache
			if err != nil {
				log.Println("Failed to set value in cache: ", err)
				http.Error(responseWriter, "Internal Server Error", http.StatusBadGateway)
				return
			}
			next.ServeHTTP(responseWriter, request) // Chama o próximo handler na cadeia
			return
		}

		// Converte o valor obtido do cache para um inteiro
		count, err := strconv.Atoi(val)
		if err != nil {
			log.Println("Failed to convert value to int: ", err)
			http.Error(responseWriter, "Internal Server Error", http.StatusBadGateway)
			return
		}

		// Define o limite máximo de requisições baseado no tipo
		maxRequest := rateLimiter.maxIPRequests
		if kind == "token" {
			maxRequest = rateLimiter.maxTokenRequests
		}

		if count+1 > maxRequest { // Verifica se o número de requisições excede o limite
			log.Println("Too many requests")
			responseWriter.WriteHeader(http.StatusTooManyRequests)
			responseWriter.Write([]byte("you have reached the maximum number of requests or actions allowed within a certain time frame"))
			return
		}

		// Incrementa o contador no cache e atualiza o tempo de bloqueio
		err = rateLimiter.cache.Set(key, strconv.Itoa(count+1), time.Duration(rateLimiter.blockDuration)*time.Second)
		if err != nil {
			log.Println("Filed to set value in cache: ", err)
			http.Error(responseWriter, "Internal Server Error", http.StatusBadGateway)
			return
		}

		next.ServeHTTP(responseWriter, request) // Chama o próximo handler na cadeia
	})
}

// extractClientIP extrai o IP do cliente a partir do request
func extractClientIP(r *http.Request) (string, error) {
	// Obtém a lista de IPs do header X-Forwarded-For
	forwardedIPs := r.Header.Get("X-Forwarded-For")
	ipList := strings.Split(forwardedIPs, ",")

	if len(ipList) > 0 { // Se a lista de IPs não for vazia
		clientIP := net.ParseIP(strings.TrimSpace(ipList[len(ipList)-1])) // Obtém o último IP da lista
		if clientIP != nil {
			return clientIP.String(), nil
		}
	}

	// Se o header X-Forwarded-For não estiver presente ou não contiver IPs válidos, usa o RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	clientIP := net.ParseIP(ip)
	if clientIP != nil {
		if ip == "::1" { // Se o IP for "::1" (loopback IPv6), retorna "127.0.0.1"
			return "127.0.0.1", nil
		}
		return clientIP.String(), nil
	}

	return "", errors.New("client IP not found") // Retorna erro se não conseguir obter o IP
}

// Factory method to create a new RateLimiter
func NewRateLimiter(maxIPRequests, maxTokenRequests, blockDuration int, cache db.Cache) *RateLimiter {
	return &RateLimiter{
		maxIPRequests:    maxIPRequests,    // Define o número máximo de requisições por IP
		maxTokenRequests: maxTokenRequests, // Define o número máximo de requisições por token
		blockDuration:    blockDuration,    // Define a duração do bloqueio em segundos
		cache:            cache,            // Define a interface do cache
	}
}

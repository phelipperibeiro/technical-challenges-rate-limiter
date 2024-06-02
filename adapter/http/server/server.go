package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"                                                                        // Importa o pacote chi para roteamento
	db "github.com/phelipperibeiro/technical-challenges-rate-limiter/adapter/db/redis"                // Importa o pacote db para usar o cache Redis
	middleware "github.com/phelipperibeiro/technical-challenges-rate-limiter/adapter/http/middleware" // Importa o pacote middleware para usar o limitador de taxa
)

// NewWebServer cria e configura um novo servidor web com limitador de taxa
func NewWebServer(maxRequestsWithoutToken, maxTokenRequests, blockDuration int, redis db.Cache) *chi.Mux {

	// Loga a configuração do limitador de taxa
	log.Printf(
		"Rate Limiter Configuration: maxRequestsWithoutToken=%d, maxTokenRequests=%d, blockDuration=%d",
		maxRequestsWithoutToken,
		maxTokenRequests,
		blockDuration,
	)

	// Cria uma nova instância do limitador de taxa com os parâmetros fornecidos
	rateLimiter := middleware.NewRateLimiter(maxRequestsWithoutToken, maxTokenRequests, blockDuration, redis)

	// Cria um novo roteador chi
	router := chi.NewRouter()

	// Adiciona o middleware do limitador de taxa ao roteador
	router.Use(rateLimiter.Limit)

	// Define uma rota GET para o caminho raiz "/"
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Service is running")) // Responde com uma mensagem simples indicando que o serviço está funcionando
	})

	return router // Retorna o roteador configurado
}

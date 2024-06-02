package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache é uma interface que define os métodos para interagir com o cache
type Cache interface {
	Get(key string) (string, error)                        // Método para obter um valor do cache
	Set(key, value string, expiration time.Duration) error // Método para definir um valor no cache com uma expiração
}

// RedisCache é uma estrutura que implementa a interface Cache usando Redis
type RedisCache struct {
	client *redis.Client // Cliente Redis para interagir com o servidor Redis
}

// NewRedis cria uma nova instância de RedisCache e retorna um erro se a conexão falhar
func NewRedis(addr string) (*RedisCache, error) {
	// Cria um novo cliente Redis com o endereço fornecido
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	// Testa a conexão com o Redis usando o comando PING
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("FAILED TO CONNECT TO REDIS: %v", err)
	}

	// Retorna a instância de RedisCache com o cliente conectado
	return &RedisCache{
		client: client,
	}, nil
}

// Get obtém um valor do Redis com a chave fornecida
func (redisCache *RedisCache) Get(key string) (string, error) {
	// Tenta obter o valor da chave do Redis
	val, err := redisCache.client.Get(context.Background(), key).Result()
	if !errors.Is(err, redis.Nil) && err != nil {
		return "", fmt.Errorf("FAILED TO GET VALUE FROM REDIS: %v", err)
	}
	return val, nil
}

// Set define um valor no Redis com a chave e expiração fornecidas
func (redisCache *RedisCache) Set(key, value string, expiration time.Duration) error {
	// Tenta definir o valor da chave no Redis
	err := redisCache.client.Set(context.Background(), key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("FAILED TO SET VALUE IN REDIS: %v", err)
	}
	return nil
}

// Delete remove um valor do Redis com a chave fornecida
func (redisCache *RedisCache) Delete(key string) error {
	// Tenta deletar a chave do Redis
	err := redisCache.client.Del(context.Background(), key).Err()
	if err != nil {
		return fmt.Errorf("FAILED TO DELETE VALUE FROM REDIS: %v", err)
	}
	return nil
}

// Clear remove todos os valores do Redis
func (redisCache *RedisCache) Clear() error {
	// Tenta remover todas as chaves do Redis
	err := redisCache.client.FlushAll(context.Background()).Err()
	if err != nil {
		return fmt.Errorf("FAILED TO FLUSH ALL KEYS: %v", err)
	}
	return nil
}

package main

import (
	"fmt"
	"log"
	"net/http"

	db "github.com/phelipperibeiro/technical-challenges-rate-limiter/adapter/db/redis"
	"github.com/phelipperibeiro/technical-challenges-rate-limiter/adapter/http/server"
	"github.com/phelipperibeiro/technical-challenges-rate-limiter/config"
)

func main() {

	configs, err := config.LoadConfig(".")

	if err != nil {
		panic(err)
	}

	log.Println("Creating redis ...")

	redis, err := db.NewRedis(fmt.Sprintf("%s:%s", configs.RedisHost, configs.RedisPort))

	if err != nil {
		panic(err)
	}

	webServer := server.NewWebServer(
		configs.MaxRequestsWithoutToken,
		configs.MaxRequestsWithToken,
		configs.TimeBlockInSecond,
		redis,
	)

	log.Println("Starting web server on port", "8080")

	err = http.ListenAndServe(fmt.Sprintf(":%s", "8080"), webServer)

	if err != nil {
		panic(err)
	}
}

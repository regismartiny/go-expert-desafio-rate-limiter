package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/regismartiny/go-expert-desafio-rate-limiter/configs"
	db "github.com/regismartiny/go-expert-desafio-rate-limiter/internal/infra/database"
	web "github.com/regismartiny/go-expert-desafio-rate-limiter/internal/infra/web"
	"github.com/regismartiny/go-expert-desafio-rate-limiter/internal/infra/web/webserver"
)

func main() {
	ctx := context.Background()
	configs, err := configs.LoadConfig(".")
	if err != nil {
		log.Fatalf("Could not load configurations: %v\n", err)
	}

	log.Println("Configurations:")
	log.Println("ServerPort:", configs.ServerPort)
	log.Println("RateLimiter:", configs.RateLimiter)
	log.Println("Persistence:", configs.Persistence)

	webserver := webserver.NewWebServer(configs.ServerPort)

	rateLimiterRepository := db.RateLimiterRepositoryStrategy(ctx, configs.Persistence, "redis")
	rateLimiterMiddleware := web.NewRateLimiterMiddleware(ctx, configs.RateLimiter, rateLimiterRepository)
	webserver.AddMiddleware(rateLimiterMiddleware.Handle)
	homeHandler := web.NewHomeHandler()
	webserver.AddHandler("/", homeHandler.Handle)
	webserver.Start()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	webserver.Stop(ctx)
}

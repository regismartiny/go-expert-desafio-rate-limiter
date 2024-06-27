package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/regismartiny/go-expert-desafio-rate-limiter/configs"
	web "github.com/regismartiny/go-expert-desafio-rate-limiter/internal/infra/web"
	"github.com/regismartiny/go-expert-desafio-rate-limiter/internal/infra/web/webserver"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		log.Fatalf("Could not load configurations: %v\n", err)
	}

	fmt.Println("Configurations:")
	fmt.Println("ReqsPerSecond:", configs.ReqsPerSecond)
	fmt.Println("Token Configs:", configs.TokenConfigs)

	webserver := webserver.NewWebServer(":8080")
	rateLimiterMiddleware := web.NewRateLimiterMiddleware(configs, nil)
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

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
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

	fmt.Println("Configurations:")
	fmt.Println("ServerPort:", configs.ServerPort)
	fmt.Println("ReqsPerSecond:", configs.ReqsPerSecond)
	fmt.Println("Token Configs:", configs.TokenConfigs)
	fmt.Println("Redis Configs:", configs.Redis)

	redisClient := getRedisClient(ctx, configs)

	webserver := webserver.NewWebServer(configs.ServerPort)
	ratelimiterMiddlewareConfigs := web.RateLimiterMiddlewareConfigs{
		ReqsPerSecond: configs.ReqsPerSecond,
		TokenConfigs:  configs.TokenConfigs,
	}

	rateLimiterRepository := db.NewRateLimiterRedisRepository(ctx, redisClient)

	rateLimiterMiddleware := web.NewRateLimiterMiddleware(ratelimiterMiddlewareConfigs, rateLimiterRepository)
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
func getRedisClient(ctx context.Context, configs *configs.Conf) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     configs.Redis.Addr,
		Password: configs.Redis.Password,
		DB:       configs.Redis.Db,
	})

	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(pong)

	return redisClient
}

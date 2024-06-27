package web

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/regismartiny/go-expert-desafio-rate-limiter/configs"
)

type RateLimiterMiddleware struct {
	Configs *configs.Conf
	Db      *sql.DB
}

func NewRateLimiterMiddleware(
	Configs *configs.Conf,
	Db *sql.DB,
) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{
		Configs: Configs,
		Db:      Db,
	}
}

func (h *RateLimiterMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implementar o rate limiter
		fmt.Println("Passou pelo rate limiter middleware")
		w.Write([]byte("Passou pelo rate limiter middleware\n"))

		next.ServeHTTP(w, r)
	})
}

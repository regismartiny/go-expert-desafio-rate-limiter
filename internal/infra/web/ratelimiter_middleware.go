package web

import (
	"fmt"
	"net"
	"net/http"

	configs "github.com/regismartiny/go-expert-desafio-rate-limiter/configs"
	db "github.com/regismartiny/go-expert-desafio-rate-limiter/internal/infra/database"
)

type RateLimiterMiddleware struct {
	RateLimiter *RateLimiter
}

func NewRateLimiterMiddleware(
	Configs configs.RateLimiterConfigs,
	Repository db.RateLimiterRepository,
) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{
		RateLimiter: &RateLimiter{
			Configs:    Configs,
			Repository: Repository,
		},
	}
}

func (h *RateLimiterMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ipAddr := getIP(r)
		apiKeyHeader := r.Header.Get("API_KEY")

		fmt.Println("ipAddr", ipAddr)
		fmt.Println("apiKeyHeader", apiKeyHeader)

		allow := h.RateLimiter.Allow(ipAddr, apiKeyHeader)

		if allow {
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("Rate limit exceeded"))
		}
	})
}

func getIP(req *http.Request) string {
	ip, port, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		fmt.Printf("userip: %q is not IP:port\n", req.RemoteAddr)
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		//return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
		fmt.Printf("userip: %q is not IP:port\n", req.RemoteAddr)
		return ""
	}

	// This will only be defined when site is accessed via non-anonymous proxy
	// and takes precedence over RemoteAddr
	// Header.Get is case-insensitive
	forward := req.Header.Get("X-Forwarded-For")

	fmt.Printf("IP: %s\n", ip)
	fmt.Printf("Port: %s\n", port)
	fmt.Printf("Forwarded for: %s\n", forward)

	return ip
}

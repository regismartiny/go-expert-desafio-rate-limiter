package web

import (
	"fmt"
	"net"
	"net/http"

	configs "github.com/regismartiny/go-expert-desafio-rate-limiter/configs"
	db "github.com/regismartiny/go-expert-desafio-rate-limiter/internal/infra/database"
	rateLimiter "github.com/regismartiny/go-expert-desafio-rate-limiter/internal/ratelimiter"
)

type RateLimiterMiddleware struct {
	RateLimiter *rateLimiter.RateLimiter
}

func NewRateLimiterMiddleware(
	Configs configs.RateLimiterConfigs,
	Repository db.RateLimiterRepository,
) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{
		RateLimiter: rateLimiter.NewRateLimiter(
			rateLimiter.RateLimiterConfigs{
				BlockingDuration:   Configs.BlockingDuration,
				IpMaxReqsPerSecond: Configs.IpMaxReqsPerSecond,
				TokenConfigs:       Configs.TokenConfigs},
			Repository),
	}
}

func (h *RateLimiterMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		apiKeyHeader := r.Header.Get("API_KEY")
		ipAddr, err := getIP(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Println("ipAddr", ipAddr)
		fmt.Println("apiKeyHeader", apiKeyHeader)

		allow := h.RateLimiter.Allow(ipAddr, apiKeyHeader)

		if allow {
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("you have reached the maximum number of requests or actions allowed within a certain time frame"))
		}
	})
}

func getIP(req *http.Request) (string, error) {
	ip, port, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		fmt.Printf("userip: %q is not IP:port\n", req.RemoteAddr)
		return "", err
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		//return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
		fmt.Printf("userip: %q is not IP:port\n", req.RemoteAddr)
		return "", err
	}

	// This will only be defined when site is accessed via non-anonymous proxy
	// and takes precedence over RemoteAddr
	// Header.Get is case-insensitive
	forward := req.Header.Get("X-Forwarded-For")

	fmt.Printf("IP: %s\n", ip)
	fmt.Printf("Port: %s\n", port)
	fmt.Printf("Forwarded for: %s\n", forward)

	return ip, nil
}

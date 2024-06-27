package web

import (
	"fmt"
	"net"
	"net/http"

	db "github.com/regismartiny/go-expert-desafio-rate-limiter/internal/infra/database"
)

type RateLimiterMiddleware struct {
	Configs    RateLimiterMiddlewareConfigs
	Repository *db.RateLimiterRepository
}

type RateLimiterMiddlewareConfigs struct {
	ReqsPerSecond int
	TokenConfigs  map[string]int
}

func NewRateLimiterMiddleware(
	Configs RateLimiterMiddlewareConfigs,
	Repository *db.RateLimiterRepository,
) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{
		Configs:    Configs,
		Repository: Repository,
	}
}

func (h *RateLimiterMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implementar o rate limiter
		fmt.Println("Passou pelo rate limiter middleware")
		w.Write([]byte("Passou pelo rate limiter middleware\n"))

		ipAddr := getIP(r)
		apiKeyHeader := r.Header.Get("API_KEY")

		fmt.Println("ipAddr", ipAddr)
		fmt.Println("apiKeyHeader", apiKeyHeader)

		ipAddrRequestCount, err := h.Repository.GetRequestCount(ipAddr)
		if err != nil {
			fmt.Println(err)
		}

		newIpAddrRequestCount := ipAddrRequestCount + 1

		h.Repository.SaveRequestCount(ipAddr, newIpAddrRequestCount)

		fmt.Println("ipAddrRequestCount", newIpAddrRequestCount)

		maxReqsPerSecond := h.Configs.ReqsPerSecond

		fmt.Println("maxReqsPerSecond", maxReqsPerSecond)

		next.ServeHTTP(w, r)
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

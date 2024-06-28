package web

import (
	"fmt"

	configs "github.com/regismartiny/go-expert-desafio-rate-limiter/configs"
	db "github.com/regismartiny/go-expert-desafio-rate-limiter/internal/infra/database"
)

type RateLimiter struct {
	Configs    configs.RateLimiterConfigs
	Repository db.RateLimiterRepository
}

func (r *RateLimiter) Allow(ipAddr string, apiKeyHeader string) bool {

	if apiKeyHeader != "" {
		tokenMaxReqsPerSecond := r.Configs.TokenConfigs[apiKeyHeader]
		if tokenMaxReqsPerSecond != 0 {

			fmt.Println("tokenMaxReqsPerSecond", tokenMaxReqsPerSecond)

			tokenRequestCount, err := r.Repository.GetRequestCount(apiKeyHeader)
			if err != nil {
				fmt.Println(err)
				return false
			}

			fmt.Println("tokenRequestCount", tokenRequestCount)

			r.Repository.SaveRequestCount(apiKeyHeader, tokenRequestCount+1)

			return tokenRequestCount < tokenMaxReqsPerSecond
		}
	}

	ipMaxReqsPerSecond := r.Configs.IpMaxReqsPerSecond
	fmt.Println("ipMaxReqsPerSecond", ipMaxReqsPerSecond)

	ipAddrRequestCount, err := r.Repository.GetRequestCount(ipAddr)
	if err != nil {
		fmt.Println(err)
		return false
	}

	fmt.Println("ipAddrRequestCount", ipAddrRequestCount)

	r.Repository.SaveRequestCount(ipAddr, ipAddrRequestCount+1)

	return ipAddrRequestCount < ipMaxReqsPerSecond
}

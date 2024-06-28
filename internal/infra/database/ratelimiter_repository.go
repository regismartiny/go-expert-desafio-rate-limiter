package database

import entity "github.com/regismartiny/go-expert-desafio-rate-limiter/internal/entity"

type RateLimiterRepository interface {
	GetActiveClients() (map[string]entity.ActiveClient, error)
	SaveActiveClients(clients map[string]entity.ActiveClient) error
}

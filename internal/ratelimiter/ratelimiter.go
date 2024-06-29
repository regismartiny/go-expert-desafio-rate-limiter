package ratelimiter

import (
	"fmt"
	"sync"
	"time"

	entity "github.com/regismartiny/go-expert-desafio-rate-limiter/internal/entity"
	db "github.com/regismartiny/go-expert-desafio-rate-limiter/internal/infra/database"
	"golang.org/x/time/rate"
)

type RateLimiterConfigs struct {
	BlockingDuration   time.Duration
	IpMaxReqsPerSecond int
	TokenConfigs       map[string]int
}

type RateLimiter struct {
	Configs       RateLimiterConfigs
	Repository    db.RateLimiterRepository
	activeClients ActiveClients
}

type ActiveClients struct {
	mu      sync.Mutex
	clients map[string]entity.ActiveClient
}

func NewRateLimiter(
	Configs RateLimiterConfigs,
	Repository db.RateLimiterRepository) *RateLimiter {

	rateLimiter := &RateLimiter{
		Configs:    Configs,
		Repository: Repository,
		activeClients: ActiveClients{
			clients: make(map[string]entity.ActiveClient)}}

	rateLimiter.loadActiveClients()

	// Unblock client after expiration time
	go func() {
		for {
			time.Sleep(1 * time.Second)

			now := time.Now()

			for k, v := range rateLimiter.activeClients.clients {
				client := rateLimiter.activeClients.clients[k]
				if v.Blocked && now.After(v.BlockedUntil) {
					client.Blocked = false
					client.BlockedUntil = time.Time{}
					rateLimiter.updateActiveClient(k, client)
				}
			}
		}
	}()

	// Remove inactive clients
	go func() {
		for {
			time.Sleep(3 * time.Minute)

			for k := range rateLimiter.activeClients.clients {
				client := rateLimiter.activeClients.clients[k]
				if time.Since(client.LastSeen) > 3*time.Minute {
					rateLimiter.removeActiveClient(client)
				}
			}
		}
	}()

	return rateLimiter
}

func (r *RateLimiter) loadActiveClients() {

	fmt.Println("Loading active clients...")

	activeClients, err := r.Repository.GetActiveClients()
	if err != nil {
		fmt.Println("Error loading active clients. Starting clean.", err)
		activeClients = make(map[string]entity.ActiveClient)
	}

	fmt.Printf("%d active clients loaded\n", len(activeClients))

	// populate limiters
	for k := range activeClients {
		activeClient := activeClients[k]
		maxReqsPerSecond := 0

		if activeClient.ClientType == entity.Ip {
			maxReqsPerSecond = r.Configs.IpMaxReqsPerSecond
		} else {
			if tokenConfig, ok := r.Configs.TokenConfigs[activeClient.ClientId]; ok {
				maxReqsPerSecond = tokenConfig
			}
		}

		if entry, ok := activeClients[k]; ok {
			entry.Limiter = getRateLimiter(maxReqsPerSecond)
			activeClients[k] = entry
		}
	}

	r.activeClients.mu.Lock()
	r.activeClients.clients = activeClients
	r.activeClients.mu.Unlock()
}

func (r *RateLimiter) saveActiveClients() {

	err := r.Repository.SaveActiveClients(r.activeClients.clients)
	if err != nil {
		panic(err)
	}
}

func (r *RateLimiter) addActiveClient(client entity.ActiveClient) {

	r.activeClients.mu.Lock()
	r.activeClients.clients[client.ClientId] = client
	r.activeClients.mu.Unlock()

	r.saveActiveClients()
}

func (r *RateLimiter) removeActiveClient(client entity.ActiveClient) {

	r.activeClients.mu.Lock()
	delete(r.activeClients.clients, client.ClientId)
	r.activeClients.mu.Unlock()

	r.saveActiveClients()
}

func (r *RateLimiter) updateActiveClient(key string, client entity.ActiveClient) {

	r.activeClients.mu.Lock()
	r.activeClients.clients[key] = client
	r.activeClients.mu.Unlock()

	r.saveActiveClients()
}

func (r *RateLimiter) Allow(ipAddr string, apiKeyHeader string) bool {

	if apiKeyHeader != "" {

		tokenMaxReqsPerSecond, ok := r.Configs.TokenConfigs[apiKeyHeader]
		if ok {

			fmt.Println("tokenMaxReqsPerSecond", tokenMaxReqsPerSecond)
			return r.verifyClientAllowed(apiKeyHeader, entity.Token, tokenMaxReqsPerSecond)
		}
	}

	ipMaxReqsPerSecond := r.Configs.IpMaxReqsPerSecond
	fmt.Println("ipMaxReqsPerSecond", ipMaxReqsPerSecond)
	return r.verifyClientAllowed(ipAddr, entity.Ip, ipMaxReqsPerSecond)
}

func (r *RateLimiter) verifyClientAllowed(id string, clientType entity.ClientType, maxReqsPerSecond int) bool {

	activeClient, exists := r.activeClients.clients[id]

	if !exists {
		activeClient = createActiveClient(id, clientType, maxReqsPerSecond)
		r.addActiveClient(activeClient)
		return true
	}

	activeClient.LastSeen = time.Now()

	if activeClient.Blocked {
		fmt.Println("Client is blocked until", activeClient.BlockedUntil)
		r.updateActiveClient(activeClient.ClientId, activeClient)
		return false
	}

	allow := activeClient.Limiter.Allow()

	if !allow {
		activeClient.Blocked = true
		activeClient.BlockedUntil = time.Now().Add(r.Configs.BlockingDuration)
		fmt.Printf("Blocking client %s until %s\n", activeClient.ClientId, activeClient.BlockedUntil)
	}

	r.updateActiveClient(activeClient.ClientId, activeClient)

	return allow
}

func createActiveClient(id string, clientType entity.ClientType, maxReqsPerSeconds int) entity.ActiveClient {
	return entity.ActiveClient{
		ClientId:     id,
		LastSeen:     time.Now(),
		ClientType:   clientType,
		BlockedUntil: time.Time{},
		Blocked:      false,
		Limiter:      getRateLimiter(maxReqsPerSeconds),
	}
}

func getRateLimiter(maxReqsPerSecond int) *rate.Limiter {
	return rate.NewLimiter(rate.Limit(maxReqsPerSecond), maxReqsPerSecond)
}

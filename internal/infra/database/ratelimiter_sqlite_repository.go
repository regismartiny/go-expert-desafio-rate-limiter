package database

import (
	"context"
	"database/sql"
	"time"

	entity "github.com/regismartiny/go-expert-desafio-rate-limiter/internal/entity"

	_ "github.com/mattn/go-sqlite3"
)

type RateLimiterSQLiteRepository struct {
	ctx    context.Context
	client *sql.DB
}

func NewRateLimiterSQLiteRepository(ctx context.Context, client *sql.DB) *RateLimiterSQLiteRepository {
	client.Exec("CREATE TABLE ActiveClient (ClientId TEXT NOT NULL, LastSeen DATETIME NOT NULL, ClientType INTEGER NOT NULL, BlockedUntil DATETIME, Blocked BOOLEAN NOT NULL);")
	return &RateLimiterSQLiteRepository{ctx: ctx, client: client}
}

func (r *RateLimiterSQLiteRepository) SaveActiveClients(clients map[string]entity.ActiveClient) error {

	for _, client := range clients {

		_, err := r.client.Exec("INSERT INTO ActiveClient (ClientId, LastSeen, ClientType, BlockedUntil, Blocked) VALUES (?, ?, ?, ?, ?)",
			client.ClientId, client.LastSeen, client.ClientType, client.BlockedUntil, client.Blocked)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *RateLimiterSQLiteRepository) GetActiveClients() (map[string]entity.ActiveClient, error) {

	activeClients := make(map[string]entity.ActiveClient, 0)

	rows, err := r.client.Query("SELECT ClientId, LastSeen, ClientType, BlockedUntil, Blocked FROM ActiveClient")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var clientId string
		var lastSeen time.Time
		var clientType int
		var blockedUntil time.Time
		var blocked bool

		err = rows.Scan(&clientId, &lastSeen, &clientType, &blockedUntil, &blocked)
		if err != nil {
			return nil, err
		}
		activeClients[clientId] = entity.ActiveClient{
			ClientId:     clientId,
			LastSeen:     lastSeen,
			ClientType:   entity.ClientType(clientType),
			BlockedUntil: blockedUntil,
			Blocked:      blocked,
		}
	}
	defer rows.Close()

	return activeClients, nil
}

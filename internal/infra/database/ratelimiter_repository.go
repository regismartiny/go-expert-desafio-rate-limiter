package database

import (
	"database/sql"
)

type RateLimiterRepository struct {
	Db *sql.DB
}

func NewRateLimiterRepository(db *sql.DB) *RateLimiterRepository {
	return &RateLimiterRepository{Db: db}
}

func (r *RateLimiterRepository) SaveRequestCount(key string, value int) error {

	return nil
}

func (r *RateLimiterRepository) GetRequestCount(key string) (int, error) {

	return 0, nil
}

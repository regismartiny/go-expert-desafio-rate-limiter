package database

type RateLimiterRepository interface {
	GetRequestCount(key string) (int, error)
	SaveRequestCount(key string, value int) error
}

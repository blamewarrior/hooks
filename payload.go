package hooks

import (
	"github.com/go-redis/redis"
)

type Payloads interface {
	Save(repositoryFullName, payload string) error
	// ListFor(repositoryFullName string, limit, offset int) ([]string, error)
	// Delete(repositoryFullName, payload string) error
}

type PayloadRepository struct {
	redisClient *redis.Client
}

func NewPayloadRepository(redisClient *redis.Client) *PayloadRepository {
	return &PayloadRepository{redisClient}
}

func (repo *PayloadRepository) Save(repositoryFullName, payload string) (err error) {

	err = repo.redisClient.LPush(repositoryFullName, payload).Err()

	return
}

// func (repo *PayloadRepository) ListFor(repositoryFullName string, limit, offset int) ([]string, error) {

// }

// func (repo *PayloadRepository) Delete(repositoryFullName, payload string) error {

// }

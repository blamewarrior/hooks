package hooks_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/blamewarrior/hooks"
	"github.com/go-redis/redis"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSavePayload(t *testing.T) {
	redisClient, teardown := setup()

	defer teardown()

	payloadRepo := hooks.NewPayloadRepository(redisClient)

	testRepo := "blamewarrior/test"

	testPayload := "test_payload"

	err := payloadRepo.Save(testRepo, testPayload)

	require.NoError(t, err)

	val, err := redisClient.LRange(testRepo, 0, 0).Result()
	require.NoError(t, err)
	require.Equal(t, 1, len(val))
	assert.Equal(t, testPayload, val[0])
}

func setup() (client *redis.Client, teardownFn func()) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "6379"
	}

	opts := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: "", // no password set
		DB:       0,  // use default DB
	}

	client = redis.NewClient(opts)

	err := client.Ping().Err()

	if err != nil {
		log.Fatalf("failed to establish connection with test redis by addr: %s:%s", host, port)
	}

	return client, func() {
		client.FlushDB()
	}

}

/*
   Copyright (C) 2016 The BlameWarrior Authors.
   This file is a part of BlameWarrior service.
   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.
   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.
   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

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

func TestListPayload(t *testing.T) {
	redisClient, teardown := setup()

	defer teardown()

	testRepo := "blamewarrior/test"
	testPayload := "test_payload"

	err := createTestPayload(redisClient, testRepo, testPayload)
	require.NoError(t, err)

	payloadRepo := hooks.NewPayloadRepository(redisClient)

	list, err := payloadRepo.List(testRepo, int64(2))
	require.NoError(t, err)
	require.Equal(t, 1, len(list))
	require.Equal(t, testPayload, list[0])
}

func TestDeletePayload(t *testing.T) {
	redisClient, teardown := setup()

	defer teardown()

	testRepo := "blamewarrior/test"
	testPayload := "test_payload"

	err := createTestPayload(redisClient, testRepo, testPayload)
	require.NoError(t, err)

	payloadRepo := hooks.NewPayloadRepository(redisClient)

	err = payloadRepo.Delete(testRepo, testPayload)
	require.NoError(t, err)

	val, err := redisClient.LRange(testRepo, 0, 0).Result()
	require.NoError(t, err)

	assert.Empty(t, val)
}

func createTestPayload(client *redis.Client, testRepo, testPayload string) error {
	return client.LPush(testRepo, testPayload).Err()
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

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

package hooks

import (
	"github.com/go-redis/redis"
)

type Payloads interface {
	Save(payload string) error
	List(limit int) ([]string, error)
	Delete(payload string) error
}

type PayloadRepository struct {
	redisClient *redis.Client
}

func NewPayloadRepository(redisClient *redis.Client) *PayloadRepository {
	return &PayloadRepository{redisClient}
}

func (repo *PayloadRepository) Save(payload string) (err error) {
	return repo.redisClient.LPush("hooks", payload).Err()
}

func (repo *PayloadRepository) List(limit int64) ([]string, error) {
	return repo.redisClient.LRange("hooks", 0, limit).Result()
}

func (repo *PayloadRepository) Delete(payload string) error {
	return repo.redisClient.LRem("hooks", 0, payload).Err()
}

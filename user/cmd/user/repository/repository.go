package repository

import (
	"github.com/redis/go-redis/v9"
)

type UserRepository struct {
	Redis *redis.Client
}

func NewUserRepository(redis *redis.Client) *UserRepository {
	return &UserRepository{
		Redis: redis,
	}
}

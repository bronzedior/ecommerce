package repository

import "github.com/redis/go-redis/v9"

type OrderRepository struct {
	Redis *redis.Client
}

func NewOrderRepository(redis *redis.Client) *OrderRepository {
	return &OrderRepository{
		Redis: redis,
	}
}

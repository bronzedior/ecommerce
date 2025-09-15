package repository

import "github.com/redis/go-redis/v9"

type ProductRepository struct {
	Redis *redis.Client
}

func NewProductRepository(redis *redis.Client) *ProductRepository {
	return &ProductRepository{
		Redis: redis,
	}
}

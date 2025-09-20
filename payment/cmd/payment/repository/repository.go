package repository

import "github.com/redis/go-redis/v9"

type PaymentRepository struct {
	Redis *redis.Client
}

func NewPaymentRepository(redis *redis.Client) *PaymentRepository {
	return &PaymentRepository{
		Redis: redis,
	}
}

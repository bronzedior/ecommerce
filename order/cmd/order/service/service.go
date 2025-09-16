package service

import "order/cmd/order/repository"

type OrderService struct {
	OrderRepo repository.OrderRepository
}

func NewOrderService(orderRepo repository.OrderRepository) *OrderService {
	return &OrderService{
		OrderRepo: orderRepo,
	}
}

package usecase

import "order/cmd/order/service"

type OrderUsecase struct {
	OrderService service.OrderService
}

func NewOrderUsecase(orderService service.OrderService) *OrderUsecase {
	return &OrderUsecase{
		OrderService: orderService,
	}
}

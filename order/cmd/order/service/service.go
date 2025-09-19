package service

import (
	"context"
	"order/cmd/order/repository"
	"order/models"

	"gorm.io/gorm"
)

type OrderService struct {
	OrderRepository repository.OrderRepository
}

func NewOrderService(orderRepository repository.OrderRepository) *OrderService {
	return &OrderService{
		OrderRepository: orderRepository,
	}
}

func (s *OrderService) CheckIdempotency(ctx context.Context, token string) (bool, error) {
	isExist, err := s.OrderRepository.CheckIdempotency(ctx, token)
	if err != nil {
		return false, err
	}

	return isExist, nil
}

func (s *OrderService) SaveIdempotencyToken(ctx context.Context, token string) error {
	err := s.OrderRepository.SaveIdempotencyToken(ctx, token)
	if err != nil {
		return err
	}

	return nil
}

func (s *OrderService) GetOrderInfoByOrderID(ctx context.Context, orderID int64) (models.Order, error) {
	orderInfo, err := s.OrderRepository.GetOrderInfoByOrderID(ctx, orderID)
	if err != nil {
		return models.Order{}, err
	}

	return orderInfo, nil
}

func (s *OrderService) GetOrderDetailByOrderDetailID(ctx context.Context, orderDetailID int64) (models.OrderDetail, error) {
	orderDetail, err := s.OrderRepository.GetOrderDetailByOrderDetailID(ctx, orderDetailID)
	if err != nil {
		return models.OrderDetail{}, err
	}

	return orderDetail, nil
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID int64, status int) error {
	err := s.OrderRepository.UpdateOrderStatus(ctx, orderID, status)
	if err != nil {
		return err
	}

	return nil
}

func (s *OrderService) SaveOrderAndOrderDetail(ctx context.Context, order *models.Order, orderDetail *models.OrderDetail) (int64, error) {
	var orderID int64

	err := s.OrderRepository.WithTransaction(ctx, func(tx *gorm.DB) error {
		err := s.OrderRepository.InsertOrderDetailTx(ctx, tx, orderDetail)
		if err != nil {
			return err
		}

		order.OrderDetailID = orderDetail.ID

		err = s.OrderRepository.InsertOrderTx(ctx, tx, order)
		if err != nil {
			return err
		}

		orderID = order.ID
		return nil
	})

	if err != nil {
		return 0, err
	}

	return orderID, nil
}

func (s *OrderService) GetOrderHistoriesByUserID(ctx context.Context, param models.OrderHistoryParam) ([]models.OrderHistoryResponse, error) {
	orderHistories, err := s.OrderRepository.GetOrderHistoriesByUserID(ctx, param)
	if err != nil {
		return nil, err
	}
	return orderHistories, nil
}

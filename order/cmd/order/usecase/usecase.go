package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"order/cmd/order/service"
	"order/infrastructure/constant"
	"order/kafka"
	"order/models"
	"time"
)

type OrderUsecase struct {
	OrderService service.OrderService
	Producer     kafka.KafkaProducer
}

func NewOrderUsecase(orderService service.OrderService, kafkaProducer kafka.KafkaProducer) *OrderUsecase {
	return &OrderUsecase{
		OrderService: orderService,
		Producer:     kafkaProducer,
	}
}

func (uc *OrderUsecase) CheckoutOrder(ctx context.Context, param *models.CheckoutRequest) (int64, error) {
	if param.IdempotencyToken != "" {
		isExist, err := uc.OrderService.CheckIdempotency(ctx, param.IdempotencyToken)
		if err != nil {
			return 0, err
		}

		if isExist {
			return 0, errors.New("order already created, please check again")
		}
	}

	if err := uc.validateProducts(ctx, param.Items); err != nil {
		return 0, err
	}

	totalQty, totalAmount := uc.calculateOrderSummary(param.Items)
	productJSON, historyJSON, err := uc.constructOrderDetail(param.Items)
	if err != nil {
		return 0, err
	}

	orderDetail := &models.OrderDetail{
		Products:     productJSON,
		OrderHistory: historyJSON,
	}

	order := &models.Order{
		UserID:          param.UserID,
		Amount:          totalAmount,
		TotalQty:        totalQty,
		Status:          constant.OrderStatusCreated,
		PaymentMethod:   param.PaymentMethod,
		ShippingAddress: param.ShippingAddress,
	}

	orderID, err := uc.OrderService.SaveOrderAndOrderDetail(ctx, order, orderDetail)
	if err != nil {
		return 0, err
	}

	if param.IdempotencyToken != "" {
		_ = uc.OrderService.SaveIdempotencyToken(ctx, param.IdempotencyToken)
	}

	orderCreatedEvent := models.OrderCreatedEvent{
		OrderID:         orderID,
		UserID:          param.UserID,
		TotalAmount:     order.Amount,
		PaymentMethod:   param.PaymentMethod,
		ShippingAddress: param.ShippingAddress,
	}
	err = uc.Producer.PublishOrderCreated(ctx, orderCreatedEvent)
	if err != nil {
		return 0, err
	}

	return orderID, nil
}

func (uc *OrderUsecase) validateProducts(ctx context.Context, items []models.CheckoutItem) error {
	for _, item := range items {
		productInfo, err := uc.OrderService.GetProductInfo(ctx, item.ProductID)
		if err != nil {
			return fmt.Errorf("failed get product info: %d, err %s", item.ProductID, err.Error())
		}

		if item.Price != productInfo.Price {
			return fmt.Errorf("invalid price for product %d (%.2f - %.2f)", item.ProductID, item.Price, productInfo.Price)
		}

		if item.Quantity <= 0 || item.Quantity > 1000 {
			return fmt.Errorf("invalid quantity product %d, maximum product quantity is 1000", item.ProductID)
		}

		if item.Quantity > productInfo.Stock {
			return fmt.Errorf("invalid prodouct quantity %d, product stock is only %d", item.ProductID, productInfo.Stock)
		}
	}

	return nil
}

func (uc *OrderUsecase) calculateOrderSummary(items []models.CheckoutItem) (int, float64) {
	var totalQty int
	var totalAmount float64
	for _, item := range items {
		totalQty += item.Quantity
		totalAmount += float64(item.Quantity) * item.Price
	}
	return totalQty, totalAmount
}

func (uc *OrderUsecase) constructOrderDetail(items []models.CheckoutItem) (string, string, error) {
	productsJSON, _ := json.Marshal(items)
	history := []map[string]interface{}{
		{"status": "created", "timestamp": time.Now()},
	}
	historyJSON, _ := json.Marshal(history)

	return string(productsJSON), string(historyJSON), nil
}

func (uc *OrderUsecase) GetOrderHistory(ctx context.Context, param models.OrderHistoryParam) ([]models.OrderHistoryResponse, error) {
	orderHistory, err := uc.OrderService.GetOrderHistoriesByUserID(ctx, param)
	if err != nil {
		return nil, err
	}

	return orderHistory, nil
}

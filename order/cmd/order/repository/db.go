package repository

import (
	"context"
	"errors"
	"order/models"
	"time"

	"gorm.io/gorm"
)

func (r *OrderRepository) InsertOrderDetailTx(ctx context.Context, tx *gorm.DB, orderDetail *models.OrderDetail) error {
	err := tx.WithContext(ctx).Table("order_detail").Create(orderDetail).Error
	return err
}

func (r *OrderRepository) InsertOrderTx(ctx context.Context, tx *gorm.DB, order *models.Order) error {
	err := tx.WithContext(ctx).Table("orders").Create(order).Error
	return err
}

func (r *OrderRepository) WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	tx := r.Database.Begin().WithContext(ctx)

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *OrderRepository) CheckIdempotency(ctx context.Context, token string) (bool, error) {
	var log models.OrderRequestLog
	err := r.Database.WithContext(ctx).Table("order_request_log").First(&log, "idempotency_token = ?", token).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return err == nil, err
}

func (r *OrderRepository) SaveIdempotencyToken(ctx context.Context, token string) error {
	log := models.OrderRequestLog{
		IdempotencyToken: token,
		CreateTime:       time.Now(),
	}
	return r.Database.WithContext(ctx).Table("order_request_log").Create(&log).Error
}

func (r *OrderRepository) GetOrderInfoByOrderID(ctx context.Context, orderID int64) (models.Order, error) {
	var result models.Order
	err := r.Database.Table("orders").WithContext(ctx).Where("order_id = ?", orderID).Find(&result).Error
	if err != nil {
		return models.Order{}, err
	}

	return result, nil
}

func (r *OrderRepository) GetOrderDetailByOrderDetailID(ctx context.Context, orderDetailID int64) (models.OrderDetail, error) {
	var result models.OrderDetail
	err := r.Database.Table("order_detail").WithContext(ctx).Where("id = ?", orderDetailID).Find(&result).Error
	if err != nil {
		return models.OrderDetail{}, err
	}

	return result, nil
}

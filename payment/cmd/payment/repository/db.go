package repository

import (
	"context"
	"payment/infrastructure/log"
	"payment/models"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PaymentDatabase interface {
	MarkPaid(ctx context.Context, orderID int64) error

	SavePayment(ctx context.Context, param models.Payment) error
}

type paymentDatabase struct {
	DB *gorm.DB
}

func NewPaymentDatabase(db *gorm.DB) PaymentDatabase {
	return &paymentDatabase{
		DB: db,
	}
}

func (r *paymentDatabase) MarkPaid(ctx context.Context, orderID int64) error {
	err := r.DB.Model(&models.Payment{}).Table("payments").WithContext(ctx).Where("order_id = ?", orderID).Update("status", "paid").Error
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"order_id": orderID,
		}).Errorf("r.DB.Update() got error: %v", err)
		return err
	}

	return nil
}

func (r *paymentDatabase) SavePayment(ctx context.Context, param models.Payment) error {
	err := r.DB.Table("payments").WithContext(ctx).Create(param).Error
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"param": param,
		}).Errorf("r.DB.Create() got error: %v", err)
		return err
	}

	return nil
}

package service

import (
	"context"
	"payment/cmd/payment/repository"
	"payment/infrastructure/log"

	"github.com/sirupsen/logrus"
)

type PaymentService interface {
	ProcessPaymentSuccess(ctx context.Context, orderID int64) error
}

type paymentService struct {
	database  repository.PaymentDatabase
	publisher repository.PaymentEventPublisher
}

func NewPaymentService(database repository.PaymentDatabase, publisher repository.PaymentEventPublisher) PaymentService {
	return &paymentService{
		database:  database,
		publisher: publisher,
	}
}

func (s *paymentService) ProcessPaymentSuccess(ctx context.Context, orderID int64) error {
	err := s.publisher.PublishPaymentSuccess(ctx, orderID)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"order_id": orderID,
		}).Errorf("s.publisher.PublishPaymentSuccess() got error: %v", err)
		return err
	}

	err = s.database.MarkPaid(ctx, orderID)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"order_id": orderID,
		}).Errorf("s.database.MarkPaid() got error: %v", err)
		return err
	}

	return nil
}

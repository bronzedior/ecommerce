package service

import (
	"context"
	"fmt"
	"payment/cmd/payment/repository"
	"payment/infrastructure/log"
	"payment/models"
	"time"

	"github.com/sirupsen/logrus"
)

type XenditService interface {
	CreateInvoice(ctx context.Context, param models.OrderCreatedEvent) error
}

type xenditService struct {
	database repository.PaymentDatabase
	xendit   repository.XenditClient
}

func NewXenditService(database repository.PaymentDatabase, xenditClient repository.XenditClient) XenditService {
	return &xenditService{
		database: database,
		xendit:   xenditClient,
	}
}

func (s *xenditService) CreateInvoice(ctx context.Context, param models.OrderCreatedEvent) error {
	externalID := fmt.Sprintf("order-%d", param.OrderID)
	req := models.XenditInvoiceRequest{
		ExternalID:  externalID,
		Amount:      param.TotalAmount,
		Description: fmt.Sprintf("[FC] Pembayaran Order %d", param.OrderID),
		PayerEmail:  fmt.Sprintf("user%d@test.com", param.UserID),
	}

	xenditInvoiceDetail, err := s.xendit.CreateInvoice(ctx, req)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"param":   param,
			"payload": req,
		}).Errorf("s.xendit.CreateInvoice() got error: %v", err)
		return err
	}

	newPayment := models.Payment{
		OrderID:     param.OrderID,
		UserID:      param.UserID,
		ExternalID:  externalID,
		Amount:      param.TotalAmount,
		Status:      "PENDING",
		CreateTime:  time.Now(),
		ExpiredTime: xenditInvoiceDetail.ExpiryDate,
	}
	err = s.database.SavePayment(ctx, newPayment)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"param":      param,
			"newPayment": newPayment,
		}).Errorf("s.database.SavePayment() got error: %v", err)
		return err
	}

	return nil
}

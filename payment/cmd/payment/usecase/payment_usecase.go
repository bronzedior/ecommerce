package usecase

import (
	"context"
	"payment/cmd/payment/service"
	"payment/infrastructure/log"
	"payment/models"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

type PaymentUsecase interface {
	ProcessPaymentWebhook(ctx context.Context, param models.XenditWebhookPayload) error
}

type paymentUsecase struct {
	Service service.PaymentService
}

func NewPaymentUsecase(svc service.PaymentService) PaymentUsecase {
	return &paymentUsecase{
		Service: svc,
	}
}

func (uc *paymentUsecase) ProcessPaymentWebhook(ctx context.Context, payload models.XenditWebhookPayload) error {
	switch payload.Status {
	case "PAID":
		orderID := extractOrderID(payload.ExternalID)
		err := uc.Service.ProcessPaymentSuccess(ctx, orderID)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"status":      payload.Status,
				"external_id": payload.ExternalID,
			}).Errorf("uc.service.ProcessPaymentSuccess() got error: %v", err)
			return err
		}
	case "FAILED":
		//
	case "PENDING":
		//
	default:
		log.Logger.WithFields(logrus.Fields{
			"status":      payload.Status,
			"external_id": payload.ExternalID,
		}).Infof("[%s] Anomaly Payment Webhook Status: %s", payload.ExternalID, payload.Status)
	}

	return nil
}

func extractOrderID(externalID string) int64 {
	// order id: 123456
	// key kafka event: "order-12345"
	idStr := strings.TrimPrefix(externalID, "order-")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	return id
}

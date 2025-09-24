package usecase

import (
	"context"
	"errors"
	"fmt"
	"payment/cmd/payment/service"
	"payment/infrastructure/constant"
	"payment/infrastructure/log"
	"payment/models"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type PaymentUsecase interface {
	ProcessPaymentRequests(ctx context.Context, payload models.OrderCreatedEvent) error

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

func (uc *paymentUsecase) ProcessPaymentRequests(ctx context.Context, payload models.OrderCreatedEvent) error {
	err := uc.Service.SavePaymentRequests(ctx, models.PaymentRequests{
		OrderID:    payload.OrderID,
		Amount:     payload.TotalAmount,
		UserID:     payload.UserID,
		Status:     "PENDING",
		CreateTime: time.Now(),
	})

	if err != nil {
		return err
	}

	return nil
}

func (uc *paymentUsecase) ProcessPaymentWebhook(ctx context.Context, payload models.XenditWebhookPayload) error {
	switch payload.Status {
	case "PAID":
		orderID := extractOrderID(payload.ExternalID)

		amount, err := uc.Service.CheckPaymentAmountByOrderID(ctx, orderID)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"order_id":       orderID,
				"status":         payload.Status,
				"external_id":    payload.ExternalID,
				"webhook_amount": payload,
			})
		}

		if amount != payload.Amount {
			errorMessage := fmt.Sprintf("Webhook amount mismatched: expected %.2f, got %.2f", amount, payload.Amount)
			paymentAnomaly := models.PaymentAnomaly{
				OrderID:     orderID,
				ExternalID:  payload.ExternalID,
				AnomalyType: constant.AnomalyTypeInvalidAmount,
				Notes:       errorMessage,
				Status:      constant.PaymentAnomalyStatusNeedToCheck,
				CreateTime:  time.Now(),
			}

			err := uc.Service.SavePaymentAnomaly(ctx, paymentAnomaly)
			if err != nil {
				log.Logger.WithFields(logrus.Fields{
					"payload":        payload,
					"paymentAnomaly": paymentAnomaly,
				}).WithError(err)
				return err
			}

			log.Logger.WithFields(logrus.Fields{
				"payload": payload,
			}).Errorf("Webhook amount mismatched: expected %.2f, got %.2f", amount, payload.Amount)
			err = errors.New(errorMessage)
			return err
		}

		err = uc.Service.ProcessPaymentSuccess(ctx, orderID)
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

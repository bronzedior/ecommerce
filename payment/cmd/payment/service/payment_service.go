package service

import (
	"context"
	"math"
	"payment/cmd/payment/repository"
	"payment/infrastructure/constant"
	"payment/infrastructure/log"
	"payment/models"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	maxRetryPublish = 5
)

type PaymentService interface {
	CheckPaymentAmountByOrderID(ctx context.Context, orderID int64) (float64, error)

	ProcessPaymentSuccess(ctx context.Context, orderID int64) error

	SavePaymentAnomaly(ctx context.Context, param models.PaymentAnomaly) error

	SavePaymentRequests(ctx context.Context, param models.PaymentRequests) error

	GetPaymentInfoByOrderID(ctx context.Context, orderID int64) (models.Payment, error)
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

func (s *paymentService) CheckPaymentAmountByOrderID(ctx context.Context, orderID int64) (float64, error) {
	amount, err := s.database.CheckPaymentAmountByOrderID(ctx, orderID)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"order_id": orderID,
		}).Errorf("s.database.CheckPaymentAmountByOrderID() got error: %v", err)
		return 0, err
	}

	return amount, nil
}

func (s *paymentService) GetPaymentInfoByOrderID(ctx context.Context, orderID int64) (models.Payment, error) {
	paymentInfo, err := s.database.GetPaymentInfoByOrderID(ctx, orderID)
	if err != nil {
		return models.Payment{}, err
	}

	return paymentInfo, nil
}

func (s *paymentService) SavePaymentAnomaly(ctx context.Context, param models.PaymentAnomaly) error {
	err := s.database.SavePaymentAnomaly(ctx, param)
	if err != nil {
		return err
	}

	return nil
}

func (s *paymentService) SavePaymentRequests(ctx context.Context, param models.PaymentRequests) error {
	err := s.database.SavePaymentRequests(ctx, param)
	if err != nil {
		return err
	}

	return nil
}

func (s *paymentService) ProcessPaymentSuccess(ctx context.Context, orderID int64) error {
	isAlreadyPaid, err := s.database.IsAlreadyPaid(ctx, orderID)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"order_id": orderID,
		}).Errorf("s.database.isAlreadyPaid() got error: %v", err)
		return err
	}

	if isAlreadyPaid {
		log.Logger.WithFields(logrus.Fields{
			"order_id": orderID,
		}).Infof("[skip - order %d] Payment status already paid!", orderID)
		return nil
	}

	err = retryPublishPayment(maxRetryPublish, func() error {
		return s.publisher.PublishPaymentSuccess(ctx, orderID)
	})
	if err != nil {
		failedEventsParam := models.FailedEvents{
			OrderID:    orderID,
			FailedType: constant.FailedPublishEventPaymentSuccess,
			Status:     constant.FailedPublishEventStatusNeedToCheck,
			Notes:      err.Error(),
			CreateTime: time.Now(),
		}

		errSaveFailedPublish := s.database.SaveFailedPublishEvent(ctx, failedEventsParam)
		if errSaveFailedPublish != nil {
			log.Logger.WithFields(logrus.Fields{
				"failedEventParam": failedEventsParam,
			}).WithError(errSaveFailedPublish)
			return errSaveFailedPublish
		}

		log.Logger.WithFields(logrus.Fields{
			"order_id": orderID,
		}).Errorf("s.publisher.PublishPaymentSuccess() got error: %v", err)
		return err
	}

	err = s.publisher.PublishPaymentSuccess(ctx, orderID)
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

func retryPublishPayment(max int, fn func() error) error {
	var err error
	for i := 0; i < max; i++ {
		err = fn()
		if err == nil {
			return nil
		}

		wait := time.Duration(math.Pow(2, float64(i))) * time.Second
		log.Logger.Printf("Retry: %d, Error: %s. Retrying in %d seconds...", i+1, err, wait)
		time.Sleep(wait)
	}

	return err
}

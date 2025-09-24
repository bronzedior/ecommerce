package service

import (
	"context"
	"fmt"
	"log"
	"payment/cmd/payment/repository"
	"payment/models"
	"time"
)

type SchedulerService struct {
	Database       repository.PaymentDatabase
	Xendit         repository.XenditClient
	Publisher      repository.PaymentEventPublisher
	PaymentService PaymentService
}

func (s *SchedulerService) StartSweepingExpiredPendingPayments() {
	go func(ctx context.Context) {
		for {
			log.Println("Scheduler StartSpeeingExpiredPendingPayments is running...")
			expiredPayments, err := s.Database.GetExpiredPendingPayments(ctx)
			if err != nil {
				log.Println("Failed get expired pending payments, err: ", err.Error())
				time.Sleep(5 * time.Minute)
				continue
			}

			for _, expiredPayment := range expiredPayments {
				err = s.Database.MarkExpired(ctx, expiredPayment.ID)
				if err != nil {
					log.Printf("[payment id: %d] Failed update expired, err: %s", expiredPayment.ID, err.Error())
				}
			}

			time.Sleep(10 * time.Minute)
		}
	}(context.Background())
}

func (s *SchedulerService) StartProcessFailedPaymentRequests() {
	go func(ctx context.Context) {
		for {
			var paymentRequests []models.PaymentRequests
			err := s.Database.GetFailedPaymentRequests(ctx, &paymentRequests)
			if err != nil {
				log.Panicln("Error get failed payment request! error: ", err.Error())

				time.Sleep(10 * time.Second)
				continue
			}

			for _, paymentRequest := range paymentRequests {
				err = s.Database.UpdatePendingPaymentRequests(ctx, paymentRequest.ID)
				if err != nil {
					log.Println("s.Database.UpdatePendingPaymentRequests() got error: ", err.Error())

					errUpdateStatus := s.Database.UpdateFailedPaymentRequests(ctx, paymentRequest.ID, err.Error())
					if errUpdateStatus != nil {
						log.Panicln("s.Database.UpdateFailedPaymentRequests() got error: ", err)
					}

					continue
				}
			}
			time.Sleep(1 * time.Minute)
		}

	}(context.Background())
}

func (s *SchedulerService) StartProcessPendingPaymentRequests() {
	go func(ctx context.Context) {
		for {
			var paymentRequests []models.PaymentRequests
			err := s.Database.GetPendingPaymentRequests(ctx, &paymentRequests)
			if err != nil {
				log.Println("s.Database.GetPendingPaymentRequests() got error: ", err.Error())
				time.Sleep(10 * time.Second)
				continue
			}

			for _, paymentRequest := range paymentRequests {
				log.Printf("[DEBUG] Processing Payment Request Order %d", paymentRequest.OrderID)

				paymentInfo, err := s.Database.GetPaymentInfoByOrderID(ctx, paymentRequest.OrderID)
				if err != nil {
					log.Println("s.Database.GetPaymentInfoByOrderID() got error ", err.Error())
					continue
				}

				externalID := fmt.Sprintf("order-%d", paymentRequest.OrderID)
				if paymentInfo.ID != 0 {
					err = s.Database.UpdateSuccessPaymentRequests(ctx, paymentRequest.ID)
					if err != nil {
						log.Printf("[req id: %d] s.Database.UpdateSuccessPaymentRequest() got error: %s", paymentRequest.ID, err.Error())
					}
					continue
				}

				xenditInvoiceDetail, err := s.Xendit.CreateInvoice(ctx, models.XenditInvoiceRequest{
					ExternalID:  externalID,
					Amount:      paymentRequest.Amount,
					Description: fmt.Sprintf("[FC] Pembayaran Order %d", paymentRequest.OrderID),
					PayerEmail:  fmt.Sprintf("user%d@test.com", paymentRequest.UserID),
				})
				if err != nil {
					log.Printf("[req id: %d] s.Xendit.CreateInvoice() got error: %v", paymentRequest.ID, err.Error())

					errSaveFailedPaymentRequest := s.Database.UpdateFailedPaymentRequests(ctx, paymentRequest.ID, err.Error())
					if errSaveFailedPaymentRequest != nil {
						log.Printf("[req id: %d] s.Database.UpdateFailedPaymentRequests() got error: %v", paymentRequest.ID, errSaveFailedPaymentRequest.Error())
					}

					continue
				}

				err = s.Database.UpdateSuccessPaymentRequests(ctx, paymentRequest.ID)
				if err != nil {
					log.Printf("[req id: %d] s.Database.UpdateSuccessPaymentRequest() got error: %s", paymentRequest.ID, err.Error())
				}

				err = s.Database.SavePayment(ctx, models.Payment{
					OrderID:     paymentRequest.OrderID,
					UserID:      paymentRequest.UserID,
					Amount:      paymentRequest.Amount,
					ExternalID:  externalID,
					Status:      "PENDING",
					CreateTime:  time.Now(),
					ExpiredTime: xenditInvoiceDetail.ExpiryDate,
				})
				if err != nil {
					log.Printf("[req id: %d] s.Database.SavePayment() got error: %s", paymentRequest.ID, err.Error())
				}
			}

			time.Sleep(5 * time.Second)
		}
	}(context.Background())
}

func (s *SchedulerService) StartCheckPendingInvoices() {
	ticker := time.NewTicker(10 * time.Minute)

	go func() {
		for range ticker.C {
			ctx := context.Background()
			listPendingInvoices, err := s.Database.GetPendingInvoices(ctx)
			if err != nil {
				log.Println("s.Database.GetPendingInvoices() got error: ", err.Error())
				continue
			}

			for _, pendingInvoice := range listPendingInvoices {
				invoiceStatus, err := s.Xendit.CheckInvoiceStatus(ctx, pendingInvoice.ExternalID)
				if err != nil {
					log.Println("s.Xendit.CheckInvoiceStatus() got error: ", err.Error())
					continue
				}

				if invoiceStatus == "PAID" {
					err = s.PaymentService.ProcessPaymentSuccess(ctx, pendingInvoice.ID)
					if err != nil {
						log.Println("s.PaymentService.ProcessPaymentSuccess() got error: ", err)
						continue
					}
				}
			}
		}
	}()
}

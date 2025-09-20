package usecase

import "payment/cmd/payment/service"

type PaymentUsecase struct {
	PaymentService service.PaymentService
}

func NewPaymentUsecase(paymentService service.PaymentService) *PaymentUsecase {
	return &PaymentUsecase{
		PaymentService: paymentService,
	}
}

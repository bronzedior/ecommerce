package handler

import (
	"net/http"
	"payment/cmd/payment/usecase"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	PaymentUsecase usecase.PaymentUsecase
}

func NewPaymentHandler(paymentUsecase usecase.PaymentUsecase) *PaymentHandler {
	return &PaymentHandler{
		PaymentUsecase: paymentUsecase,
	}
}

func (h *PaymentHandler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}

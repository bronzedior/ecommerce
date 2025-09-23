package handler

import (
	"net/http"
	"payment/cmd/payment/usecase"
	"payment/infrastructure/log"
	"payment/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type PaymentHandler interface {
	HandleXenditWebhook(c *gin.Context)
}

type paymentHandler struct {
	Usecase usecase.PaymentUsecase
}

func NewPaymentHandler(usecase usecase.PaymentUsecase) PaymentHandler {
	return &paymentHandler{
		Usecase: usecase,
	}
}

func (h *paymentHandler) HandleXenditWebhook(c *gin.Context) {
	var payload models.XenditWebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Logger.WithFields(logrus.Fields{
			"payload": payload,
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload", "error_message": err.Error()})
		return
	}

	err := h.Usecase.ProcessPaymentWebhook(c.Request.Context(), payload)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"payload": payload,
		})
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success!",
	})
	return
}

package controllers

import (
	"net/http"

	"github.com/Rabiann/weather-mailer/pkg/go/api/internal/dto"
	"github.com/gin-gonic/gin"
)

type (
	SubscriptionController struct {
		SubscriptionService SubscriptionService
	}

	MailingService interface {
		SendConfirmationLetter(string, string) error
	}

	SubscriptionService interface {
		Subscribe(dto.SubscriptionRequest, *gin.Context) error
		Confirm(*gin.Context) error
		Unsubscribe(*gin.Context) error
	}
)

func NewSubscriptionController(subscriptionService SubscriptionService) SubscriptionController {
	return SubscriptionController{SubscriptionService: subscriptionService}
}

func (s *SubscriptionController) Subscribe(ctx *gin.Context) {
	var subscription dto.SubscriptionRequest

	if err := ctx.ShouldBind(&subscription); err != nil {
		ctx.JSON(400, gin.H{"status": "bad request"})
		return
	}

	if err := s.SubscriptionService.Subscribe(subscription, ctx); err != nil {
		ctx.JSON(300, gin.H{"status": err.Error()})
		return
	}

	ctx.HTML(http.StatusOK, "needconfirmation.html", gin.H{})
}

func (s *SubscriptionController) Confirm(ctx *gin.Context) {

	if err := s.SubscriptionService.Confirm(ctx); err != nil {
		ctx.HTML(400, "registrationfailed.html", gin.H{})
		return
	}

	ctx.HTML(http.StatusOK, "registration.html", gin.H{})
}

func (s SubscriptionController) Unsubscribe(ctx *gin.Context) {
	if err := s.SubscriptionService.Unsubscribe(ctx); err != nil {
		ctx.JSON(400, gin.H{"status": "invalid params"})
		return
	}

	ctx.HTML(http.StatusOK, "unsubscription.html", gin.H{})
}

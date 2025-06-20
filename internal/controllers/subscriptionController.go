package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/Rabiann/weather-mailer/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	SubscriptionController struct {
		SubscriptionService SubscriptionService
	}

	MailingService interface {
		SendConfirmationLetter(string, string) error
	}

	TokenService interface {
		CreateToken(uint, context.Context, context.CancelFunc) (uuid.UUID, error)
		GetSubscriptionOfToken(uuid.UUID, context.Context, context.CancelFunc) (uint, error)
		UseToken(uuid.UUID, context.Context, context.CancelFunc) error
	}

	SubscriptionService interface {
		Subscribe(models.Subscription, *gin.Context, context.Context, context.CancelFunc) error
		Confirm(*gin.Context, context.Context, context.CancelFunc) error
		Unsubscribe(*gin.Context, context.Context, context.CancelFunc) error
	}
)

func NewSubscriptionController(subscriptionService SubscriptionService) SubscriptionController {
	return SubscriptionController{SubscriptionService: subscriptionService}
}

func (s *SubscriptionController) Subscribe(ctx *gin.Context) {
	var subscription models.Subscription
	ctx_, cancel := context.WithTimeout(ctx.Request.Context(), 2*time.Second)
	defer cancel()

	if err := ctx.ShouldBind(&subscription); err != nil {
		ctx.JSON(400, gin.H{"status": "bad request"})
		return
	}

	if err := s.SubscriptionService.Subscribe(subscription, ctx, ctx_, cancel); err != nil {
		ctx.JSON(400, gin.H{"status": "bad request"})
		return
	}

	ctx.HTML(http.StatusOK, "needconfirmation.html", gin.H{})
}

func (s *SubscriptionController) Confirm(ctx *gin.Context) {
	ctx_, cancel := context.WithTimeout(ctx.Request.Context(), 2*time.Second)
	defer cancel()
	handleTokenErr := func(ctx *gin.Context, err error, code int) {
		ctx.HTML(code, "registrationfailed.html", gin.H{})
	}

	if err := s.SubscriptionService.Confirm(ctx, ctx_, cancel); err != nil {
		handleTokenErr(ctx, err, 400)
		return
	}

	ctx.HTML(http.StatusOK, "registration.html", gin.H{})
}

func (s SubscriptionController) Unsubscribe(ctx *gin.Context) {
	ctx_, cancel := context.WithTimeout(ctx.Request.Context(), 2*time.Second)
	defer cancel()
	if err := s.SubscriptionService.Unsubscribe(ctx, ctx_, cancel); err != nil {
		ctx.JSON(400, gin.H{"status": "invalid params"})
		return
	}

	ctx.HTML(http.StatusOK, "unsubscription.html", gin.H{})
}

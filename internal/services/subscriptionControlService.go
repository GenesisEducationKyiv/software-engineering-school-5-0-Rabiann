package services

import (
	"context"
	"fmt"

	"github.com/Rabiann/weather-mailer/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	SubscriptionControlService struct {
		subscriptionDataService SubscriptionDataServer
		tokenService            TokenServer
		emailService            EmailServer
		baseUrl                 string
	}

	SubscriptionDataServer interface {
		AddSubscription(models.Subscription, context.Context, context.CancelFunc) (uint, error)
		ActivateSubscription(uint, context.Context, context.CancelFunc) (string, error)
		DeleteSubscription(uint, context.Context, context.CancelFunc) error
	}

	TokenServer interface {
		CreateToken(uint, context.Context, context.CancelFunc) (uuid.UUID, error)
		GetSubscriptionOfToken(uuid.UUID, context.Context, context.CancelFunc) (uint, error)
		UseToken(uuid.UUID, context.Context, context.CancelFunc) error
	}

	EmailServer interface {
		SendConfirmationLetter(recipient string, confirmationUrl string) error
	}
)

func NewSubscriptionBusinessService(subscriptionService SubscriptionDataServer, tokenService TokenServer, emailService EmailServer, baseUrl string) *SubscriptionControlService {
	return &SubscriptionControlService{subscriptionService, tokenService, emailService, baseUrl}
}

func (s *SubscriptionControlService) Subscribe(subscription models.Subscription, ctx *gin.Context, ctx_ context.Context, cancel context.CancelFunc) error {
	id, err := s.subscriptionDataService.AddSubscription(MapSubscription(subscription), ctx_, cancel)
	if err != nil {
		return err
	}

	token, err := s.tokenService.CreateToken(id, ctx_, cancel)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/confirm/%s", s.baseUrl, token)

	if err := s.emailService.SendConfirmationLetter(subscription.Email, url); err != nil {
		return err
	}

	return nil
}

func (s *SubscriptionControlService) Confirm(ctx *gin.Context, ctx_ context.Context, cancel context.CancelFunc) error {
	defer cancel()
	token, err := uuid.Parse(ctx.Param("token"))
	if err != nil {
		return err
	}

	subscriberId, err := s.tokenService.GetSubscriptionOfToken(token, ctx_, cancel)
	if err != nil {
		return err
	}

	if err := s.tokenService.UseToken(token, ctx_, cancel); err != nil {
		return err
	}

	_, err = s.subscriptionDataService.ActivateSubscription(subscriberId, ctx_, cancel)
	if err != nil {
		return err
	}

	return nil
}

func (s *SubscriptionControlService) Unsubscribe(ctx *gin.Context, ctx_ context.Context, cancel context.CancelFunc) error {
	token, err := uuid.Parse(ctx.Param("token"))
	if err != nil {
		return err
	}
	subscriberId, err := s.tokenService.GetSubscriptionOfToken(token, ctx_, cancel)
	if err != nil {
		return err
	}

	if err := s.tokenService.UseToken(token, ctx_, cancel); err != nil {
		return err
	}

	if err := s.subscriptionDataService.DeleteSubscription(subscriberId, ctx_, cancel); err != nil {
		return err
	}

	return err
}

package services

import (
	"context"
	"fmt"

	"github.com/Rabiann/weather-mailer/internal/models"
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
		AddSubscription(models.Subscription, context.Context) (uint, error)
		ActivateSubscription(uint, context.Context) (string, error)
		DeleteSubscription(uint, context.Context) error
	}

	TokenServer interface {
		CreateToken(uint, context.Context) (uuid.UUID, error)
		GetSubscriptionOfToken(uuid.UUID, context.Context) (uint, error)
		UseToken(uuid.UUID, context.Context) error
	}

	EmailServer interface {
		SendConfirmationLetter(recipient string, confirmationUrl string) error
	}
)

func NewSubscriptionBusinessService(subscriptionService SubscriptionDataServer, tokenService TokenServer, emailService EmailServer, baseUrl string) *SubscriptionControlService {
	return &SubscriptionControlService{subscriptionService, tokenService, emailService, baseUrl}
}

func (s *SubscriptionControlService) Subscribe(subscription models.Subscription, ctx context.Context) error {
	id, err := s.subscriptionDataService.AddSubscription(MapSubscription(subscription), ctx)
	if err != nil {
		return err
	}

	token, err := s.tokenService.CreateToken(id, ctx)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/confirm/%s", s.baseUrl, token)

	if err := s.emailService.SendConfirmationLetter(subscription.Email, url); err != nil {
		return err
	}

	return nil
}

func (s *SubscriptionControlService) Confirm(token uuid.UUID, ctx context.Context) error {
	subscriberId, err := s.tokenService.GetSubscriptionOfToken(token, ctx)
	if err != nil {
		return err
	}

	if err := s.tokenService.UseToken(token, ctx); err != nil {
		return err
	}

	_, err = s.subscriptionDataService.ActivateSubscription(subscriberId, ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *SubscriptionControlService) Unsubscribe(token uuid.UUID, ctx context.Context) error {
	subscriberId, err := s.tokenService.GetSubscriptionOfToken(token, ctx)
	if err != nil {
		return err
	}

	if err := s.tokenService.UseToken(token, ctx); err != nil {
		return err
	}

	if err := s.subscriptionDataService.DeleteSubscription(subscriberId, ctx); err != nil {
		return err
	}

	return err
}

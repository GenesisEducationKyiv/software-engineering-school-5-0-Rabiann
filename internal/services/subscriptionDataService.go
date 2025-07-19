package services

import (
	"context"
	"errors"
	"net/mail"

	"github.com/Rabiann/weather-mailer/internal/models"
)

type (
	SubscriptionDataService struct {
		subscriptionRepository SubscriptionRepository
	}

	SubscriptionRepository interface {
		AddSubscription(subscription models.Subscription, ctx context.Context) (uint, error)
		ActivateSubscription(id uint, ctx context.Context) (string, error)
		GetActiveSubscriptions(per string, ctx context.Context) ([]models.Subscription, error)
		UpdateSubscription(id uint, new_subscription models.Subscription, ctx context.Context) error
		DeleteSubscription(id uint, ctx context.Context) error
	}
)

func NewSubscriptionService(subscriptionRepository SubscriptionRepository) *SubscriptionDataService {
	return &SubscriptionDataService{subscriptionRepository}
}

func MapSubscription(subscriptionRequest models.Subscription) (models.Subscription, error) {
	var mapped models.Subscription

	_, err := mail.ParseAddress(subscriptionRequest.Email)
	if err != nil {
		return mapped, err
	}

	frequency := subscriptionRequest.Frequency
	if frequency != "daily" && frequency != "hourly" {
		return mapped, errors.New("Invalid frequency")
	}

	return models.Subscription{
		Email:     subscriptionRequest.Email,
		Frequency: frequency,
		City:      subscriptionRequest.City,
		Confirmed: false,
	}, nil
}

func (s *SubscriptionDataService) AddSubscription(subscription models.Subscription, ctx context.Context) (uint, error) {
	return s.subscriptionRepository.AddSubscription(subscription, ctx)
}

func (s SubscriptionDataService) ActivateSubscription(id uint, ctx context.Context) (string, error) {
	return s.subscriptionRepository.ActivateSubscription(id, ctx)
}

func (s SubscriptionDataService) GetActiveSubscriptions(per string, ctx context.Context) ([]models.Subscription, error) {
	return s.subscriptionRepository.GetActiveSubscriptions(per, ctx)
}

func (s SubscriptionDataService) UpdateSubscription(id uint, new_subscription models.Subscription, ctx context.Context) error {
	return s.subscriptionRepository.UpdateSubscription(id, new_subscription, ctx)
}

func (s SubscriptionDataService) DeleteSubscription(id uint, ctx context.Context) error {
	return s.subscriptionRepository.DeleteSubscription(id, ctx)
}

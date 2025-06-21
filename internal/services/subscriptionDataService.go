package services

import (
	"context"

	"github.com/Rabiann/weather-mailer/internal/models"
)

type (
	SubscriptionDataService struct {
		subscriptionRepository SubscriptionRepository
	}

	SubscriptionRepository interface {
		AddSubscription(subscription models.Subscription, ctx context.Context, cancel context.CancelFunc) (uint, error)
		ActivateSubscription(id uint, ctx context.Context, cancel context.CancelFunc) (string, error)
		GetActiveSubscriptions(per string, ctx context.Context, cancel context.CancelFunc) ([]models.Subscription, error)
		UpdateSubscription(id uint, new_subscription models.Subscription, ctx context.Context, cancel context.CancelFunc) error
		DeleteSubscription(id uint, ctx context.Context, cancel context.CancelFunc) error
	}
)

func NewSubscriptionService(subscriptionRepository SubscriptionRepository) *SubscriptionDataService {
	return &SubscriptionDataService{subscriptionRepository}
}

func MapSubscription(subscriptionRequest models.Subscription) models.Subscription {
	return models.Subscription{
		Email:     subscriptionRequest.Email,
		Frequency: subscriptionRequest.Frequency,
		City:      subscriptionRequest.City,
		Confirmed: false,
	}
}

func (s *SubscriptionDataService) AddSubscription(subscription models.Subscription, ctx context.Context, cancel context.CancelFunc) (uint, error) {
	return s.subscriptionRepository.AddSubscription(subscription, ctx, cancel)
}

func (s SubscriptionDataService) ActivateSubscription(id uint, ctx context.Context, cancel context.CancelFunc) (string, error) {
	return s.subscriptionRepository.ActivateSubscription(id, ctx, cancel)
}

func (s SubscriptionDataService) GetActiveSubscriptions(per string, ctx context.Context, cancel context.CancelFunc) ([]models.Subscription, error) {
	return s.subscriptionRepository.GetActiveSubscriptions(per, ctx, cancel)
}

func (s SubscriptionDataService) UpdateSubscription(id uint, new_subscription models.Subscription, ctx context.Context, cancel context.CancelFunc) error {
	return s.subscriptionRepository.UpdateSubscription(id, new_subscription, ctx, cancel)
}

func (s SubscriptionDataService) DeleteSubscription(id uint, ctx context.Context, cancel context.CancelFunc) error {
	return s.subscriptionRepository.DeleteSubscription(id, ctx, cancel)
}

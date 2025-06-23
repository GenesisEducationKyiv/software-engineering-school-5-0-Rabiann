package persistance

import (
	"context"
	"errors"
	"github.com/Rabiann/weather-mailer/internal/models"
	"gorm.io/gorm"
)

type (
	SubscriptionRepository struct {
		Db *gorm.DB
	}
)

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db}
}

func (s *SubscriptionRepository) GetSubscriptions(ctx context.Context, cancel context.CancelFunc) ([]models.Subscription, error) {
	defer cancel()
	var subscriptions []models.Subscription
	result := s.Db.WithContext(ctx).Find(&subscriptions)
	return subscriptions, result.Error
}

func (s *SubscriptionRepository) GetSubscriptionById(id uint, ctx context.Context, cancel context.CancelFunc) (models.Subscription, error) {
	defer cancel()
	subscription := models.Subscription{ID: id}
	result := s.Db.WithContext(ctx).First(&subscription)
	return subscription, result.Error
}

func (s *SubscriptionRepository) AddSubscription(subscription models.Subscription, ctx context.Context, cancel context.CancelFunc) (uint, error) {
	if s.Db == nil {
		return 0, nil
	}
	result := s.Db.WithContext(ctx).Create(&subscription)
	return subscription.ID, result.Error
}

func (s *SubscriptionRepository) ActivateSubscription(id uint, ctx context.Context, cancel context.CancelFunc) (string, error) {
	var subscription models.Subscription
	subscription.ID = id

	result := s.Db.Find(&subscription)
	if result.Error != nil {
		return "", result.Error
	}

	if subscription.Confirmed {
		return "", errors.New("subscription already confirmed")
	}

	subscription.Confirmed = true
	result = s.Db.WithContext(ctx).Save(subscription)
	return subscription.Email, result.Error
}

func (s *SubscriptionRepository) GetActiveSubscriptions(per string, ctx context.Context, cancel context.CancelFunc) ([]models.Subscription, error) {
	var subscribers []models.Subscription
	result := s.Db.WithContext(ctx).Where("frequency = ? and confirmed = true", per).Find(&subscribers)

	if result.Error != nil {
		return nil, result.Error
	}

	return subscribers, nil
}

func (s *SubscriptionRepository) UpdateSubscription(id uint, new_subscription models.Subscription, ctx context.Context, cancel context.CancelFunc) error {
	subscription := models.Subscription{ID: id}

	if id != new_subscription.ID {
		return errors.New("IDs differ")
	}

	result := s.Db.WithContext(ctx).Find(&subscription)

	if result.Error != nil {
		return result.Error
	}

	subscription.City = new_subscription.City
	subscription.Confirmed = new_subscription.Confirmed
	subscription.CreatedAt = new_subscription.CreatedAt
	subscription.Email = new_subscription.Email
	subscription.Frequency = new_subscription.Frequency

	result = s.Db.WithContext(ctx).Save(subscription)
	return result.Error
}

func (s *SubscriptionRepository) DeleteSubscription(id uint, ctx context.Context, cancel context.CancelFunc) error {
	result := s.Db.WithContext(ctx).Delete(&models.Subscription{}, id)
	return result.Error
}

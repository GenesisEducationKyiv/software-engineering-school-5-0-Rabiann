package persistance

import (
	"context"
	"errors"
	"time"

	"github.com/Rabiann/weather-mailer/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	TokenRepository struct {
		Db *gorm.DB
	}
)

func NewTokenRepository(db *gorm.DB) *TokenRepository {
	return &TokenRepository{db}
}

func (t *TokenRepository) CreateToken(subscriptionId uint, ctx context.Context, cancel context.CancelFunc) (uuid.UUID, error) {
	defer cancel()
	id := uuid.New()

	token := models.Token{
		ID:             id,
		SubscriptionID: subscriptionId,
		Expires:        time.Now().Add(time.Hour * 24),
	}

	result := t.Db.WithContext(ctx).Create(&token)
	return id, result.Error
}

func (t *TokenRepository) GetSubscriptionOfToken(id uuid.UUID, ctx context.Context, cancel context.CancelFunc) (uint, error) {
	var token models.Token
	token.ID = id
	defer cancel()

	result := t.Db.WithContext(ctx).Find(&token)
	return token.SubscriptionID, result.Error
}

func (t *TokenRepository) UseToken(id uuid.UUID, ctx context.Context, cancel context.CancelFunc) error {
	token := models.Token{
		ID: id,
	}
	defer cancel()

	result := t.Db.WithContext(ctx).First(&token)
	if result.Error != nil {
		return result.Error
	}

	if time.Now().Compare(token.Expires) > 0 {
		if result := t.Db.WithContext(ctx).Delete(token); result.Error != nil {
			return result.Error
		}
		return errors.New("token already expired")
	}

	if result := t.Db.WithContext(ctx).Delete(token); result.Error != nil {
		return result.Error
	}

	return nil
}

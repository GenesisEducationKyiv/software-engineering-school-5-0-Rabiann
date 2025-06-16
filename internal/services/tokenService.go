package services

import (
	"errors"
	"time"

	"github.com/Rabiann/weather-mailer/internal/services/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TokenServer interface {
	CreateToken(subscriptionId uint) (uuid.UUID, error)
	GetSubscriptionOfToken(id uuid.UUID) (uint, error)
	UseToken(id uuid.UUID) error
}

type TokenService struct {
	Db *gorm.DB
}

func NewTokenService(db *gorm.DB) TokenService {
	return TokenService{db}
}

func (t TokenService) CreateToken(subscriptionId uint) (uuid.UUID, error) {
	id := uuid.New()

	token := models.Token{
		ID:             id,
		SubscriptionID: subscriptionId,
		Expires:        time.Now().Add(time.Hour * 24),
	}

	result := t.Db.Create(&token)
	return id, result.Error
}

func (t TokenService) GetSubscriptionOfToken(id uuid.UUID) (uint, error) {
	var token models.Token
	token.ID = id

	result := t.Db.Find(&token)
	return token.SubscriptionID, result.Error
}

func (t TokenService) UseToken(id uuid.UUID) error {
	token := models.Token{
		ID: id,
	}

	result := t.Db.First(&token)
	if result.Error != nil {
		return result.Error
	}

	if time.Now().Compare(token.Expires) > 0 {
		if result := t.Db.Delete(token); result.Error != nil {
			return result.Error
		}
		return errors.New("token already expired")
	}

	if result := t.Db.Delete(token); result.Error != nil {
		return result.Error
	}

	return nil
}

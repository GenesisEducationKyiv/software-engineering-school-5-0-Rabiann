package services

import (
	"context"

	"github.com/google/uuid"
)

type (
	TokenService struct {
		tokenRepository TokenRepository
	}

	TokenRepository interface {
		CreateToken(subscriptionId uint, ctx context.Context) (uuid.UUID, error)
		GetSubscriptionOfToken(id uuid.UUID, ctx context.Context) (uint, error)
		UseToken(id uuid.UUID, ctx context.Context) error
	}
)

func NewTokenService(tokenRepository TokenRepository) *TokenService {
	return &TokenService{tokenRepository}
}

func (t TokenService) CreateToken(subscriptionId uint, ctx context.Context) (uuid.UUID, error) {
	return t.tokenRepository.CreateToken(subscriptionId, ctx)
}

func (t TokenService) GetSubscriptionOfToken(id uuid.UUID, ctx context.Context) (uint, error) {
	return t.tokenRepository.GetSubscriptionOfToken(id, ctx)
}

func (t TokenService) UseToken(id uuid.UUID, ctx context.Context) error {
	return t.tokenRepository.UseToken(id, ctx)
}

package services

import (
	"context"
	"fmt"

	"github.com/Rabiann/weather-mailer/pkg/go/api/internal/dto"
	pUser "github.com/Rabiann/weather-mailer/protos/user"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type (
	SubscriptionControlService struct {
		userClient   pUser.UserServiceClient
		emailService EmailServer
		baseUrl      string
	}

	EmailServer interface {
		SendConfirmationLetter(recipient string, confirmationUrl string) error
	}
)

func NewSubscriptionBusinessService(conn *grpc.ClientConn, emailService EmailServer, baseUrl string) *SubscriptionControlService {
	c := pUser.NewUserServiceClient(conn)
	return &SubscriptionControlService{c, emailService, baseUrl}
}

func (s *SubscriptionControlService) Subscribe(subscription dto.SubscriptionRequest, ctx *gin.Context) error {
	userClient := s.userClient
	reply, err := userClient.AddSubscription(ctx, &pUser.SubscriptionRequest{
		Email:  subscription.Email,
		City:   subscription.City,
		Period: subscription.Frequency,
	})

	if err != nil {
		return err
	}

	token, err := userClient.CreateToken(ctx, &pUser.CreateTokenRequest{SubscriptionId: reply.Id})
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/confirm/%s", s.baseUrl, token)

	if err := s.emailService.SendConfirmationLetter(subscription.Email, url); err != nil {
		return err
	}

	return nil
}

func (s *SubscriptionControlService) Confirm(ctx *gin.Context) error {
	userClient := s.userClient
	token, err := uuid.Parse(ctx.Param("token"))
	if err != nil {
		return err
	}

	reply, err := userClient.GetSubscriptionOfToken(ctx, &pUser.GetSubOfToken{Uuid: token.String()})
	if err != nil {
		return err
	}

	subscriberId := reply.Sub

	_, err = userClient.UseToken(ctx, &pUser.UseTokenRequest{Uuid: token.String()})
	if err != nil {
		return err
	}

	_, err = userClient.ActivateSubscription(ctx, &pUser.ActivateRequest{Id: subscriberId})
	if err != nil {
		return err
	}

	return nil
}

func (s *SubscriptionControlService) Unsubscribe(ctx *gin.Context) error {
	userClient := s.userClient
	token, err := uuid.Parse(ctx.Param("token"))
	if err != nil {
		return err
	}
	reply, err := userClient.GetSubscriptionOfToken(ctx, &pUser.GetSubOfToken{Uuid: token.String()})
	if err != nil {
		return err
	}

	subscriberId := reply.Sub

	if _, err := userClient.UseToken(ctx, &pUser.UseTokenRequest{Uuid: token.String()}); err != nil {
		return err
	}

	if _, err := userClient.DeleteSubscription(ctx, &pUser.DeleteRequest{Id: subscriberId}); err != nil {
		return err
	}

	return err
}

func (s *SubscriptionControlService) GetActiveSubscriptions(period string, ctx context.Context) ([]dto.Subscription, error) {
	perRequest := pUser.PeriodRequest{
		Period: period,
	}

	var subscribers []dto.Subscription

	reply, err := s.userClient.GetActiveSubscriptions(ctx, &perRequest)
	if err != nil {
		return nil, err
	}

	for _, sub := range reply.Subscriptions {
		subscribers = append(subscribers, dto.Subscription{
			City:   sub.City,
			Email:  sub.Email,
			Period: sub.Period,
		})
	}

	return subscribers, nil

}

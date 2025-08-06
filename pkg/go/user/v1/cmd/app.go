package cmd

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/Rabiann/weather-mailer/lib/go/config"
	"github.com/Rabiann/weather-mailer/lib/go/logger"
	localConf "github.com/Rabiann/weather-mailer/pkg/go/user/internal/config"
	persistance "github.com/Rabiann/weather-mailer/pkg/go/user/internal/database"
	"github.com/Rabiann/weather-mailer/pkg/go/user/internal/models"
	"github.com/Rabiann/weather-mailer/pkg/go/user/internal/repositories"
	"github.com/Rabiann/weather-mailer/pkg/go/user/internal/services"
	pb "github.com/Rabiann/weather-mailer/protos/user"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type App struct{}

type Server struct {
	pb.UnimplementedUserServiceServer
	subscriptionService *services.SubscriptionDataService
	tokenService        *services.TokenService
}

func (s *Server) AddSubscription(ctx context.Context, in *pb.SubscriptionRequest) (*pb.SubscriptionReply, error) {
	log.Printf("Add subscription")
	subscriptionService := s.subscriptionService
	subscription := models.Subscription{
		Email:     in.Email,
		City:      in.City,
		Frequency: in.Period,
	}
	id, err := subscriptionService.AddSubscription(subscription, ctx)
	return &pb.SubscriptionReply{
		Id: int32(id),
	}, err
}

func (s *Server) ActivateSubscription(ctx context.Context, in *pb.ActivateRequest) (*pb.ActivateReply, error) {
	log.Printf("Activate subscription")
	subscriptionService := s.subscriptionService
	id := uint(in.Id)
	token, err := subscriptionService.ActivateSubscription(id, ctx)
	return &pb.ActivateReply{
		Id: token,
	}, err
}

func (s *Server) GetActiveSubscriptions(ctx context.Context, in *pb.PeriodRequest) (*pb.ActiveReply, error) {
	log.Printf("Get active subscriptions")
	subscriptionService := s.subscriptionService
	per := in.Period
	subscriptions, err := subscriptionService.GetActiveSubscriptions(per, ctx)

	var mappedSubscriptions []*pb.ActiveSubscriber
	for _, sub := range subscriptions {
		mappedSub := &pb.ActiveSubscriber{
			Id:     int32(sub.ID),
			Email:  sub.Email,
			City:   sub.City,
			Period: sub.Email,
		}
		mappedSubscriptions = append(mappedSubscriptions, mappedSub)
	}

	return &pb.ActiveReply{
		Subscriptions: mappedSubscriptions,
	}, err
}

func (s *Server) UpdateSubscriptions(ctx context.Context, in *pb.UpdateRequest) error {
	log.Printf("Update subscriptions")
	subscriptionService := s.subscriptionService
	id := uint(in.Id)
	new := models.Subscription{
		Email:     in.Subscriber.Email,
		City:      in.Subscriber.City,
		Frequency: in.Subscriber.Period,
	}

	return subscriptionService.UpdateSubscription(id, new, ctx)
}

func (s *Server) DeleteSubscription(ctx context.Context, in *pb.DeleteRequest) (*empty.Empty, error) {
	log.Printf("Delete subscription")
	subscriptionService := s.subscriptionService
	id := uint(in.Id)

	err := subscriptionService.DeleteSubscription(id, ctx)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *Server) CreateToken(ctx context.Context, in *pb.CreateTokenRequest) (*pb.TokenReply, error) {
	log.Printf("Create token")
	tokenService := s.tokenService
	id := uint(in.SubscriptionId)

	uuid, err := tokenService.CreateToken(id, ctx)
	if err != nil {
		return nil, err
	}

	return &pb.TokenReply{
		Uuid: uuid.String(),
	}, nil
}

func (s *Server) GetSubscriptionOfToken(ctx context.Context, in *pb.GetSubOfToken) (*pb.SubReply, error) {
	log.Printf("Get sub of token")
	tokenService := s.tokenService
	id, err := uuid.Parse(in.Uuid)
	if err != nil {
		return nil, err
	}

	sub, err := tokenService.GetSubscriptionOfToken(id, ctx)
	if err != nil {
		return nil, err
	}

	return &pb.SubReply{
		Sub: int32(sub),
	}, nil
}

func (s *Server) UseToken(ctx context.Context, in *pb.UseTokenRequest) (*empty.Empty, error) {
	log.Printf("Use token")
	tokenService := s.tokenService
	id, err := uuid.Parse(in.Uuid)
	if err != nil {
		return nil, err
	}

	err = tokenService.UseToken(id, ctx)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func bootstrapDatabase(configuration *localConf.Configuration) (*gorm.DB, error) {
	db := persistance.ConnectToDatabase(configuration)
	if err := persistance.Migrate(db); err != nil {
		return nil, err
	}

	return db, nil
}

func (a *App) Run() error {
	configuration := &localConf.Configuration{}
	err := config.LoadEnvironment(configuration)
	if err != nil {
		return err
	}

	db, err := bootstrapDatabase(configuration)
	if err != nil {
		return err
	}

	os.Mkdir("logs", 0700)
	logger.SetupLogger("logs/app.log")
	subscriptionRepository := repositories.NewSubscriptionRepository(db)
	tokenRepository := repositories.NewTokenRepository(db)
	subscriptionDataService := services.NewSubscriptionService(subscriptionRepository)
	tokenService := services.NewTokenService(tokenRepository)

	server := &Server{
		subscriptionService: subscriptionDataService,
		tokenService:        tokenService,
	}

	s := grpc.NewServer()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 5555))
	pb.RegisterUserServiceServer(s, server)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve %v", err)
	}

	return err
}

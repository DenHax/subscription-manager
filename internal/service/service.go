package service

import (
	"github.com/DenHax/subscription-manager/internal/domain/models"
	"github.com/DenHax/subscription-manager/internal/repo"
	"github.com/DenHax/subscription-manager/internal/service/subscription"
)

type Subscriptions interface {
	CreateSubscrition(serviceName string, price int, userID string, startDate string, endDate *string) (*models.Subscription, error)
	Subscription(id string) (*models.Subscription, error)
	DeleteSubscription(id string) error
	GetAllSubscriptions(userID *string, serviceName *string, limit, offset int) ([]*models.Subscription, int, error)
	UpdateSubscription(id, serviceName *string, price *int, userID *string, startDate *string, endDate *string) (*models.Subscription, error)
	SummarySubscription(startDate, endDate string, userID *string, serviceName *string) (int, error)
}

type Service struct {
	Subscriptions
}

func NewService(repos *repo.Repository) *Service {
	subService := subscription.NewSubService(repos.Subscriptions)
	return &Service{
		Subscriptions: subService,
	}
}

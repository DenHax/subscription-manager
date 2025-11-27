package repo

import (
	"github.com/DenHax/subscription-manager/internal/domain/models"
	"github.com/DenHax/subscription-manager/internal/repo/subscription"
	storage "github.com/DenHax/subscription-manager/internal/storage/postgres"
)

type Subscriptions interface {
	CreateSubscrition(serviceName string, price int, userID string, startDate string, endDate *string) (*models.Subscription, error)
	Subscription(id string) (*models.Subscription, error)
	DeleteSubscription(id string) error
	GetAllSubscriptions(userID *string, serviceName *string, limit, offset int) ([]*models.Subscription, int, error)
	UpdateSubscription(id, serviceName *string, price *int, userID *string, startDate *string, endDate *string) (*models.Subscription, error)
	SummarySubscription(startDate, endDate string, userID *string, serviceName *string) (int, error)
}

type Repository struct {
	Subscriptions
}

func NewRepository(s *storage.Storage) *Repository {
	return &Repository{
		Subscriptions: subscription.NewSubStorage(s),
	}
}

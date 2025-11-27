package subscription

import (
	"fmt"
	"log/slog"

	"github.com/DenHax/subscription-manager/internal/domain/models"
	"github.com/DenHax/subscription-manager/internal/repo"
)

type SubService struct {
	repo repo.Subscriptions
}

func NewSubService(repo repo.Subscriptions) *SubService {
	return &SubService{repo: repo}
}

func (s *SubService) CreateSubscrition(serviceName string, price int, userID string, startDate string, endDate *string) (*models.Subscription, error) {
	const op = "service.subscription.CreateSubscrition"

	// Basic validation
	if serviceName == "" {
		return nil, fmt.Errorf("%s: service name cannot be empty", op)
	}
	if price < 0 {
		return nil, fmt.Errorf("%s: price cannot be negative", op)
	}
	if userID == "" {
		return nil, fmt.Errorf("%s: user ID cannot be empty", op)
	}
	if startDate == "" {
		return nil, fmt.Errorf("%s: start date cannot be empty", op)
	}

	// Create the subscription via repository
	sub, err := s.repo.CreateSubscrition(serviceName, price, userID, startDate, endDate)
	if err != nil {
		slog.Error("Failed to create subscription",
			slog.String("operation", op),
			slog.String("user_id", userID),
			slog.String("service_name", serviceName),
			slog.Any("error", err))
		return nil, fmt.Errorf("%s: failed to create subscription: %w", op, err)
	}

	slog.Info("Subscription created",
		slog.String("operation", op),
		slog.Int("subscription_id", sub.Id),
		slog.String("user_id", userID),
		slog.String("service_name", serviceName))

	return sub, nil
}

func (s *SubService) Subscription(id string) (*models.Subscription, error) {
	const op = "service.subscription.Subscription"

	// Validate input
	if id == "" {
		return nil, fmt.Errorf("%s: subscription ID cannot be empty", op)
	}

	// Get subscription via repository
	sub, err := s.repo.Subscription(id)
	if err != nil {
		slog.Error("Failed to get subscription",
			slog.String("operation", op),
			slog.String("subscription_id", id),
			slog.Any("error", err))
		return nil, fmt.Errorf("%s: failed to get subscription: %w", op, err)
	}

	slog.Debug("Subscription retrieved",
		slog.String("operation", op),
		slog.String("subscription_id", id))

	return sub, nil
}

func (s *SubService) DeleteSubscription(id string) error {
	const op = "service.subscription.DeleteSubscription"

	// Validate input
	if id == "" {
		return fmt.Errorf("%s: subscription ID cannot be empty", op)
	}

	// Check if subscription exists before deleting
	_, err := s.repo.Subscription(id)
	if err != nil {
		slog.Warn("Attempt to delete non-existent subscription",
			slog.String("operation", op),
			slog.String("subscription_id", id))
		return fmt.Errorf("%s: subscription not found: %w", op, err)
	}

	// Delete subscription via repository
	err = s.repo.DeleteSubscription(id)
	if err != nil {
		slog.Error("Failed to delete subscription",
			slog.String("operation", op),
			slog.String("subscription_id", id),
			slog.Any("error", err))
		return fmt.Errorf("%s: failed to delete subscription: %w", op, err)
	}

	slog.Info("Subscription deleted",
		slog.String("operation", op),
		slog.String("subscription_id", id))

	return nil
}

func (s *SubService) GetAllSubscriptions(userID *string, serviceName *string, limit, offset int) ([]*models.Subscription, int, error) {
	const op = "service.subscription.GetAllSubscriptions"

	// Validate limit and offset
	if limit < 0 {
		return nil, 0, fmt.Errorf("%s: limit cannot be negative", op)
	}
	if offset < 0 {
		return nil, 0, fmt.Errorf("%s: offset cannot be negative", op)
	}

	// Fetch subscriptions via repository
	subscriptions, totalCount, err := s.repo.GetAllSubscriptions(userID, serviceName, limit, offset)
	if err != nil {
		slog.Error("Failed to get all subscriptions",
			slog.String("operation", op),
			slog.Any("error", err))
		return nil, 0, fmt.Errorf("%s: failed to get subscriptions: %w", op, err)
	}

	slog.Debug("Subscriptions retrieved",
		slog.String("operation", op),
		slog.Int("count", len(subscriptions)),
		slog.Int("total_count", totalCount))

	return subscriptions, totalCount, nil
}

func (s *SubService) UpdateSubscription(id, serviceName *string, price *int, userID *string, startDate *string, endDate *string) (*models.Subscription, error) {
	const op = "service.subscription.UpdateSubscription"

	// Validate subscription ID
	if id == nil || *id == "" {
		return nil, fmt.Errorf("%s: subscription ID cannot be empty", op)
	}

	// Validate that at least one field should be updated
	if serviceName == nil && price == nil && userID == nil && startDate == nil && endDate == nil {
		return nil, fmt.Errorf("%s: at least one field must be provided for update", op)
	}

	// Check if subscription exists before updating
	_, err := s.repo.Subscription(*id)
	if err != nil {
		slog.Warn("Attempt to update non-existent subscription",
			slog.String("operation", op),
			slog.String("subscription_id", *id))
		return nil, fmt.Errorf("%s: subscription not found: %w", op, err)
	}

	// Update subscription via repository
	updatedSub, err := s.repo.UpdateSubscription(id, serviceName, price, userID, startDate, endDate)
	if err != nil {
		slog.Error("Failed to update subscription",
			slog.String("operation", op),
			slog.String("subscription_id", *id),
			slog.Any("error", err))
		return nil, fmt.Errorf("%s: failed to update subscription: %w", op, err)
	}

	slog.Info("Subscription updated",
		slog.String("operation", op),
		slog.Int("subscription_id", updatedSub.Id))

	return updatedSub, nil
}

func (s *SubService) SummarySubscription(startDate, endDate string, userID *string, serviceName *string) (int, error) {
	const op = "service.subscription.SummarySubscription"

	// Validate required dates
	if startDate == "" {
		return 0, fmt.Errorf("%s: start date cannot be empty", op)
	}
	if endDate == "" {
		return 0, fmt.Errorf("%s: end date cannot be empty", op)
	}

	// Calculate summary via repository
	total, err := s.repo.SummarySubscription(startDate, endDate, userID, serviceName)
	if err != nil {
		slog.Error("Failed to calculate subscription summary",
			slog.String("operation", op),
			slog.String("start_date", startDate),
			slog.String("end_date", endDate),
			slog.Any("error", err))
		return 0, fmt.Errorf("%s: failed to calculate summary: %w", op, err)
	}

	slog.Debug("Subscription summary calculated",
		slog.String("operation", op),
		slog.String("start_date", startDate),
		slog.String("end_date", endDate),
		slog.Int("summary", total))

	return total, nil
}

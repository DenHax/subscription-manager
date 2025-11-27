package subscription

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/DenHax/subscription-manager/internal/domain/models"
	storage "github.com/DenHax/subscription-manager/internal/storage/postgres"
)

type SubStore struct {
	storage *storage.Storage
}

func NewSubStorage(s *storage.Storage) *SubStore {
	return &SubStore{storage: s}
}

func (s *SubStore) CreateSubscrition(serviceName string, price int, userID string, startDate string, endDate *string) (*models.Subscription, error) {
	const op = "repo.subscription.CreateSubscrition"

	// Parse start date from format like "07-2025" to time.Time
	startTime, err := time.Parse("01-2006", startDate)
	if err != nil {
		return nil, fmt.Errorf("%s: invalid start date format: %w", op, err)
	}

	var endTime *time.Time
	if endDate != nil {
		parsedEndDate, err := time.Parse("01-2006", *endDate)
		if err != nil {
			return nil, fmt.Errorf("%s: invalid end date format: %w", op, err)
		}
		endTime = &parsedEndDate
	}

	query := `
		INSERT INTO subscriptions.subscriptions (user_id, service_name, price, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING subscription_id, user_id, service_name, price, start_date, end_date
	`

	var sub models.Subscription
	err = s.storage.DB.QueryRow(query, userID, serviceName, price, startTime, endTime).Scan(
		&sub.Id,
		&sub.UserId,
		&sub.ServiceName,
		&sub.Price,
		&sub.StartDate,
		&sub.EndDate,
	)
	if err != nil {
		slog.Error("Failed to create subscription",
			slog.String("operation", op),
			slog.String("user_id", userID),
			slog.String("service_name", serviceName),
			slog.Int("price", price),
			slog.Any("error", err))
		return nil, fmt.Errorf("%s: failed to create subscription: %w", op, err)
	}

	slog.Debug("Subscription created",
		slog.String("operation", op),
		slog.Int("subscription_id", sub.Id))

	return &sub, nil
}

func (s *SubStore) Subscription(id string) (*models.Subscription, error) {
	const op = "repo.subscription.Subscription"

	query := `
		SELECT subscription_id, user_id, service_name, price, start_date, end_date
		FROM subscriptions.subscriptions
		WHERE subscription_id = $1
	`

	var sub models.Subscription
	err := s.storage.DB.QueryRow(query, id).Scan(
		&sub.Id,
		&sub.UserId,
		&sub.ServiceName,
		&sub.Price,
		&sub.StartDate,
		&sub.EndDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s: subscription not found", op)
		}
		slog.Error("Failed to get subscription",
			slog.String("operation", op),
			slog.String("subscription_id", id),
			slog.Any("error", err))
		return nil, fmt.Errorf("%s: failed to get subscription: %w", op, err)
	}

	slog.Debug("Subscription retrieved",
		slog.String("operation", op),
		slog.Int("subscription_id", sub.Id))

	return &sub, nil
}

func (s *SubStore) DeleteSubscription(id string) error {
	const op = "repo.subscription.DeleteSubscription"

	query := `DELETE FROM subscriptions.subscriptions WHERE subscription_id = $1`

	result, err := s.storage.DB.Exec(query, id)
	if err != nil {
		slog.Error("Failed to delete subscription",
			slog.String("operation", op),
			slog.String("subscription_id", id),
			slog.Any("error", err))
		return fmt.Errorf("%s: failed to delete subscription: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("Failed to get affected rows when deleting subscription",
			slog.String("operation", op),
			slog.String("subscription_id", id),
			slog.Any("error", err))
		return fmt.Errorf("%s: failed to get affected rows: %w", op, err)
	}

	if rowsAffected == 0 {
		slog.Warn("Attempt to delete non-existent subscription",
			slog.String("operation", op),
			slog.String("subscription_id", id))
		return fmt.Errorf("%s: subscription not found", op)
	}

	slog.Debug("Subscription deleted",
		slog.String("operation", op),
		slog.String("subscription_id", id))

	return nil
}

func (s *SubStore) GetAllSubscriptions(userID *string, serviceName *string, limit, offset int) ([]*models.Subscription, int, error) {
	const op = "repo.subscription.GetAllSubscriptions"

	// Count query to get total number of subscriptions matching the filters
	countQuery := `SELECT COUNT(*) FROM subscriptions.subscriptions WHERE true`
	countArgs := []interface{}{}
	countArgCount := 0

	// Main query with filters
	query := `SELECT subscription_id, user_id, service_name, price, start_date, end_date FROM subscriptions.subscriptions WHERE true`
	args := []interface{}{}
	argCount := 0

	if userID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argCount+1)
		args = append(args, *userID)
		argCount++

		countQuery += fmt.Sprintf(" AND user_id = $%d", countArgCount+1)
		countArgs = append(countArgs, *userID)
		countArgCount++
	}

	if serviceName != nil {
		query += fmt.Sprintf(" AND service_name = $%d", argCount+1)
		args = append(args, *serviceName)
		argCount++

		countQuery += fmt.Sprintf(" AND service_name = $%d", countArgCount+1)
		countArgs = append(countArgs, *serviceName)
		countArgCount++
	}

	query += fmt.Sprintf(" ORDER BY subscription_id LIMIT $%d OFFSET $%d", argCount+1, argCount+2)
	args = append(args, limit, offset)

	// Execute count query
	var totalCount int
	err := s.storage.DB.QueryRow(countQuery, countArgs...).Scan(&totalCount)
	if err != nil {
		slog.Error("Failed to count subscriptions",
			slog.String("operation", op),
			slog.Any("error", err))
		return nil, 0, fmt.Errorf("%s: failed to count subscriptions: %w", op, err)
	}

	// Execute main query
	rows, err := s.storage.DB.Query(query, args...)
	if err != nil {
		slog.Error("Failed to query subscriptions",
			slog.String("operation", op),
			slog.Any("error", err))
		return nil, 0, fmt.Errorf("%s: failed to query subscriptions: %w", op, err)
	}
	defer rows.Close()

	var subscriptions []*models.Subscription
	for rows.Next() {
		var sub models.Subscription
		err := rows.Scan(
			&sub.Id,
			&sub.UserId,
			&sub.ServiceName,
			&sub.Price,
			&sub.StartDate,
			&sub.EndDate,
		)
		if err != nil {
			slog.Error("Failed to scan subscription",
				slog.String("operation", op),
				slog.Any("error", err))
			return nil, 0, fmt.Errorf("%s: failed to scan subscription: %w", op, err)
		}
		subscriptions = append(subscriptions, &sub)
	}

	slog.Debug("Fetched subscriptions",
		slog.String("operation", op),
		slog.Int("count", len(subscriptions)),
		slog.Int("total_count", totalCount))

	return subscriptions, totalCount, nil
}

func (s *SubStore) UpdateSubscription(id, serviceName *string, price *int, userID *string, startDate *string, endDate *string) (*models.Subscription, error) {
	const op = "repo.subscription.UpdateSubscription"

	// Build the dynamic query and arguments
	query := "UPDATE subscriptions.subscriptions SET "
	args := []interface{}{}
	argCount := 0

	if serviceName != nil {
		query += fmt.Sprintf("service_name = $%d, ", argCount+1)
		args = append(args, *serviceName)
		argCount++
	}
	if price != nil {
		query += fmt.Sprintf("price = $%d, ", argCount+1)
		args = append(args, *price)
		argCount++
	}
	if userID != nil {
		query += fmt.Sprintf("user_id = $%d, ", argCount+1)
		args = append(args, *userID)
		argCount++
	}
	if startDate != nil {
		parsedStartDate, err := time.Parse("01-2006", *startDate)
		if err != nil {
			return nil, fmt.Errorf("%s: invalid start date format: %w", op, err)
		}
		query += fmt.Sprintf("start_date = $%d, ", argCount+1)
		args = append(args, parsedStartDate)
		argCount++
	}
	if endDate != nil {
		var parsedEndDate *time.Time
		if *endDate != "" {
			t, err := time.Parse("01-2006", *endDate)
			if err != nil {
				return nil, fmt.Errorf("%s: invalid end date format: %w", op, err)
			}
			parsedEndDate = &t
		} else {
			parsedEndDate = nil
		}
		query += fmt.Sprintf("end_date = $%d, ", argCount+1)
		args = append(args, parsedEndDate)
		argCount++
	}

	// Remove trailing comma and space
	if len(args) > 0 {
		query = query[:len(query)-2]
	} else {
		// If no fields to update, return early
		sub, err := s.Subscription(id)
		if err != nil {
			return nil, fmt.Errorf("%s: subscription not found: %w", op, err)
		}
		return sub, nil
	}

	// Add WHERE clause
	query += fmt.Sprintf(" WHERE subscription_id = $%d RETURNING subscription_id, user_id, service_name, price, start_date, end_date", argCount+1)
	args = append(args, id)

	var updatedSub models.Subscription
	err := s.storage.DB.QueryRow(query, args...).Scan(
		&updatedSub.Id,
		&updatedSub.UserId,
		&updatedSub.ServiceName,
		&updatedSub.Price,
		&updatedSub.StartDate,
		&updatedSub.EndDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Warn("Attempt to update non-existent subscription",
				slog.String("operation", op),
				slog.String("subscription_id", *id))
			return nil, fmt.Errorf("%s: subscription not found", op)
		}
		slog.Error("Failed to update subscription",
			slog.String("operation", op),
			slog.String("subscription_id", *id),
			slog.Any("error", err))
		return nil, fmt.Errorf("%s: failed to update subscription: %w", op, err)
	}

	slog.Debug("Subscription updated",
		slog.String("operation", op),
		slog.Int("subscription_id", updatedSub.Id))

	return &updatedSub, nil
}

func (s *SubStore) SummarySubscription(startDate, endDate string, userID *string, serviceName *string) (int, error) {
	const op = "repo.subscription.SummarySubscription"

	// Parse start and end dates
	startTime, err := time.Parse("01-2006", startDate)
	if err != nil {
		return 0, fmt.Errorf("%s: invalid start date format: %w", op, err)
	}
	endTime, err := time.Parse("01-2006", endDate)
	if err != nil {
		return 0, fmt.Errorf("%s: invalid end date format: %w", op, err)
	}

	query := `SELECT COALESCE(SUM(price), 0) FROM subscriptions.subscriptions WHERE start_date <= $1 AND ($2 IS NULL OR end_date >= $3 OR end_date IS NULL)`
	args := []interface{}{endTime, endTime, startTime}
	argCount := 3

	if userID != nil {
		argCount++
		query += fmt.Sprintf(" AND user_id = $%d", argCount)
		args = append(args, *userID)
	}

	if serviceName != nil {
		argCount++
		query += fmt.Sprintf(" AND service_name = $%d", argCount)
		args = append(args, *serviceName)
	}

	var total int
	err = s.storage.DB.QueryRow(query, args...).Scan(&total)
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

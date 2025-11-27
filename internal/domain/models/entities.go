package models

import "time"

type Subscription struct {
	Id          int       `json:"subscription_id" db:"subscription_id"`
	UserId      string    `json:"user_id" db:"user_id"`
	ServiceName string    `json:"service_name" db:"service_name"`
	Price       int       `json:"price" db:"price"`
	StartDate   time.Time `json:"start_date" db:"start_date"`
	EndDate     time.Time `json:"end_date" db:"end_date"`
}

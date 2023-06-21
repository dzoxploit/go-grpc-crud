package entities

import "time"

type Contact struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Gender     string   `json:"gender"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	Created_at time.Time 
	Updated_at time.Time
}
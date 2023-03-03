package models

import (
	"time"

	"gorm.io/gorm"
)

type Routine struct {
	gorm.Model
	Name   string `json:"name"`
	Days   string `json:"days"`
	UserID string `json:"-"`
}

type RoutineRecord struct {
	ID        uint      `json:"id"`
	Timestamp time.Time `json:"timestamp"`
}

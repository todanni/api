package models

import "gorm.io/gorm"

type Routine struct {
	gorm.Model
	Name   string `json:"name"`
	Days   string `json:"days"`
	UserID string `json:"-"`
}

type RoutineRecord struct {
	Timestamp string `json:"timestamp"`
	ID        uint   `gorm:"primarykey"`
}

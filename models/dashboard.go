package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Dashboard struct {
	ID        uuid.UUID      `gorm:"primarykey" json:"id"`
	Status    Status         `json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Members   []User         `json:"members" gorm:"many2many:user_dashboards;"`
}

type Status string

const (
	PendingStatus  Status = "PENDING"
	AcceptedStatus Status = "ACCEPTED"
	RejectedStatus Status = "REJECTED"
)

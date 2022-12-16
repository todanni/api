package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID          string         `json:"id" gorm:"primarykey"`
	DisplayName string         `json:"display_name"`
	Email       string         `json:"email"`
	ProfilePic  string         `json:"profile_pic"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index" `

	Dashboards []Dashboard `json:"-" gorm:"many2many:user_dashboards;"`
	Projects   []Project   `json:"-" gorm:"many2many:user_projects;"`
}

package models

import "gorm.io/gorm"

type Project struct {
	gorm.Model
	Name    string `json:"name"`
	Owner   string `json:"owner"`
	Members []User `json:"members" gorm:"many2many:user_projects;"`
}

// TODO: Make status an enum

type ProjectInvite struct {
	ID        uint   `json:"id" gorm:"primarykey"`
	ProjectID uint   `json:"project_id"`
	UserID    string `json:"user_id"`
	Status    string `json:"status"`
}

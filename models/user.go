package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Email      string      `json:"email"`
	FirstName  string      `json:"first_name"`
	LastName   string      `json:"last_name"`
	ProfilePic string      `json:"profile_pic"`
	Dashboards []Dashboard `json:"-" gorm:"many2many:user_dashboards;"`
	Projects   []Project   `json:"-" gorm:"many2many:user_projects;"`
}

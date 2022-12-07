package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	FirstName  string      `json:"first_name"`
	LastName   string      `json:"last_name"`
	Email      string      `json:"email"`
	ProfilePic string      `json:"profile_pic"`
	Dashboards []Dashboard `json:"-" gorm:"many2many:user_dashboards;"`
	Projects   []Project   `json:"-" gorm:"many2many:user_projects;"`
}

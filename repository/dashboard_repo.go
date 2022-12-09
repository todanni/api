package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/todanni/api/models"
)

type DashboardRepository interface {
	CreateDashboard(dashboard models.Dashboard) (models.Dashboard, error)
	DeleteDashboard(id uuid.UUID) (models.Dashboard, error)
	ListDashboardsByUser(userID uint) ([]models.Dashboard, error)
}

type dashboardRepo struct {
	db *gorm.DB
}

func (d dashboardRepo) CreateDashboard(dashboard models.Dashboard) (models.Dashboard, error) {
	//TODO implement me
	panic("implement me")
}

func (d dashboardRepo) DeleteDashboard(id uuid.UUID) (models.Dashboard, error) {
	//TODO implement me
	panic("implement me")
}

func (d dashboardRepo) ListDashboardsByUser(userID uint) ([]models.Dashboard, error) {
	var user models.User
	result := d.db.Model(&models.User{}).Preload("Dashboard.Members").First(&user, userID)
	return user.Dashboards, result.Error
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepo{
		db: db,
	}
}

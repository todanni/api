package repository

type DashboardRepository interface {
}

type dashboardRepo struct {
}

func NewDashboardRepository() DashboardRepository {
	return &dashboardRepo{}
}

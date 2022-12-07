package repository

type ProjectRepository interface {
}

type projectRepo struct {
}

func NewProjectRepository() ProjectRepository {
	return &projectRepo{}
}

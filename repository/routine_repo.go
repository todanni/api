package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/todanni/api/models"
)

type RoutineRepository interface {
	CreateRoutine(routine models.Routine) (models.Routine, error)
	UpdateRoutine(routine models.Routine) (models.Routine, error)
	ListRoutinesByUser(userID string) ([]models.Routine, error)
	GetRoutineByID(routineID string) (models.Routine, error)
	DeleteRoutine(routineID string) error
	CreateRoutineRecord(record models.RoutineRecord) (models.RoutineRecord, error)
	DeleteRoutineRecord(record models.RoutineRecord) error
}

type routineRepo struct {
	db *gorm.DB
}

func (r *routineRepo) CreateRoutineRecord(record models.RoutineRecord) (models.RoutineRecord, error) {
	result := r.db.Create(&record)
	return record, result.Error
}

func (r *routineRepo) DeleteRoutineRecord(record models.RoutineRecord) error {
	result := r.db.Delete(record)
	return result.Error
}

func (r *routineRepo) GetRoutineByID(routineID string) (models.Routine, error) {
	var routine models.Routine
	result := r.db.First(&routine, routineID)
	return routine, result.Error
}

func (r *routineRepo) DeleteRoutine(routineID string) error {
	result := r.db.Delete(&models.Routine{}, routineID)
	return result.Error
}

func (r *routineRepo) ListRoutinesByUser(userID string) ([]models.Routine, error) {
	var routines []models.Routine
	result := r.db.Where("user_id = ?", userID).Find(&routines)
	return routines, result.Error
}

func (r *routineRepo) CreateRoutine(routine models.Routine) (models.Routine, error) {
	result := r.db.Create(&routine)
	return routine, result.Error
}

func (r *routineRepo) UpdateRoutine(routine models.Routine) (models.Routine, error) {
	result := r.db.Model(&routine).Clauses(clause.Returning{}).Updates(routine)
	return routine, result.Error
}

func NewRoutineRepository(db *gorm.DB) RoutineRepository {
	return &routineRepo{
		db: db,
	}
}

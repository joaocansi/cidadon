package repository

import (
	"cidadon/internal/adapters/persistence/postgres/model"
	"cidadon/internal/domain/entity"
	"cidadon/internal/domain/repository"
	"cidadon/internal/platform/database"
	"context"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

type DemandRepositoryImpl struct {
	*database.BaseRepository
}

func NewDemandRepository(baseRepository *database.BaseRepository) *DemandRepositoryImpl {
	return &DemandRepositoryImpl{
		BaseRepository: baseRepository,
	}
}

func (r *DemandRepositoryImpl) Create(ctx context.Context, demand repository.CreateDemandData) (*entity.Demand, error) {
	demandModel := model.Demand{
		Protocol:     demand.Protocol,
		CitizenID:    demand.CitizenID,
		Title:        demand.Title,
		Description:  demand.Description,
		Category:     demand.Category,
		Street:       demand.Street,
		Number:       demand.Number,
		Neighborhood: demand.Neighborhood,
		City:         demand.City,
		State:        demand.State,
		Latitude:     demand.Latitude, Longitude: demand.Longitude,
		Images:           func() json.RawMessage { b, _ := json.Marshal(demand.Images); return b }(),
		DirectedOfficeID: demand.DirectedOfficeID,
		Status:           string(demand.Status),
	}

	if err := r.GetDB(ctx).Create(&demandModel).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, repository.NewDBError(repository.DBErrorConflict, err)
		}
		return nil, repository.NewDBError(repository.DBErrorInternal, err)
	}

	return demandModel.ToDomain(), nil
}

func (r *DemandRepositoryImpl) AssignOffices(ctx context.Context, demandID uint, officeIDs []uint) error {
	if len(officeIDs) == 0 {
		return nil
	}
	assignments := make([]model.DemandAssignment, 0, len(officeIDs))
	for _, officeID := range officeIDs {
		assignments = append(assignments, model.DemandAssignment{DemandID: demandID, OfficeID: officeID})
	}
	if err := r.GetDB(ctx).Create(&assignments).Error; err != nil {
		return repository.NewDBError(repository.DBErrorInternal, err)
	}
	return nil
}

func (r *DemandRepositoryImpl) ListOfficeDemands(ctx context.Context, officeID uint, filters repository.DemandFilters) ([]entity.Demand, error) {
	var models []model.Demand
	query := r.GetDB(ctx).Joins("JOIN demand_assignments ON demand_assignments.demand_id = demands.id").Where("demand_assignments.office_id = ?", officeID).Order("demands.created_at desc").Limit(100)
	if filters.Status != "" {
		query = query.Where("demands.status = ?", filters.Status)
	}
	if err := query.Find(&models).Error; err != nil {
		return nil, repository.NewDBError(repository.DBErrorInternal, err)
	}
	result := make([]entity.Demand, 0, len(models))
	for _, item := range models {
		result = append(result, *item.ToDomain())
	}
	return result, nil
}

func (r *DemandRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.Demand, error) {
	var demandModel model.Demand
	if err := r.GetDB(ctx).First(&demandModel, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.NewDBError(repository.DBErrorNotFound, err)
		}
		return nil, repository.NewDBError(repository.DBErrorInternal, err)
	}

	return demandModel.ToDomain(), nil
}

func (r *DemandRepositoryImpl) List(ctx context.Context, filters repository.DemandFilters) ([]entity.Demand, error) {
	var demandModels []model.Demand
	query := r.GetDB(ctx).Order("created_at desc").Limit(50)

	if filters.City != "" {
		query = query.Where("city = ?", filters.City)
	}
	if filters.State != "" {
		query = query.Where("state = ?", filters.State)
	}
	if filters.Neighborhood != "" {
		query = query.Where("neighborhood = ?", filters.Neighborhood)
	}
	if filters.Category != "" {
		query = query.Where("category = ?", filters.Category)
	}
	if filters.Status != "" {
		query = query.Where("status = ?", string(filters.Status))
	}

	if err := query.Find(&demandModels).Error; err != nil {
		return nil, repository.NewDBError(repository.DBErrorInternal, err)
	}

	demands := make([]entity.Demand, 0, len(demandModels))
	for _, demandModel := range demandModels {
		demands = append(demands, *demandModel.ToDomain())
	}

	return demands, nil
}

func (r *DemandRepositoryImpl) ListByCitizen(ctx context.Context, citizenID uint) ([]entity.Demand, error) {
	var models []model.Demand
	if err := r.GetDB(ctx).Where("citizen_id = ?", citizenID).Order("created_at desc").Find(&models).Error; err != nil {
		return nil, repository.NewDBError(repository.DBErrorInternal, err)
	}
	result := make([]entity.Demand, 0, len(models))
	for _, item := range models {
		result = append(result, *item.ToDomain())
	}
	return result, nil
}

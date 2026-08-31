package repository

import (
	"cidadon/internal/adapters/persistence/postgres/model"
	"cidadon/internal/domain/entity"
	"cidadon/internal/domain/repository"
	"cidadon/internal/platform/database"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OfficeRepositoryImpl struct {
	*database.BaseRepository
}

func NewOfficeRepository(baseRepository *database.BaseRepository) *OfficeRepositoryImpl {
	return &OfficeRepositoryImpl{
		BaseRepository: baseRepository,
	}
}

func (o *OfficeRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.Office, error) {
	var officeModel model.Office
	if err := o.GetDB(ctx).Preload("Councillor.User").First(&officeModel, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.NewDBError(repository.DBErrorNotFound, err)
		}
		return nil, repository.NewDBError(repository.DBErrorInternal, err)
	}
	return officeModel.ToDomain(), nil
}

func (o *OfficeRepositoryImpl) Search(ctx context.Context, queryText, city, state string) ([]entity.Office, error) {
	var offices []model.Office
	query := o.GetDB(ctx).Preload("Councillor.User").Joins("JOIN councillors ON councillors.user_id = offices.councillor_id").Joins("JOIN users ON users.id = councillors.user_id")
	if queryText != "" {
		query = query.Where("users.name ILIKE ? OR councillors.party ILIKE ?", "%"+queryText+"%", "%"+queryText+"%")
	}
	if city != "" {
		query = query.Where("councillors.city = ?", city)
	}
	if state != "" {
		query = query.Where("councillors.state = ?", state)
	}
	if err := query.Order("users.name asc").Limit(30).Find(&offices).Error; err != nil {
		return nil, repository.NewDBError(repository.DBErrorInternal, err)
	}
	result := make([]entity.Office, 0, len(offices))
	for _, office := range offices {
		result = append(result, *office.ToDomain())
	}
	return result, nil
}

func (o *OfficeRepositoryImpl) FindByCouncillorID(ctx context.Context, councillorID uint) (*entity.Office, error) {
	var officeModel model.Office
	if err := o.GetDB(ctx).First(&officeModel, "councillor_id = ?", councillorID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.NewDBError(repository.DBErrorNotFound, err)
		}
		return nil, repository.NewDBError(repository.DBErrorInternal, err)
	}
	return officeModel.ToDomain(), nil
}

func (o *OfficeRepositoryImpl) ListByCityState(ctx context.Context, city, state string) ([]entity.Office, error) {
	var offices []model.Office
	if err := o.GetDB(ctx).Preload("Councillor.User").Joins("JOIN councillors ON councillors.user_id = offices.councillor_id").Where("councillors.city = ? AND councillors.state = ?", city, state).Find(&offices).Error; err != nil {
		return nil, repository.NewDBError(repository.DBErrorInternal, err)
	}
	result := make([]entity.Office, 0, len(offices))
	for _, office := range offices {
		result = append(result, *office.ToDomain())
	}
	return result, nil
}

func (o *OfficeRepositoryImpl) ListActiveOfficeIDsNear(ctx context.Context, latitude, longitude, radiusKM float64, limit int) ([]uint, error) {
	type officeMatch struct {
		OfficeID uint
	}
	matches := make([]officeMatch, 0)
	query := `
		WITH office_activity AS (
			SELECT
				da.office_id,
				COUNT(*) AS activity_count,
				AVG(d.latitude) AS latitude,
				AVG(d.longitude) AS longitude
			FROM demand_assignments da
			JOIN demands d ON d.id = da.demand_id
			WHERE d.deleted_at IS NULL AND d.status IN ('in_progress', 'resolved')
			GROUP BY da.office_id
		), ranked AS (
			SELECT
				office_id,
				activity_count,
				6371 * 2 * ASIN(LEAST(1, SQRT(
					POWER(SIN(RADIANS((? - latitude) / 2)), 2) +
					COS(RADIANS(?)) * COS(RADIANS(latitude)) *
					POWER(SIN(RADIANS((? - longitude) / 2)), 2)
				))) AS distance_km
			FROM office_activity
		)
		SELECT office_id
		FROM ranked
		WHERE distance_km <= ?
		ORDER BY activity_count DESC, distance_km ASC
		LIMIT ?`
	if err := o.GetDB(ctx).Raw(query, latitude, latitude, longitude, radiusKM, limit).Scan(&matches).Error; err != nil {
		return nil, repository.NewDBError(repository.DBErrorInternal, err)
	}
	ids := make([]uint, 0, len(matches))
	for _, match := range matches {
		ids = append(ids, match.OfficeID)
	}
	return ids, nil
}

func (o *OfficeRepositoryImpl) Create(ctx context.Context, data repository.CreateOfficeData) (*entity.Office, error) {
	officeModel := model.Office{
		CouncillorID:   data.CouncillorID,
		Description:    data.Description,
		Contacts:       o.toRawMessage(data.Contacts),
		SocialNetworks: o.toRawMessage(data.SocialNetworks),
	}
	if err := o.GetDB(ctx).Create(&officeModel).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, repository.NewDBError(repository.DBErrorConflict, err)
		}
		return nil, repository.NewDBError(repository.DBErrorInternal, err)
	}
	return officeModel.ToDomain(), nil
}

func (o *OfficeRepositoryImpl) UpdateByCouncillorID(ctx context.Context, councillorID uint, data repository.UpdateOfficeData) (*entity.Office, error) {
	officeModel := model.Office{
		Description:    data.Description,
		SocialNetworks: o.toRawMessage(data.SocialNetworks),
		Contacts:       o.toRawMessage(data.Contacts),
	}

	result := o.GetDB(ctx).Model(&officeModel).
		Clauses(clause.Returning{}).
		Where("councillor_id = ?", councillorID).
		Updates(&officeModel)

	if result.Error != nil {
		return nil, repository.NewDBError(repository.DBErrorInternal, result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, repository.NewDBError(repository.DBErrorNotFound, fmt.Errorf("office not found"))
	}

	return officeModel.ToDomain(), nil
}

func (o *OfficeRepositoryImpl) toRawMessage(v any) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		return json.RawMessage("[]")
	}
	return b
}

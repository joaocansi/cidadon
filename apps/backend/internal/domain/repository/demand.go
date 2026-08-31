package repository

import (
	"cidadon/internal/domain/entity"
	"context"
)

type CreateDemandData struct {
	Protocol         string
	CitizenID        uint
	Title            string
	Description      string
	Category         string
	Street           string
	Number           string
	Neighborhood     string
	City             string
	State            string
	Latitude         float64
	Longitude        float64
	Images           []string
	DirectedOfficeID *uint
	Status           entity.DemandStatus
	OfficeID         uint
}

type DemandFilters struct {
	City         string
	State        string
	Neighborhood string
	Category     string
	Status       entity.DemandStatus
}

type DemandRepository interface {
	Create(ctx context.Context, demand CreateDemandData) (*entity.Demand, error)
	FindByID(ctx context.Context, id uint) (*entity.Demand, error)
	List(ctx context.Context, filters DemandFilters) ([]entity.Demand, error)
	ListByCitizen(ctx context.Context, citizenID uint) ([]entity.Demand, error)
	AssignOffices(ctx context.Context, demandID uint, officeIDs []uint) error
	ListOfficeDemands(ctx context.Context, officeID uint, filters DemandFilters) ([]entity.Demand, error)
	UpdateStatus(ctx context.Context, demandID, officeID uint, status entity.DemandStatus) (*entity.Demand, error)
}

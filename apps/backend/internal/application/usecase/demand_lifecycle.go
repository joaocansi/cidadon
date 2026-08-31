package usecase

import (
	"cidadon/internal/adapters/persistence/postgres/model"
	"cidadon/internal/domain/entity"
	"context"
	"encoding/json"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// DemandLifecycleService owns background state transitions for demands.
type DemandLifecycleService struct{ db *gorm.DB }

func NewDemandLifecycleService(db *gorm.DB) *DemandLifecycleService {
	return &DemandLifecycleService{db: db}
}

func (s *DemandLifecycleService) CompleteExpiredConfirmations(ctx context.Context, now time.Time) error {
	var ids []uint
	if err := s.db.WithContext(ctx).Model(&model.Demand{}).Where("status = ? AND confirmation_due_at <= ?", entity.DemandStatusAwaitingConfirmation, now).Pluck("id", &ids).Error; err != nil {
		return err
	}
	for _, id := range ids {
		if err := s.completeExpired(ctx, id, now); err != nil {
			return err
		}
	}
	return nil
}

func (s *DemandLifecycleService) completeExpired(ctx context.Context, id uint, now time.Time) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var demand model.Demand
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).First(&demand, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil
			}
			return err
		}
		if demand.Status != string(entity.DemandStatusAwaitingConfirmation) || demand.ConfirmationDueAt == nil || demand.ConfirmationDueAt.After(now) {
			return nil
		}
		if err := tx.Model(&demand).Updates(map[string]any{"status": entity.DemandStatusCompleted, "confirmation_due_at": nil}).Error; err != nil {
			return err
		}
		metadata, _ := json.Marshal(map[string]string{"reason": "confirmation_expired"})
		if err := tx.Create(&model.DemandEvent{DemandID: demand.ID, Type: "automatically_completed", Metadata: metadata}).Error; err != nil {
			return err
		}
		return s.notifyCompletion(tx, demand)
	})
}

func (s *DemandLifecycleService) notifyCompletion(tx *gorm.DB, demand model.Demand) error {
	userIDs := []uint{demand.CitizenID}
	if demand.ResponsibleOfficeID != nil {
		var councillorID uint
		if err := tx.Table("offices").Select("councillor_id").Where("id = ?", *demand.ResponsibleOfficeID).Scan(&councillorID).Error; err != nil {
			return err
		}
		var memberIDs []uint
		if err := tx.Table("office_members").Select("user_id").Where("office_id = ?", *demand.ResponsibleOfficeID).Scan(&memberIDs).Error; err != nil {
			return err
		}
		userIDs = append(userIDs, councillorID)
		userIDs = append(userIDs, memberIDs...)
	}
	seen := map[uint]bool{}
	notifications := make([]model.Notification, 0, len(userIDs))
	for _, userID := range userIDs {
		if userID != 0 && !seen[userID] {
			seen[userID] = true
			notifications = append(notifications, model.Notification{UserID: userID, DemandID: demand.ID, Type: "demand_automatically_completed"})
		}
	}
	if len(notifications) == 0 {
		return nil
	}
	return tx.Create(&notifications).Error
}

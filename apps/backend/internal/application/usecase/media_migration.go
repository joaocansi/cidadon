package usecase

import (
	"cidadon/internal/adapters/persistence/postgres/model"
	"context"
	"encoding/json"
	"strings"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// MediaMigrationService incrementally replaces legacy database data URLs with
// storage URLs. A value is changed only after its file is saved successfully,
// therefore retries are safe and naturally idempotent.
type MediaMigrationService struct {
	db    *gorm.DB
	media *MediaService
	log   *zap.SugaredLogger
}

func NewMediaMigrationService(db *gorm.DB, media *MediaService, log *zap.SugaredLogger) *MediaMigrationService {
	return &MediaMigrationService{db: db, media: media, log: log.Named("MediaMigration")}
}

func (s *MediaMigrationService) MigrateLegacyMedia(ctx context.Context) error {
	if err := s.migrateDemandImages(ctx); err != nil {
		return err
	}
	if err := s.migrateCommentImages(ctx); err != nil {
		return err
	}
	if err := s.migrateEventImages(ctx); err != nil {
		return err
	}
	if err := s.migrateAvatars(ctx); err != nil {
		return err
	}
	return nil
}

func (s *MediaMigrationService) migrateDemandImages(ctx context.Context) error {
	var demands []model.Demand
	if err := s.db.WithContext(ctx).Where("images::text LIKE ?", "%data:image/%").Limit(100).Find(&demands).Error; err != nil {
		return err
	}
	for _, demand := range demands {
		var images []string
		if json.Unmarshal(demand.Images, &images) != nil {
			continue
		}
		updated, changed := s.migrateURLs(ctx, "demands", images)
		if !changed {
			continue
		}
		bytes, _ := json.Marshal(updated)
		if err := s.db.WithContext(ctx).Model(&model.Demand{}).Where("id = ?", demand.ID).Update("images", bytes).Error; err != nil {
			s.log.Warnw("failed to migrate demand media", "demandID", demand.ID, "error", err)
		}
	}
	return nil
}

func (s *MediaMigrationService) migrateCommentImages(ctx context.Context) error {
	var comments []model.DemandComment
	if err := s.db.WithContext(ctx).Where("images::text LIKE ?", "%data:image/%").Limit(100).Find(&comments).Error; err != nil {
		return err
	}
	for _, comment := range comments {
		var images []string
		if json.Unmarshal(comment.Images, &images) != nil {
			continue
		}
		updated, changed := s.migrateURLs(ctx, "comments", images)
		if !changed {
			continue
		}
		bytes, _ := json.Marshal(updated)
		if err := s.db.WithContext(ctx).Model(&model.DemandComment{}).Where("id = ?", comment.ID).Update("images", bytes).Error; err != nil {
			s.log.Warnw("failed to migrate comment media", "commentID", comment.ID, "error", err)
		}
	}
	return nil
}

func (s *MediaMigrationService) migrateEventImages(ctx context.Context) error {
	var events []model.DemandEvent
	if err := s.db.WithContext(ctx).Where("metadata::text LIKE ?", "%data:image/%").Limit(100).Find(&events).Error; err != nil {
		return err
	}
	for _, event := range events {
		metadata := map[string]any{}
		if json.Unmarshal(event.Metadata, &metadata) != nil {
			continue
		}
		rawImages, ok := metadata["images"].([]any)
		if !ok {
			continue
		}
		images := make([]string, 0, len(rawImages))
		for _, raw := range rawImages {
			if image, ok := raw.(string); ok {
				images = append(images, image)
			}
		}
		updated, changed := s.migrateURLs(ctx, "timeline", images)
		if !changed {
			continue
		}
		metadata["images"] = updated
		bytes, _ := json.Marshal(metadata)
		if err := s.db.WithContext(ctx).Model(&model.DemandEvent{}).Where("id = ?", event.ID).Update("metadata", bytes).Error; err != nil {
			s.log.Warnw("failed to migrate timeline media", "eventID", event.ID, "error", err)
		}
	}
	return nil
}

func (s *MediaMigrationService) migrateAvatars(ctx context.Context) error {
	var councillors []model.Councillor
	if err := s.db.WithContext(ctx).Where("image_url LIKE ?", "data:image/%").Limit(100).Find(&councillors).Error; err != nil {
		return err
	}
	for _, councillor := range councillors {
		item, err := s.media.StoreDataURL(ctx, "avatars/councillors", councillor.ImageURL)
		if err != nil {
			s.log.Warnw("failed to migrate councillor avatar", "userID", councillor.UserID, "error", err)
			continue
		}
		if err := s.db.WithContext(ctx).Model(&model.Councillor{}).Where("user_id = ?", councillor.UserID).Update("image_url", item.URL).Error; err != nil {
			s.log.Warnw("failed to persist councillor avatar", "userID", councillor.UserID, "error", err)
		}
	}
	var members []model.OfficeMember
	if err := s.db.WithContext(ctx).Where("image_url LIKE ?", "data:image/%").Limit(100).Find(&members).Error; err != nil {
		return err
	}
	for _, member := range members {
		item, err := s.media.StoreDataURL(ctx, "avatars/members", member.ImageURL)
		if err != nil {
			s.log.Warnw("failed to migrate office member avatar", "userID", member.UserID, "error", err)
			continue
		}
		if err := s.db.WithContext(ctx).Model(&model.OfficeMember{}).Where("user_id = ?", member.UserID).Update("image_url", item.URL).Error; err != nil {
			s.log.Warnw("failed to persist office member avatar", "userID", member.UserID, "error", err)
		}
	}
	return nil
}

func (s *MediaMigrationService) migrateURLs(ctx context.Context, prefix string, images []string) ([]string, bool) {
	changed := false
	updated := append([]string(nil), images...)
	for index, image := range images {
		if !strings.HasPrefix(image, "data:image/") {
			continue
		}
		item, err := s.media.StoreDataURL(ctx, prefix, image)
		if err != nil {
			s.log.Warnw("failed to migrate media value", "prefix", prefix, "error", err)
			continue
		}
		updated[index] = item.URL
		changed = true
	}
	return updated, changed
}

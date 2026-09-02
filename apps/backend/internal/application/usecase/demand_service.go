package usecase

import (
	"cidadon/internal/adapters/persistence/postgres/model"
	service "cidadon/internal/application/contract"
	"cidadon/internal/domain/entity"
	"cidadon/internal/domain/repository"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/url"
	"strings"
	"time"

	"go.uber.org/zap"
)

type DemandServiceImpl struct {
	demandRepo repository.DemandRepository
	officeRepo repository.OfficeRepository
	logger     *zap.SugaredLogger
	db         *gorm.DB
}

func NewDemandService(demandRepo repository.DemandRepository, officeRepo repository.OfficeRepository, db *gorm.DB, logger *zap.SugaredLogger) *DemandServiceImpl {
	return &DemandServiceImpl{
		demandRepo: demandRepo,
		officeRepo: officeRepo,
		logger:     logger.Named("DemandService"),
		db:         db,
	}
}

func (s *DemandServiceImpl) Claim(ctx context.Context, id, officeID, userID uint, input service.DemandTimelineInput) (*service.DemandOutput, error) {
	return s.transition(ctx, id, officeID, userID, entity.DemandStatusRegistered, entity.DemandStatusReview, "claimed", input)
}
func (s *DemandServiceImpl) Start(ctx context.Context, id, officeID, userID uint, input service.DemandTimelineInput) (*service.DemandOutput, error) {
	return s.transition(ctx, id, officeID, userID, entity.DemandStatusReview, entity.DemandStatusInProgress, "execution_started", input)
}
func (s *DemandServiceImpl) RequestConfirmation(ctx context.Context, id, officeID, userID uint, input service.DemandTimelineInput) (*service.DemandOutput, error) {
	d, e := s.transition(ctx, id, officeID, userID, entity.DemandStatusInProgress, entity.DemandStatusAwaitingConfirmation, "confirmation_requested", input)
	return d, e
}
func (s *DemandServiceImpl) Confirm(ctx context.Context, id, citizenID uint) (*service.DemandOutput, error) {
	var m model.Demand
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&m, id).Error; err != nil {
			return err
		}
		if m.CitizenID != citizenID || m.Status != string(entity.DemandStatusAwaitingConfirmation) {
			return service.Forbidden("cannot confirm this demand")
		}
		if err := tx.Model(&m).Updates(map[string]any{"status": entity.DemandStatusCompleted, "confirmation_due_at": nil}).Error; err != nil {
			return err
		}
		metadata, _ := json.Marshal(statusMetadata(nil, entity.DemandStatusAwaitingConfirmation, entity.DemandStatusCompleted))
		return tx.Create(&model.DemandEvent{DemandID: id, Type: "confirmed", ActorUserID: &citizenID, Metadata: metadata}).Error
	})
	if err != nil {
		if _, ok := service.From(err); ok {
			return nil, err
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service.NotFound("demand not found")
		}
		return nil, service.Internal(err)
	}
	m.Status = string(entity.DemandStatusCompleted)
	m.ConfirmationDueAt = nil
	_ = s.notifyParticipants(ctx, m, citizenID, "demand_confirmed")
	return demandToOutput(m.ToDomain()), nil
}
func (s *DemandServiceImpl) Reopen(ctx context.Context, id, citizenID uint, input service.DemandTimelineInput) (*service.DemandOutput, error) {
	if err := validateTimelineInput(input); err != nil {
		return nil, err
	}
	var m model.Demand
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&m, id).Error; err != nil {
			return err
		}
		if m.CitizenID != citizenID || m.Status != string(entity.DemandStatusAwaitingConfirmation) {
			return service.Forbidden("cannot reopen this demand")
		}
		if err := tx.Model(&m).Updates(map[string]any{"status": entity.DemandStatusReview, "confirmation_due_at": nil}).Error; err != nil {
			return err
		}
		metadata, _ := json.Marshal(statusMetadata(&input, entity.DemandStatusAwaitingConfirmation, entity.DemandStatusReview))
		return tx.Create(&model.DemandEvent{DemandID: id, Type: "reopened", ActorUserID: &citizenID, Metadata: metadata}).Error
	})
	if err != nil {
		if _, ok := service.From(err); ok {
			return nil, err
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service.NotFound("demand not found")
		}
		return nil, service.Internal(err)
	}
	m.Status = string(entity.DemandStatusReview)
	m.ConfirmationDueAt = nil
	_ = s.notifyParticipants(ctx, m, citizenID, "demand_reopened")
	return demandToOutput(m.ToDomain()), nil
}
func (s *DemandServiceImpl) transition(ctx context.Context, id, officeID, userID uint, from, to entity.DemandStatus, event string, input service.DemandTimelineInput) (*service.DemandOutput, error) {
	if err := validateTimelineInput(input); err != nil {
		return nil, err
	}
	var m model.Demand
	e := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if e := tx.First(&m, id).Error; e != nil {
			return e
		}
		var a model.DemandAssignment
		if e := tx.Where("demand_id=? AND office_id=?", id, officeID).First(&a).Error; e != nil {
			return service.Forbidden("demand not assigned")
		}
		if m.Status != string(from) || (m.ResponsibleOfficeID != nil && *m.ResponsibleOfficeID != officeID) {
			return service.Conflict("invalid demand transition")
		}
		updates := map[string]any{"status": to}
		if from == entity.DemandStatusRegistered {
			updates["responsible_office_id"] = officeID
			updates["claimed_by_user_id"] = userID
		}
		if to == entity.DemandStatusAwaitingConfirmation {
			updates["confirmation_due_at"] = time.Now().Add(120 * time.Hour)
		}
		if e := tx.Model(&m).Updates(updates).Error; e != nil {
			return e
		}
		metadata, _ := json.Marshal(statusMetadata(&input, from, to))
		return tx.Create(&model.DemandEvent{DemandID: id, Type: event, ActorUserID: &userID, Metadata: metadata}).Error
	})
	if e != nil {
		if _, ok := service.From(e); ok {
			return nil, e
		}
		return nil, service.Internal(e)
	}
	m.Status = string(to)
	if to == entity.DemandStatusAwaitingConfirmation {
		due := time.Now().Add(120 * time.Hour)
		m.ConfirmationDueAt = &due
	}
	if from == entity.DemandStatusRegistered {
		m.ResponsibleOfficeID = &officeID
	}
	_ = s.notifyParticipants(ctx, m, userID, "demand_"+event)
	return demandToOutput(m.ToDomain()), nil
}

func statusMetadata(input *service.DemandTimelineInput, from, to entity.DemandStatus) map[string]any {
	metadata := map[string]any{"from_status": from, "to_status": to}
	if input != nil {
		metadata["message"] = strings.TrimSpace(input.Message)
		metadata["images"] = normalizedImages(input.Images)
	}
	return metadata
}

func (s *DemandServiceImpl) CreateMilestone(ctx context.Context, id, officeID, userID uint, input service.DemandTimelineInput) error {
	if err := validateTimelineInput(input); err != nil {
		return err
	}
	var demand model.Demand
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&demand, id).Error; err != nil {
			return err
		}
		if demand.ResponsibleOfficeID == nil || *demand.ResponsibleOfficeID != officeID {
			return service.Forbidden("only the responsible office can create a milestone")
		}
		metadata, _ := json.Marshal(map[string]any{
			"message": strings.TrimSpace(input.Message),
			"images":  normalizedImages(input.Images),
		})
		if err := tx.Create(&model.DemandEvent{DemandID: id, Type: "milestone", ActorUserID: &userID, Metadata: metadata}).Error; err != nil {
			return err
		}
		return tx.Model(&demand).Update("updated_at", time.Now()).Error
	})
	if err != nil {
		if _, ok := service.From(err); ok {
			return err
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service.NotFound("demand not found")
		}
		return service.Internal(err)
	}
	return s.notifyParticipants(ctx, demand, userID, "demand_milestone")
}
func (s *DemandServiceImpl) Comment(ctx context.Context, id, userID uint, role entity.UserRole, input service.DemandCommentInput) (*entity.DemandComment, error) {
	if err := validateComment(input); err != nil {
		return nil, err
	}
	var demand model.Demand
	if e := s.db.WithContext(ctx).First(&demand, id).Error; e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			return nil, service.NotFound("demand not found")
		}
		return nil, service.Internal(e)
	}
	var u model.User
	if e := s.db.WithContext(ctx).First(&u, userID).Error; e != nil {
		return nil, service.Internal(e)
	}
	images := normalizedImages(input.Images)
	b, _ := json.Marshal(images)
	parentID, e := s.resolveCommentParent(ctx, id, input.ParentID)
	if e != nil {
		return nil, e
	}
	m := model.DemandComment{DemandID: id, ParentID: parentID, AuthorID: userID, Body: strings.TrimSpace(input.Body), Images: b}
	if e := s.db.WithContext(ctx).Create(&m).Error; e != nil {
		return nil, service.Internal(e)
	}
	s.db.WithContext(ctx).Model(&model.Demand{}).Where("id=?", id).UpdateColumn("comment_count", gorm.Expr("comment_count + 1"))
	_ = s.notifyParticipants(ctx, demand, userID, "demand_commented")
	return &entity.DemandComment{ID: m.ID, DemandID: id, ParentID: parentID, AuthorID: userID, AuthorName: u.Name, AuthorRole: role, AuthorImageURL: s.commentAuthorImage(ctx, userID, role), Body: m.Body, Images: images, CreatedAt: m.CreatedAt}, nil
}

// resolveCommentParent keeps one root comment plus a single reply column. Replies to an
// existing reply are attached to that reply's root, preventing deeper nesting.
func (s *DemandServiceImpl) resolveCommentParent(ctx context.Context, demandID uint, requested *uint) (*uint, error) {
	if requested == nil {
		return nil, nil
	}
	var parent model.DemandComment
	if err := s.db.WithContext(ctx).Where("id = ? AND demand_id = ?", *requested, demandID).First(&parent).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service.InvalidInput("invalid comment parent")
		}
		return nil, service.Internal(err)
	}
	if parent.ParentID != nil {
		return parent.ParentID, nil
	}
	return &parent.ID, nil
}

func (s *DemandServiceImpl) commentAuthorImage(ctx context.Context, userID uint, role entity.UserRole) string {
	switch role {
	case entity.CouncillorUser:
		var councillor model.Councillor
		if s.db.WithContext(ctx).First(&councillor, "user_id = ?", userID).Error == nil {
			return councillor.ImageURL
		}
	case entity.OfficeMemberUser:
		var member model.OfficeMember
		if s.db.WithContext(ctx).Preload("Office.Councillor").First(&member, "user_id = ?", userID).Error == nil {
			if member.ImageURL != "" {
				return member.ImageURL
			}
			if member.Office != nil && member.Office.Councillor != nil {
				return member.Office.Councillor.ImageURL
			}
		}
	}
	return ""
}

func (s *DemandServiceImpl) notifyParticipants(ctx context.Context, demand model.Demand, actorID uint, kind string) error {
	users := []uint{demand.CitizenID}
	if demand.ResponsibleOfficeID != nil {
		var officeUsers []uint
		if err := s.db.WithContext(ctx).Table("offices").Select("councillor_id").Where("id = ?", *demand.ResponsibleOfficeID).Scan(&officeUsers).Error; err != nil {
			return err
		}
		var members []uint
		if err := s.db.WithContext(ctx).Table("office_members").Select("user_id").Where("office_id = ?", *demand.ResponsibleOfficeID).Scan(&members).Error; err != nil {
			return err
		}
		users = append(users, officeUsers...)
		users = append(users, members...)
	}
	var commenters []uint
	if err := s.db.WithContext(ctx).Table("demand_comments").Distinct("author_id").Where("demand_id = ?", demand.ID).Scan(&commenters).Error; err != nil {
		return err
	}
	users = append(users, commenters...)
	return s.createNotifications(ctx, demand.ID, actorID, kind, users)
}
func (s *DemandServiceImpl) Activity(ctx context.Context, id uint) (*service.DemandActivityOutput, error) {
	var demand model.Demand
	if e := s.db.WithContext(ctx).First(&demand, id).Error; e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			return nil, service.NotFound("demand not found")
		}
		return nil, service.Internal(e)
	}
	var es []model.DemandEvent
	var cs []model.DemandComment
	if e := s.db.WithContext(ctx).Where("demand_id=? AND type IN ?", id, publicTimelineEventTypes()).Order("created_at asc").Find(&es).Error; e != nil {
		return nil, service.Internal(e)
	}
	if e := s.db.WithContext(ctx).Where("demand_id=?", id).Order("created_at asc").Find(&cs).Error; e != nil {
		return nil, service.Internal(e)
	}
	out := &service.DemandActivityOutput{Events: make([]entity.DemandEvent, 0, len(es)), Comments: make([]entity.DemandComment, 0, len(cs))}
	for _, e := range es {
		var md map[string]any
		_ = json.Unmarshal(e.Metadata, &md)
		images, _ := md["images"].([]any)
		imageURLs := make([]string, 0, len(images))
		for _, image := range images {
			if value, ok := image.(string); ok {
				imageURLs = append(imageURLs, value)
			}
		}
		if imageURLs == nil {
			imageURLs = []string{}
		}
		actorName, actorRole, actorImage := "", entity.UserRole(""), ""
		if e.ActorUserID != nil {
			actorName, actorRole, actorImage = s.eventActor(ctx, *e.ActorUserID)
		}
		message, _ := md["message"].(string)
		out.Events = append(out.Events, entity.DemandEvent{
			ID: e.ID, DemandID: e.DemandID, Type: e.Type, ActorUserID: e.ActorUserID,
			ActorName: actorName, ActorRole: actorRole, ActorImageURL: actorImage,
			Message: message, Images: imageURLs, Metadata: md, CreatedAt: e.CreatedAt,
		})
	}
	for _, c := range cs {
		var u model.User
		_ = s.db.WithContext(ctx).First(&u, c.AuthorID).Error
		imgs := make([]string, 0)
		_ = json.Unmarshal(c.Images, &imgs)
		imgs = normalizedImages(imgs)
		body, images := c.Body, imgs
		if c.HiddenAt != nil {
			body, images = "", []string{}
		}
		role := entity.UserRole(u.Role)
		out.Comments = append(out.Comments, entity.DemandComment{ID: c.ID, DemandID: c.DemandID, ParentID: c.ParentID, AuthorID: c.AuthorID, AuthorName: u.Name, AuthorRole: role, AuthorImageURL: s.commentAuthorImage(ctx, c.AuthorID, role), Body: body, Images: images, HiddenAt: c.HiddenAt, Hidden: c.HiddenAt != nil, CreatedAt: c.CreatedAt})
	}
	return out, nil
}

func publicTimelineEventTypes() []string {
	return []string{
		"created", "claimed", "execution_started", "confirmation_requested", "confirmed",
		"reopened", "automatically_completed", "milestone",
	}
}

func (s *DemandServiceImpl) eventActor(ctx context.Context, userID uint) (string, entity.UserRole, string) {
	var user model.User
	if err := s.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		return "", "", ""
	}
	role := entity.UserRole(user.Role)
	return user.Name, role, s.commentAuthorImage(ctx, userID, role)
}

func validateComment(input service.DemandCommentInput) error {
	if strings.TrimSpace(input.Body) == "" && len(input.Images) == 0 {
		return service.InvalidInput("comment content required")
	}
	return validateImages(input.Images)
}

func validateTimelineInput(input service.DemandTimelineInput) error {
	if len(strings.TrimSpace(input.Message)) < 3 {
		return service.InvalidInput("timeline message required")
	}
	return validateImages(input.Images)
}

func validateImages(images []string) error {
	if len(images) > 5 {
		return service.InvalidInput("too many images")
	}
	for _, image := range images {
		url, err := url.Parse(image)
		if err != nil || (url.Scheme != "http" && url.Scheme != "https") || url.Host == "" {
			return service.InvalidInput("invalid image URL")
		}
	}
	return nil
}

func (s *DemandServiceImpl) GetSupport(ctx context.Context, demandID, citizenID uint) (*service.DemandSupportOutput, error) {
	var demand model.Demand
	if err := s.db.WithContext(ctx).First(&demand, demandID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service.NotFound("demand not found")
		}
		return nil, service.Internal(err)
	}
	return s.supportState(ctx, demand, citizenID)
}

func (s *DemandServiceImpl) AddSupport(ctx context.Context, demandID, citizenID uint) (*service.DemandSupportOutput, error) {
	var demand model.Demand
	var added bool
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&demand, demandID).Error; err != nil {
			return err
		}
		if demand.CitizenID == citizenID {
			return service.Forbidden("demand author cannot support")
		}
		result := tx.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "demand_id"}, {Name: "citizen_id"}}, DoNothing: true}).Create(&model.DemandSupport{DemandID: demandID, CitizenID: citizenID})
		if result.Error != nil {
			return result.Error
		}
		added = result.RowsAffected > 0
		if added {
			return tx.Model(&demand).UpdateColumn("support_count", gorm.Expr("support_count + 1")).Error
		}
		return nil
	}); err != nil {
		if _, ok := service.From(err); ok {
			return nil, err
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service.NotFound("demand not found")
		}
		return nil, service.Internal(err)
	}
	if added {
		demand.SupportCount++
	}
	return &service.DemandSupportOutput{SupportCount: demand.SupportCount, Supported: true, CanSupport: true}, nil
}

func (s *DemandServiceImpl) RemoveSupport(ctx context.Context, demandID, citizenID uint) (*service.DemandSupportOutput, error) {
	var demand model.Demand
	var removed bool
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&demand, demandID).Error; err != nil {
			return err
		}
		result := tx.Unscoped().Where("demand_id = ? AND citizen_id = ?", demandID, citizenID).Delete(&model.DemandSupport{})
		if result.Error != nil {
			return result.Error
		}
		removed = result.RowsAffected > 0
		if removed {
			return tx.Model(&demand).UpdateColumn("support_count", gorm.Expr("GREATEST(support_count - 1, 0)")).Error
		}
		return nil
	}); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service.NotFound("demand not found")
		}
		return nil, service.Internal(err)
	}
	if removed && demand.SupportCount > 0 {
		demand.SupportCount--
	}
	return &service.DemandSupportOutput{SupportCount: demand.SupportCount, Supported: false, CanSupport: demand.CitizenID != citizenID}, nil
}

func (s *DemandServiceImpl) supportState(ctx context.Context, demand model.Demand, citizenID uint) (*service.DemandSupportOutput, error) {
	output := &service.DemandSupportOutput{SupportCount: demand.SupportCount, CanSupport: demand.CitizenID != citizenID}
	if !output.CanSupport {
		return output, nil
	}
	var support model.DemandSupport
	err := s.db.WithContext(ctx).Where("demand_id = ? AND citizen_id = ?", demand.ID, citizenID).First(&support).Error
	if err == nil {
		output.Supported = true
		return output, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return output, nil
	}
	return nil, service.Internal(err)
}

func (s *DemandServiceImpl) DeleteComment(ctx context.Context, commentID, userID uint) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var comment model.DemandComment
		if err := tx.First(&comment, commentID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return service.NotFound("comment not found")
			}
			return service.Internal(err)
		}
		if comment.AuthorID != userID {
			return service.Forbidden("only the comment author can delete it")
		}

		commentIDs := []uint{comment.ID}
		var replyIDs []uint
		if err := tx.Model(&model.DemandComment{}).Where("parent_id = ?", comment.ID).Pluck("id", &replyIDs).Error; err != nil {
			return service.Internal(err)
		}
		commentIDs = append(commentIDs, replyIDs...)

		if err := tx.Where("comment_id IN ?", commentIDs).Delete(&model.DemandCommentReport{}).Error; err != nil {
			return service.Internal(err)
		}
		if err := tx.Where("id IN ?", commentIDs).Delete(&model.DemandComment{}).Error; err != nil {
			return service.Internal(err)
		}
		if err := tx.Model(&model.Demand{}).Where("id = ?", comment.DemandID).UpdateColumn("comment_count", gorm.Expr("GREATEST(comment_count - ?, 0)", len(commentIDs))).Error; err != nil {
			return service.Internal(err)
		}
		metadata, _ := json.Marshal(map[string]any{"comment_id": comment.ID, "deleted_comment_count": len(commentIDs)})
		if err := tx.Create(&model.DemandEvent{DemandID: comment.DemandID, Type: "comment_deleted", ActorUserID: &userID, Metadata: metadata}).Error; err != nil {
			return service.Internal(err)
		}
		return nil
	})
}

func normalizedImages(images []string) []string {
	if images == nil {
		return []string{}
	}
	return images
}

func (s *DemandServiceImpl) ReportComment(ctx context.Context, commentID, userID uint, reason string) error {
	if strings.TrimSpace(reason) == "" {
		return service.InvalidInput("report reason required")
	}
	var comment model.DemandComment
	if err := s.db.WithContext(ctx).First(&comment, commentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service.NotFound("comment not found")
		}
		return service.Internal(err)
	}
	report := model.DemandCommentReport{CommentID: commentID, ReporterID: userID, Reason: strings.TrimSpace(reason)}
	if err := s.db.WithContext(ctx).Create(&report).Error; err != nil {
		return service.Conflict("comment already reported")
	}
	return s.record(ctx, comment.DemandID, "comment_reported", &userID, map[string]any{"comment_id": commentID})
}

func (s *DemandServiceImpl) HideComment(ctx context.Context, commentID, officeID, userID uint) error {
	var comment model.DemandComment
	if err := s.db.WithContext(ctx).First(&comment, commentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service.NotFound("comment not found")
		}
		return service.Internal(err)
	}
	var demand model.Demand
	if err := s.db.WithContext(ctx).First(&demand, comment.DemandID).Error; err != nil {
		return service.Internal(err)
	}
	if demand.ResponsibleOfficeID == nil || *demand.ResponsibleOfficeID != officeID {
		return service.Forbidden("only responsible office can moderate")
	}
	now := time.Now()
	if err := s.db.WithContext(ctx).Model(&comment).Updates(map[string]any{"hidden_at": now, "hidden_by_user_id": userID}).Error; err != nil {
		return service.Internal(err)
	}
	return s.record(ctx, demand.ID, "comment_hidden", &userID, map[string]any{"comment_id": commentID})
}

func (s *DemandServiceImpl) ListNotifications(ctx context.Context, userID uint) ([]service.NotificationOutput, error) {
	var rows []model.Notification
	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at desc").Limit(50).Find(&rows).Error; err != nil {
		return nil, service.Internal(err)
	}
	out := make([]service.NotificationOutput, 0, len(rows))
	for _, n := range rows {
		out = append(out, service.NotificationOutput{ID: n.ID, Type: n.Type, DemandID: n.DemandID, ReadAt: n.ReadAt, CreatedAt: n.CreatedAt})
	}
	return out, nil
}
func (s *DemandServiceImpl) ListNotificationsAfter(ctx context.Context, userID, afterID uint) ([]service.NotificationOutput, error) {
	var rows []model.Notification
	if err := s.db.WithContext(ctx).Where("user_id = ? AND id > ?", userID, afterID).Order("id asc").Limit(100).Find(&rows).Error; err != nil {
		return nil, service.Internal(err)
	}
	out := make([]service.NotificationOutput, 0, len(rows))
	for _, n := range rows {
		out = append(out, service.NotificationOutput{ID: n.ID, Type: n.Type, DemandID: n.DemandID, ReadAt: n.ReadAt, CreatedAt: n.CreatedAt})
	}
	return out, nil
}
func (s *DemandServiceImpl) ReadNotifications(ctx context.Context, userID uint, ids []uint) error {
	now := time.Now()
	query := s.db.WithContext(ctx).Model(&model.Notification{}).Where("user_id = ? AND read_at IS NULL", userID)
	if len(ids) > 0 {
		query = query.Where("id IN ?", ids)
	}
	if err := query.Update("read_at", now).Error; err != nil {
		return service.Internal(err)
	}
	return nil
}
func (s *DemandServiceImpl) record(ctx context.Context, id uint, t string, actor *uint, meta map[string]any) error {
	b, _ := json.Marshal(meta)
	return s.db.WithContext(ctx).Create(&model.DemandEvent{DemandID: id, Type: t, ActorUserID: actor, Metadata: b}).Error
}

func (s *DemandServiceImpl) Create(ctx context.Context, input service.CreateDemandInput) (*service.DemandOutput, error) {
	if input.DirectedOfficeID != nil {
		if _, err := s.officeRepo.FindByID(ctx, *input.DirectedOfficeID); err != nil {
			var dbErr *repository.DBError
			if errors.As(err, &dbErr) && dbErr.Code == repository.DBErrorNotFound {
				return nil, service.NotFound("directed office not found")
			}
			return nil, service.Internal(err)
		}
	}
	protocol, err := generateDemandProtocol(time.Now())
	if err != nil {
		s.logger.Error("failed to generate demand protocol", "error", err)
		return nil, service.Internal(err)
	}

	demand, err := s.demandRepo.Create(ctx, repository.CreateDemandData{
		Protocol:     protocol,
		CitizenID:    input.CitizenID,
		Title:        strings.TrimSpace(input.Title),
		Description:  strings.TrimSpace(input.Description),
		Category:     strings.TrimSpace(input.Category),
		Street:       strings.TrimSpace(input.Street),
		Number:       strings.TrimSpace(input.Number),
		Neighborhood: strings.TrimSpace(input.Neighborhood),
		City:         strings.TrimSpace(input.City),
		State:        strings.ToUpper(strings.TrimSpace(input.State)),
		Status:       entity.DemandStatusRegistered,
		Latitude:     input.Latitude, Longitude: input.Longitude, Images: input.Images, DirectedOfficeID: input.DirectedOfficeID,
	})
	if err != nil {
		var dbErr *repository.DBError
		if errors.As(err, &dbErr) {
			if dbErr.Code == repository.DBErrorConflict {
				return nil, service.Conflict("protocol already exists")
			}
			s.logger.Error("failed to create demand", "citizenID", input.CitizenID, "error", err)
			return nil, service.Internal(err)
		}
		return nil, err
	}

	officeIDs := make([]uint, 0)
	if input.DirectedOfficeID != nil {
		officeIDs = append(officeIDs, *input.DirectedOfficeID)
	} else {
		officeIDs = s.matchRecipientOffices(ctx, input)
	}
	if err := s.demandRepo.AssignOffices(ctx, demand.ID, officeIDs); err != nil {
		s.logger.Error("failed to assign demand", "error", err)
		return nil, service.Internal(err)
	}
	if err := s.record(ctx, demand.ID, "created", &input.CitizenID, nil); err != nil {
		return nil, service.Internal(err)
	}
	if err := s.notifyOffices(ctx, demand.ID, officeIDs, input.CitizenID, "demand_created"); err != nil {
		s.logger.Warn("failed to notify offices", "error", err)
	}
	return demandToOutput(demand), nil
}

func (s *DemandServiceImpl) notifyOffices(ctx context.Context, demandID uint, officeIDs []uint, actorID uint, kind string) error {
	if len(officeIDs) == 0 {
		return nil
	}
	var users []uint
	if err := s.db.WithContext(ctx).Table("offices").Select("offices.councillor_id").Where("offices.id IN ?", officeIDs).Scan(&users).Error; err != nil {
		return err
	}
	var members []uint
	if err := s.db.WithContext(ctx).Table("office_members").Select("user_id").Where("office_id IN ?", officeIDs).Scan(&members).Error; err != nil {
		return err
	}
	users = append(users, members...)
	return s.createNotifications(ctx, demandID, actorID, kind, users)
}

func (s *DemandServiceImpl) createNotifications(ctx context.Context, demandID, actorID uint, kind string, userIDs []uint) error {
	seen := map[uint]bool{}
	rows := make([]model.Notification, 0, len(userIDs))
	for _, id := range userIDs {
		if id != 0 && id != actorID && !seen[id] {
			seen[id] = true
			rows = append(rows, model.Notification{UserID: id, DemandID: demandID, Type: kind})
		}
	}
	if len(rows) == 0 {
		return nil
	}
	return s.db.WithContext(ctx).Create(&rows).Error
}

func (s *DemandServiceImpl) matchRecipientOffices(ctx context.Context, input service.CreateDemandInput) []uint {
	for _, radiusKM := range []float64{5, 15, 30, 60} {
		officeIDs, err := s.officeRepo.ListActiveOfficeIDsNear(ctx, input.Latitude, input.Longitude, radiusKM, 3)
		if err != nil {
			s.logger.Error("failed to rank offices by regional activity", "radiusKM", radiusKM, "error", err)
			break
		}
		if len(officeIDs) > 0 {
			return officeIDs
		}
	}

	offices, err := s.officeRepo.ListByCityState(ctx, strings.TrimSpace(input.City), strings.ToUpper(strings.TrimSpace(input.State)))
	if err != nil {
		s.logger.Error("failed to find fallback recipient offices", "error", err)
		return []uint{}
	}
	officeIDs := make([]uint, 0, len(offices))
	for _, office := range offices {
		officeIDs = append(officeIDs, office.ID)
	}
	return officeIDs
}

func (s *DemandServiceImpl) ListForOffice(ctx context.Context, officeID uint, filters service.DemandListFilters) ([]service.DemandOutput, error) {
	demands, err := s.demandRepo.ListOfficeDemands(ctx, officeID, repository.DemandFilters{Status: filters.Status})
	if err != nil {
		return nil, service.Internal(err)
	}
	result := make([]service.DemandOutput, 0, len(demands))
	for _, item := range demands {
		result = append(result, *demandToOutput(&item))
	}
	return result, nil
}

func (s *DemandServiceImpl) FindByID(ctx context.Context, id uint) (*service.DemandOutput, error) {
	demand, err := s.demandRepo.FindByID(ctx, id)
	if err != nil {
		var dbErr *repository.DBError
		if errors.As(err, &dbErr) {
			if dbErr.Code == repository.DBErrorNotFound {
				return nil, service.NotFound("demand not found")
			}
			s.logger.Error("failed to find demand", "demandID", id, "error", err)
			return nil, service.Internal(err)
		}
		return nil, err
	}

	return demandToOutput(demand), nil
}

func (s *DemandServiceImpl) List(ctx context.Context, filters service.DemandListFilters) ([]service.DemandOutput, error) {
	if filters.Status != "" && !isValidDemandStatus(filters.Status) {
		return nil, service.InvalidInput("invalid demand status")
	}

	demands, err := s.demandRepo.List(ctx, repository.DemandFilters{
		City:         strings.TrimSpace(filters.City),
		State:        strings.ToUpper(strings.TrimSpace(filters.State)),
		Neighborhood: strings.TrimSpace(filters.Neighborhood),
		Category:     strings.TrimSpace(filters.Category),
		Status:       filters.Status,
	})
	if err != nil {
		var dbErr *repository.DBError
		if errors.As(err, &dbErr) {
			s.logger.Error("failed to list demands", "error", err)
			return nil, service.Internal(err)
		}
		return nil, err
	}

	output := make([]service.DemandOutput, 0, len(demands))
	for _, demand := range demands {
		output = append(output, *demandToOutput(&demand))
	}

	return output, nil
}

func (s *DemandServiceImpl) ListMine(ctx context.Context, citizenID uint) ([]service.DemandOutput, error) {
	demands, err := s.demandRepo.ListByCitizen(ctx, citizenID)
	if err != nil {
		return nil, service.Internal(err)
	}
	output := make([]service.DemandOutput, 0, len(demands))
	for _, demand := range demands {
		output = append(output, *demandToOutput(&demand))
	}
	return output, nil
}

func generateDemandProtocol(now time.Time) (string, error) {
	bytes := make([]byte, 3)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return fmt.Sprintf("CDD-%s-%s", now.Format("20060102"), strings.ToUpper(hex.EncodeToString(bytes))), nil
}

func demandToOutput(demand *entity.Demand) *service.DemandOutput {
	if demand == nil {
		return nil
	}

	return &service.DemandOutput{
		ID:           demand.ID,
		Protocol:     demand.Protocol,
		Title:        demand.Title,
		Description:  demand.Description,
		Category:     demand.Category,
		Street:       demand.Street,
		Number:       demand.Number,
		Neighborhood: demand.Neighborhood,
		City:         demand.City,
		State:        demand.State,
		Status:       demand.Status,
		SupportCount: demand.SupportCount,
		CommentCount: demand.CommentCount,
		CreatedAt:    demand.CreatedAt,
		UpdatedAt:    demand.UpdatedAt,
		Latitude:     demand.Latitude, Longitude: demand.Longitude, Images: normalizedImages(demand.Images), DirectedOfficeID: demand.DirectedOfficeID, ResponsibleOfficeID: demand.ResponsibleOfficeID, ConfirmationDueAt: demand.ConfirmationDueAt,
	}
}

func isValidDemandStatus(status entity.DemandStatus) bool {
	switch status {
	case entity.DemandStatusRegistered, entity.DemandStatusReview, entity.DemandStatusInProgress, entity.DemandStatusAwaitingConfirmation, entity.DemandStatusCompleted:
		return true
	default:
		return false
	}
}

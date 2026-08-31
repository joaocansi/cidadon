package usecase

import (
	"cidadon/internal/adapters/persistence/postgres/model"
	service "cidadon/internal/application/contract"
	"cidadon/internal/domain/entity"
	"cidadon/internal/domain/repository"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
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

func (s *DemandServiceImpl) Claim(ctx context.Context, id, officeID, userID uint) (*service.DemandOutput, error) {
	return s.transition(ctx, id, officeID, userID, entity.DemandStatusRegistered, entity.DemandStatusReview, "claimed", nil)
}
func (s *DemandServiceImpl) Start(ctx context.Context, id, officeID, userID uint) (*service.DemandOutput, error) {
	return s.transition(ctx, id, officeID, userID, entity.DemandStatusReview, entity.DemandStatusInProgress, "execution_started", nil)
}
func (s *DemandServiceImpl) RequestConfirmation(ctx context.Context, id, officeID, userID uint, input service.DemandCommentInput) (*service.DemandOutput, error) {
	if err := validateComment(input); err != nil {
		return nil, err
	}
	d, e := s.transition(ctx, id, officeID, userID, entity.DemandStatusInProgress, entity.DemandStatusAwaitingConfirmation, "confirmation_requested", &input)
	return d, e
}
func (s *DemandServiceImpl) Confirm(ctx context.Context, id, citizenID uint) (*service.DemandOutput, error) {
	var m model.Demand
	if e := s.db.WithContext(ctx).First(&m, id).Error; e != nil || m.CitizenID != citizenID || m.Status != string(entity.DemandStatusAwaitingConfirmation) {
		return nil, service.Forbidden("cannot confirm this demand")
	}
	if e := s.db.WithContext(ctx).Model(&m).Updates(map[string]any{"status": entity.DemandStatusCompleted, "confirmation_due_at": nil}).Error; e != nil {
		return nil, service.Internal(e)
	}
	_ = s.record(ctx, id, "confirmed", &citizenID, nil)
	_ = s.notifyParticipants(ctx, m, citizenID, "demand_confirmed")
	return demandToOutput(m.ToDomain()), nil
}
func (s *DemandServiceImpl) Reopen(ctx context.Context, id, citizenID uint, input service.DemandCommentInput) (*service.DemandOutput, error) {
	if err := validateComment(input); err != nil {
		return nil, err
	}
	var m model.Demand
	if e := s.db.WithContext(ctx).First(&m, id).Error; e != nil || m.CitizenID != citizenID || m.Status != string(entity.DemandStatusAwaitingConfirmation) {
		return nil, service.Forbidden("cannot reopen this demand")
	}
	if e := s.db.WithContext(ctx).Model(&m).Updates(map[string]any{"status": entity.DemandStatusReview, "confirmation_due_at": nil}).Error; e != nil {
		return nil, service.Internal(e)
	}
	_, e := s.Comment(ctx, id, citizenID, entity.CitizenUser, input)
	if e != nil {
		return nil, e
	}
	_ = s.record(ctx, id, "reopened", &citizenID, nil)
	_ = s.notifyParticipants(ctx, m, citizenID, "demand_reopened")
	m.Status = string(entity.DemandStatusReview)
	return demandToOutput(m.ToDomain()), nil
}
func (s *DemandServiceImpl) transition(ctx context.Context, id, officeID, userID uint, from, to entity.DemandStatus, event string, input *service.DemandCommentInput) (*service.DemandOutput, error) {
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
		if input != nil {
			b, _ := json.Marshal(input.Images)
			if e := tx.Create(&model.DemandComment{DemandID: id, AuthorID: userID, Body: strings.TrimSpace(input.Body), Images: b}).Error; e != nil {
				return e
			}
			tx.Model(&m).UpdateColumn("comment_count", gorm.Expr("comment_count + 1"))
		}
		return tx.Create(&model.DemandEvent{DemandID: id, Type: event, ActorUserID: &userID, Metadata: json.RawMessage("{}")}).Error
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
	b, _ := json.Marshal(input.Images)
	m := model.DemandComment{DemandID: id, AuthorID: userID, Body: strings.TrimSpace(input.Body), Images: b}
	if e := s.db.WithContext(ctx).Create(&m).Error; e != nil {
		return nil, service.Internal(e)
	}
	s.db.WithContext(ctx).Model(&model.Demand{}).Where("id=?", id).UpdateColumn("comment_count", gorm.Expr("comment_count + 1"))
	_ = s.record(ctx, id, "commented", &userID, nil)
	_ = s.notifyParticipants(ctx, demand, userID, "demand_commented")
	return &entity.DemandComment{ID: m.ID, DemandID: id, AuthorID: userID, AuthorName: u.Name, AuthorRole: role, Body: m.Body, Images: input.Images, CreatedAt: m.CreatedAt}, nil
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
	if e := s.db.WithContext(ctx).Where("demand_id=?", id).Order("created_at asc").Find(&es).Error; e != nil {
		return nil, service.Internal(e)
	}
	if e := s.db.WithContext(ctx).Where("demand_id=?", id).Order("created_at asc").Find(&cs).Error; e != nil {
		return nil, service.Internal(e)
	}
	out := &service.DemandActivityOutput{Events: make([]entity.DemandEvent, 0, len(es)), Comments: make([]entity.DemandComment, 0, len(cs))}
	for _, e := range es {
		var md map[string]any
		_ = json.Unmarshal(e.Metadata, &md)
		out.Events = append(out.Events, entity.DemandEvent{ID: e.ID, DemandID: e.DemandID, Type: e.Type, ActorUserID: e.ActorUserID, Metadata: md, CreatedAt: e.CreatedAt})
	}
	for _, c := range cs {
		var u model.User
		_ = s.db.WithContext(ctx).First(&u, c.AuthorID).Error
		var imgs []string
		_ = json.Unmarshal(c.Images, &imgs)
		body, images := c.Body, imgs
		if c.HiddenAt != nil {
			body, images = "", nil
		}
		out.Comments = append(out.Comments, entity.DemandComment{ID: c.ID, DemandID: c.DemandID, AuthorID: c.AuthorID, AuthorName: u.Name, AuthorRole: entity.UserRole(u.Role), Body: body, Images: images, HiddenAt: c.HiddenAt, Hidden: c.HiddenAt != nil, CreatedAt: c.CreatedAt})
	}
	return out, nil
}

func validateComment(input service.DemandCommentInput) error {
	if strings.TrimSpace(input.Body) == "" && len(input.Images) == 0 {
		return service.InvalidInput("comment content required")
	}
	if len(input.Images) > 5 {
		return service.InvalidInput("too many images")
	}
	for _, image := range input.Images {
		if !(strings.HasPrefix(image, "data:image/jpeg;base64,") || strings.HasPrefix(image, "data:image/png;base64,") || strings.HasPrefix(image, "data:image/webp;base64,")) {
			return service.InvalidInput("invalid image type")
		}
		parts := strings.SplitN(image, ",", 2)
		if len(parts) != 2 {
			return service.InvalidInput("invalid image")
		}
		decoded, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil || len(decoded) > 2<<20 {
			return service.InvalidInput("invalid image size")
		}
	}
	return nil
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

func (s *DemandServiceImpl) UpdateStatus(ctx context.Context, demandID, officeID uint, input service.UpdateDemandStatusInput) (*service.DemandOutput, error) {
	if !isValidDemandStatus(input.Status) {
		return nil, service.InvalidInput("invalid demand status")
	}
	demand, err := s.demandRepo.UpdateStatus(ctx, demandID, officeID, input.Status)
	if err != nil {
		var dbErr *repository.DBError
		if errors.As(err, &dbErr) && dbErr.Code == repository.DBErrorNotFound {
			return nil, service.NotFound("demand is not assigned to this office")
		}
		return nil, service.Internal(err)
	}
	return demandToOutput(demand), nil
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
		Latitude:     demand.Latitude, Longitude: demand.Longitude, Images: demand.Images, DirectedOfficeID: demand.DirectedOfficeID, ResponsibleOfficeID: demand.ResponsibleOfficeID, ConfirmationDueAt: demand.ConfirmationDueAt,
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

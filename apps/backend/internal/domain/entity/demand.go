package entity

import "time"

type DemandStatus string

const (
	DemandStatusRegistered           DemandStatus = "registered"
	DemandStatusReview               DemandStatus = "under_review"
	DemandStatusInProgress           DemandStatus = "in_progress"
	DemandStatusAwaitingConfirmation DemandStatus = "awaiting_confirmation"
	DemandStatusCompleted            DemandStatus = "completed"
)

type Demand struct {
	ID                  uint         `json:"id"`
	Protocol            string       `json:"protocol"`
	CitizenID           uint         `json:"citizen_id"`
	Citizen             *Citizen     `json:"citizen,omitempty"`
	Title               string       `json:"title"`
	Description         string       `json:"description"`
	Category            string       `json:"category"`
	Street              string       `json:"street"`
	Number              string       `json:"number"`
	Neighborhood        string       `json:"neighborhood"`
	City                string       `json:"city"`
	State               string       `json:"state"`
	Latitude            float64      `json:"latitude"`
	Longitude           float64      `json:"longitude"`
	Images              []string     `json:"images"`
	DirectedOfficeID    *uint        `json:"directed_office_id,omitempty"`
	ResponsibleOfficeID *uint        `json:"responsible_office_id,omitempty"`
	ClaimedByUserID     *uint        `json:"claimed_by_user_id,omitempty"`
	ConfirmationDueAt   *time.Time   `json:"confirmation_due_at,omitempty"`
	Status              DemandStatus `json:"status"`
	SupportCount        int          `json:"support_count"`
	CommentCount        int          `json:"comment_count"`
	CreatedAt           time.Time    `json:"created_at"`
	UpdatedAt           time.Time    `json:"updated_at"`
}

type DemandEvent struct {
	ID          uint           `json:"id"`
	DemandID    uint           `json:"demand_id"`
	Type        string         `json:"type"`
	ActorUserID *uint          `json:"actor_user_id,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
}
type DemandComment struct {
	ID         uint       `json:"id"`
	DemandID   uint       `json:"demand_id"`
	AuthorID   uint       `json:"author_id"`
	AuthorName string     `json:"author_name"`
	AuthorRole UserRole   `json:"author_role"`
	Body       string     `json:"body"`
	Images     []string   `json:"images"`
	HiddenAt   *time.Time `json:"hidden_at,omitempty"`
	Hidden     bool       `json:"hidden"`
	CreatedAt  time.Time  `json:"created_at"`
}

type DemandAssignment struct {
	ID        uint      `json:"id"`
	DemandID  uint      `json:"demand_id"`
	OfficeID  uint      `json:"office_id"`
	CreatedAt time.Time `json:"created_at"`
}

package model

import (
	"time"

	"github.com/Agmer17/srd_lab_creative/internal/db/sqlcgen"
	"github.com/google/uuid"
)

type ProjectMember struct {
	ID        uuid.UUID   `json:"id"`
	ProjectID uuid.UUID   `json:"project_id"`
	User      User        `json:"user"`
	Role      ProjectRole `json:"role"`
	IsOwner   bool        `json:"is_owner"`
	JoinedAt  time.Time   `json:"joined_at"`
	LeftAt    *time.Time  `json:"left_at,omitempty"`
}

type ProjectProgress struct {
	ID            uuid.UUID     `json:"id"`
	ProjectID     uuid.UUID     `json:"project_id"`
	ProjectMember ProjectMember `json:"member"` // Lebih enak dibaca 'member' aja di JSON
	Title         string        `json:"title"`
	Weight        float64       `json:"weight"`
	IsCompleted   bool          `json:"is_completed"`
	CreatedAt     time.Time     `json:"created_at"`
}

type Project struct {
	ID                   uuid.UUID         `json:"id"`
	OrderID              uuid.UUID         `json:"order_id"`
	Name                 string            `json:"name"`
	Description          *string           `json:"description"`
	Status               string            `json:"status"`
	AllowedRevisionCount int32             `json:"allowed_revision_count"`
	ProjectMembers       []ProjectMember   `json:"project_members"`
	Progress             []ProjectProgress `json:"progress"`
	ProjectRevision      []ProjectRevision `json:"project_revision,omitempty"`
	OrderData            *Order            `json:"order,omitempty"`
	ActualStartDate      *time.Time        `json:"actual_start_date"`
	EndDate              *time.Time        `json:"end_date"`
	CreatedAt            time.Time         `json:"created_at"`
	UpdatedAt            time.Time         `json:"updated_at"`
}

type ProjectRevision struct {
	Id        uuid.UUID `json:"id"`
	ProjectId uuid.UUID `json:"project_id"`
	Title     string    `json:"title"`
	Reason    string    `json:"reason"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func MapProjectDataToModel(dbProject sqlcgen.Project) Project {
	return Project{
		ID:                   dbProject.ID,
		OrderID:              dbProject.OrderID,
		Name:                 dbProject.Name,
		Description:          dbProject.Description,
		Status:               dbProject.Status,
		AllowedRevisionCount: dbProject.AllowedRevisionCount,
		ActualStartDate:      dbProject.ActualStartDate,
		EndDate:              dbProject.EndDate,
		CreatedAt:            dbProject.CreatedAt,
		UpdatedAt:            dbProject.UpdatedAt,
	}
}

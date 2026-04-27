package project

import (
	"time"
)

type createProjectRequest struct {
	OrderId         string     `json:"order_id" binding:"required,uuid"`
	Name            string     `json:"name" binding:"required,max=255"`
	Description     *string    `json:"description,omitempty"`
	Status          string     `json:"status" binding:"required,oneof=in_progress in_review completed archive"`
	AllowedRevision *int32     `json:"allowed_revision_count" binding:"omitempty,gte=3,lte=100"`
	Deadline        *time.Time `json:"end_date" binding:"omitempty"`
	CreatorRoleId   string     `json:"creator_role_id" binding:"required,uuid"`
}

type AddNewMemberDto struct {
	ProjectId string `json:"project_id" binding:"required,uuid"`
	UserId    string `json:"user_id" binding:"required,uuid"`
	RoleId    string `json:"role_id" binding:"required,uuid"`
	IsOwner   bool   `json:"is_owner" binding:"omitempty"`
}

type updateProjectRequest struct {
	Name            *string    `json:"name,omitempty" binding:"omitempty"`
	Description     *string    `json:"description,omitempty" binding:"omitempty"`
	Status          *string    `json:"status,omitempty" binding:"omitempty"`
	AllowedRevision *int32     `json:"allowed_revision,omitempty" binding:"omitempty,min=1,max=100"`
	EndDate         *time.Time `json:"end_date" binding:"omitempty"`
}

type UpdateMemberDataRequest struct {
	MemberId string `json:"member_id" binding:"required,uuid"`
	NewRole  string `json:"role_id" binding:"uuid"`
	IsOwner  *bool  `json:"is_owner,omitempty" binding:"omitempty"`
}

type createProgressRequests struct {
	Title           string  `json:"title" binding:"required,alphanumspace"`
	Weight          float64 `json:"weight" binding:"required,min=1,max=100"`
	ProjectMemberID string  `json:"project_member_id" binding:"required,uuid"`
}

type updateProgressRequest struct {
	Title           *string  `json:"title" binding:"omitempty,min=3,max=255,alphanumspace"`
	Weight          *float64 `json:"weight" binding:"omitempty,min=1,max=99"`
	IsComplete      *bool    `json:"is_completed" binding:"omitempty"`
	ProjectMemberId string   `json:"project_member_id" binding:"required,uuid"`
}

type createRevisionRequest struct {
	Title  string `json:"title" binding:"required,alphanumspace"`
	Reason string `json:"reason" binding:"required"`
}

type updateRevisionRequest struct {
	Status string `json:"status" binding:"required,oneof=pending accepted rejected"`
}

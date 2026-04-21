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

type addNewMemberDto struct {
	ProjectId string `json:"project_id" binding:"required,uuid"`
	UserId    string `json:"user_id" binding:"required,uuid"`
	RoleId    string `json:"role_id" binding:"required,uuid"`
	IsOwner   bool   `json:"is_owner" binding:"omitempty"`
}

package project

import (
	"context"
	"encoding/json"

	"github.com/Agmer17/srd_lab_creative/internal/db/sqlcgen"
	"github.com/Agmer17/srd_lab_creative/internal/shared/model"
	"github.com/google/uuid"
)

type ProjectMemberRepository struct {
	db *sqlcgen.Queries
}

func NewProjectMemberRepository(q *sqlcgen.Queries) *ProjectMemberRepository {

	return &ProjectMemberRepository{
		db: q,
	}
}

func (pmr *ProjectMemberRepository) CreateProjectMember(ctx context.Context, md model.ProjectMember) error {
	_, err := pmr.db.AddProjectMember(ctx, sqlcgen.AddProjectMemberParams{
		ProjectID: md.ProjectID,
		UserID:    md.User.ID,
		RoleID:    md.Role.Id,
		IsOwner:   md.IsOwner,
	})
	return err
}

func (pmr *ProjectMemberRepository) GetMemberFromProject(ctx context.Context, projectId uuid.UUID) ([]model.ProjectMember, error) {

	data, err := pmr.db.ListProjectMembersWithUser(ctx, projectId)
	if err != nil {
		return []model.ProjectMember{}, err
	}

	var listData []model.ProjectMember = make([]model.ProjectMember, len(data))
	for i, v := range data {

		var userData model.User
		err := json.Unmarshal(v.User, &userData)
		if err != nil {
			return []model.ProjectMember{}, err
		}

		var roleData model.ProjectRole
		umsErr := json.Unmarshal(v.Role, &roleData)
		if umsErr != nil {
			return []model.ProjectMember{}, err
		}

		listData[i] = model.ProjectMember{
			ID:        v.ID,
			ProjectID: v.ProjectID,
			IsOwner:   v.IsOwner,
			JoinedAt:  v.JoinedAt,
			LeftAt:    v.LeftAt,
			User:      userData,
			Role:      roleData,
		}

	}

	return listData, nil
}

package project

import (
	"context"

	"github.com/Agmer17/srd_lab_creative/internal/shared"
	"github.com/Agmer17/srd_lab_creative/internal/shared/model"
	"github.com/google/uuid"
)

type ProjectMemberService struct {
	memberRepo *ProjectMemberRepository
}

func NewProjectMemberService(repo *ProjectMemberRepository) *ProjectMemberService {
	return &ProjectMemberService{
		memberRepo: repo,
	}
}

func (pms *ProjectMemberService) addNewMember(ctx context.Context, req addNewMemberDto) ([]model.ProjectMember, *shared.ErrorResponse) {
	projectId, err := uuid.Parse(req.ProjectId)
	if err != nil {
		return []model.ProjectMember{}, shared.NewErrorResponse(400, "invalid projectId! please provide a valid uuid")
	}

	roleId, err := uuid.Parse(req.RoleId)
	if err != nil {
		return []model.ProjectMember{}, shared.NewErrorResponse(400, "invalid projectId! please provide a valid uuid")
	}

	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return []model.ProjectMember{}, shared.NewErrorResponse(400, "invalid projectId! please provide a valid uuid")
	}

	member := model.ProjectMember{
		ProjectID: projectId,
		Role: model.ProjectRole{
			Id: roleId,
		},
		User: model.User{
			ID: userId,
		},
		IsOwner: req.IsOwner,
	}

	insertErr := pms.memberRepo.CreateProjectMember(ctx, member)
	if insertErr != nil {
		return []model.ProjectMember{}, shared.NewErrorResponse(500, "something wrong while trying to add new project member")
	}

	newData, err := pms.memberRepo.GetMemberFromProject(ctx, projectId)
	if err != nil {
		return []model.ProjectMember{}, shared.NewErrorResponse(500, "something wrong while trying to create new member")
	}

	return newData, nil
}

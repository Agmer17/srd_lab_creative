package project

import (
	"context"
	"errors"
	"fmt"

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
		return []model.ProjectMember{}, shared.NewErrorResponse(400, "invalid projectId format")
	}

	roleId, err := uuid.Parse(req.RoleId)
	if err != nil {
		return []model.ProjectMember{}, shared.NewErrorResponse(400, "invalid roleId format")
	}

	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return []model.ProjectMember{}, shared.NewErrorResponse(400, "invalid userId format")
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
		switch {
		case errors.Is(insertErr, ErrUserNotFound):
			return []model.ProjectMember{}, shared.NewErrorResponse(404, "target user does not exist")

		case errors.Is(insertErr, ErrRoleNotFound):
			return []model.ProjectMember{}, shared.NewErrorResponse(404, "target role does not exist")
		case errors.Is(insertErr, ErrMemberExists):
			return []model.ProjectMember{}, shared.NewErrorResponse(409, "this user is already a member of the project")

		default:
			return []model.ProjectMember{}, shared.NewErrorResponse(500, "internal server error: failed to add member")
		}
	}

	newData, err := pms.memberRepo.GetMemberFromProject(ctx, projectId)
	if err != nil {
		return []model.ProjectMember{}, shared.NewErrorResponse(500, "member added but failed to fetch updated list")
	}

	return newData, nil
}

func (pms *ProjectMemberService) UpdateMemberRole(ctx context.Context, memberId uuid.UUID, newRole uuid.UUID, isOwner *bool) (model.ProjectMember, *shared.ErrorResponse) {
	err := pms.memberRepo.UpdateUserRoleFromProject(ctx, newRole, memberId, isOwner)
	if err != nil {
		if errors.Is(err, memberNotFound) {
			return model.ProjectMember{}, shared.NewErrorResponse(404, "member with this id not found")
		}
		return model.ProjectMember{}, shared.NewErrorResponse(500, "something wrong while trying to update role")
	}

	newData, err := pms.memberRepo.GetMemberDataById(ctx, memberId)
	if err != nil {
		if errors.Is(err, memberNotFound) {
			return model.ProjectMember{}, shared.NewErrorResponse(404, "member with this id not found")
		}
		return model.ProjectMember{}, shared.NewErrorResponse(500, "something wrong while trying to update role")
	}

	return newData, nil
}

func (pms *ProjectMemberService) GetMemberDataById(ctx context.Context, projectId uuid.UUID, memberId uuid.UUID) (model.ProjectMember, *shared.ErrorResponse) {
	data, err := pms.memberRepo.GetMemberDataById(ctx, memberId)
	if err != nil {
		if errors.Is(err, memberNotFound) {
			return model.ProjectMember{}, shared.NewErrorResponse(404, "member with this id not found")
		}
		return model.ProjectMember{}, shared.NewErrorResponse(500, "something wrong while trying to update data")
	}

	return data, nil
}

func (pms *ProjectMemberService) GetAllMemberFromProjectId(ctx context.Context, projectId uuid.UUID, userId uuid.UUID) ([]model.ProjectMember, *shared.ErrorResponse) {
	data, err := pms.memberRepo.GetMemberFromProject(ctx, projectId)
	if err != nil {
		return []model.ProjectMember{}, shared.NewErrorResponse(500, "something wrong while trying to getting user data")
	}

	allowed := false
	for _, v := range data {
		if v.User.ID == userId {
			allowed = true
		}
	}
	if !allowed {
		return []model.ProjectMember{}, shared.NewErrorResponse(403, "you can't access this data")
	}

	return data, nil
}

// owner, member, error
func (ps *ProjectMemberService) validateOwnerOrMember(ctx context.Context, userId uuid.UUID, projectId uuid.UUID) (bool, bool, error) {
	tempData, err := ps.memberRepo.GetMemberDataByUserId(ctx, userId, projectId)
	if err != nil {
		if errors.Is(err, memberNotFound) {
			fmt.Println("ERROR VALIDASI : ", err)
			return false, false, nil
		}
		return false, false, err
	}

	if tempData.ProjectID != projectId {
		fmt.Println(tempData)
		return false, false, nil
	}

	if !tempData.IsOwner {
		return false, true, nil
	}

	return true, true, nil
}

func (ps *ProjectMemberService) RemoveUserFromProject(ctx context.Context, toRemove uuid.UUID) *shared.ErrorResponse {

	err := ps.memberRepo.RemoveFromProject(ctx, toRemove)
	if err != nil {
		if errors.Is(err, memberNotFound) {
			return shared.NewErrorResponse(404, "no member found with this id")
		}
		return shared.NewErrorResponse(500, "something wrong while trying to remove user from project")
	}

	return nil
}

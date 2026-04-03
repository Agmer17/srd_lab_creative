package projectrole

import (
	"context"
	"errors"

	"github.com/Agmer17/srd_lab_creative/internal/shared"
	"github.com/Agmer17/srd_lab_creative/internal/shared/model"
	"github.com/google/uuid"
)

type ProjectRoleService struct {
	repo *ProjectRoleRepository
}

func NewProjectRoleService(rp *ProjectRoleRepository) *ProjectRoleService {
	return &ProjectRoleService{
		repo: rp,
	}
}

func (ps *ProjectRoleService) GetAllProjectRoles(ctx context.Context) ([]model.ProjectRole, *shared.ErrorResponse) {

	data, err := ps.repo.GetAllRoles(ctx)
	if err != nil {
		return []model.ProjectRole{}, shared.NewErrorResponse(500, "something wrong with the server while trying to get role data")
	}

	return data, nil
}

func (ps *ProjectRoleService) SearchRole(ctx context.Context, query string) ([]model.ProjectRole, *shared.ErrorResponse) {

	data, err := ps.repo.SearchRole(ctx, query)
	if err != nil {
		return []model.ProjectRole{}, shared.NewErrorResponse(500, "something wrong with the server while trying to get role data")
	}

	return data, nil
}

func (ps *ProjectRoleService) CreateRole(ctx context.Context, name string) (model.ProjectRole, *shared.ErrorResponse) {

	exist, err := ps.repo.ExistByName(ctx, name)
	if err != nil {
		return model.ProjectRole{}, shared.NewErrorResponse(500, "something wrong with the server while trying to insert new role")
	}

	if exist {
		return model.ProjectRole{}, shared.NewErrorResponse(409, "Project role with "+name+" already exist")
	}

	data, err := ps.repo.CreateRole(ctx, name)
	if err != nil {
		return model.ProjectRole{}, shared.NewErrorResponse(500, "something wrong with the server while trying to insert new role")
	}
	return data, nil
}

func (ps *ProjectRoleService) DeleteRoles(ctx context.Context, id uuid.UUID) *shared.ErrorResponse {

	_, err := ps.repo.GetRoleById(ctx, id)
	if err != nil {
		if errors.Is(err, errRoleNotFound) {
			return shared.NewErrorResponse(404, "role with this id is not found!")
		}
	}
	err = ps.repo.DeleteRole(ctx, id)
	if err != nil {
		return shared.NewErrorResponse(409, "you can't delete role that already being used in projects!")
	}
	return nil
}

func (ps *ProjectRoleService) UpdateRole(ctx context.Context, name string, id uuid.UUID) (model.ProjectRole, *shared.ErrorResponse) {

	exist, err := ps.repo.ExistByName(ctx, name)
	if err != nil {
		return model.ProjectRole{}, shared.NewErrorResponse(500, "something wrong with the server while trying to update new role")
	}

	if exist {
		return model.ProjectRole{}, shared.NewErrorResponse(409, "Project role with "+name+" already exist")
	}

	newData, err := ps.repo.UpdateRole(ctx, name, id)

	if err != nil {
		return model.ProjectRole{}, shared.NewErrorResponse(500, "something wrong with the server while trying to update new role")
	}

	return newData, nil

}

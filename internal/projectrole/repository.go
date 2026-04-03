package projectrole

import (
	"context"
	"errors"

	"github.com/Agmer17/srd_lab_creative/internal/db/sqlcgen"
	"github.com/Agmer17/srd_lab_creative/internal/shared/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var errRoleNotFound = errors.New("role not found")

type ProjectRoleRepository struct {
	db *sqlcgen.Queries
}

func NewProjectRoleRepository(q *sqlcgen.Queries) *ProjectRoleRepository {
	return &ProjectRoleRepository{
		db: q,
	}
}

func (pr *ProjectRoleRepository) GetAllRoles(ctx context.Context) ([]model.ProjectRole, error) {

	gen, err := pr.db.GetAllRoles(ctx)
	if err != nil {
		return []model.ProjectRole{}, err
	}
	return model.GenListToRoleModel(gen), nil
}

func (pr *ProjectRoleRepository) GetRoleById(ctx context.Context, id uuid.UUID) (model.ProjectRole, error) {

	gen, err := pr.db.GetRoleById(ctx, id)
	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return model.ProjectRole{}, errRoleNotFound
		}
		return model.ProjectRole{}, err
	}

	return model.MapToProjectRoleModel(gen), nil
}

func (pr *ProjectRoleRepository) SearchRole(ctx context.Context, query string) ([]model.ProjectRole, error) {

	gen, err := pr.db.SearchRoles(ctx, sqlcgen.SearchRolesParams{
		Keyword:   &query,
		OffsetVal: 0,
		LimitVal:  10000,
	})

	if err != nil {
		return []model.ProjectRole{}, err
	}

	return model.GenListToRoleModel(gen), nil
}

func (pr *ProjectRoleRepository) CreateRole(ctx context.Context, name string) (model.ProjectRole, error) {
	gen, err := pr.db.CreateRole(ctx, name)
	if err != nil {
		return model.ProjectRole{}, err
	}
	return model.MapToProjectRoleModel(gen), nil
}

func (pr *ProjectRoleRepository) ExistByName(ctx context.Context, name string) (bool, error) {

	_, err := pr.db.GetRoleByName(ctx, name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (pr *ProjectRoleRepository) DeleteRole(ctx context.Context, id uuid.UUID) error {
	err := pr.db.DeleteRole(ctx, id)
	return err
}

func (pr *ProjectRoleRepository) UpdateRole(ctx context.Context, newName string, id uuid.UUID) (model.ProjectRole, error) {

	data, err := pr.db.UpdateRole(ctx, sqlcgen.UpdateRoleParams{
		Name: newName,
		ID:   id,
	})

	if err != nil {
		return model.ProjectRole{}, err
	}

	return model.MapToProjectRoleModel(data), nil
}

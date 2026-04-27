package project

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/Agmer17/srd_lab_creative/internal/db/sqlcgen"
	"github.com/Agmer17/srd_lab_creative/internal/shared/model"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	memberNotFound  = errors.New("member not found")
	ErrUserNotFound = errors.New("user not found")
	ErrRoleNotFound = errors.New("role not found")
	ErrMemberExists = errors.New("user is already a member of this project")
	ErrInternal     = errors.New("internal server error")
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

	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			switch pgErr.Code {

			case pgerrcode.ForeignKeyViolation:
				detail := pgErr.Detail
				if strings.Contains(detail, "user_id") {
					return ErrUserNotFound
				}
				if strings.Contains(detail, "role_id") {
					return ErrRoleNotFound
				}

			case pgerrcode.UniqueViolation:
				return ErrMemberExists
			}
		}

		return err
	}

	return nil
}

func (pmr *ProjectMemberRepository) GetMemberFromProject(ctx context.Context, projectId uuid.UUID) ([]model.ProjectMember, error) {

	data, err := pmr.db.GetActiveProjectMembers(ctx, projectId)
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
			User:      userData,
			Role:      roleData,
		}

	}

	return listData, nil
}

func (pmr *ProjectMemberRepository) UpdateUserRoleFromProject(ctx context.Context, newRole uuid.UUID, memberId uuid.UUID, isOwner *bool) error {

	_, err := pmr.db.UpdateProjectMemberRole(ctx, sqlcgen.UpdateProjectMemberRoleParams{
		RoleID:   newRole,
		MemberID: memberId,
		IsOwner:  isOwner,
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return memberNotFound
		}
		return err
	}
	return nil
}

func (pmr *ProjectMemberRepository) GetMemberDataById(ctx context.Context, memberId uuid.UUID) (model.ProjectMember, error) {

	data, err := pmr.db.GetProjectMemberByID(ctx, memberId)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ProjectMember{}, memberNotFound
		}
		return model.ProjectMember{}, err
	}

	return MapGetProjectMemberByIDRowToModel(data)
}

func MapGetProjectMemberByIDRowToModel(row sqlcgen.GetProjectMemberByIDRow) (model.ProjectMember, error) {
	var user model.User
	if err := json.Unmarshal(row.User, &user); err != nil {
		return model.ProjectMember{}, err
	}

	var role model.ProjectRole
	if err := json.Unmarshal(row.Role, &role); err != nil {
		return model.ProjectMember{}, err
	}

	return model.ProjectMember{
		ID:        row.ID,
		ProjectID: row.ProjectID,
		User:      user,
		Role:      role,
		IsOwner:   row.IsOwner,
		JoinedAt:  row.JoinedAt,
		LeftAt:    row.LeftAt,
	}, nil
}

func (pmr *ProjectMemberRepository) GetMemberDataByUserId(ctx context.Context, userId uuid.UUID, projectId uuid.UUID) (model.ProjectMember, error) {

	row, err := pmr.db.GetMemberDataByUserId(ctx, sqlcgen.GetMemberDataByUserIdParams{
		UserID:    userId,
		ProjectID: projectId,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			fmt.Println(err)
			return model.ProjectMember{}, memberNotFound
		}
		return model.ProjectMember{}, ErrInternal
	}

	var user model.User
	if err := json.Unmarshal(row.User, &user); err != nil {
		return model.ProjectMember{}, err
	}

	var role model.ProjectRole
	if err := json.Unmarshal(row.Role, &role); err != nil {
		return model.ProjectMember{}, err
	}

	return model.ProjectMember{
		ID:        row.ID,
		ProjectID: row.ProjectID,
		User:      user,
		Role:      role,
		IsOwner:   row.IsOwner,
		JoinedAt:  row.JoinedAt,
		LeftAt:    row.LeftAt,
	}, nil
}

func (pmr *ProjectMemberRepository) RemoveFromProject(ctx context.Context, toRemove uuid.UUID) (uuid.UUID, uuid.UUID, error) {

	rem, err := pmr.db.RemoveProjectMember(ctx, toRemove)
	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, uuid.Nil, memberNotFound
		}
		return uuid.Nil, uuid.Nil, err
	}

	return rem.UserID, rem.ProjectID, nil
}

func (pmr *ProjectMemberRepository) getAllMembers(ctx context.Context) ([]model.ProjectMember, error) {

	data, err := pmr.db.GetAllMember(ctx)
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

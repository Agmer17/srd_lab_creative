package user

import (
	"context"
	"errors"
	"time"

	"github.com/Agmer17/srd_lab_creative/internal/db/sqlcgen"
	"github.com/Agmer17/srd_lab_creative/internal/shared/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var errNoUserFound = errors.New("user not found!")

type UserRepository struct {
	db *sqlcgen.Queries
}

func NewRepository(q *sqlcgen.Queries) *UserRepository {
	return &UserRepository{
		db: q,
	}
}

func (ur *UserRepository) SoftDeleteUser(ctx context.Context, id uuid.UUID) (*time.Time, error) {
	exec, err := ur.db.SoftDeleteUser(ctx, id)
	if err != nil {
		return nil, err
	}
	return exec.DeletedAt, nil
}

func (ur *UserRepository) GetAllUser(ctx context.Context) ([]model.User, error) {
	data, err := ur.db.ListUsers(ctx, sqlcgen.ListUsersParams{
		OffsetVal: 0,
		LimitVal:  10000,
	})

	if err != nil {
		return []model.User{}, err
	}

	userMap := model.GenListToUserMap(data)

	return userMap, nil
}

func (ur *UserRepository) SearchUser(ctx context.Context, query string) ([]model.User, error) {

	data, err := ur.db.SearchUsers(ctx, sqlcgen.SearchUsersParams{
		Keyword:   &query,
		LimitVal:  10000,
		OffsetVal: 0,
	})

	if err != nil {
		return []model.User{}, err
	}

	userList := model.GenListToUserMap(data)
	return userList, nil
}

func (ur *UserRepository) UpdateUserGlobalRole(ctx context.Context, role string, id uuid.UUID) (model.User, error) {
	data, err := ur.db.UpdateUserGlobalRole(ctx, sqlcgen.UpdateUserGlobalRoleParams{
		GlobalRole: role,
		ID:         id,
	})

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, errNoUserFound
		}
		return model.User{}, err
	}

	return model.MapToUserModel(data), nil
}

func (ur *UserRepository) UpdateUser(ctx context.Context, data UpdateUserDto, id uuid.UUID) (model.User, error) {
	updatedData, err := ur.db.UpdateUser(ctx, sqlcgen.UpdateUserParams{
		FullName:    data.FullName,
		Gender:      data.Gender,
		PhoneNumber: data.PhoneNumber,
		ID:          id,
	})
	if err != nil {
		return model.User{}, err
	}

	return model.MapToUserModel(updatedData), nil
}

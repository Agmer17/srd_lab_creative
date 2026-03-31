package user

import (
	"context"
	"errors"
	"time"

	"github.com/Agmer17/srd_lab_creative/internal/shared"
	"github.com/Agmer17/srd_lab_creative/internal/shared/model"
	"github.com/google/uuid"
)

type UserService struct {
	repo *UserRepository
}

func NewUserService(rp *UserRepository) *UserService {
	return &UserService{
		repo: rp,
	}
}

func (us *UserService) GetAllUser(ctx context.Context) ([]model.User, *shared.ErrorResponse) {

	data, err := us.repo.GetAllUser(ctx)
	if err != nil {
		return []model.User{}, shared.NewErrorResponse(500, "something wrong with the server rightnow! try again later")
	}

	return data, nil

}

func (us *UserService) SearchUser(ctx context.Context, query string) ([]model.User, *shared.ErrorResponse) {

	data, err := us.repo.SearchUser(ctx, query)
	if err != nil {
		return []model.User{}, shared.NewErrorResponse(500, "something wrong with the server rightnow! try again later")
	}

	return data, nil
}

func (us *UserService) UpdateUserRole(ctx context.Context, targetId uuid.UUID, newRole string) (time.Time, *shared.ErrorResponse) {
	data, err := us.repo.UpdateUserGlobalRole(ctx, newRole, targetId)
	if err != nil {
		if errors.Is(err, errNoUserFound) {
			return time.Now(), shared.NewErrorResponse(404, "no user with this id was found")

		}
		return time.Now(), shared.NewErrorResponse(500, "something went wrong while updating roles! try again later!")
	}

	return data.UpdatedAt, nil
}

func (us *UserService) UpdateUserData(ctx context.Context, data UpdateUserDto, userId uuid.UUID) (time.Time, *shared.ErrorResponse) {
	updatedData, err := us.repo.UpdateUser(ctx, data, userId)
	if err != nil {
		if errors.Is(err, errNoUserFound) {
			return time.Now(), shared.NewErrorResponse(404, "no user with this id was found")

		}
		return time.Now(), shared.NewErrorResponse(500, "something went wrong while updating roles! try again later!")
	}

	return updatedData.UpdatedAt, nil
}

func (us *UserService) DeleteUser(ctx context.Context, userId uuid.UUID) *shared.ErrorResponse {
	_, err := us.repo.SoftDeleteUser(ctx, userId)
	if err != nil {
		if errors.Is(err, errNoUserFound) {
			return shared.NewErrorResponse(404, "no user with this id was found")

		}
		return shared.NewErrorResponse(500, "something went wrong while updating roles! try again later!")
	}
	return nil
}

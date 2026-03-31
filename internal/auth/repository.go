package auth

import (
	"context"
	"errors"

	"github.com/Agmer17/srd_lab_creative/internal/db/sqlcgen"
	"github.com/Agmer17/srd_lab_creative/internal/shared/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var errAccountNotFound = errors.New("Accound doesn't exist!")

// ====================================================

type createUser struct {
	FullName       string
	Email          string
	ProfilePicture *string
	Provider       string
	ProviderUserID string
}

type userAuthInfo struct {
	UserId uuid.UUID
	Role   string
}

// ====================================================

type AuthRepository struct {
	db *sqlcgen.Queries
}

func NewAuthRepository(q *sqlcgen.Queries) *AuthRepository {
	return &AuthRepository{
		db: q,
	}
}

func (ar *AuthRepository) GetByEmail(ctx context.Context, email string) (model.User, error) {
	data, err := ar.db.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, errAccountNotFound
		}
	}
	return model.MapToUserModel(data), nil
}

func (ar *AuthRepository) CreateNewUser(ctx context.Context, newData createUser) (model.User, error) {

	d, err := ar.db.CreateUser(ctx, sqlcgen.CreateUserParams{
		FullName:       newData.FullName,
		Email:          newData.Email,
		ProfilePicture: newData.ProfilePicture,
		Provider:       newData.Provider,
		ProviderUserID: newData.ProviderUserID,
	})

	if err != nil {
		return model.User{}, err
	}

	return model.MapToUserModel(d), nil
}

func (ar *AuthRepository) ExistByProviderId(ctx context.Context, providerUserId string, providerName string) (userAuthInfo, bool, error) {

	existingData, err := ar.db.GetUserAuthInfoByProviderID(ctx, sqlcgen.GetUserAuthInfoByProviderIDParams{
		Provider:       providerName,
		ProviderUserID: providerUserId,
	})

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return userAuthInfo{}, false, errAccountNotFound
		}

		return userAuthInfo{}, false, err
	}

	return userAuthInfo{
		UserId: existingData.ID,
		Role:   existingData.GlobalRole,
	}, true, nil
}

func (ar *AuthRepository) GetUserById(ctx context.Context, id uuid.UUID) (model.User, error) {

	data, err := ar.db.GetUserById(ctx, id)
	if err != nil {

		if errors.Is(err, errAccountNotFound) {
			return model.User{}, errAccountNotFound
		}

		return model.User{}, err
	}

	return model.MapToUserModel(data), nil
}

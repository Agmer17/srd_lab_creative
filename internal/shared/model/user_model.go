package model

import (
	"time"

	"github.com/Agmer17/srd_lab_creative/internal/db/sqlcgen"
	"github.com/google/uuid"
)

const RoleAdmin = "ADMIN"
const RoleUser = "USER"

type User struct {
	ID             uuid.UUID  `json:"id"`
	GlobalRole     string     `json:"global_role"`
	FullName       string     `json:"full_name"`
	Email          string     `json:"email"`
	PhoneNumber    *string    `json:"phone_number"`
	ProfilePicture *string    `json:"profile_picture"`
	Gender         string     `json:"gender"`
	Provider       string     `json:"oauth_provider"`
	ProviderUserID string     `json:"oauth_provider_user_id"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at"`
}

func MapToUserModel(gen sqlcgen.User) User {
	return User{
		ID:             gen.ID,
		GlobalRole:     gen.GlobalRole,
		FullName:       gen.FullName,
		Email:          gen.Email,
		PhoneNumber:    gen.PhoneNumber,
		ProfilePicture: gen.ProfilePicture,
		Provider:       gen.Provider,
		ProviderUserID: gen.ProviderUserID,
		CreatedAt:      gen.CreatedAt,
		UpdatedAt:      gen.UpdatedAt,
		DeletedAt:      gen.DeletedAt,
	}
}

func GenListToUserMap(list []sqlcgen.User) []User {

	s := make([]User, len(list))

	for i, v := range list {
		s[i] = MapToUserModel(v)
	}

	return s
}

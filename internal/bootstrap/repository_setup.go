package bootstrap

import (
	"github.com/Agmer17/srd_lab_creative/internal/auth"
	"github.com/Agmer17/srd_lab_creative/internal/db/sqlcgen"
)

type RepositoryConfigs struct {
	AuthRepository *auth.AuthRepository
}

func NewRepositoryConfigs(q *sqlcgen.Queries) *RepositoryConfigs {
	authRepo := auth.NewAuthRepository(q)
	return &RepositoryConfigs{
		AuthRepository: authRepo,
	}

}

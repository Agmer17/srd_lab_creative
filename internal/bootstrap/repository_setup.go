package bootstrap

import (
	"github.com/Agmer17/srd_lab_creative/internal/auth"
	"github.com/Agmer17/srd_lab_creative/internal/db/sqlcgen"
	"github.com/Agmer17/srd_lab_creative/internal/projectrole"
	"github.com/Agmer17/srd_lab_creative/internal/user"
)

type RepositoryConfigs struct {
	AuthRepository        *auth.AuthRepository
	UserRepository        *user.UserRepository
	ProjectRoleRepository *projectrole.ProjectRoleRepository
}

func NewRepositoryConfigs(q *sqlcgen.Queries) *RepositoryConfigs {
	authRepo := auth.NewAuthRepository(q)
	userRepo := user.NewRepository(q)
	projectRoleRepo := projectrole.NewProjectRoleRepository(q)

	return &RepositoryConfigs{
		AuthRepository:        authRepo,
		UserRepository:        userRepo,
		ProjectRoleRepository: projectRoleRepo,
	}

}

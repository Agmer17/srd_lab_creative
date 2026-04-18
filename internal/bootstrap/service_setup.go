package bootstrap

import (
	"github.com/Agmer17/srd_lab_creative/internal/auth"
	"github.com/Agmer17/srd_lab_creative/internal/category"
	"github.com/Agmer17/srd_lab_creative/internal/projectrole"
	"github.com/Agmer17/srd_lab_creative/internal/user"
	"github.com/Agmer17/srd_lab_creative/internal/ws"
	"github.com/olahol/melody"
)

type ServiceConfigs struct {
	AuthService        *auth.AuthService
	UserService        *user.UserService
	ProjectRoleService *projectrole.ProjectRoleService
	CategoryService    *category.CategoryService
	WebsocketHub       *ws.WebsocketHub
}

func NewServiceConfigs(googleClientId string, googleSecret string, rpf *RepositoryConfigs, mel *melody.Melody) *ServiceConfigs {

	authService := auth.NewAuthService(googleClientId, googleSecret, rpf.AuthRepository)
	userService := user.NewUserService(rpf.UserRepository)
	projectRoleService := projectrole.NewProjectRoleService(rpf.ProjectRoleRepository)
	categoryService := category.NewCategoryService(rpf.CategoryRepository)

	wshub := ws.NewWebsocketHub(mel)

	return &ServiceConfigs{
		AuthService:        authService,
		UserService:        userService,
		ProjectRoleService: projectRoleService,
		CategoryService:    categoryService,
		WebsocketHub:       wshub,
	}
}

package bootstrap

import (
	"github.com/Agmer17/srd_lab_creative/internal/auth"
	"github.com/Agmer17/srd_lab_creative/internal/db/sqlcgen"
	"github.com/Agmer17/srd_lab_creative/internal/projectrole"
	"github.com/Agmer17/srd_lab_creative/internal/user"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Router       *gin.Engine
	Services     *ServiceConfigs
	Repositories *RepositoryConfigs
}

func NewApp(router *gin.Engine, googleClient string, googleSecret string, pool *pgxpool.Pool) *App {
	db := sqlcgen.New(pool)

	// setup reposiotry
	repoConfigs := NewRepositoryConfigs(db)

	// setup service
	serviceConfigs := NewServiceConfigs(googleClient, googleSecret, repoConfigs)

	// generate sama daftarin ke router
	authHandler := auth.NewAuthHandler(serviceConfigs.AuthService)
	userHandler := user.NewUserHandler(serviceConfigs.UserService)
	projectRoleHandler := projectrole.NewProjectRoleHandler(serviceConfigs.ProjectRoleService)

	SetupRoutes(
		router,
		authHandler,
		userHandler,
		projectRoleHandler,
	)

	return &App{
		Router:       router,
		Repositories: repoConfigs,
		Services:     serviceConfigs,
	}
}

func (a *App) Run() {
	a.Router.Run("0.0.0.0:80")
}

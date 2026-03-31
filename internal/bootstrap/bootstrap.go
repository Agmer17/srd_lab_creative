package bootstrap

import (
	"fmt"

	"github.com/Agmer17/srd_lab_creative/internal/auth"
	"github.com/Agmer17/srd_lab_creative/internal/db/sqlcgen"
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

	// generate all handler
	authHandler := auth.NewAuthHandler(serviceConfigs.AuthService)
	SetupRoutes(router, authHandler)

	return &App{
		Router:       router,
		Repositories: repoConfigs,
		Services:     serviceConfigs,
	}
}

func (a *App) Run() {
	fmt.Println("Server berjalan di port 80\n\n\n\n")
	a.Router.Run("0.0.0.0:80")
}

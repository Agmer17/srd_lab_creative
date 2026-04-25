package bootstrap

import (
	"context"
	"net/http"

	"github.com/Agmer17/srd_lab_creative/internal/auth"
	"github.com/Agmer17/srd_lab_creative/internal/category"
	"github.com/Agmer17/srd_lab_creative/internal/db/sqlcgen"
	"github.com/Agmer17/srd_lab_creative/internal/order"
	"github.com/Agmer17/srd_lab_creative/internal/product"
	"github.com/Agmer17/srd_lab_creative/internal/project"
	"github.com/Agmer17/srd_lab_creative/internal/projectrole"
	"github.com/Agmer17/srd_lab_creative/internal/user"
	"github.com/Agmer17/srd_lab_creative/internal/ws"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/olahol/melody"
	"github.com/redis/go-redis/v9"
)

type App struct {
	Router       *gin.Engine
	Services     *ServiceConfigs
	Repositories *RepositoryConfigs
}

func NewApp(ctx context.Context, router *gin.Engine, googleClient string, googleSecret string, pool *pgxpool.Pool, redisCli *redis.Client) *App {
	db := sqlcgen.New(pool)

	// setup melody
	mel := melody.New()

	mel.Upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	// setup reposiotry
	repoConfigs := NewRepositoryConfigs(db)

	// setup service
	serviceConfigs := NewServiceConfigs(
		ctx,
		googleClient,
		googleSecret,
		repoConfigs,
		mel,
		redisCli,
	)

	// generate sama daftarin ke router
	authHandler := auth.NewAuthHandler(serviceConfigs.AuthService)
	userHandler := user.NewUserHandler(serviceConfigs.UserService)
	projectRoleHandler := projectrole.NewProjectRoleHandler(serviceConfigs.ProjectRoleService)
	categoryHandler := category.NewCategoryHandler(serviceConfigs.CategoryService)
	productHandler := product.NewProductHandler(serviceConfigs.ProductService)

	orderHandler := order.NewOrderHandler(serviceConfigs.OrderService)
	projectHandler := project.NewProjectHandler(serviceConfigs.ProjectService)

	// ws
	wsHandler := ws.NewWebsocketHandler(mel)

	SetupRoutes(
		router,
		authHandler,
		userHandler,
		projectRoleHandler,
		categoryHandler,
		wsHandler,
		productHandler,
		orderHandler,
		projectHandler,
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

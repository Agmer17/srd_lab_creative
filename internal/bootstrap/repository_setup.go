package bootstrap

import (
	"github.com/Agmer17/srd_lab_creative/internal/auth"
	"github.com/Agmer17/srd_lab_creative/internal/category"
	"github.com/Agmer17/srd_lab_creative/internal/db/sqlcgen"
	"github.com/Agmer17/srd_lab_creative/internal/order"
	"github.com/Agmer17/srd_lab_creative/internal/product"
	"github.com/Agmer17/srd_lab_creative/internal/project"
	"github.com/Agmer17/srd_lab_creative/internal/projectrole"
	"github.com/Agmer17/srd_lab_creative/internal/user"
)

type RepositoryConfigs struct {
	AuthRepository          *auth.AuthRepository
	UserRepository          *user.UserRepository
	ProjectRoleRepository   *projectrole.ProjectRoleRepository
	CategoryRepository      *category.CategoryRepository
	ProductRepository       *product.ProductRepository
	OrderRepository         *order.OrderRepository
	ProjectRepository       *project.ProjectRepository
	ProjectMemberRepository *project.ProjectMemberRepository
	ProgressRepository      *project.ProgresRepository
}

func NewRepositoryConfigs(q *sqlcgen.Queries) *RepositoryConfigs {
	authRepo := auth.NewAuthRepository(q)
	userRepo := user.NewRepository(q)
	projectRoleRepo := projectrole.NewProjectRoleRepository(q)
	categoryRepo := category.NewCategoryRepository(q)
	productRepo := product.NewProductRepository(q)
	orderRepo := order.NewOrderRepositories(q)

	projectRepo := project.NewProjectRepository(q)
	projectMemberRepo := project.NewProjectMemberRepository(q)
	progressRepo := project.NewProgresRepository(q)

	return &RepositoryConfigs{
		AuthRepository:          authRepo,
		UserRepository:          userRepo,
		ProjectRoleRepository:   projectRoleRepo,
		CategoryRepository:      categoryRepo,
		ProductRepository:       productRepo,
		OrderRepository:         orderRepo,
		ProjectRepository:       projectRepo,
		ProjectMemberRepository: projectMemberRepo,
		ProgressRepository:      progressRepo,
	}

}

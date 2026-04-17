package bootstrap

import (
	"github.com/Agmer17/srd_lab_creative/internal/auth"
	"github.com/Agmer17/srd_lab_creative/internal/category"
	"github.com/Agmer17/srd_lab_creative/internal/product"
	"github.com/Agmer17/srd_lab_creative/internal/projectrole"
	"github.com/Agmer17/srd_lab_creative/internal/user"
)

type ServiceConfigs struct {
	AuthService        *auth.AuthService
	UserService        *user.UserService
	ProjectRoleService *projectrole.ProjectRoleService
	CategoryService    *category.CategoryService
	ProductService     *product.ProductService
}

func NewServiceConfigs(googleClientId string, googleSecret string, rpf *RepositoryConfigs) *ServiceConfigs {

	authService := auth.NewAuthService(googleClientId, googleSecret, rpf.AuthRepository)
	userService := user.NewUserService(rpf.UserRepository)
	projectRoleService := projectrole.NewProjectRoleService(rpf.ProjectRoleRepository)
	categoryService := category.NewCategoryService(rpf.CategoryRepository)
	productService := product.NewProductService(rpf.ProductRepository)

	return &ServiceConfigs{
		AuthService:        authService,
		UserService:        userService,
		ProjectRoleService: projectRoleService,
		CategoryService:    categoryService,
		ProductService: productService,
	}
}

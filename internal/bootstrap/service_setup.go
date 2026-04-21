package bootstrap

import (
	"github.com/Agmer17/srd_lab_creative/internal/auth"
	"github.com/Agmer17/srd_lab_creative/internal/category"
	"github.com/Agmer17/srd_lab_creative/internal/order"
	"github.com/Agmer17/srd_lab_creative/internal/product"
	"github.com/Agmer17/srd_lab_creative/internal/project"
	"github.com/Agmer17/srd_lab_creative/internal/projectrole"
	"github.com/Agmer17/srd_lab_creative/internal/storage"
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
	ProductService     *product.ProductService
	OrderService       *order.OrderService

	ProjectService       *project.ProjectService
	ProjectMemberService *project.ProjectMemberService
}

func NewServiceConfigs(googleClientId string, googleSecret string, rpf *RepositoryConfigs, mel *melody.Melody) *ServiceConfigs {

	authService := auth.NewAuthService(googleClientId, googleSecret, rpf.AuthRepository)
	userService := user.NewUserService(rpf.UserRepository)
	projectRoleService := projectrole.NewProjectRoleService(rpf.ProjectRoleRepository)
	categoryService := category.NewCategoryService(rpf.CategoryRepository)
	myStorage := storage.NewFileStorage(5);
	productService := product.NewProductService(rpf.ProductRepository,myStorage);


	orderService := order.NewOrderService(rpf.OrderRepository, productService)

	memberService := project.NewProjectMemberService(rpf.ProjectMemberRepository)
	projectService := project.NewProjectService(
		rpf.ProjectRepository,
		orderService,
		memberService,
	)

	wshub := ws.NewWebsocketHub(mel)

	return &ServiceConfigs{
		AuthService:        authService,
		UserService:        userService,
		ProjectRoleService: projectRoleService,
		CategoryService:    categoryService,
		WebsocketHub:       wshub,
		ProductService:     productService,
		OrderService:       orderService,
		ProjectService:     projectService,
	}
}

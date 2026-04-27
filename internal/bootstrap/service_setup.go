package bootstrap

import (
	"context"

	"github.com/Agmer17/srd_lab_creative/internal/auth"
	"github.com/Agmer17/srd_lab_creative/internal/category"
	"github.com/Agmer17/srd_lab_creative/internal/chat"
	"github.com/Agmer17/srd_lab_creative/internal/order"
	"github.com/Agmer17/srd_lab_creative/internal/product"
	"github.com/Agmer17/srd_lab_creative/internal/project"
	"github.com/Agmer17/srd_lab_creative/internal/projectrole"
	"github.com/Agmer17/srd_lab_creative/internal/storage"
	"github.com/Agmer17/srd_lab_creative/internal/user"
	"github.com/Agmer17/srd_lab_creative/internal/ws"
	"github.com/olahol/melody"
	"github.com/redis/go-redis/v9"
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
	ProgressService      *project.ProgressService

	ChatroomService *chat.ChatroomService
	ChatService     *chat.ChatService
	MessaginService *chat.MessagingService
}

func NewServiceConfigs(ctx context.Context, googleClientId string, googleSecret string, rpf *RepositoryConfigs, mel *melody.Melody, rdb *redis.Client) *ServiceConfigs {

	authService := auth.NewAuthService(googleClientId, googleSecret, rpf.AuthRepository)
	userService := user.NewUserService(rpf.UserRepository)
	projectRoleService := projectrole.NewProjectRoleService(rpf.ProjectRoleRepository)
	categoryService := category.NewCategoryService(rpf.CategoryRepository)
	myStorage := storage.NewFileStorage(5)
	productService := product.NewProductService(rpf.ProductRepository, myStorage)

	orderService := order.NewOrderService(rpf.OrderRepository, productService)

	wshub := ws.NewWebsocketHub(mel)
	chatroomSvc := chat.NewChatroomService(rpf.ChatroomRepository)
	chatSvc := chat.NewChatService(rpf.ChatRepository, rpf.ChatMediaRepository, myStorage)
	messagingService := chat.NewMessagingService(chatSvc, chatroomSvc, wshub, rdb)

	memberService := project.NewProjectMemberService(ctx, rpf.ProjectMemberRepository, rdb)
	progressService := project.NewProgressService(rpf.ProgressRepository)
	revisonService := project.NewRevisionService(rpf.RevisionRepository)

	projectService := project.NewProjectService(
		rpf.ProjectRepository,
		orderService,
		memberService,
		progressService,
		revisonService,
		chatroomSvc,
	)

	return &ServiceConfigs{
		AuthService:          authService,
		UserService:          userService,
		ProjectRoleService:   projectRoleService,
		CategoryService:      categoryService,
		WebsocketHub:         wshub,
		ProductService:       productService,
		OrderService:         orderService,
		ProjectService:       projectService,
		ProgressService:      progressService,
		ProjectMemberService: memberService,
		ChatroomService:      chatroomSvc,
		ChatService:          chatSvc,
		MessaginService:      messagingService,
	}
}

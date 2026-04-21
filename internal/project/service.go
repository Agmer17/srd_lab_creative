package project

import (
	"context"
	"errors"
	"fmt"

	"github.com/Agmer17/srd_lab_creative/internal/order"
	"github.com/Agmer17/srd_lab_creative/internal/shared"
	"github.com/Agmer17/srd_lab_creative/internal/shared/model"
	"github.com/google/uuid"
)

type ProjectService struct {
	projectRepo   *ProjectRepository
	orderService  *order.OrderService
	memberService *ProjectMemberService
}

func NewProjectService(repo *ProjectRepository, orderSvc *order.OrderService, memSvc *ProjectMemberService) *ProjectService {
	return &ProjectService{
		projectRepo:   repo,
		orderService:  orderSvc,
		memberService: memSvc,
	}
}

func (ps *ProjectService) GetAllProjects(ctx context.Context) ([]model.Project, *shared.ErrorResponse) {

	data, err := ps.projectRepo.GetAllProjects(ctx)
	if err != nil {
		fmt.Println(err)
		return []model.Project{}, shared.NewErrorResponse(500, "something wrong while trying to get project data")
	}

	return data, nil
}

func (ps *ProjectService) CreateProject(ctx context.Context, dto createProjectRequest, userId uuid.UUID) (model.Project, *shared.ErrorResponse) {
	orderId, err := uuid.Parse(dto.OrderId)
	if err != nil {
		return model.Project{}, shared.NewErrorResponse(400, "invalid order id!")
	}

	roleId, err := uuid.Parse(dto.CreatorRoleId)
	if err != nil {
		return model.Project{}, shared.NewErrorResponse(400, "invalid order id!")
	}

	orderData, getErr := ps.orderService.GetOrderById(ctx, orderId)
	if getErr != nil {
		return model.Project{}, getErr
	}

	if orderData.Status != model.TypeOrderStatusCompleted {
		return model.Project{}, shared.NewErrorResponse(409, "you can't make project from incompleted order!")
	}

	data, err := ps.projectRepo.CreateProjects(ctx, dto)
	if err != nil {
		return model.Project{}, shared.NewErrorResponse(500, "something wrong while trying to create project")
	}

	creatorData := addNewMemberDto{
		ProjectId: data.ID.String(),
		UserId:    userId.String(),
		RoleId:    roleId.String(),
		IsOwner:   true,
	}
	memberData, insMemErr := ps.memberService.addNewMember(ctx, creatorData)
	if insMemErr != nil {
		fmt.Println(memberData)
		return model.Project{}, insMemErr
	}

	data.ProjectMembers = memberData
	return data, nil
}

func (ps *ProjectService) DeleteProjects(ctx context.Context, id uuid.UUID) *shared.ErrorResponse {

	err := ps.projectRepo.DeleteProjects(ctx, id)
	if err != nil {
		if errors.Is(err, projectNotFound) {
			return shared.NewErrorResponse(404, "no projects with this id found")
		}
		return shared.NewErrorResponse(500, "something went wrong while trying to delete projects!")
	}

	return nil
}

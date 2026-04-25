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
	projectRepo     *ProjectRepository
	orderService    *order.OrderService
	memberService   *ProjectMemberService
	progressService *ProgressService
	revisionService *RevisionService
}

func NewProjectService(
	repo *ProjectRepository,
	orderSvc *order.OrderService,
	memSvc *ProjectMemberService,
	progSvc *ProgressService,
	revSvc *RevisionService,
) *ProjectService {

	return &ProjectService{
		projectRepo:     repo,
		orderService:    orderSvc,
		memberService:   memSvc,
		progressService: progSvc,
		revisionService: revSvc,
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
		if errors.Is(err, errOrderAlreadyUse) {
			return model.Project{}, shared.NewErrorResponse(409, "project with this order id already exist!")

		}
		return model.Project{}, shared.NewErrorResponse(500, "something wrong while trying to create project")
	}

	creatorData := addNewMemberDto{
		ProjectId: data.ID.String(),
		UserId:    userId.String(),
		RoleId:    roleId.String(),
		IsOwner:   true,
	}
	memberData, insMemErr := ps.memberService.addNewMember(ctx, creatorData)
	ps.memberService.setOneMembersRedis(ctx, memberData[0])
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

func (ps *ProjectService) UpdateProjectData(ctx context.Context, id uuid.UUID, dto updateProjectRequest) (model.Project, *shared.ErrorResponse) {

	data, err := ps.projectRepo.UpdateProject(ctx, id, dto)
	if err != nil {
		if errors.Is(err, projectNotFound) {
			return model.Project{}, shared.NewErrorResponse(404, "no project with this id found")
		}
		return model.Project{}, shared.NewErrorResponse(500, "something went wronf with the server while trying to update project data")
	}

	return data, nil
}

func (ps *ProjectService) GetDetailById(ctx context.Context, id uuid.UUID, userId uuid.UUID) (model.Project, *shared.ErrorResponse) {

	data, err := ps.projectRepo.GetProjectDetailById(ctx, id)
	if err != nil {
		if errors.Is(err, projectNotFound) {
			return model.Project{}, shared.NewErrorResponse(404, "project with this id not found! please try again")
		}

		return model.Project{}, shared.NewErrorResponse(500, "something wronf while trying to get project details, try again later")
	}

	allowed := false

	for _, v := range data.ProjectMembers {
		if v.User.ID == userId {
			allowed = true
			break
		}
	}

	if !allowed {
		return model.Project{}, shared.NewErrorResponse(403, "you can't access this data")
	}

	return data, nil
}

func (ps *ProjectService) GetMemberFromProject(ctx context.Context, projectId uuid.UUID, userId uuid.UUID) ([]model.ProjectMember, *shared.ErrorResponse) {

	mem, err := ps.memberService.GetAllMemberFromProjectId(ctx, projectId, userId)
	if err != nil {
		return []model.ProjectMember{}, err
	}

	return mem, nil
}

func (ps *ProjectService) AddNewMember(ctx context.Context, projectId uuid.UUID, userId uuid.UUID, req addNewMemberDto) ([]model.ProjectMember, *shared.ErrorResponse) {

	own, _, err := ps.memberService.validateOwnerOrMember(ctx, userId, projectId)
	if err != nil {
		return []model.ProjectMember{}, shared.NewErrorResponse(500, "something wrong with the server")
	}
	if !own {
		return []model.ProjectMember{}, shared.NewErrorResponse(403, "permision denied")
	}

	newData, addErr := ps.memberService.addNewMember(ctx, req)
	if addErr != nil {
		return []model.ProjectMember{}, addErr
	}

	return newData, nil

}

func (ps *ProjectService) UpdateProjectMemberRole(ctx context.Context, curr uuid.UUID, memberId uuid.UUID, projectId uuid.UUID, role uuid.UUID, isOwner *bool) (model.ProjectMember, *shared.ErrorResponse) {

	owner, _, err := ps.memberService.validateOwnerOrMember(ctx, curr, projectId)
	if err != nil {
		return model.ProjectMember{}, shared.NewErrorResponse(500, "something wrong with the server")
	}

	if !owner {
		return model.ProjectMember{}, shared.NewErrorResponse(403, "permision denied")
	}

	data, uptErr := ps.memberService.UpdateMemberRole(ctx, memberId, role, isOwner)
	if uptErr != nil {
		return model.ProjectMember{}, uptErr
	}

	return data, nil

}

func (ps *ProjectService) RemoveUserFromProject(ctx context.Context, curr uuid.UUID, rmf uuid.UUID, projectId uuid.UUID) *shared.ErrorResponse {

	owner, _, err := ps.memberService.validateOwnerOrMember(ctx, curr, projectId)
	if err != nil {
		return shared.NewErrorResponse(500, "something wrong with the server")
	}

	if !owner {
		return shared.NewErrorResponse(403, "permision denied")
	}

	delErr := ps.memberService.RemoveUserFromProject(ctx, rmf)
	if delErr != nil {
		return delErr
	}

	return nil

}

func (ps *ProjectService) GetProgressFromProject(ctx context.Context, id uuid.UUID, userId uuid.UUID) ([]model.ProjectProgress, *shared.ErrorResponse) {

	own, mem, err := ps.memberService.validateOwnerOrMember(ctx, userId, id)
	if err != nil {
		return []model.ProjectProgress{}, shared.NewErrorResponse(500, "something wrong while trying to get user data")
	}

	if !own && !mem {
		return []model.ProjectProgress{}, shared.NewErrorResponse(403, "access denied! you can't access this data")
	}

	return ps.progressService.GetProgressFromProject(ctx, id)
}

func (ps *ProjectService) CreateProjectProgress(ctx context.Context, curr uuid.UUID, projectId uuid.UUID, dto createProgressRequests) ([]model.ProjectProgress, *shared.ErrorResponse) {

	own, _, err := ps.memberService.validateOwnerOrMember(ctx, curr, projectId)
	if err != nil {
		return []model.ProjectProgress{}, shared.NewErrorResponse(500, "something wrong while trying to get user data")
	}

	if !own {
		return []model.ProjectProgress{}, shared.NewErrorResponse(403, "access denied! you can't access this data")
	}

	newData, insErr := ps.progressService.AddNewProgress(ctx, projectId, dto)
	if insErr != nil {
		return []model.ProjectProgress{}, insErr
	}

	return newData, nil
}

func (ps *ProjectService) RemoveProjectProgress(ctx context.Context, id uuid.UUID) *shared.ErrorResponse {

	return ps.progressService.DeleteProgress(ctx, id)
}

func (ps *ProjectService) UpdateProjectProgress(
	ctx context.Context,
	id uuid.UUID,
	dto updateProgressRequest,
) (model.ProjectProgress, *shared.ErrorResponse) {

	return ps.progressService.UpdateProgressData(ctx, id, dto)
}

func (ps *ProjectService) GetRevisionFromProject(ctx context.Context, curr uuid.UUID, projectId uuid.UUID) ([]model.ProjectRevision, *shared.ErrorResponse) {

	own, mem, err := ps.memberService.validateOwnerOrMember(ctx, curr, projectId)
	if err != nil {
		return []model.ProjectRevision{}, shared.NewErrorResponse(500, "something wrong while trying to get user data")
	}

	if !own && !mem {
		return []model.ProjectRevision{}, shared.NewErrorResponse(403, "access denied! you can't access this data")
	}

	return ps.revisionService.GetRevisionFromProject(ctx, projectId)
}

func (ps *ProjectService) CreateProjectRevision(
	ctx context.Context,
	projectId uuid.UUID,
	curr uuid.UUID,
	dto createRevisionRequest,
) (model.ProjectRevision, *shared.ErrorResponse) {

	projectData, getErr := ps.GetDetailById(ctx, projectId, curr)
	if getErr != nil {
		return model.ProjectRevision{}, getErr
	}

	if projectData.OrderData.ID == curr {
		return model.ProjectRevision{}, shared.NewErrorResponse(403, "you are forbidden to do this operation!")
	}

	if projectData.Status != "in_progress" || projectData.Status == "in_review" {

		return model.ProjectRevision{}, shared.NewErrorResponse(409, "you can't request revision to complete or archieved project")
	}

	return ps.revisionService.CreateNewRevisionn(ctx, projectId, dto)
}

func (ps *ProjectService) UpdateRevisionStatus(ctx context.Context, curr uuid.UUID, projectId uuid.UUID, id uuid.UUID, status string) (model.ProjectRevision, *shared.ErrorResponse) {
	own, mem, err := ps.memberService.validateOwnerOrMember(ctx, curr, projectId)
	if err != nil {
		return model.ProjectRevision{}, shared.NewErrorResponse(500, "something wrong while trying to get user data")
	}

	if !own && !mem {
		return model.ProjectRevision{}, shared.NewErrorResponse(403, "access denied! you can't access this data")
	}

	return ps.revisionService.UpdateProjectRevision(ctx, id, status)
}

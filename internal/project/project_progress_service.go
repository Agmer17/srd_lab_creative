package project

import (
	"context"
	"errors"

	"github.com/Agmer17/srd_lab_creative/internal/shared"
	"github.com/Agmer17/srd_lab_creative/internal/shared/model"
	"github.com/google/uuid"
)

type ProgressService struct {
	repo *ProgresRepository
}

func NewProgressService(rp *ProgresRepository) *ProgressService {
	return &ProgressService{
		repo: rp,
	}
}

func (pss *ProgressService) GetProgressFromProject(ctx context.Context, projId uuid.UUID) ([]model.ProjectProgress, *shared.ErrorResponse) {

	data, err := pss.repo.GetProgressFromProject(ctx, projId)
	if err != nil {
		return []model.ProjectProgress{}, shared.NewErrorResponse(500, "something wrong while trying to get progress data")
	}
	return data, nil
}

func (pss *ProgressService) AddNewProgress(ctx context.Context, projectId uuid.UUID, dto createProgressRequests) ([]model.ProjectProgress, *shared.ErrorResponse) {

	memberId, err := uuid.Parse(dto.ProjectMemberID)
	if err != nil {
		return []model.ProjectProgress{}, shared.NewErrorResponse(400, "invalid member id")
	}
	_, err = pss.repo.CreateProgress(
		ctx,
		projectId,
		dto.Title,
		dto.Weight,
		memberId,
	)

	if err != nil {
		if errors.Is(err, memberNotFound) {
			return []model.ProjectProgress{}, shared.NewErrorResponse(409, "member id not found! you cant assign task to this id")
		}

		if errors.Is(err, projectNotFound) {
			return []model.ProjectProgress{}, shared.NewErrorResponse(409, "project id not found! you can't make task for invalid project id")
		}

		return []model.ProjectProgress{}, shared.NewErrorResponse(500, "something wrong with the server while trying to create progress")
	}

	newData, getErr := pss.GetProgressFromProject(ctx, projectId)
	if getErr != nil {
		return []model.ProjectProgress{}, getErr
	}
	return newData, nil
}

func (pss *ProgressService) DeleteProgress(ctx context.Context, id uuid.UUID) *shared.ErrorResponse {

	err := pss.repo.DeleteProgress(ctx, id)
	if err != nil {
		if errors.Is(err, errProgressNotFound) {
			return shared.NewErrorResponse(404, "progress id not found")
		}

		return shared.NewErrorResponse(500, "something goes wrong whle trying to delete progress")
	}

	return nil
}

func (pss *ProgressService) UpdateProgressData(
	ctx context.Context,
	progId uuid.UUID,
	dto updateProgressRequest,
) (model.ProjectProgress, *shared.ErrorResponse) {

	memberId, err := uuid.Parse(dto.ProjectMemberId)
	if err != nil {
		return model.ProjectProgress{}, shared.NewErrorResponse(400, "invalid uuid for member! please provide a valid uuid")
	}

	data, err := pss.repo.UpdateProgress(ctx, progId, dto, memberId)
	if err != nil {
		if errors.Is(err, errProgressNotFound) {
			return model.ProjectProgress{}, shared.NewErrorResponse(404, "progress not found")
		}
		return model.ProjectProgress{}, shared.NewErrorResponse(500, "failed to update progress")
	}

	newData, err := pss.repo.getProgressDataById(ctx, data.ID)
	if err != nil {
		return model.ProjectProgress{}, shared.NewErrorResponse(500, "failed to get new data, but the update was success : "+err.Error())

	}

	return newData, nil
}

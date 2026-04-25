package project

import (
	"context"
	"errors"

	"github.com/Agmer17/srd_lab_creative/internal/shared"
	"github.com/Agmer17/srd_lab_creative/internal/shared/model"
	"github.com/google/uuid"
)

type RevisionService struct {
	repo *RevisionRepository
}

func NewRevisionService(rp *RevisionRepository) *RevisionService {

	return &RevisionService{
		repo: rp,
	}
}

func (rps *RevisionService) GetRevisionFromProject(ctx context.Context, projectId uuid.UUID) ([]model.ProjectRevision, *shared.ErrorResponse) {

	data, err := rps.repo.GetRevisionFromProject(ctx, projectId)
	if err != nil {
		return []model.ProjectRevision{}, shared.NewErrorResponse(500, "somehting wrong while trying to get revision data")
	}
	return data, nil
}

func (rps *RevisionService) CreateNewRevisionn(ctx context.Context, projectId uuid.UUID, dto createRevisionRequest) (model.ProjectRevision, *shared.ErrorResponse) {

	newData, err := rps.repo.CreateRevision(ctx, projectId, dto.Title, dto.Reason)
	if err != nil {
		if errors.Is(err, projectNotFound) {
			return model.ProjectRevision{}, shared.NewErrorResponse(409, "project id not found! you can't make revision from unexisting project")
		}
	}
	return newData, nil
}

func (rps *RevisionService) UpdateProjectRevision(ctx context.Context, id uuid.UUID, status string) (model.ProjectRevision, *shared.ErrorResponse) {

	data, err := rps.repo.UpdateRevisionStatus(ctx, id, status)
	if err != nil {
		if errors.Is(err, errRevisionNotFound) {
			return model.ProjectRevision{}, shared.NewErrorResponse(404, "revision id not found")
		}

		return model.ProjectRevision{}, shared.NewErrorResponse(500, "somehting wrong while trying to update revision data")
	}

	return data, nil
}

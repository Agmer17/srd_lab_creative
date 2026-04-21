package project

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Agmer17/srd_lab_creative/internal/db/sqlcgen"
	"github.com/Agmer17/srd_lab_creative/internal/shared/model"
	"github.com/google/uuid"
)

var projectNotFound = errors.New("projects not found")

type ProjectRepository struct {
	db *sqlcgen.Queries
}

func NewProjectRepository(q *sqlcgen.Queries) *ProjectRepository {

	return &ProjectRepository{
		db: q,
	}
}

func (pr *ProjectRepository) CreateProjects(ctx context.Context, dto createProjectRequest) (model.Project, error) {

	orderId, _ := uuid.Parse(dto.OrderId)
	startDate := time.Now()

	data, err := pr.db.CreateProject(ctx, sqlcgen.CreateProjectParams{
		OrderID:              orderId,
		Name:                 dto.Name,
		Description:          dto.Description,
		Status:               dto.Status,
		AllowedRevisionCount: dto.AllowedRevision,
		StartDate:            &startDate,
		EndDate:              dto.Deadline,
	})

	if err != nil {
		return model.Project{}, err
	}

	dataModel := model.MapProjectDataToModel(data)
	dataModel.ProjectMembers = make([]model.ProjectMember, 0)
	dataModel.Progress = make([]model.ProjectProgress, 0)

	return dataModel, nil
}

func (pr *ProjectRepository) DeleteProjects(ctx context.Context, id uuid.UUID) error {

	aff, err := pr.db.DeleteProject(ctx, id)

	if err != nil {
		return err
	}

	if aff == 0 {
		return projectNotFound
	}

	return nil
}

func (pr *ProjectRepository) GetAllProjects(ctx context.Context) ([]model.Project, error) {

	data, err := pr.db.ListProjects(ctx)
	if err != nil {
		return []model.Project{}, err
	}

	var listProjects []model.Project = make([]model.Project, len(data))

	for i, v := range data {
		tempProject := model.Project{
			ID:                   v.ID,
			OrderID:              v.OrderID,
			Name:                 v.Name,
			Description:          v.Description,
			Status:               v.Status,
			AllowedRevisionCount: v.AllowedRevisionCount,
			ActualStartDate:      v.ActualStartDate,
			EndDate:              v.EndDate,
			UpdatedAt:            v.UpdatedAt,
		}

		var members []model.ProjectMember = make([]model.ProjectMember, 0)
		err := json.Unmarshal(v.ProjectMembers, &members)
		if err != nil {
			fmt.Println(err)
			return []model.Project{}, err
		}

		var progresses []model.ProjectProgress
		err = json.Unmarshal(v.Progress, &progresses)
		if err != nil {
			fmt.Println(err)
			return []model.Project{}, err
		}

		tempProject.Progress = progresses
		tempProject.ProjectMembers = members

		listProjects[i] = tempProject
	}

	return listProjects, nil
}

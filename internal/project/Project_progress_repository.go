package project

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/Agmer17/srd_lab_creative/internal/db/sqlcgen"
	"github.com/Agmer17/srd_lab_creative/internal/shared/model"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var errProgressNotFound = errors.New("progres with this id not found")

type ProgresRepository struct {
	db *sqlcgen.Queries
}

func NewProgresRepository(q *sqlcgen.Queries) *ProgresRepository {
	return &ProgresRepository{
		db: q,
	}
}

func (prp *ProgresRepository) GetProgressFromProject(ctx context.Context, prj uuid.UUID) ([]model.ProjectProgress, error) {

	data, err := prp.db.GetProgressByProject(ctx, prj)
	if err != nil {
		return []model.ProjectProgress{}, err
	}

	var listModel []model.ProjectProgress = make([]model.ProjectProgress, len(data))
	for i, v := range data {
		listModel[i] = model.ProjectProgress{
			ID:          v.ProgressID,
			ProjectID:   v.ProjectID,
			Title:       v.Title,
			Weight:      v.Weight,
			IsCompleted: v.IsCompleted,
			CreatedAt:   v.ProgressCreatedAt,
		}

		var projectMember model.ProjectMember
		err := json.Unmarshal(v.ProjectMember, &projectMember)
		if err != nil {
			return []model.ProjectProgress{}, err
		}

		listModel[i].ProjectMember = projectMember
	}
	return listModel, nil
}

func (prp *ProgresRepository) CreateProgress(
	ctx context.Context,
	projectId uuid.UUID,
	title string,
	weight float64,
	memberId uuid.UUID,
) (model.ProjectProgress, error) {

	data, err := prp.db.CreateProgress(ctx, sqlcgen.CreateProgressParams{
		ProjectID:       projectId,
		Title:           title,
		Weight:          weight,
		ProjectMemberID: memberId,
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {

			case pgerrcode.ForeignKeyViolation:
				detail := pgErr.Detail
				if strings.Contains(detail, "project_member_id") {
					return model.ProjectProgress{}, memberNotFound
				}
				if strings.Contains(detail, "project_id") {
					return model.ProjectProgress{}, projectNotFound
				}
			}
		}

		return model.ProjectProgress{}, err
	}

	return model.ProjectProgress{
		ID:        data.ID,
		ProjectID: data.ProjectID,
		ProjectMember: model.ProjectMember{
			ID: data.ProjectMemberID,
		},
		Title:       data.Title,
		Weight:      data.Weight,
		IsCompleted: data.IsCompleted,
		CreatedAt:   data.CreatedAt,
	}, nil

}

func (prp *ProgresRepository) DeleteProgress(ctx context.Context, progresId uuid.UUID) error {

	aff, err := prp.db.DeleteProgress(ctx, progresId)
	if err != nil {
		return err
	}

	if aff == 0 {
		return errProgressNotFound
	}

	return nil

}

func (prp *ProgresRepository) UpdateProgress(
	ctx context.Context,
	id uuid.UUID,
	dto updateProgressRequest,
	memberId uuid.UUID,
) (model.ProjectProgress, error) {

	data, err := prp.db.UpdateProgress(ctx, sqlcgen.UpdateProgressParams{
		ID:              id,
		Title:           dto.Title,
		Weight:          dto.Weight,
		IsCompleted:     dto.IsComplete,
		ProjectMemberID: memberId,
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ProjectProgress{}, errProgressNotFound
		}

		return model.ProjectProgress{}, err
	}

	return model.ProjectProgress{
		ID:          data.ID,
		ProjectID:   data.ProjectID,
		Title:       data.Title,
		Weight:      data.Weight,
		IsCompleted: data.IsCompleted,
		CreatedAt:   data.CreatedAt,
		ProjectMember: model.ProjectMember{
			ID: data.ProjectMemberID,
		},
	}, nil
}

func (prp *ProgresRepository) getProgressDataById(ctx context.Context, id uuid.UUID) (model.ProjectProgress, error) {

	data, err := prp.db.GetProgressById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ProjectProgress{}, errProgressNotFound
		}
		return model.ProjectProgress{}, err
	}

	var memberData model.ProjectMember
	err = json.Unmarshal(data.ProjectMember, &memberData)
	if err != nil {
		return model.ProjectProgress{}, err
	}

	return model.ProjectProgress{
		ID:            data.ProgressID,
		ProjectID:     data.ProjectID,
		ProjectMember: memberData,
		Title:         data.Title,
		Weight:        data.Weight,
		IsCompleted:   data.IsCompleted,
		CreatedAt:     data.ProgressCreatedAt,
	}, nil
}

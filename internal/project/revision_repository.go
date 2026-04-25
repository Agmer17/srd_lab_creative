package project

import (
	"context"
	"errors"

	"github.com/Agmer17/srd_lab_creative/internal/db/sqlcgen"
	"github.com/Agmer17/srd_lab_creative/internal/shared/model"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var errRevisionNotFound = errors.New("revision not found")

type RevisionRepository struct {
	db *sqlcgen.Queries
}

func NewRevisionRepository(q *sqlcgen.Queries) *RevisionRepository {
	return &RevisionRepository{
		db: q,
	}
}

func (rvp *RevisionRepository) GetRevisionFromProject(ctx context.Context, projectId uuid.UUID) ([]model.ProjectRevision, error) {
	data, err := rvp.db.GetRevisionsByProject(ctx, projectId)
	if err != nil {
		return []model.ProjectRevision{}, err
	}
	model := model.ListGenToRevision(data)
	return model, nil
}

func (rvp *RevisionRepository) CreateRevision(ctx context.Context, projectId uuid.UUID, title string, reason string) (model.ProjectRevision, error) {

	data, err := rvp.db.CreateRevisionRequest(ctx, sqlcgen.CreateRevisionRequestParams{
		ProjectID: projectId,
		Title:     title,
		Reason:    reason,
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {

			if pgErr.Code == pgerrcode.ForeignKeyViolation {
				return model.ProjectRevision{}, projectNotFound
			}

		}
		return model.ProjectRevision{}, err
	}

	return model.MapRevisionToModel(data), nil
}

func (rvp *RevisionRepository) UpdateRevisionStatus(ctx context.Context, id uuid.UUID, status string) (model.ProjectRevision, error) {

	data, err := rvp.db.UpdateRevisionStatus(ctx, sqlcgen.UpdateRevisionStatusParams{
		ID: id,
		Status: sqlcgen.NullRevisionStatusEnum{
			RevisionStatusEnum: sqlcgen.RevisionStatusEnum(status),
			Valid:              true,
		},
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ProjectRevision{}, errRevisionNotFound
		}

		return model.ProjectRevision{}, err
	}

	return model.MapRevisionToModel(data), nil
}

package category

import (
	"context"
	"errors"

	"github.com/Agmer17/srd_lab_creative/internal/db/sqlcgen"
	"github.com/Agmer17/srd_lab_creative/internal/shared/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var errCategoryNotFound = errors.New("no category found")

type CategoryRepository struct {
	db *sqlcgen.Queries
}

func NewCategoryRepository(q *sqlcgen.Queries) *CategoryRepository {

	return &CategoryRepository{
		db: q,
	}
}

func (cr *CategoryRepository) GetAllCategories(ctx context.Context) ([]model.Category, error) {
	data, err := cr.db.ListCategories(ctx, sqlcgen.ListCategoriesParams{
		PageOffset: 0,
		PageLimit:  10000,
	})
	if err != nil {
		return []model.Category{}, err
	}
	return model.MapListToCategoryModel(data), nil
}

func (cr *CategoryRepository) CreateCategories(ctx context.Context, c createCategoryRequest) (model.Category, error) {

	data, err := cr.db.CreateCategory(ctx, sqlcgen.CreateCategoryParams{
		Name:        c.Name,
		Slug:        c.Slug,
		Description: c.Desc,
	})

	if err != nil {
		return model.Category{}, err
	}

	return model.MapToCategoryModel(data), nil
}

func (cr *CategoryRepository) GetCategoryById(ctx context.Context, id uuid.UUID) (model.Category, error) {

	data, err := cr.db.GetCategoryByID(ctx, id)
	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return model.Category{}, errCategoryNotFound

		}
		return model.Category{}, err
	}

	return model.MapToCategoryModel(data), nil
}

func (cr *CategoryRepository) GetCategoryBySlug(ctx context.Context, slug string) (model.Category, error) {

	data, err := cr.db.GetCategoryBySlug(ctx, slug)
	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return model.Category{}, errCategoryNotFound

		}
		return model.Category{}, err
	}

	return model.MapToCategoryModel(data), nil
}

func (cr *CategoryRepository) UpdateCategory(ctx context.Context, id uuid.UUID, u updateCategoryRequest) (model.Category, error) {

	data, err := cr.db.UpdateCategory(ctx, sqlcgen.UpdateCategoryParams{
		Name:        u.Name,
		Slug:        u.Slug,
		Description: u.Desc,
		ID:          id,
	})

	if err != nil {
		return model.Category{}, err
	}

	return model.MapToCategoryModel(data), nil
}

func (cr *CategoryRepository) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	return cr.db.DeleteCategory(ctx, id)
}

func (cr *CategoryRepository) ExistByCategorySlyg(ctx context.Context, slug string) (bool, error) {
	return cr.db.ExistCategoryBySlug(ctx, slug)
}

func (cr *CategoryRepository) SearchCategory(ctx context.Context, query *string) ([]model.Category, error) {
	data, err := cr.db.SearchCategories(ctx, sqlcgen.SearchCategoriesParams{
		Query:      query,
		PageOffset: 0,
		PageLimit:  10000,
	})

	if err != nil {
		return []model.Category{}, err
	}

	return model.MapListToCategoryModel(data), nil
}

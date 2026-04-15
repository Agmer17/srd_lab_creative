package category

import (
	"context"
	"errors"

	"github.com/Agmer17/srd_lab_creative/internal/shared"
	"github.com/Agmer17/srd_lab_creative/internal/shared/model"
	"github.com/google/uuid"
)

type CategoryService struct {
	repo *CategoryRepository
}

func NewCategoryService(rp *CategoryRepository) *CategoryService {
	return &CategoryService{
		repo: rp,
	}
}

func (cs *CategoryService) GetAllCategories(ctx context.Context) ([]model.Category, *shared.ErrorResponse) {
	data, err := cs.repo.GetAllCategories(ctx)
	if err != nil {
		return []model.Category{}, shared.NewErrorResponse(500, "something wrong with the server right now! try again later")
	}
	return data, nil
}

func (cs *CategoryService) GetCategoryById(ctx context.Context, id uuid.UUID) (model.Category, *shared.ErrorResponse) {
	data, err := cs.repo.GetCategoryById(ctx, id)
	if err != nil {
		if errors.Is(err, errCategoryNotFound) {
			return model.Category{}, shared.NewErrorResponse(404, "no category with this id was found")
		}
		return model.Category{}, shared.NewErrorResponse(500, "something went wrong while getting category! try again later")
	}
	return data, nil
}

func (cs *CategoryService) GetCategoryBySlug(ctx context.Context, slug string) (model.Category, *shared.ErrorResponse) {
	data, err := cs.repo.GetCategoryBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, errCategoryNotFound) {
			return model.Category{}, shared.NewErrorResponse(404, "no category with this slug was found")
		}
		return model.Category{}, shared.NewErrorResponse(500, "something went wrong while getting category! try again later")
	}
	return data, nil
}

func (cs *CategoryService) CreateCategory(ctx context.Context, c createCategoryRequest) (model.Category, *shared.ErrorResponse) {
	// cek dulu slug udah ada apa belum
	exists, err := cs.repo.ExistByCategorySlyg(ctx, c.Slug)
	if err != nil {
		return model.Category{}, shared.NewErrorResponse(500, "something went wrong while checking slug! try again later")
	}
	if exists {
		return model.Category{}, shared.NewErrorResponse(409, "category with this slug already exists")
	}

	data, err := cs.repo.CreateCategories(ctx, c)
	if err != nil {
		return model.Category{}, shared.NewErrorResponse(500, "something went wrong while creating category! try again later")
	}
	return data, nil
}

func (cs *CategoryService) UpdateCategory(ctx context.Context, id uuid.UUID, u updateCategoryRequest) (model.Category, *shared.ErrorResponse) {
	_, err := cs.repo.GetCategoryById(ctx, id)
	if err != nil {
		if errors.Is(err, errCategoryNotFound) {
			return model.Category{}, shared.NewErrorResponse(404, "no category with this id was found")
		}
		return model.Category{}, shared.NewErrorResponse(500, "something went wrong while checking category! try again later")
	}

	if u.Slug != nil {
		exists, err := cs.repo.ExistByCategorySlyg(ctx, *u.Slug)
		if err != nil {
			return model.Category{}, shared.NewErrorResponse(500, "something went wrong while checking slug! try again later")
		}
		if exists {
			return model.Category{}, shared.NewErrorResponse(409, "category with this slug already exists")
		}
	}

	data, err := cs.repo.UpdateCategory(ctx, id, u)
	if err != nil {
		// fmt.Println(err)
		return model.Category{}, shared.NewErrorResponse(500, "something went wrong while updating category! try again later")
	}
	return data, nil
}

func (cs *CategoryService) DeleteCategory(ctx context.Context, id uuid.UUID) *shared.ErrorResponse {
	_, err := cs.repo.GetCategoryById(ctx, id)
	if err != nil {
		if errors.Is(err, errCategoryNotFound) {
			return shared.NewErrorResponse(404, "no category with this id was found")
		}
		return shared.NewErrorResponse(500, "something went wrong while checking category! try again later")
	}

	err = cs.repo.DeleteCategory(ctx, id)
	if err != nil {
		return shared.NewErrorResponse(500, "something went wrong while deleting category! try again later")
	}
	return nil
}

func (cs *CategoryService) SearchCategory(ctx context.Context, query *string) ([]model.Category, *shared.ErrorResponse) {
	data, err := cs.repo.SearchCategory(ctx, query)
	if err != nil {
		return []model.Category{}, shared.NewErrorResponse(500, "something went wrong while searching categories! try again later")
	}
	return data, nil
}

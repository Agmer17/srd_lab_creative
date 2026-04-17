package product

import (
	"context"
	"errors"

	"github.com/Agmer17/srd_lab_creative/internal/shared"
	"github.com/Agmer17/srd_lab_creative/internal/shared/model"
	"github.com/google/uuid"
)


type ProductService struct{
	repo *ProductRepository
}

func NewProductService(pr *ProductRepository) *ProductService{
	return &ProductService{
		repo: pr,
	}
}

func (ps *ProductService) GetAllProduct (ctx context.Context) ([]model.Product, *shared.ErrorResponse){
	data, err := ps.repo.GetAllProduct(ctx);
	if err != nil{
		return []model.Product{},shared.NewErrorResponse(500,"something wrong with the server while trying to get product data");
	}
	return data,nil;
}

func (ps *ProductService) GetProductById (ctx context.Context, id uuid.UUID) (model.Product, *shared.ErrorResponse){
	data, err := ps.repo.GetProduct(ctx,id);
	if err != nil{
		if errors.Is(err,errProductNotFound){
			return model.Product{},shared.NewErrorResponse(404, "no product with this id was found");
		}
		return model.Product{},shared.NewErrorResponse(500,"something wrong with the server while trying to get product data");
	}
	return data,nil;
}

func (ps *ProductService) DeleteProduct(ctx context.Context, id uuid.UUID) *shared.ErrorResponse{
	_, err := ps.repo.GetProduct(ctx,id);
	if err != nil{
		return shared.NewErrorResponse(404,"product data with this id is not found")
	}
	err = ps.repo.DeleteProduct(ctx,id);
	if err != nil{
		return shared.NewErrorResponse(500, "something wrong with the server while trying to delete the product data");
	}
	return nil;

}

func (ps *ProductService) CreateProduct(ctx context.Context, req createProductRequest) (model.Product, *shared.ErrorResponse){
	isUsed , err := ps.repo.CheckProductSlugExists(ctx,req.Slug);

	if err != nil {
		return model.Product{},shared.NewErrorResponse(500, "something went wrong while checking slug! try again later");
	}
	if isUsed {
		return model.Product{}, shared.NewErrorResponse(409, "product with this slug already exists")
	}

	data,err := ps.repo.CreateProduct(ctx,req);

	if err != nil{
		return model.Product{},shared.NewErrorResponse(500, "something went wrong while creating product! try again later");
	}

	return data,nil;
}

func (ps *ProductService) UpdateProduct (ctx context.Context,id uuid.UUID,req updateProductRequest) (model.Product, *shared.ErrorResponse){
	
	data_lama,err := ps.repo.GetProduct(ctx,id);
	
	if err != nil {
		if errors.Is(err,errProductNotFound) {
			return model.Product{}, shared.NewErrorResponse(404, "no product with this id was found");
		}
		return model.Product{}, shared.NewErrorResponse(500, "something went wrong while checking product! try again later");
	}

	if req.Slug != nil && *req.Slug != data_lama.Slug {
		exists, errSlug := ps.repo.CheckProductSlugExists(ctx,*req.Slug);
		if errSlug != nil {
			return model.Product{}, shared.NewErrorResponse(500, "something went wrong while checking slug! try again later");
		}
		if exists {
			return model.Product{}, shared.NewErrorResponse(409, "product with this slug already exists");
		}
	}
	
	
	data,err := ps.repo.UpdateProduct(ctx,id,req);
	if err != nil{
		return model.Product{},shared.NewErrorResponse(500, "something went wrong while updating the product, please try again");
	}
	return data,nil;
}


	

package product

import (
	"context"
	"errors"

	"github.com/Agmer17/srd_lab_creative/internal/db/sqlcgen"
	"github.com/Agmer17/srd_lab_creative/internal/shared/model"
	"github.com/google/uuid"
)


type ProductRepository struct {
	db *sqlcgen.Queries
}

func NewProductRepository(q *sqlcgen.Queries) *ProductRepository{
	return &ProductRepository{
		db: q,
	}
}

var errProductNotFound = errors.New("product not found")

func (pr *ProductRepository) GetAllProduct(ctx context.Context) ([]model.Product,error){
	data,err := pr.db.GetAllProduct(ctx);
	if err != nil{
		return []model.Product{},err;
	}
	return model.MapListToProductModel(data),nil;
}

func (pr *ProductRepository) GetProduct(ctx context.Context, id uuid.UUID) (model.Product,error){
	data,err := pr.db.GetProductById(ctx,id);
	if err != nil{
		return model.Product{},err;
	}
	return model.MapToProductModel(data),nil;
}

func (pr *ProductRepository) DeleteProduct(ctx context.Context, id uuid.UUID) error{
	err := pr.db.DeleteProduct(ctx,id);
	return err;
}

func (pr *ProductRepository) CreateProduct (ctx context.Context, req createProductRequest) (model.Product,error){
	data,err := pr.db.CreateProduct(ctx,sqlcgen.CreateProductParams{
		Name: req.Name,
		Slug: req.Slug,
		Description: &req.Description,
		Price: req.Price,
		Status: req.Status,
		IsFeatured: *req.IsFeatured,
	});
	if err != nil{
		return model.Product{},err;
	}
	return model.MapToProductModel(data),nil;
}

func (pr *ProductRepository) UpdateProduct (ctx context.Context,id uuid.UUID,req updateProductRequest) (model.Product,error){
	data,err := pr.db.UpdateProduct(ctx,sqlcgen.UpdateProductParams{
		ID: id,
		Name: req.Name,
		Slug: req.Slug,
		Description: req.Description,
		Price: req.Price,
		Status: req.Status,
		IsFeatured: req.IsFeatured,
	});
	if err != nil{
		return model.Product{},err;
	}
	return model.MapToProductModel(data),nil;
}

func (pr *ProductRepository) CheckProductSlugExists(ctx context.Context, slug string) (bool, error) {
	// Menghasilkan 'true' kalau slug kembar 100% sudah beneran kepake
	exists, err := pr.db.CheckProductSlugExists(ctx, slug);
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (pr *ProductRepository) GetProductBySlug (ctx context.Context, slug string) (model.Product,error){
	data,err := pr.db.GetProductBySlug(ctx,slug);
	if err != nil {
		return model.Product{},err;
	}
	return model.MapToProductModel(data),nil;
}


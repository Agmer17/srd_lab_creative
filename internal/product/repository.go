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

var errProductNotFound = errors.New("product not found");
var errProductIdNotFound = errors.New("product id not found");

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



// berhubungan sama product image

func(pr *ProductRepository) CountImageofProductId(ctx context.Context, productId uuid.UUID) (int,error){
	count,err := pr.db.GetTotalImageOfProductId(ctx,productId);
	if err != nil{
		return int(count),err;
	}
	return int(count),nil;
}

func (pr *ProductRepository) CreateProductImage (ctx context.Context, productID uuid.UUID, urllist []string, primarylist []bool, sortorderlist []int32) ([]model.ProductImage,error){
	result,err := pr.db.CreateProductImage(ctx,sqlcgen.CreateProductImageParams{
		ProductID: productID,
		Column2: urllist,
		Column3: primarylist,
		Column4: sortorderlist,
	})
	if err != nil{
		return []model.ProductImage{},err;
	}
	return model.MapToProductImageListModel(result), nil;
}

func (pr *ProductRepository) GetAllImageIdAndOrderByProductId(ctx context.Context, productID uuid.UUID)([]UpdateProductImageSort,error){
	data,err := pr.db.GetImageIdsAndOrderByProductId(ctx,productID);
	if err != nil{
		return []UpdateProductImageSort{},err;
	}
	// SORRY AGMER CONVERSION DI SINI GA BIKIN FUNCTION :(
    result := make([]UpdateProductImageSort, len(data))
    for i, row := range data {
        idStr := row.ID.String()
        order := int(row.SortOrder)
        result[i] = UpdateProductImageSort{
            ImageId: idStr,   
            SortOrder: order,
        }
    }
    return result, nil
}

func (pr *ProductRepository) BulkUpdateImageSortOrder(ctx context.Context, ids []uuid.UUID, orders []int32)error{
	err := pr.db.ImageIdOrderChange(ctx,sqlcgen.ImageIdOrderChangeParams{
		Column1: ids,
		Column2: orders,
	});
	if err != nil{
		return err;
	}
	return nil;

}


func (pr *ProductRepository) DeleteImageById(ctx context.Context, imageId uuid.UUID)error{
	err := pr.db.DeleteProductImageByImageId(ctx,imageId);
	if err != nil{
		return err;
	}
	return nil;
}

func (pr *ProductRepository) GetProductImageById(ctx context.Context, imageId uuid.UUID)(model.ProductImage,error){
	data,err := pr.db.GetProductImageByImageId(ctx,imageId);
	if err != nil {
		return model.ProductImage{},err;
	}
	return model.MapToProductImageModel(data),nil;
}

func (pr *ProductRepository) GetAllProductImage(ctx context.Context, productId uuid.UUID) ([]model.ProductImage, error){
	data,err := pr.db.GetAllProductImageByProductId(ctx,productId);
	if err != nil{
		return []model.ProductImage{},err;
	}
	return model.MapToProductImageListModel(data),nil;
	
}

func (pr *ProductRepository) DeleteAllProductImage (ctx context.Context, productId uuid.UUID) (error){
	err := pr.db.DeleteProductImageByProductId(ctx,productId);
	if err != nil{
		return err;
	}
	return nil;
}
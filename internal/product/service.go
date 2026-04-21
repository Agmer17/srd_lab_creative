package product

import (
	"context"
	"errors"
	"log"
	"mime/multipart"
	"path/filepath"

	"github.com/Agmer17/srd_lab_creative/internal/shared"
	"github.com/Agmer17/srd_lab_creative/internal/shared/model"
	"github.com/Agmer17/srd_lab_creative/internal/storage"
	"github.com/google/uuid"
)


type ProductService struct{
	repo *ProductRepository
	storage *storage.FileStorage
}

func NewProductService(pr *ProductRepository, stg *storage.FileStorage) *ProductService{
	return &ProductService{
		repo: pr,
		storage: stg,
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

func (ps *ProductService) GetProductBySlug (ctx context.Context, slug string) (model.Product, *shared.ErrorResponse){
	data,err := ps.repo.GetProductBySlug(ctx,slug);
	if err != nil{
		if errors.Is(err,errProductNotFound){
			return model.Product{},shared.NewErrorResponse(404,"No product with this slug was found");
		}
		return model.Product{},shared.NewErrorResponse(500, "something went wrong while checking product! try again later");
	}
	return data,nil;
}


// product image related

func (ps *ProductService) AddImage(ctx context.Context, productId uuid.UUID, images []*multipart.FileHeader) ([]model.ProductImage, *shared.ErrorResponse){
	// verify dulu product idnya ada apa engga
	_,err := ps.repo.GetProduct(ctx,productId);
	if err != nil{
		if errors.Is(err,errProductIdNotFound){
			return []model.ProductImage{},shared.NewErrorResponse(404, "no product with this id was found");
		}
		return []model.ProductImage{},shared.NewErrorResponse(500,"something wrong with the server while trying to get product id");
	}

	//cek udah berapa banyak gambarnya
	count,err := ps.repo.CountImageofProductId(ctx,productId);
	if err != nil{
		return []model.ProductImage{},shared.NewErrorResponse(500, "something wrong with the server while trying to count image");
	}

	// verify,compress dan save
	filenames,err := ps.storage.SaveAllPublicFile(ctx,images,"products");
	if err != nil{
		return []model.ProductImage{},shared.NewErrorResponse(500,"something wrong when saving and processing the image");
	}

	// bikin array tiap column
	imageUrls  := make([]string, len(filenames))
	isPrimaries := make([]bool, len(filenames))
	sortOrders  := make([]int32, len(filenames))
	
	
	// ubah nama file biasa jadi relatif
	for i, file := range filenames{
		imageUrl := "/uploads/public/products/" + file;
		isPrimary := count == 0 && i == 0;
		sortOrder := int(count+i);
		imageUrls[i] = imageUrl;
		isPrimaries[i] = isPrimary;
		sortOrders[i] = int32(sortOrder);
	}


	// add image ke product id tersebut
	data,err := ps.repo.CreateProductImage(ctx,productId,imageUrls,isPrimaries,sortOrders);
	if err != nil{
		return []model.ProductImage{}, shared.NewErrorResponse(500, "Something wrong when saving product image to database");
	}
	// return
		return data,nil;

}


func (ps *ProductService) SortImageOrder(ctx context.Context, productId uuid.UUID, req []UpdateProductImageSort) ([]UpdateProductImageSort, *shared.ErrorResponse){
	// validasi product id
	// verify dulu product idnya ada apa engga
	_,err := ps.repo.GetProduct(ctx,productId);
	if err != nil{
		if errors.Is(err,errProductIdNotFound){
			return []UpdateProductImageSort{},shared.NewErrorResponse(404, "no product with this id was found");
		}
		return []UpdateProductImageSort{},shared.NewErrorResponse(500,"something wrong with the server while trying to get product id");
	}
	
	// dapetin order image yang ada sekarang
	oldData,err := ps.repo.GetAllImageIdAndOrderByProductId(ctx,productId);
	if err != nil{
		return []UpdateProductImageSort{},shared.NewErrorResponse(500,"something wrong with the server while trying to get image id and order");
	}

	// sorting
    for _, change := range req {
        indexGanti := change.SortOrder;
        idGanti    := change.ImageId;
        for j := range oldData {
            if oldData[j].ImageId == idGanti {
                oldPos := oldData[j].SortOrder;
                oldData[j].SortOrder = indexGanti;
                for k := range oldData {
                    if k != j && oldData[k].SortOrder == indexGanti {
                        oldData[k].SortOrder = oldPos;
                        break;
                    }
                }
                break;
            }
        }
    }

	// mecah jadi 2 buat query
	ids    := make([]uuid.UUID, len(oldData))
	orders := make([]int32, len(oldData))
	for i, item := range oldData {
    	ids[i],_    = uuid.Parse(item.ImageId)
    	orders[i]   = int32(item.SortOrder)
	}

    // query ke db — update sort_order semua gambar yang berubah
    updateErr := ps.repo.BulkUpdateImageSortOrder(ctx,ids,orders);
    if updateErr != nil {
        return []UpdateProductImageSort{}, shared.NewErrorResponse(500, "something wrong while updating image order")
    }
    // return hasil final ke user
    return oldData, nil

}

func (ps *ProductService) DeleteSpecificProductImage(ctx context.Context,imageId uuid.UUID) (*shared.ErrorResponse){
	
	// ambil dulu url image
	data,err := ps.repo.GetProductImageById(ctx,imageId);
	if err != nil {
		return shared.NewErrorResponse(500, "Something wrong with the server while geting the product image data");
	}
	
	// delete di db
	err = ps.repo.DeleteImageById(ctx,imageId);
	if err != nil{
		return shared.NewErrorResponse(500, "Something wrong with the server while deleting the product image");
	}

	fileErr := ps.storage.DeletePublicFile(filepath.Base(data.ImageUrl),"products");
	if fileErr != nil{
		log.Println("Warning: Failed to delete product image file in disk", fileErr);
	}

	return nil;
}


func (ps *ProductService) DeleteAllProductImage (ctx context.Context, productId uuid.UUID) (*shared.ErrorResponse){
	// ambil slice url image
	data,err := ps.repo.GetAllProductImage(ctx,productId);
	if err != nil {
		return shared.NewErrorResponse(500, "Something wrong with the server while geting the product image data");
	}
	// delete di db 
	err = ps.repo.DeleteAllProductImage(ctx,productId);
	if err != nil{
		return shared.NewErrorResponse(500, "Something wrong with the server while deleting the product image");
	}
	for _, img := range data {
		fileErr := ps.storage.DeletePublicFile(filepath.Base(img.ImageUrl), "products")
		if fileErr != nil {
			log.Printf("Warning: failed to delete file %s: %v", img.ImageUrl, fileErr)
			// lanjut terus meski 1 gagal
		}
	}
	return nil

}

func (ps *ProductService) GetProductImage (ctx context.Context, imageId uuid.UUID) (model.ProductImage,*shared.ErrorResponse){
	data,err := ps.repo.GetProductImageById(ctx,imageId);
	if err != nil{
		return model.ProductImage{},shared.NewErrorResponse(500,"Something wrong with the server while getting the product image");
	}
	return data,nil;
}

func (ps *ProductService) GetAllProductImage (ctx context.Context, productId uuid.UUID) ([]model.ProductImage,*shared.ErrorResponse){
	data, err := ps.repo.GetAllProductImage(ctx,productId);
	if err != nil{
		return []model.ProductImage{},shared.NewErrorResponse(500,"Something wrong with the server while getting the product image");
	}
	return data,nil;
}
	

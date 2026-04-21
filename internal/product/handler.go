package product

import (
	"github.com/Agmer17/srd_lab_creative/internal/shared"
	"github.com/Agmer17/srd_lab_creative/internal/shared/middleware"
	"github.com/Agmer17/srd_lab_creative/pkg"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductHandler struct {
	svc *ProductService
}

func NewProductHandler(s *ProductService) *ProductHandler {
	return &ProductHandler{
		svc: s,
	}
}

// Product

func (pth *ProductHandler) HandleGetAllProducts(c *gin.Context) {
	data, err := pth.svc.GetAllProduct(c)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}
	c.JSON(200, shared.NewSuccessResponse(200, "Successfully getting the product data", data))
}

func (pth *ProductHandler) HandleGetProductById(c *gin.Context) {
	path := c.Param("id")
	id, err := uuid.Parse(path)
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid id params"))
		return
	}

	data, getErr := pth.svc.GetProductById(c, id)
	if getErr != nil {
		c.JSON(getErr.Code, getErr)
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "succesffuly getting the product data", data))
}

func (pth *ProductHandler) PostCreateProduct(c *gin.Context) {
	var req createProductRequest


	if err := c.ShouldBindJSON(&req); err != nil {
		vldMsg, ok := pkg.ParseValidationErrors(err)
		if !ok {
			c.JSON(400, shared.NewErrorResponse(400, "invalid request body! please provide valid body for create product request"))
			return
		}
		c.JSON(400, shared.NewErrorResponse(400, vldMsg))
		return
	}

	data, insErr := pth.svc.CreateProduct(c, req)
	if insErr != nil {
		c.JSON(insErr.Code, insErr)
		return
	}
	c.JSON(200, shared.NewSuccessResponse(200, "New product successfully created", data))
}

func (pth *ProductHandler) PatchUpdateProduct(c *gin.Context) {
	var req updateProductRequest

	path := c.Param("id")
	id, err := uuid.Parse(path)
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid id params"))
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		vldMsg, ok := pkg.ParseValidationErrors(err)
		if !ok {
			c.JSON(400, shared.NewErrorResponse(400, "Invalid request body! please provide valid body for update product request"))
			return
		}
		c.JSON(400, shared.NewErrorResponse(400, vldMsg))
		return
	}

	data, updErr := pth.svc.UpdateProduct(c, id, req)
	if updErr != nil {
		c.JSON(updErr.Code, updErr)
		return
	}
	c.JSON(200, shared.NewSuccessResponse(200, "product successfully updated", data))

}

func (pth *ProductHandler) HandleGetProductBySlug(c *gin.Context) {
	slug := c.Param("slug")

	data, err := pth.svc.GetProductBySlug(c, slug)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "succesffuly getting the product data", data))

}

func (pth *ProductHandler) DeleteProductHandler(c *gin.Context) {
	path := c.Param("id")
	id, err := uuid.Parse(path)
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid id params"))
		return
	}

	delErr := pth.svc.DeleteProduct(c, id)
	if delErr != nil {
		c.JSON(delErr.Code, delErr)
		return
	}
	c.JSON(200, shared.NewSuccessResponse(200, "product successfully deleted", nil))

}

// Product Images

func (pth *ProductHandler) PostAddProductImage(c *gin.Context){

	// nerima data dari multipart form
	form,err := c.MultipartForm()
	if err != nil{
		c.JSON(400, shared.NewErrorResponse(400, "Failed to process form data"));
		return;
	}

	// ngambil specific yang ada key images nya
	files := form.File["images"]
	if len(files) == 0 {
		c.JSON(400, shared.NewErrorResponse(400, "No image was received"));
		return;
	}
	// validasi input id
	path := c.Param("product_id")
	id, err := uuid.Parse(path)
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid id params"))
		return
	}
	// query nambahin image
	data, insErr := pth.svc.AddImage(c,id,files);
	if insErr != nil {
		c.JSON(insErr.Code, insErr);
		return;
	}
	
	// response
	c.JSON(200, shared.NewSuccessResponse(200, "New image successfully added", data));
	
}

func (pth *ProductHandler) PatchUpdateProductImageOrder(c *gin.Context){
	// terima id
	idparams := c.Param("product_id");
	
	// validasi id
	productid, err := uuid.Parse(idparams)
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid id params"))
		return
	}

	// nerima json yang isinya pasangan antara id image dan urutannya
	var req []UpdateProductImageSort;

	// validasi
	if err := c.ShouldBindJSON(&req); err != nil {
		vldMsg, ok := pkg.ParseValidationErrors(err)
		if !ok {
			c.JSON(400, shared.NewErrorResponse(400, "invalid request body! please provide valid body for sorting image request"))
			return
		}
		c.JSON(400, shared.NewErrorResponse(400, vldMsg))
		return
	}

	// nama image
	data,err_sort := pth.svc.SortImageOrder(c,productid,req);
	if err_sort != nil{
		c.JSON(err_sort.Code,err_sort);
		return
	}

	// response
	c.JSON(200, shared.NewSuccessResponse(200, "Image order successfully sorted", data));
}

func (pth *ProductHandler) DeleteProductImageHandler(c *gin.Context){
	// terima id
	idParams := c.Param("id");
	
	// validasi id
	imageId, err := uuid.Parse(idParams);
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid id params"))
		return
	}

	errDel := pth.svc.DeleteSpecificProductImage(c,imageId);
	if errDel != nil{
		c.JSON(errDel.Code, errDel);
		return
	}

	// response
	c.JSON(200, shared.NewSuccessResponse(200, "Product image successfully deleted", nil));

}

func (pth *ProductHandler) DeleteAllProductImagesHandler(c *gin.Context){
	// terima id
	productIdParams := c.Param("product_id");
	
	// validasi id
	productId, err := uuid.Parse(productIdParams);
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid id params"))
		return
	}

	errDel := pth.svc.DeleteAllProductImage(c,productId);
	if errDel != nil{
		c.JSON(errDel.Code, errDel);
		return
	}

	// response
	c.JSON(200, shared.NewSuccessResponse(200, "All product image successfully deleted", nil));
}

func (pth *ProductHandler) HandleGetProductImageById(c *gin.Context){
	IdParams := c.Param("id");
	// validasi id
	imageId, err := uuid.Parse(IdParams);
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid id params"))
		return
	}
	// tarik dataaa
	data,errGet := pth.svc.GetProductImage(c,imageId);
	
	if errGet != nil{
		c.JSON(errGet.Code,errGet);
		return
	}
	c.JSON(200, shared.NewSuccessResponse(200, "Product image successfully retrieved", data));
}

func (pth *ProductHandler) HandleGetAllProductImages(c *gin.Context){
	productIdParams := c.Param("product_id");
	// validasi id
	productId, err := uuid.Parse(productIdParams);
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid id params"))
		return
	}
	data,errGet := pth.svc.GetAllProductImage(c,productId);
	
	if errGet != nil{
		c.JSON(errGet.Code,errGet);
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "Product image successfully retrieved", data));

}

// Product Categories
func (pth *ProductHandler) PostAssignProductCategory(c *gin.Context){}

func (pth *ProductHandler) DeleteRemoveProductCategoryHandler(c *gin.Context){}

func (pth *ProductHandler) DeleteRemoveAllProductCategoriesHandler(c *gin.Context){}

func (pth *ProductHandler) HandleGetProductsByCategoryFilter(c *gin.Context){}

func (pth *ProductHandler) HandleGetProductCategories(c *gin.Context){}



func (pth *ProductHandler) RegisterRoutes(r gin.IRouter) {
	productApi := r.Group("/product");
	productApi.GET("/get-all", pth.HandleGetAllProducts);
	// kurang nyari with slug
	productApi.GET("/details/:slug", pth.HandleGetProductBySlug);

	productAdmin := productApi.Group("/");

	productAdmin.Use(middleware.AuthMiddleware());
	productAdmin.Use(middleware.RoleMiddleware(middleware.RoleAdmin));

	productAdmin.POST("/add", pth.PostCreateProduct);
	productAdmin.PATCH("/update/:id", pth.PatchUpdateProduct);
	productAdmin.DELETE("/delete/:id", pth.DeleteProductHandler);
	productAdmin.GET("/id/:id", pth.HandleGetProductById);


	// PRODUCT IMAGES

	productImages := productAdmin.Group("/images");
	productImagesPublic := productApi.Group("/images");

	// Create product_image yang di attach ke suatu produk
	productImages.POST("/add/:product_id", pth.PostAddProductImage);

	// Ganti susunan product_image dari suatu produk
	productImages.PATCH("/order/:product_id", pth.PatchUpdateProductImageOrder);

	// Hapus gambar dari suatu produk
	productImages.DELETE("/delete/:id", pth.DeleteProductImageHandler);

	// hapus semua gambar dari suatu produk
	productImages.DELETE("/delete-all/:product_id", pth.DeleteAllProductImagesHandler);

	// Ngambil suatu gambar spesifik
	productImagesPublic.GET("/id/:id", pth.HandleGetProductImageById);

	// Ngambil seluruh gambar dari suatu produk
	productImages.GET("/all/:product_id", pth.HandleGetAllProductImages);
	

	// PRODUCT Categories

	productCategories := productAdmin.Group("/categories");
	productCategoriesPublic := productApi.Group("/categories");

	// Attach suatu product ke dalam suatu category
	productCategories.POST("/add/:product_id", pth.PostAssignProductCategory);

	// hapus suatu product dari specific satu category
	productCategories.DELETE("/remove/:product_id/:category_id", pth.DeleteRemoveProductCategoryHandler);

	// reset product dari semua category
	productCategories.DELETE("/remove-all/:product_id", pth.DeleteRemoveAllProductCategoriesHandler);

	// get semua product dari category tertentu
	productCategoriesPublic.GET("/filter", pth.HandleGetProductsByCategoryFilter);

	// get semua category yang terikat pada suatu product tertentu
	productCategoriesPublic.GET("/all/:product_id", pth.HandleGetProductCategories);


}

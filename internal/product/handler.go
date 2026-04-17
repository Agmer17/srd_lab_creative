package product

import (
	"github.com/Agmer17/srd_lab_creative/internal/shared"
	"github.com/Agmer17/srd_lab_creative/internal/shared/middleware"
	"github.com/Agmer17/srd_lab_creative/pkg"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductHandler struct{
	svc *ProductService
}

func NewProductHandler (s *ProductService) *ProductHandler{
	return &ProductHandler{
		svc:s,
	}
}

func(pth *ProductHandler) HandleGetAllProducts(c *gin.Context){
	data,err := pth.svc.GetAllProduct(c);
	if err != nil{
		c.JSON(err.Code,err)
		return
	}
	c.JSON(200,shared.NewSuccessResponse(200,"Successfully getting the product data",data))
}


func(pth *ProductHandler) HandleGetProductById(c *gin.Context){
	path := c.Param("id");
	id,err := uuid.Parse(path);
	if err != nil{
		c.JSON(400, shared.NewErrorResponse(400, "invalid id params"));
		return
	}
	
	data,getErr := pth.svc.GetProductById(c,id);
	if getErr != nil {
		c.JSON(getErr.Code, getErr)
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "succesffuly getting the product data", data))
}

func (pth *ProductHandler) PostCreateProduct(c *gin.Context){
	var req createProductRequest;

	if err := c.ShouldBindJSON(&req); err != nil{
		vldMsg,ok := pkg.ParseValidationErrors(err)
		if !ok{
			c.JSON(400, shared.NewErrorResponse(400, "invalid request body! please provide valid body for create product request"))
			return
		}
		c.JSON(400, shared.NewErrorResponse(400, vldMsg));
		return
	}

	data, insErr := pth.svc.CreateProduct(c,req);
	if insErr != nil{
		c.JSON(insErr.Code,insErr);
		return
	}
	c.JSON(200, shared.NewSuccessResponse(200,"New product successfully created",data));
}

func(pth *ProductHandler) PatchUpdateProduct(c *gin.Context){
	var req updateProductRequest;

	path := c.Param("id");
	id,err := uuid.Parse(path);
	if err != nil{
		c.JSON(400, shared.NewErrorResponse(400, "invalid id params"));
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil{
		vldMsg,ok := pkg.ParseValidationErrors(err);
		if !ok{
			c.JSON(400, shared.NewErrorResponse(400,"Invalid request body! please provide valid body for update product request"));
			return
		}
		c.JSON(400, shared.NewErrorResponse(400,vldMsg));
		return
	}

	data, updErr := pth.svc.UpdateProduct(c,id,req);
	if updErr != nil{
		c.JSON(updErr.Code,updErr);
		return
	}
	c.JSON(200, shared.NewSuccessResponse(200,"product successfully updated",data));

}

func(pth *ProductHandler) DeleteProductHandler(c *gin.Context){
	path := c.Param("id");
	id,err := uuid.Parse(path);
	if err != nil{
		c.JSON(400, shared.NewErrorResponse(400, "invalid id params"));
		return
	}

	delErr := pth.svc.DeleteProduct(c,id);
	if delErr != nil{
		c.JSON(delErr.Code,delErr);
		return
	}
	c.JSON(200, shared.NewSuccessResponse(200,"product successfully deleted",nil));

}


func(pth *ProductHandler) RegisterRoutes(r gin.IRouter){
	productApi := r.Group("/product");
	productApi.GET("/get-all",pth.HandleGetAllProducts);
	productApi.GET("/id/:id",pth.HandleGetProductById);

	productAdmin := productApi.Group("/");
	
	productAdmin.Use(middleware.AuthMiddleware());
	productAdmin.Use(middleware.RoleMiddleware(middleware.RoleAdmin));
	
	productAdmin.POST("/add",pth.PostCreateProduct);
	productAdmin.PATCH("/update/:id",pth.PatchUpdateProduct);
	productAdmin.DELETE("/delete/:id",pth.DeleteProductHandler);

}
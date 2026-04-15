package category

import (
	"github.com/Agmer17/srd_lab_creative/internal/shared"
	"github.com/Agmer17/srd_lab_creative/internal/shared/middleware"
	"github.com/Agmer17/srd_lab_creative/pkg"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CategoryHandler struct {
	svc *CategoryService
}

func NewCategoryHandler(sv *CategoryService) *CategoryHandler {
	return &CategoryHandler{
		svc: sv,
	}
}

func (cth *CategoryHandler) HandleGetAllCategories(c *gin.Context) {
	data, err := cth.svc.GetAllCategories(c.Request.Context())
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "successfully getting the categories data", data))
}

func (cth *CategoryHandler) HandleSearchCategories(c *gin.Context) {

	query := c.Query("q")

	data, err := cth.svc.SearchCategory(c.Request.Context(), &query)

	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "succesffuly searching the category data", data))

}

func (cth *CategoryHandler) HandleGetCategorybyId(c *gin.Context) {

	path := c.Param("id")
	id, err := uuid.Parse(path)
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid id params"))
		return
	}

	data, getErr := cth.svc.GetCategoryById(c.Request.Context(), id)
	if getErr != nil {
		c.JSON(getErr.Code, getErr)
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "succesffuly getting the category data", data))

}

func (cth *CategoryHandler) HandleGetCategorybySlug(c *gin.Context) {

	slug := c.Param("slug")

	data, err := cth.svc.GetCategoryBySlug(c.Request.Context(), slug)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "succesffuly getting the category data", data))

}

func (cth *CategoryHandler) PostCreateCategory(c *gin.Context) {

	var req createCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		vldMsg, ok := pkg.ParseValidationErrors(err)
		if !ok {
			c.JSON(400, shared.NewErrorResponse(400, "invalid request body! please provide valid valid body for create category request"))
			return
		}

		c.JSON(400, shared.NewErrorResponse(400, vldMsg))
		return
	}

	data, insErr := cth.svc.CreateCategory(c.Request.Context(), req)
	if insErr != nil {
		c.JSON(insErr.Code, insErr)
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "successfully creating new category", data))
}

func (cth *CategoryHandler) PatchUpdateCategory(c *gin.Context) {

	path := c.Param("id")
	id, err := uuid.Parse(path)
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid id params"))
		return
	}

	var req updateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {

		vldMsg, ok := pkg.ParseValidationErrors(err)
		if !ok {
			c.JSON(400, shared.NewErrorResponse(400, "invalid request body! please provide valid valid body for create category request"))
			return
		}

		c.JSON(400, shared.NewErrorResponse(400, vldMsg))
		return
	}

	data, uptErr := cth.svc.UpdateCategory(c.Request.Context(), id, req)
	if uptErr != nil {
		c.JSON(uptErr.Code, uptErr)
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "successfully updating category data!", data))
}

func (cth *CategoryHandler) DeleteCategoryHandler(c *gin.Context) {

	path := c.Param("id")
	id, err := uuid.Parse(path)
	if err != nil {
		c.JSON(400, shared.NewErrorResponse(400, "invalid id params"))
		return
	}

	delErr := cth.svc.DeleteCategory(c.Request.Context(), id)
	if delErr != nil {
		c.JSON(delErr.Code, delErr)
		return
	}

	c.JSON(200, shared.NewSuccessResponse(200, "successfully deleting the category data", nil))
}

func (cth *CategoryHandler) RegisterRoutes(r gin.IRouter) {

	categoryApi := r.Group("/category")

	categoryApi.GET("/get-all", cth.HandleGetAllCategories)
	categoryApi.GET("/search", cth.HandleSearchCategories)
	categoryApi.GET("/id/:id", cth.HandleGetCategorybyId)
	categoryApi.GET("slug/:slug", cth.HandleGetCategorybySlug)

	catAdminOnly := categoryApi.Group("/")
	catAdminOnly.Use(middleware.AuthMiddleware())
	catAdminOnly.Use(middleware.RoleMiddleware(middleware.RoleAdmin))

	catAdminOnly.POST("/add", cth.PostCreateCategory)
	catAdminOnly.PATCH("/update/:id", cth.PatchUpdateCategory)
	catAdminOnly.DELETE("/delete/:id", cth.DeleteCategoryHandler)

}

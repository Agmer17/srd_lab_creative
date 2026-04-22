package product

type createProductRequest struct{
	Name string `json:"name" binding:"required,min=3,max=255"`
	Slug string `json:"slug" binding:"required,min=3"`
	Description string `json:"description" binding:"omitempty,max=2000"`
	Price float64 `json:"price" binding:"required,min=0"`
	Status string `json:"status" binding:"required,oneof=active inactive"`
	IsFeatured *bool `json:"is_featured" binding:"required"`
}
type updateProductRequest struct {
	Name        *string  `json:"name" binding:"omitempty,min=3"`
	Slug        *string  `json:"slug" binding:"omitempty,min=3"`
	Description *string  `json:"description" binding:"omitempty,max=2000"`
	Price       *float64 `json:"price" binding:"omitempty,min=0"`
	Status      *string  `json:"status" binding:"omitempty,oneof=active inactive"`
	IsFeatured  *bool    `json:"is_featured" binding:"omitempty"`
}


type UpdateProductImageSort struct{
	ImageId string `json:"image_id" binding:"required,uuid"`
	SortOrder int `json:"sort_order" binding:"required,min=0"`
}

type updateProductImageRequest struct{
	ProductId *string `json:"product_id" binding:"omitempty,uuid"`
	ImageUrl *string `json:"image_url" binding:"omitempty,min=3"`
	IsPrimary *bool `json:"is_primary" binding:"omitempty"`
	SortOrder *int `json:"sort_order" binding:"omitempty,min=0"`
}

type assignProductCategoryRequest struct {
	CategoryIds []string `json:"category_ids" binding:"required,min=1,dive,uuid"`
}
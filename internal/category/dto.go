package category

type createCategoryRequest struct {
	Name string  `json:"name" binding:"required,min=3,max=100"`
	Slug string  `json:"slug" binding:"required,min=3,max=100"`
	Desc *string `json:"desc" binding:"omitempty,max=255"`
}

type updateCategoryRequest struct {
	Name *string `json:"name" binding:"omitempty,min=3,max=100"`
	Slug *string `json:"slug" binding:"omitempty,min=3,max=100"`
	Desc *string `json:"desc" binding:"omitempty,max=255"`
}

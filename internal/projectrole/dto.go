package projectrole

type createRoleRequest struct {
	Name string `json:"name" binding:"required,min=3,max=255,alphanumspace"`
}
type updateRoleRequest struct {
	Name string `json:"name" binding:"required,min=3,max=255,alphanumspace"`
	ID   string `json:"id" binding:"required,uuid"`
}

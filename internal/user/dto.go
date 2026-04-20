package user

type UpdateUserDto struct {
	FullName    *string `json:"full_name" binding:"omitempty,min=3,max=100"`
	PhoneNumber *string `json:"phone_number" binding:"omitempty,e164"`
	Gender      *string `json:"gender" binding:"omitempty,oneof=male female"`
}

type updateUserRoleDTO struct {
	Role string `json:"role" binding:"required,oneof=ADMIN USER"`
}

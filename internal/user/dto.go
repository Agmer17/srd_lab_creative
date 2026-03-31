package user

type UpdateUserDto struct {
	FullName    *string `json:"full_name"`
	PhoneNumber *string `json:"phone_number"`
	Gender      *string `json:"gender"`
}

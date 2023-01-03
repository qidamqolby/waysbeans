package userdto

type UpdateUserRequest struct {
	Name     string `json:"name" form:"name"`
	Password string `json:"password" form:"password"`
	Image    string `json:"image" form:"image"`
	Phone    string `json:"phone" form:"phone"`
	Address  string `json:"address" form:"address"`
}

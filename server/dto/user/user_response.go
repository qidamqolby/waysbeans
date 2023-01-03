package userdto

type UserResponse struct {
	ID      int32  `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Image   string `json:"image"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

type UpdateUserResponse struct {
	Name    string `json:"name"`
	Image   string `json:"image"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

type DeleteUserResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

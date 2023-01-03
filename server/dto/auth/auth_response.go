package authdto

type RegisterResponse struct {
	ID    int32  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"Email"`
}

type LoginResponse struct {
	ID    int32  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
	Token string `json:"token"`
}

type CheckAuthResponse struct {
	ID    int32  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

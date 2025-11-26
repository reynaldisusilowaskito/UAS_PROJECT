package model

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type RegisterRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	RoleID   string `json:"role_id"`
}

package requests

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ChangePasswordRequest struct {
	ID              string `json:"id"`
	Password        string `json:"password1"`
	ConfirmPassword string `json:"password"`
}

type CreateAdmin struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

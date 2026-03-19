package auth

type RegisterDto struct {
	Email    string `json:"email" format:"email" example:"user@example.com"`
	Name     string `json:"name" validate:"required,min=3,max=100,notblank"`
	Password string `json:"password" validate:"required,min=6,notblank"`
}

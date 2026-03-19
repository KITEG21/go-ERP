package auth

type LoginDto struct {
	Email    string `json:"email" format:"email" example:"user@example.com"`
	Password string `json:"password" validate:"required,min=6,notblank"`
}

package auth

type LoginDto struct {
	Email    string `json:"email" validate:"required,email,notblank"`
	Password string `json:"password" validate:"required,min=6,notblank"`
}

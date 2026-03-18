package auth

type RegisterDto struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=3,max=100,notblank"`
	Password string `json:"password" validate:"required,min=6,notblank"`
}

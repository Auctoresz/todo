package dto

type UserSignupRequest struct {
	Email    string `validate:"required,min=1,email" json:"email"`
	Password string `validate:"required,min=10" json:"password"`
	Name     string `validate:"required,min=1" json:"name"`
}

type UserLoginRequest struct {
	Email    string `validate:"required,min=1,email" json:"email"`
	Password string `validate:"required,min=10" json:"password"`
}

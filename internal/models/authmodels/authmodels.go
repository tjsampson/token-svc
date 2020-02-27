package authmodels

import "net/http"

// UserCreds represents the user credentials used for Auth/Login
type UserCreds struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UserRegistration is the data required to register a user
type UserRegistration struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"eqfield=ConfirmPassword,required,max=50,min=12"`
	ConfirmPassword string `json:"confirm_password" validate:"eqfield=Password,required,max=50,min=12"`
}

// LoginResponse represents the login response object
type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	HTTPCookie   *http.Cookie `json:"cookie"`
}

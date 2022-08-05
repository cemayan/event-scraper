package model

// LoginInput is representation of the login payload
type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

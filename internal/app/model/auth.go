package model

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

package dto

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RegisterRequest struct {
	Email    string `json:"email" example:"test@example.com"`
	Password string `json:"password" example:"secret123"`
}

type RegisterResponse struct {
	ID    uint   `json:"id" example:"1"`
	Email string `json:"email" example:"test@example.com"`
}

type LoginRequest struct {
	Email    string `json:"email" example:"test@example.com"`
	Password string `json:"password" example:"secret123"`
}

type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type CreateAPIKeyResponse struct {
	APIKey string `json:"api_key" example:"mws_sk_1234567890abcdef"`
}

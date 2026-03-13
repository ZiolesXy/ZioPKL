package respons

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Role string `json:"role"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	Role string `json:"role"`
}
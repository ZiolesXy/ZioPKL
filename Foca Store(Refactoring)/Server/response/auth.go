package response

type AuthResponse struct {
	User UserResponse `json:"user"`
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenRefreshResponse struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

func BuildAuthResponse(user UserResponse, accessToken, refreshToken string) AuthResponse {
	return AuthResponse{
		User: user,
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	}
}

func BuildTokenRefreshResponse(accessToken, refreshToken string) TokenRefreshResponse {
	return TokenRefreshResponse{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	}
}

func BuildToken(accessToken string) TokenResponse {
	return TokenResponse{
		AccessToken: accessToken,
	}
}
package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"
	"voca-store/internal/domain/models"
	"voca-store/internal/domain/repository"
	"voca-store/internal/domain/dto/request"
	"voca-store/internal/domain/dto/response"
	"voca-store/internal/helper"
)

type AuthService interface {
	Register(req request.RegisterRequest) (response.UserResponse, error)
	Login(req request.LoginRequest) (response.AuthResponse, error)
	RefreshToken(req request.RefreshTokenRequest) (response.TokenResponse, error)
	Logout(userID uint, token string) error
	ChangePassword(userID uint, req request.ChangePasswordRequest) error
	ForgotPassword(req request.ForgotPasswordRequest) error
	VerifyOTP(req request.VerifyOTPRequest) error
	GetProfile(userID uint) (response.UserProfileResponse, error)
	UpdateProfile(userID uint, req request.UpdateProfileRequest) (response.UserResponse, error)
}

type authService struct {
	userRepo repository.UserRepository
	authRepo repository.AuthRepository
}

func NewAuthService(userRepo repository.UserRepository, authRepo repository.AuthRepository) AuthService {
	return &authService{userRepo, authRepo}
}

func (s *authService) Register(req request.RegisterRequest) (response.UserResponse, error) {
	_, err := s.userRepo.FindByEmail(req.Email)
	if err == nil {
		return response.UserResponse{}, errors.New("email already registered")
	}

	hashedPassword, _ := helper.HashPassword(req.Password)
	userRole, err := s.userRepo.FindRoleByName("User")
	if err != nil {
		return response.UserResponse{}, errors.New("user role not found")
	}

	user := models.User{
		Name:            req.Name,
		Email:           req.Email,
		Password:        hashedPassword,
		TelephoneNumber: req.TelephoneNumber,
		RoleID:          userRole.ID,
	}

	if err := s.userRepo.Create(&user); err != nil {
		return response.UserResponse{}, err
	}

	// Create cart
	s.userRepo.CreateCart(&models.Cart{UserID: user.ID})

	return response.BuildUserResponse(user.ID, user.Name, user.Email, "User", user.TelephoneNumber), nil
}

func (s *authService) Login(req request.LoginRequest) (response.AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return response.AuthResponse{}, errors.New("invalid email or password")
	}

	if err := helper.VerifyPassword(user.Password, req.Password); err != nil {
		return response.AuthResponse{}, errors.New("invalid email or password")
	}

	accessToken, _ := helper.GenerateAccessToken(user.ID, user.Role.Name)
	refreshToken, _ := helper.GenerateRefreshToken(user.ID, user.Role.Name)

	s.authRepo.SetRefreshToken(user.ID, refreshToken, 7*24*time.Hour)

	userResp := response.BuildUserResponse(user.ID, user.Name, user.Email, user.Role.Name, user.TelephoneNumber)
	return response.BuildAuthResponse(userResp, accessToken, refreshToken), nil
}

func (s *authService) RefreshToken(req request.RefreshTokenRequest) (response.TokenResponse, error) {
	claims, err := helper.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return response.TokenResponse{}, errors.New("invalid refresh token")
	}

	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return response.TokenResponse{}, errors.New("user not found")
	}

	storedToken, err := s.authRepo.GetRefreshToken(user.ID)
	if err != nil || storedToken != req.RefreshToken {
		return response.TokenResponse{}, errors.New("invalid refresh token")
	}

	accessToken, _ := helper.GenerateAccessToken(user.ID, user.Role.Name)
	return response.BuildToken(accessToken), nil
}

func (s *authService) Logout(userID uint, token string) error {
	claims, err := helper.ValidateAccessToken(token)
	if err == nil {
		ttl := time.Until(claims.ExpiresAt.Time)
		s.authRepo.BlacklistToken(token, ttl)
	}

	s.authRepo.DeleteRefreshToken(userID)
	return nil
}

func (s *authService) ChangePassword(userID uint, req request.ChangePasswordRequest) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if err := helper.VerifyPassword(user.Password, req.OldPassword); err != nil {
		return errors.New("old password incorrect")
	}

	if req.NewPassword != req.ConfirmPassword {
		return errors.New("password confirmation does not match")
	}

	hashed, _ := helper.HashPassword(req.NewPassword)
	return s.userRepo.Update(&user, map[string]interface{}{"password": hashed})
}

func (s *authService) ForgotPassword(req request.ForgotPasswordRequest) error {
	limit, _ := s.authRepo.GetOTPLimit(req.Email)
	if limit >= 5 {
		return errors.New("too many otp attempts, try again later")
	}

	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		s.authRepo.IncrementOTPLimit(req.Email, 5*time.Minute)
		return nil // Return nil to avoid email harvesting
	}

	n, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	otp := fmt.Sprintf("%06d", n.Int64())

	s.authRepo.SetOTP(req.Email, otp, 5*time.Minute)
	s.authRepo.IncrementOTPLimit(req.Email, 5*time.Minute)

	go func() {
		helper.SendOTPEmail(user.Email, otp)
	}()

	return nil
}

func (s *authService) VerifyOTP(req request.VerifyOTPRequest) error {
	storedOTP, err := s.authRepo.GetOTP(req.Email)
	if err != nil || storedOTP != req.OTP {
		return errors.New("invalid or expired otp")
	}

	hashed, _ := helper.HashPassword(req.NewPassword)
	user, _ := s.userRepo.FindByEmail(req.Email)
	
	err = s.userRepo.Update(&user, map[string]interface{}{"password": hashed})
	if err == nil {
		s.authRepo.DeleteOTP(req.Email)
	}
	return err
}

func (s *authService) GetProfile(userID uint) (response.UserProfileResponse, error) {
    user, err := s.userRepo.FindByID(userID)
    if err != nil {
        return response.UserProfileResponse{}, err
    }

    // Convert addresses ke response
    var addresses []response.AddressResponse
    for _, addr := range user.Addresses {
        addresses = append(addresses, response.BuildAddressResponse(addr))
    }

    return response.BuildUserProfileResponse(
        user.ID,
        user.Name,
        user.Email,
        user.TelephoneNumber,
        user.Role.Name,
        user.ProfileImageURL,
        addresses,
        user.CreatedAt,
        user.UpdatedAt,
    ), nil
}

func (s *authService) UpdateProfile(userID uint, req request.UpdateProfileRequest) (response.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return response.UserResponse{}, err
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Email != nil {
		// Check if email already exists
		existing, _ := s.userRepo.FindByEmail(*req.Email)
		if existing.ID != 0 && existing.ID != userID {
			return response.UserResponse{}, errors.New("email already taken")
		}
		updates["email"] = *req.Email
	}
	if req.TelephoneNumber != nil {
		updates["telephone_number"] = *req.TelephoneNumber
	}

	if err := s.userRepo.Update(&user, updates); err != nil {
		return response.UserResponse{}, err
	}

	// Reload
	user, _ = s.userRepo.FindByID(userID)
	return response.BuildUserResponse(user.ID, user.Name, user.Email, user.Role.Name, user.TelephoneNumber), nil
}

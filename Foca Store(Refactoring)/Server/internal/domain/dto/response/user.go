package response

import "time"

type UserMiniResponse struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	TelephoneNumber string `json:"telephone_number"`
}

type UserResponse struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	TelephoneNumber string `json:"telephone_number"`
	ProfileImageURL string `json:"profile_image_url,omitempty"`
	Role            string `json:"role"`
}

type AllUsers struct {
	Entries []UserResponse `json:"entries"`
}

type UserProfileResponse struct {
	ID              uint              `json:"id"`
	Name            string            `json:"name"`
	Email           string            `json:"email"`
	TelephoneNumber string            `json:"telephone_number"`
	ProfileImageURL string            `json:"profile_image_url,omitempty"`
	Role            string            `json:"role"`
	Address         []AddressResponse `json:"address"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

func BuildUserResponse(id uint, name, email, role, telephoneNumber, profileImageURL string) UserResponse {
	return UserResponse{
		ID:              id,
		Name:            name,
		Email:           email,
		TelephoneNumber: telephoneNumber,
		ProfileImageURL: profileImageURL,
		Role:            role,
	}
}

func BuildAllUser(users []UserResponse) AllUsers {
	return AllUsers{
		Entries: users,
	}
}

func BuildUserProfileResponse(id uint, name, email, telephoneNumber, profileImageURL, role string, address []AddressResponse, createdat, updatedat time.Time) UserProfileResponse {
	return UserProfileResponse{
		ID:              id,
		Name:            name,
		Email:           email,
		TelephoneNumber: telephoneNumber,
		ProfileImageURL: profileImageURL,
		Role:            role,
		Address:         address,
		CreatedAt:       createdat,
		UpdatedAt:       updatedat,
	}
}

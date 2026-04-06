package request

type UpdateProfileRequest struct {
	Name            *string `json:"name"`
	Email           *string `json:"email"`
	TelephoneNumber *string `json:"telephone_number"`
}

type UpdateProfileWithImageRequest struct {
	Name            *string `json:"name"`
	Email           *string `json:"email"`
	TelephoneNumber *string `json:"telephone_number"`
}

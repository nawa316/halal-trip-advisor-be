package domain

import "context"

type Profile struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateProfileRequest struct {
	Name string `form:"name" binding:"required"`
}

type UpdateProfileResponse struct {
	Message string `json:"message"`
}

type ProfileUsecase interface {
	GetProfileByID(c context.Context, userID string) (*Profile, error)
	UpdateProfile(c context.Context, userID string, user *User) error
}

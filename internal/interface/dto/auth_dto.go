package dto

import "challenge-travel-api/internal/domain/enums"

type RegisterRequestDTO struct {
	Name     string         `json:"name" binding:"required"`
	Email    string         `json:"email" binding:"required,email"`
	Password string         `json:"password" binding:"required,min=6"`
	Role     enums.UserType `json:"role" binding:"required"`
}

type LoginRequestDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginResponseDTO struct {
	AccessToken string `json:"access_token"`
}

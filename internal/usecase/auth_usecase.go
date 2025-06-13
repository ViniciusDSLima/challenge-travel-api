package usecase

import (
	"challenge-travel-api/internal/domain/entity"
	"challenge-travel-api/internal/domain/gateway"
	"challenge-travel-api/internal/interface/dto"
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserAlreadyExists = errors.New("usuário já existe no sistema.")
)

type AuthUseCase interface {
	Register(ctx context.Context, input dto.RegisterRequestDTO) error
	Login(ctx context.Context, input dto.LoginRequestDTO) (dto.LoginResponseDTO, error)
}

type AuthUseCaseImpl struct {
	repo gateway.UserGateway
}

func NewAUthUseCase(repo gateway.UserGateway) *AuthUseCaseImpl {
	return &AuthUseCaseImpl{
		repo: repo,
	}
}

func (uc *AuthUseCaseImpl) Register(ctx context.Context, input dto.RegisterRequestDTO) error {
	existingUser, err := uc.repo.FindByEmail(ctx, input.Email)

	if err == nil && existingUser != nil {
		return ErrUserAlreadyExists
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	hashedPassword, err := uc.hashPassword(input.Password)

	if err != nil {
		return err
	}

	newUser := &entity.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
		Role:     input.Role,
	}

	err = uc.repo.Create(ctx, newUser)

	if err != nil {
		return err
	}

	return nil
}

func (uc *AuthUseCaseImpl) Login(ctx context.Context, input dto.LoginRequestDTO) (dto.LoginResponseDTO, error) {
	user, err := uc.repo.FindByEmail(ctx, input.Email)

	if err != nil {
		return dto.LoginResponseDTO{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return dto.LoginResponseDTO{}, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.Id,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")

	if jwtSecretKey == "" {
		return dto.LoginResponseDTO{}, errors.New("JWT_SECRET_KEY environment not found")
	}

	tokenString, err := token.SignedString([]byte(jwtSecretKey))

	if err != nil {
		return dto.LoginResponseDTO{}, err
	}

	return dto.LoginResponseDTO{
		AccessToken: tokenString,
	}, nil
}

func (uc *AuthUseCaseImpl) hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		return "", nil
	}

	return string(hashedPassword), nil
}

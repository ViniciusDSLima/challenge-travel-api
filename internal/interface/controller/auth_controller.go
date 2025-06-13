package controller

import (
	"challenge-travel-api/internal/interface/dto"
	"challenge-travel-api/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authUseCase usecase.AuthUseCase
}

func NewAuthController(authUseCase usecase.AuthUseCase) *AuthController {
	return &AuthController{
		authUseCase: authUseCase,
	}
}

// Register godoc
// @Summary Registrar um novo usuário
// @Description Registra um novo usuário no sistema
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequestDTO true "Dados do usuário"
// @Success 201 {object} nil
// @Failure 400 {object} map[string]string
// @Router /auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var request dto.RegisterRequestDTO
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao decodificar requisição"})
		return
	}

	err := c.authUseCase.Register(ctx.Request.Context(), request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusCreated)
}

// Login godoc
// @Summary Autenticar usuário
// @Description Autentica um usuário e retorna um token JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequestDTO true "Credenciais do usuário"
// @Success 200 {object} dto.LoginResponseDTO
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var request dto.LoginRequestDTO
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao decodificar requisição"})
		return
	}

	response, err := c.authUseCase.Login(ctx.Request.Context(), request)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

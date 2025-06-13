package controller

import (
	"challenge-travel-api/internal/domain/enums"
	"challenge-travel-api/internal/interface/dto"
	"challenge-travel-api/internal/usecase"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TravelController struct {
	travelUseCase usecase.TravelUseCase
}

func NewTravelController(travelUseCase usecase.TravelUseCase) *TravelController {
	return &TravelController{
		travelUseCase: travelUseCase,
	}
}

// CreateTravelRequest godoc
// @Summary Criar uma nova solicitação de viagem
// @Description Cria uma nova solicitação de viagem para o usuário autenticado
// @Tags travels
// @Accept json
// @Produce json
// @Param request body dto.CreateTravelRequestDTO true "Dados da solicitação de viagem"
// @Success 201 {object} entity.TravelRequest
// @Failure 400 {object} map[string]string
// @Security Bearer
// @Router /travels [post]
func (c *TravelController) CreateTravelRequest(ctx *gin.Context) {
	var request dto.CreateTravelRequestDTO
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao decodificar requisição"})
		return
	}

	userID := ctx.MustGet("user_id").(uuid.UUID)

	travel, err := c.travelUseCase.CreateTravelRequest(
		ctx.Request.Context(),
		userID,
		request,
	)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, travel)
}

// UpdateTravelRequest godoc
// @Summary Atualizar uma solicitação de viagem
// @Description Atualiza os dados de uma solicitação de viagem existente
// @Tags travels
// @Accept json
// @Produce json
// @Param id path string true "ID da solicitação de viagem"
// @Param request body dto.UpdateTravelRequestDTO true "Dados atualizados da solicitação"
// @Success 200 {object} entity.TravelRequest
// @Failure 400 {object} map[string]string
// @Security Bearer
// @Router /travels/{id} [put]
func (c *TravelController) UpdateTravelRequest(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var request dto.UpdateTravelRequestDTO
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao decodificar requisição"})
		return
	}

	userID := ctx.MustGet("user_id").(uuid.UUID)

	travel, err := c.travelUseCase.UpdateTravelRequest(
		ctx.Request.Context(),
		id,
		userID,
		request,
	)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, travel)
}

// UpdateStatusTravelRequest godoc
// @Summary Atualizar status da solicitação de viagem
// @Description Atualiza o status de uma solicitação de viagem (aprovado/cancelado)
// @Tags travels
// @Accept json
// @Produce json
// @Param id path string true "ID da solicitação de viagem"
// @Param request body dto.UpdateStatusTravelRequestDTO true "Novo status da solicitação"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Security Bearer
// @Router /travels/{id}/status [patch]
func (c *TravelController) UpdateStatusTravelRequest(ctx *gin.Context) {
	travelID := ctx.Param("id")
	userID := ctx.MustGet("user_id").(uuid.UUID)

	var request dto.UpdateStatusTravelRequestDTO
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao decodificar requisição"})
		return
	}

	request.TravelRequestId = travelID
	err := c.travelUseCase.UpdateStatusTravelRequest(
		ctx.Request.Context(),
		userID.String(),
		request,
	)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// GetTravelRequest godoc
// @Summary Obter detalhes de uma solicitação de viagem
// @Description Retorna os detalhes de uma solicitação de viagem específica
// @Tags travels
// @Accept json
// @Produce json
// @Param id path string true "ID da solicitação de viagem"
// @Success 200 {object} entity.TravelRequest
// @Failure 404 {object} map[string]string
// @Security Bearer
// @Router /travels/{id} [get]
func (c *TravelController) GetTravelRequest(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	userID := ctx.MustGet("user_id").(uuid.UUID)

	travel, err := c.travelUseCase.GetByID(ctx.Request.Context(), id, userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, travel)
}

// ListTravelRequests godoc
// @Summary Listar solicitações de viagem
// @Description Retorna uma lista de solicitações de viagem com filtros opcionais
// @Tags travels
// @Accept json
// @Produce json
// @Param status query string false "Filtrar por status (PENDING, APPROVED, CANCELED)"
// @Param start_date query string false "Data inicial (YYYY-MM-DD)"
// @Param end_date query string false "Data final (YYYY-MM-DD)"
// @Param destination query string false "Nome do destino"
// @Param page query int false "Número da página" default(1)
// @Param page_size query int false "Tamanho da página" default(10)
// @Success 200 {array} entity.TravelRequest
// @Security Bearer
// @Router /travels [get]
func (c *TravelController) ListTravelRequests(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(uuid.UUID)

	statusStr := ctx.Query("status")
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")
	destinationName := ctx.Query("destination_name")
	pageStr := ctx.Query("page")
	pageSizeStr := ctx.Query("page_size")

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(pageSizeStr)
	if pageSize < 1 {
		pageSize = 10
	}

	var status *enums.TravelRequestStatus
	if statusStr != "" {
		s := enums.TravelRequestStatus(statusStr)
		status = &s
	}

	var startDate, endDate *time.Time
	if startDateStr != "" {
		t, err := time.Parse(time.RFC3339, startDateStr)
		if err == nil {
			startDate = &t
		}
	}

	if endDateStr != "" {
		t, err := time.Parse(time.RFC3339, endDateStr)
		if err == nil {
			endDate = &t
		}
	}

	var destName *string
	if destinationName != "" {
		destName = &destinationName
	}

	travels, err := c.travelUseCase.ListTravelRequests(
		ctx.Request.Context(),
		userID,
		status,
		startDate,
		endDate,
		destName,
		page,
		pageSize,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, travels)
}

package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/ManuelMassora/servicoJa-api/internal/usecases"
	"github.com/gin-gonic/gin"
)

type AvaliacaoUseCase interface {
	Criar(ctx context.Context, req usecases.AvaliacaoRequest, idAvaliador uint) error
	ListarPorCliente(ctx context.Context, idCliente uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]usecases.AvaliacaoResponse, error)
	ListarPorPrestador(ctx context.Context, idPrestador uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]usecases.AvaliacaoResponse, error)
	MediaPorPrestador(ctx context.Context, idPrestador uint) (float64, error)
}

type AvaliacaoHandler struct {
	uc AvaliacaoUseCase
}

func NewAvaliacaoHandler(uc AvaliacaoUseCase) *AvaliacaoHandler {
	return &AvaliacaoHandler{uc: uc}
}

func (h *AvaliacaoHandler) CriarAvaliacao(c *gin.Context) {
	var req usecases.AvaliacaoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	idAvaliador, err := getUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	err = h.uc.Criar(c.Request.Context(), req, idAvaliador)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Avaliação criada com sucesso"})
}

func (h *AvaliacaoHandler) ListarAvaliacoesPorCliente(c *gin.Context) {
	idUsuario, err := getUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	filters := ExtractFilters(c)
	limit, offset, _, _ := ExtractPagination(c)
	orderBy := c.Query("orderBy")
	orderDir := c.Query("orderDir")

	avaliacoes, err := h.uc.ListarPorCliente(c.Request.Context(), idUsuario, filters, orderBy, orderDir, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": avaliacoes})
}

func (h *AvaliacaoHandler) ListarAvaliacoesPorPrestador(c *gin.Context) {
	idPrestadorStr := c.Param("id")
	idPrestador, err := strconv.ParseUint(idPrestadorStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do prestador inválido"})
		return
	}

	filters := ExtractFilters(c)
	limit, offset, _, _ := ExtractPagination(c)
	orderBy := c.Query("orderBy")
	orderDir := c.Query("orderDir")

	avaliacoes, err := h.uc.ListarPorPrestador(c.Request.Context(), uint(idPrestador), filters, orderBy, orderDir, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": avaliacoes})
}

func (h *AvaliacaoHandler) MediaPorPrestador(c *gin.Context) {
	idPrestadorStr := c.Param("id")
	idPrestador, err := strconv.ParseUint(idPrestadorStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do prestador inválido"})
		return
	}

	media, err := h.uc.MediaPorPrestador(c.Request.Context(), uint(idPrestador))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"media": media})
}

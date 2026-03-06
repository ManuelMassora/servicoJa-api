package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/ManuelMassora/servicoJa-api/internal/usecases"
	"github.com/gin-gonic/gin"
)

type CategoriaPrestadorUsecase interface {
	Criar(ctx context.Context, req usecases.CategoriaPrestadorRequest) (*usecases.CategoriaPrestadorResponse, error)
	Editar(ctx context.Context, id uint, campos map[string]interface{}) (*usecases.CategoriaPrestadorResponse, error)
	Listar(ctx context.Context, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]usecases.CategoriaPrestadorResponse, error)
}

type CategoriaPrestadorHandler struct {
	uc CategoriaPrestadorUsecase
}

func NewCategoriaPrestadorHandler(uc CategoriaPrestadorUsecase) *CategoriaPrestadorHandler {
	return &CategoriaPrestadorHandler{uc: uc}
}

func (h *CategoriaPrestadorHandler) Criar(c *gin.Context) {
	var req usecases.CategoriaPrestadorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.uc.Criar(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *CategoriaPrestadorHandler) Editar(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var campos map[string]interface{}
	if err := c.ShouldBindJSON(&campos); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.uc.Editar(c.Request.Context(), uint(id), campos)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *CategoriaPrestadorHandler) Listar(c *gin.Context) {
	filters := ExtractFilters(c)
	orderBy := c.Query("orderBy")
	orderDir := c.Query("orderDir")
	limit, offset, page, pageSize := ExtractPagination(c)

	propostas, err := h.uc.Listar(
		c.Request.Context(),
		filters,
		orderBy,
		orderDir,
		limit,
		offset,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      propostas,
		"page":      page,
		"pageSize":  pageSize,
		"orderBy":   orderBy,
		"direction": orderDir,
		"filters":   filters,
	})
}

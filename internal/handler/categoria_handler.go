package handler

import (
	"github.com/ManuelMassora/servicoJa-api/internal/usecases"
	"github.com/gin-gonic/gin"
	"strconv"
)

type CategoriaHandler struct {
	uc usecases.CategoriaUseCase
}

func NewCategoriaHandler(uc usecases.CategoriaUseCase) *CategoriaHandler {
	return &CategoriaHandler{uc: uc}
}

func (h *CategoriaHandler) Criar(c *gin.Context) {
	var request usecases.CategoriaRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.uc.Criar(c.Request.Context(), request); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Status(201)
}

func (h *CategoriaHandler) Listar(c *gin.Context) {
	filters := make(map[string]interface{})
	orderBy := c.Query("orderBy")
	orderDir := c.Query("orderDir")
	
	// Parse pagination parameters
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	
	// Calculate offset
	offset := (page - 1) * pageSize

	categorias, err := h.uc.Listar(c.Request.Context(), filters, orderBy, orderDir, pageSize, offset)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(200, gin.H{
		"data": categorias,
		"page": page,
		"pageSize": pageSize,
	})
}

func (h *CategoriaHandler) BuscarPorID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	categoria, err := h.uc.BuscarPorID(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if categoria == nil {
		c.JSON(404, gin.H{"error": "Categoria not found"})
		return
	}
	c.JSON(200, categoria)
}
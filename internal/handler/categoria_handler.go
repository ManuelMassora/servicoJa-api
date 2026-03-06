package handler

import (
	"context"
	"strconv"

	"github.com/ManuelMassora/servicoJa-api/internal/usecases"
	"github.com/gin-gonic/gin"
)

type CategoriaUseCase interface {
	Criar(ctx context.Context, request usecases.CategoriaRequest) (uint, error)
	Editar(ctx context.Context, id uint, campos map[string]interface{}) error
	Listar(ctx context.Context, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]usecases.CategoriaResponse, error)
	BuscarPorID(ctx context.Context, id uint) (*usecases.CategoriaResponse, error)
}

type CategoriaHandler struct {
	uc CategoriaUseCase
}

func NewCategoriaHandler(uc CategoriaUseCase) *CategoriaHandler {
	return &CategoriaHandler{uc: uc}
}

func (h *CategoriaHandler) Criar(c *gin.Context) {
	var request usecases.CategoriaRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	idCategoria, err := h.uc.Criar(c.Request.Context(), request)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "Categoria criada com sucesso", "id": idCategoria})
}

func (h *CategoriaHandler) Editar(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}

	var campos map[string]interface{}
	if err := c.ShouldBindJSON(&campos); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if len(campos) == 0 {
		c.JSON(400, gin.H{"error": "no fields to update"})
		return
	}

	delete(campos, "id")

	if err := h.uc.Editar(c.Request.Context(), uint(id), campos); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Status(204)
}

func (h *CategoriaHandler) Listar(c *gin.Context) {
	filters := make(map[string]interface{})

	// Collect filters from query params except pagination and ordering controls
	for key, vals := range c.Request.URL.Query() {
		if key == "orderBy" || key == "orderDir" || key == "page" || key == "pageSize" {
			continue
		}
		if len(vals) == 0 || vals[0] == "" {
			continue
		}
		v := vals[0]

		// Try to convert common types
		if key == "id" {
			if id, err := strconv.ParseUint(v, 10, 64); err == nil {
				filters[key] = id
				continue
			}
		}
		if key == "ativo" {
			if b, err := strconv.ParseBool(v); err == nil {
				filters[key] = b
				continue
			}
		}

		// default: keep as string
		filters[key] = v
	}

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
		"data":      categorias,
		"page":      page,
		"pageSize":  pageSize,
		"orderBy":   orderBy,
		"direction": orderDir,
		"filters":   filters,
	})
}

func (h *CategoriaHandler) BuscarPorID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	categoria, err := h.uc.BuscarPorID(c.Request.Context(), uint(id))
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

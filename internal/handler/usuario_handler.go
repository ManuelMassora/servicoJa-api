package handler

import (
	"strconv"

	"github.com/ManuelMassora/servicoJa-api/internal/usecases"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UsuarioHandler struct {
	uc usecases.UsuarioUseCase
}

func NewUsuarioHandler(uc usecases.UsuarioUseCase) *UsuarioHandler {
	return &UsuarioHandler{uc: uc}
}

func (h *UsuarioHandler) CriarAdmin(c *gin.Context) {
	var request usecases.UsuarioRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"erro": "ao converter para JSON"})
		return
	}
	if err := h.uc.CriarAdmin(c.Request.Context(), request); err != nil {
		c.JSON(400, gin.H{"erro": "ao salvar admin, " + err.Error()})
		return
	}
	c.JSON(201, nil)
}

func (h *UsuarioHandler) CriarCliente(c *gin.Context) {
	var request usecases.UsuarioRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"erro": "ao converter para JSON"})
		return
	}
	if err := h.uc.CriarCliente(c.Request.Context(), request); err != nil {
		c.JSON(400, gin.H{"erro": "ao salvar cliente, " + err.Error()})
		return
	}
	c.JSON(201, nil)
}

func (h *UsuarioHandler) CriarPrestador(c *gin.Context) {
	var request usecases.PrestadorRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"erro": "ao converter para JSON"})
		return
	}
	if err := h.uc.CriarPrestador(c.Request.Context(), request); err != nil {
		c.JSON(400, gin.H{"erro": "ao salvar prestador, " + err.Error()})
		return
	}
	c.JSON(201, nil)
}

func (h *UsuarioHandler) ListarTodosUsuarios(c *gin.Context) {
	filters := make(map[string]interface{})

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
			if id, err := strconv.ParseInt(v, 10, 64); err == nil {
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

	usuarios, err := h.uc.ListarTodosUsuarios(c.Request.Context(), filters, orderBy, orderDir, pageSize, offset)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"data":      usuarios,
		"page":      page,
		"pageSize":  pageSize,
		"orderBy":   orderBy,
		"direction": orderDir,
		"filters":   filters,
	})
}

func (h *UsuarioHandler) ListarPrestadores(c *gin.Context) {
	filters := make(map[string]interface{})
	var statusDisponivel interface{}

	for key, vals := range c.Request.URL.Query() {
		// Parâmetros de controle ignorados
		if key == "orderBy" || key == "orderDir" || key == "page" || key == "pageSize" {
			continue
		}
		if len(vals) == 0 || vals[0] == "" {
			continue
		}
		v := vals[0]

		// --- Tratamento especial para booleanos e inteiros ---
		if key == "id" {
			if id, err := strconv.ParseInt(v, 10, 64); err == nil {
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

		// NOVO: Trata o parâmetro status_disponivel
		if key == "status_disponivel" {
			if b, err := strconv.ParseBool(v); err == nil {
				// Atribui o valor booleano à variável separada
				statusDisponivel = b
				continue
			}
			// Se falhar o parse, statusDisponivel continuará nil (ou você pode querer retornar um erro)
		}

		// default: keep as string (para filtros LIKE)
		// NOTA: Certifique-se de que os nomes das chaves aqui são os nomes das colunas no banco de dados
		filters[key] = v
	}

	orderBy := c.Query("orderBy")
	orderDir := c.Query("orderDir")

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	prestadores, err := h.uc.ListarPrestadores(
		c.Request.Context(),
		filters,
		statusDisponivel,
		orderBy,
		orderDir,
		pageSize,
		offset,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Adiciona o statusDisponivel ao JSON de resposta, se foi fornecido
	filters["status_disponivel"] = statusDisponivel

	c.JSON(200, gin.H{
		"data":      prestadores,
		"page":      page,
		"pageSize":  pageSize,
		"orderBy":   orderBy,
		"direction": orderDir,
		"filters":   filters,
	})
}

func (h *UsuarioHandler) ListarPrestadoresPorLocalizacao(c *gin.Context) {
	latitude, err := strconv.ParseFloat(c.Query("latitude"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Latitude inválida"})
		return
	}
	longitude, err := strconv.ParseFloat(c.Query("longitude"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Longitude inválida"})
		return
	}
	radius, err := strconv.ParseFloat(c.DefaultQuery("radius", "10"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Raio inválido"})
		return
	}

	filters := make(map[string]interface{})
	for key, vals := range c.Request.URL.Query() {
		if key == "latitude" || key == "longitude" || key == "radius" || key == "orderBy" || key == "orderDir" || key == "page" || key == "pageSize" {
			continue
		}
		if len(vals) > 0 && vals[0] != "" {
			filters[key] = vals[0]
		}
	}

	orderBy := c.Query("orderBy")
	orderDir := c.Query("orderDir")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	offset := (page - 1) * pageSize

	prestadores, err := h.uc.ListarPrestadoresPorLocalizacao(c.Request.Context(), latitude, longitude, radius, filters, orderBy, orderDir, pageSize, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      prestadores,
		"page":      page,
		"pageSize":  pageSize,
		"orderBy":   orderBy,
		"direction": orderDir,
		"filters":   filters,
	})
}

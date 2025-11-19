package handler

import (
	"net/http"
	"strconv"

	"github.com/ManuelMassora/servicoJa-api/internal/usecases"
	"github.com/gin-gonic/gin"
)

type AgendamentoHandler struct {
	uc usecases.AgendamentoUC 
}

func NewAgendamentoHandler(uc usecases.AgendamentoUC) *AgendamentoHandler {
	return &AgendamentoHandler{uc: uc}
}

func (h *AgendamentoHandler) Criar(c *gin.Context) {
	var req usecases.AgendamentoRequest
	
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	idCliente, err := getUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := h.uc.Criar(c.Request.Context(), &req, idCliente); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao criar agendamento: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"message": "Agendamento criado com sucesso!"})
}

func (h *AgendamentoHandler) Buscar(c *gin.Context) {
	
	idParam := c.Param("id")
	agendamentoID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de agendamento inválido"})
		return
	}

	idUsuario, err := getUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.uc.Buscar(c.Request.Context(), uint(agendamentoID), idUsuario)
	if err != nil {
		if err.Error() == "acesso negado: você não é o cliente nem o prestador deste agendamento" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Agendamento não encontrado ou falha na busca: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, resp)
}

func (h *AgendamentoHandler) Aceitar(c *gin.Context) {
	
	idParam := c.Param("id")
	agendamentoID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de agendamento inválido"})
		return
	}

	idUsuario, err := getUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	err = h.uc.Aceitar(c.Request.Context(), uint(agendamentoID), idUsuario)
	if err != nil {
		if err.Error() == "acesso negado: você não é o cliente nem o prestador deste agendamento" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Agendamento não encontrado ou falha na busca: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, nil)
}

func (h *AgendamentoHandler) Recusar(c *gin.Context) {
	
	idParam := c.Param("id")
	agendamentoID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de agendamento inválido"})
		return
	}

	idUsuario, err := getUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	err = h.uc.Recusar(c.Request.Context(), uint(agendamentoID), idUsuario)
	if err != nil {
		if err.Error() == "acesso negado: você não é o cliente nem o prestador deste agendamento" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Agendamento não encontrado ou falha na busca: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, nil)
}

func (h *AgendamentoHandler) Cancelar(c *gin.Context) {
	
	idParam := c.Param("id")
	agendamentoID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de agendamento inválido"})
		return
	}

	idUsuario, err := getUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	err = h.uc.Cancelar(c.Request.Context(), uint(agendamentoID), idUsuario)
	if err != nil {
		if err.Error() == "acesso negado: você não é o cliente nem o prestador deste agendamento" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Agendamento não encontrado ou falha na busca: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, nil)
}

func (h *AgendamentoHandler) Listar(c *gin.Context) {
	filters := make(map[string]interface{})

	for key, vals := range c.Request.URL.Query() {
		if key == "orderBy" || key == "orderDir" || key == "page" || key == "pageSize" {
			continue
		}
		if len(vals) == 0 || vals[0] == "" {
			continue
		}

		v := vals[0]

		if key == "id" {
			if id, err := strconv.ParseUint(v, 10, 64); err == nil {
				filters[key] = id
				continue
			}
		}

		if key == "status" {
			if b, err := strconv.ParseBool(v); err == nil {
				filters[key] = b
				continue
			}
		}

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

	agendamentos, err := h.uc.Listar(
		c.Request.Context(),
		filters,
		orderBy,
		orderDir,
		pageSize,
		offset,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      agendamentos,
		"page":      page,
		"pageSize":  pageSize,
		"orderBy":   orderBy,
		"direction": orderDir,
		"filters":   filters,
	})
}

func (h *AgendamentoHandler) ListarPorClienteID(c *gin.Context) {
	idUsuario, err := getUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	filters := make(map[string]interface{})

	for key, vals := range c.Request.URL.Query() {
		if key == "orderBy" || key == "orderDir" || key == "page" || key == "pageSize" {
			continue
		}
		if len(vals) == 0 || vals[0] == "" {
			continue
		}

		v := vals[0]

		if key == "id" {
			if id, err := strconv.ParseUint(v, 10, 64); err == nil {
				filters[key] = id
				continue
			}
		}

		if key == "status" {
			if b, err := strconv.ParseBool(v); err == nil {
				filters[key] = b
				continue
			}
		}

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

	agendamentos, err := h.uc.ListarPorClienteID(
		c.Request.Context(),
		uint(idUsuario),
		filters,
		orderBy,
		orderDir,
		pageSize,
		offset,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      agendamentos,
		"page":      page,
		"pageSize":  pageSize,
		"orderBy":   orderBy,
		"direction": orderDir,
		"filters":   filters,
	})
}

func (h *AgendamentoHandler) ListarPorCatalogID(c *gin.Context) {
	idUsuario, err := getUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	idParam := c.Param("catalogoID")
	idCatalogo, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de catálogo inválido"})
		return
	}

	filters := make(map[string]interface{})

	for key, vals := range c.Request.URL.Query() {
		if key == "orderBy" || key == "orderDir" || key == "page" || key == "pageSize" {
			continue
		}
		if len(vals) == 0 || vals[0] == "" {
			continue
		}

		v := vals[0]

		if key == "id" {
			if id, err := strconv.ParseUint(v, 10, 64); err == nil {
				filters[key] = id
				continue
			}
		}

		if key == "status" {
			if b, err := strconv.ParseBool(v); err == nil {
				filters[key] = b
				continue
			}
		}

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

	agendamentos, err := h.uc.ListarPorCatalogID(
		c.Request.Context(),
		uint(idUsuario),
		uint(idCatalogo),
		filters,
		orderBy,
		orderDir,
		pageSize,
		offset,
	)
	if err != nil {
		if err.Error() == "acesso negado: você não é o prestador deste catálogo" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      agendamentos,
		"page":      page,
		"pageSize":  pageSize,
		"orderBy":   orderBy,
		"direction": orderDir,
		"filters":   filters,
	})
}

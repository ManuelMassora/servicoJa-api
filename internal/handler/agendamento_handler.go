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

func (h *AgendamentoHandler) Listar(c *gin.Context) {
	
	idUsuario, err := getUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	filters := make(map[string]interface{})
	
	if c.Query("id_catalogo") != "" {
		if catID, err := strconv.ParseUint(c.Query("id_catalogo"), 10, 32); err == nil {
			filters["id_catalogo"] = uint(catID)
		}
	}

	resp, err := h.uc.Listar(c.Request.Context(), idUsuario, filters, "", "", 0, 0) 
	if err != nil {
		if err.Error() == "acesso negado: usuário não é cliente nem prestador" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao listar agendamentos: " + err.Error()})
		return
	}
	if len(resp) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Nenhum agendamento encontrado."})
		return
	}
	c.JSON(http.StatusOK, resp)
}
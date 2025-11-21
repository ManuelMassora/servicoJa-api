package handler

import (
	"net/http"
	"strconv"
	"errors" 
	
	"github.com/ManuelMassora/servicoJa-api/internal/usecases"
	"github.com/gin-gonic/gin"
)

type PropostaHandler struct {
	uc usecases.PropostaUseCase
}

func NewPropostaHandler(uc usecases.PropostaUseCase) *PropostaHandler {
	return &PropostaHandler{uc: uc}
}

func (h *PropostaHandler) Criar(c *gin.Context) {
	var req usecases.PropostaRequest	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados de requisição inválidos", "details": err.Error()})
		return
	}
	idUsuario, err := getUsuarioID(c) 
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}
	err = h.uc.Criar(c.Request.Context(), req, uint(idUsuario))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Proposta criada com sucesso"})
}

func (h *PropostaHandler) Responder(c *gin.Context) {
	
	idPropostaStr := c.Param("id")
	idProposta, err := strconv.ParseUint(idPropostaStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID da proposta inválido"})
		return
	}

	idUsuario, err := getUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}
	
	action := c.Query("status")
	var aceitar bool
	
	switch action {
	case "aceitar":
		aceitar = true
	case "rejeitar":
		aceitar = false
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ação de 'status' inválida. Use 'aceitar' ou 'rejeitar'"})
		return
	}
	
	err = h.uc.Responder(c.Request.Context(), uint(idProposta), uint(idUsuario), aceitar)
	if err != nil {
		
		if errors.Is(err, errors.New("acesso negado: apenas o cliente dono da vaga pode responder a proposta")) || 
		   errors.Is(err, errors.New("acesso negado: apenas propostas pendentes podem ser respondidas")) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	msg := "Proposta aceita com sucesso e serviço iniciado"
	if !aceitar {
		msg = "Proposta rejeitada com sucesso"
	}

	c.JSON(http.StatusOK, gin.H{"message": msg})
}

func (h *PropostaHandler) Cancelar(c *gin.Context) {	
	idPropostaStr := c.Param("id")
	idProposta, err := strconv.ParseUint(idPropostaStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID da proposta inválido"})
		return
	}
	
	idUsuario, err := getUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	err = h.uc.Cancelar(c.Request.Context(), uint(idProposta), uint(idUsuario))
	if err != nil {
		
		if errors.Is(err, errors.New("acesso negado: apenas o prestador que fez a proposta pode cancelá-la")) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Proposta cancelada com sucesso"})
}

func (h *PropostaHandler) ListarPorVaga(c *gin.Context) {	
	idVagaStr := c.Param("idVaga")
	idVaga, err := strconv.ParseUint(idVagaStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID da vaga inválido"})
		return
	}

	idUsuario, err := getUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	filters := ExtractFilters(c)
	orderBy := c.Query("orderBy")
	orderDir := c.Query("orderDir")
	limit, offset, page, pageSize := ExtractPagination(c)

	propostas, err := h.uc.ListarPorVaga(
		c.Request.Context(),
		uint(idUsuario),
		uint(idVaga),
		filters,
		orderBy,
		orderDir,
		limit,
		offset,
	)

	if err != nil {
		if errors.Is(err, errors.New("acesso negado: apenas o cliente que criou a vaga pode ver as propostas")) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
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

func (h *PropostaHandler) ListarPorPrestador(c *gin.Context) {	
	idUsuario, err := getUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	
	filters := ExtractFilters(c)
	orderBy := c.Query("orderBy")
	orderDir := c.Query("orderDir")
	limit, offset, page, pageSize := ExtractPagination(c)
		
	propostas, err := h.uc.ListarPorPrestador(
		c.Request.Context(),
		uint(idUsuario),
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
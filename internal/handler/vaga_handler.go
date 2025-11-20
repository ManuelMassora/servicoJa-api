package handler

import (
	"net/http"
	"strconv"
	
	"github.com/ManuelMassora/servicoJa-api/internal/usecases"
	"github.com/gin-gonic/gin"
)

type VagaHandler struct {
	uc usecases.VagaUseCase
}

func NewVagaHandler(uc usecases.VagaUseCase) *VagaHandler {
	return &VagaHandler{uc: uc}
}

func (h *VagaHandler) CriarVaga(c *gin.Context) {
	var req usecases.VagaRequest	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados de requisição inválidos", "details": err.Error()})
		return
	}
	
	idUsuario, err := getUsuarioID(c) 
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}
	
	err = h.uc.CriarVaga(c.Request.Context(), req, uint(idUsuario))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Vaga criada com sucesso"})
}

func (h *VagaHandler) CancelarVaga(c *gin.Context) {	
	idVagaStr := c.Param("id")
	idVaga, err := strconv.ParseUint(idVagaStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID da vaga inválido"})
		return
	}
	
	idUsuario, err := getUsuarioID(c) 
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	err = h.uc.CancelarVaga(c.Request.Context(), uint(idVaga), uint(idUsuario))
	if err != nil {
		if err.Error() == "vaga não pertence ao cliente" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Vaga cancelada com sucesso"})
}

func (h *VagaHandler) ListarVagasDisponiveis(c *gin.Context) {	
	filters := ExtractFilters(c)
	orderBy := c.Query("orderBy")
	orderDir := c.Query("orderDir")
	limit, offset, page, pageSize := ExtractPagination(c)

	vagas, err := h.uc.ListarVagasDisponiveis(
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
		"data":      vagas,
		"page":      page,
		"pageSize":  pageSize,
		"orderBy":   orderBy,
		"direction": orderDir,
		"filters":   filters,
	})
}

func (h *VagaHandler) ListarPorCliente(c *gin.Context) {
	idUsuario, err := getUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	filters := ExtractFilters(c)
	orderBy := c.Query("orderBy")
	orderDir := c.Query("orderDir")
	limit, offset, page, pageSize := ExtractPagination(c)
		
	vagas, err := h.uc.ListarPorCliente(
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
		"data":      vagas,
		"page":      page,
		"pageSize":  pageSize,
		"orderBy":   orderBy,
		"direction": orderDir,
		"filters":   filters,
	})
}
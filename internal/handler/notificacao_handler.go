package handler

import (
	"net/http"
	"strconv"

	"github.com/ManuelMassora/servicoJa-api/internal/usecases"
	"github.com/gin-gonic/gin"	
)

type NotificacaoResponse usecases.NotificacaoResponse

type NotificacaoHandler struct {
	uc usecases.NotificacaoUseCase 
}

func NewNotificacaoHandler(uc usecases.NotificacaoUseCase) *NotificacaoHandler {
	return &NotificacaoHandler{uc: uc} 
}

func (h *NotificacaoHandler) ListarPorUsuario(c *gin.Context) {
	idUsuario, err := getUsuarioID(c) 
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	filters := ExtractFilters(c)
	orderBy := c.Query("orderBy")
	orderDir := c.Query("orderDir")
	limit, offset, page, pageSize := ExtractPagination(c)
	
	notificacoes, err := h.uc.ListarPorUsuario(
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
		"data":      notificacoes,
		"page":      page,
		"pageSize":  pageSize,
		"orderBy":   orderBy,
		"direction": orderDir,
		"filters":   filters,
	})
}

func (h *NotificacaoHandler) MarcarComoLida(c *gin.Context) {
	idNotificacaoStr := c.Param("id")
	idNotificacao, err := strconv.ParseUint(idNotificacaoStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID da notificação inválido"})
		return
	}

	idUsuario, err := getUsuarioID(c) 
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}	
	err = h.uc.MarcarComoLida(c.Request.Context(), uint(idNotificacao), uint(idUsuario))
	if err != nil {
		
		if err.Error() == "acesso negado: nao pode marcar essa notificacao como lida" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Notificação marcada como lida com sucesso"})
}
package handler

import (
	"errors"
	"net/http"

	"github.com/ManuelMassora/servicoJa-api/internal/services"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	uc services.AuthUSer
}

func NewAuthHandler(uc services.AuthUSer) *AuthHandler {
	return &AuthHandler{uc: uc}
}

func (h *AuthHandler) Authenticate(ctx *gin.Context) {
	var input services.AuthRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Erro": "Requisição inválida:"})
		return
	}
	
	token, err := h.uc.Authenticate(ctx, input)
	if err != nil {
		if errors.Is(err, errors.New("credenciais inválidas")) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusAccepted, token)
}
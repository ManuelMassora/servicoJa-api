package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func getUsuarioID(c *gin.Context) (uint, error) {
	prestadorIDVal, exists := c.Get("userID")
	if !exists {
		return 0, errors.New("User ID não encontrado no contexto. Middleware de autenticação ausente?")
	}
	id, ok := prestadorIDVal.(uint)
	if !ok {

		return 0, errors.New("userID no contexto com formato inválido. Esperado: uint")
	}
	return id, nil
}
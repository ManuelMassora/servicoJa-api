package routes

import (
	"net/http"

	"github.com/ManuelMassora/servicoJa-api/internal/di"
	"github.com/gin-gonic/gin"
)

func SetRoutes(server *gin.Engine, container *di.Container) {
	server.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"Mensagem": "Rota nao encontrada scu scu"})
	})

	iniciar := server.Group("iniciar")
	{
		iniciar.POST("/admin", container.UsuarioH.CriarAdmin)
		iniciar.POST("/cliente", container.UsuarioH.CriarCliente)
		iniciar.POST("/prestador", container.UsuarioH.CriarPrestador)
	}

	usuario := server.Group("usuario")
	{
		usuario.GET("", container.UsuarioH.ListarTodosUsuarios)
		usuario.GET("/prestador", container.UsuarioH.ListarPrestadores)
	}

	categoria := server.Group("categoria")
	{
		categoria.POST("", container.CategoriaH.Criar)
		categoria.PATCH(":id", container.CategoriaH.Editar)
		categoria.GET("", container.CategoriaH.Listar)
		categoria.GET(":id", container.CategoriaH.BuscarPorID)
	}
}
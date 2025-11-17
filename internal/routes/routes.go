package routes

import (
	"net/http"

	"github.com/ManuelMassora/servicoJa-api/internal/di"
	"github.com/ManuelMassora/servicoJa-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetRoutes(server *gin.Engine, container *di.Container) {
	server.Use(middleware.RateLimitMiddleware())
	server.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"Mensagem": "Rota nao encontrada scu scu"})
	})
	server.POST("/login", container.AuthHandler.Authenticate)

	iniciar := server.Group("iniciar")
	{
		iniciar.POST("/admin", container.UsuarioH.CriarAdmin)
		iniciar.POST("/cliente", container.UsuarioH.CriarCliente)
		iniciar.POST("/prestador", container.UsuarioH.CriarPrestador)
	}

	usuario := server.Group("usuario", middleware.Auth())
	{
		usuario.GET("", container.UsuarioH.ListarTodosUsuarios)
		usuario.GET("/prestador", container.UsuarioH.ListarPrestadores)
	}

	categoria := server.Group("categoria", middleware.Auth())
	{
		categoria.POST("", middleware.HasRole("ADMIN"), container.CategoriaH.Criar)
		categoria.PATCH(":id", middleware.HasRole("ADMIN"), container.CategoriaH.Editar)
		categoria.GET("", container.CategoriaH.Listar)
		categoria.GET(":id", container.CategoriaH.BuscarPorID)
	}
	catalogo := server.Group("catalogo", middleware.Auth())
	{
		catalogo.POST("", middleware.HasRole("PRESTADOR"), container.CatalogoH.Criar)
		catalogo.PATCH("/:id", middleware.HasRole("PRESTADOR"), container.CatalogoH.Editar)
		catalogo.DELETE("/:id", middleware.HasRole("PRESTADOR"), container.CatalogoH.Apagar)
		catalogo.GET("", container.CatalogoH.Listar)
		catalogo.GET("/:prestadorID", container.CatalogoH.ListarPorPrestador)
	}
}
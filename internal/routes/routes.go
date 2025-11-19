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
	agendamento := server.Group("agendamento", middleware.Auth())
	{
		agendamento.POST("", middleware.HasRole("CLIENTE"), container.AgendamentoH.Criar)
		agendamento.GET("", middleware.HasRole("ADMIN"), container.AgendamentoH.Listar)
		agendamento.POST("/aceitar/:id", middleware.HasRole("PRESTADOR"), container.AgendamentoH.Aceitar)
		agendamento.POST("/recusar/:id", middleware.HasRole("PRESTADOR"), container.AgendamentoH.Recusar)
		agendamento.POST("/cancelar/:id", middleware.HasRole("CLIENTE"), container.AgendamentoH.Cancelar)
		agendamento.GET("/cliente", middleware.HasRole("CLIENTE"), container.AgendamentoH.ListarPorClienteID)
		agendamento.GET("/:catalogoID", middleware.HasRole("PRESTADOR"), container.AgendamentoH.ListarPorCatalogID)
	}
}
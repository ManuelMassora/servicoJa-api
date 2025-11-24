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
		usuario.GET("/prestador/location", container.UsuarioH.ListarPrestadoresPorLocalizacao)
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
		catalogo.GET("/location", container.CatalogoH.ListarPorLocalizacao)
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
		agendamento.GET("/location", container.AgendamentoH.ListarPorLocalizacao)
		agendamento.GET("/:catalogoID", middleware.HasRole("PRESTADOR"), container.AgendamentoH.ListarPorCatalogID)
	}
	servico := server.Group("servico", middleware.Auth())
	{
		servico.POST("/finalizar/:id", middleware.HasRole("PRESTADOR"), container.ServicoH.FinalizarServico)
		servico.POST("/cancelar/:id", middleware.HasRole("CLIENTE", "PRESTADOR"), container.ServicoH.CancelarServico)
		servico.GET("/cliente", middleware.HasRole("CLIENTE"), container.ServicoH.ListarPorCliente)
		servico.GET("/prestador", middleware.HasRole("PRESTADOR"), container.ServicoH.ListarPorPrestador)
		servico.GET("/location", container.ServicoH.ListarPorLocalizacao)
	}
	vagas := server.Group("vagas", middleware.Auth())
	{
		vagas.POST("", middleware.HasRole("CLIENTE"), container.VagaH.CriarVaga)
		vagas.POST("cancelar/:id", middleware.HasRole("CLIENTE"), container.VagaH.CancelarVaga)
		vagas.GET("", container.VagaH.ListarVagasDisponiveis)
		vagas.GET("/location", container.VagaH.ListarPorLocalizacao)
		vagas.GET("/cliente", container.VagaH.ListarPorCliente)
	}
	propostas := server.Group("propostas", middleware.Auth())
	{
		propostas.POST("", middleware.HasRole("PRESTADOR"), container.PropostaH.Criar)
		propostas.POST("/responder/:id", middleware.HasRole("CLIENTE"), container.PropostaH.Responder)
		propostas.POST("/cancelar/:id", middleware.HasRole("PRESTADOR"), container.PropostaH.Cancelar)
		propostas.GET("/prestador", middleware.HasRole("PRESTADOR"), container.PropostaH.ListarPorPrestador)
		propostas.GET("/cliente/:idVaga", middleware.HasRole("CLIENTE"), container.PropostaH.ListarPorVaga)
	}
	notificacao := server.Group("notificacao", middleware.Auth())
	{
		notificacao.GET("", container.NotificacaoH.ListarPorUsuario)
		notificacao.POST("/lida/:id", container.NotificacaoH.MarcarComoLida)
	}
	avaliacao := server.Group("avaliacao", middleware.Auth())
	{
		avaliacao.POST("/cliente/:id", middleware.HasRole("CLIENTE"), container.AvaliacaoH.CriarAvaliacao)
		avaliacao.GET("/cliente", middleware.HasRole("CLIENTE"), container.AvaliacaoH.ListarAvaliacoesPorCliente)
		avaliacao.GET("", container.AvaliacaoH.ListarAvaliacoesPorPrestador)
	}
}
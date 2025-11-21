package di

import (
	"github.com/ManuelMassora/servicoJa-api/internal/handler"
	"github.com/ManuelMassora/servicoJa-api/internal/infra/repo"
	"github.com/ManuelMassora/servicoJa-api/internal/middleware"
	"github.com/ManuelMassora/servicoJa-api/internal/services"
	"github.com/ManuelMassora/servicoJa-api/internal/usecases"
	"gorm.io/gorm"
)

type Container struct {
	//Repositories
	CategoriaRepo *repo.CategoriaRepository
	UsuarioRepo *repo.UsuarioRepository
	ClienteRepo *repo.ClienteRepository
	PrestadorRepo *repo.PrestadorRepository
	CatalogoRepo *repo.CatalogoRepository
	AgendamentoRepo *repo.AgendamentoRepo
	ServicoRepo *repo.ServicoRepository
	VagaRepo *repo.VagaRepository
	PropostaRepo *repo.PropostaRepository
	NotificacaoRepo *repo.NotificacaoRepo

	//Services
	AuthService *services.AuthUSer
	JwtService 	*middleware.JwtService

	//UseCases
	CategoriaUC *usecases.CategoriaUseCase
	UsuarioUC *usecases.UsuarioUseCase
	CatalogoUC *usecases.CatalogoUseCase
	AgendamentoUC *usecases.AgendamentoUC
	ServicoUC *usecases.ServicoUseCase
	PropostaUC *usecases.PropostaUseCase
	VagaUC *usecases.VagaUseCase
	NotificacaoUC *usecases.NotificacaoUseCase
		
	//Handler
	CategoriaH	*handler.CategoriaHandler
	UsuarioH *handler.UsuarioHandler
	AuthHandler *handler.AuthHandler
	CatalogoH *handler.CatalogoHandler
	AgendamentoH *handler.AgendamentoHandler
	ServicoH *handler.ServicoHandler
	PropostaH *handler.PropostaHandler
	VagaH *handler.VagaHandler
	NotificacaoH *handler.NotificacaoHandler
}

func NewContainer(db *gorm.DB) *Container {
	c := &Container{}

	//Init Repositories
	c.CategoriaRepo = repo.NewCategoriaRepository(db).(*repo.CategoriaRepository)
	c.UsuarioRepo = repo.NewUsuarioRepository(db).(*repo.UsuarioRepository)
	c.ClienteRepo = repo.NewClienteRepository(db).(*repo.ClienteRepository)
	c.PrestadorRepo = repo.NewPrestadorRepository(db).(*repo.PrestadorRepository)
	c.CatalogoRepo = repo.NewCatalogoRepository(db).(*repo.CatalogoRepository)
	c.AgendamentoRepo = repo.NewAgendamentoRepo(db).(*repo.AgendamentoRepo)
	c.ServicoRepo = repo.NewServicoRepository(db).(*repo.ServicoRepository)
	c.VagaRepo = repo.NewVagaRepository(db).(*repo.VagaRepository)
	c.PropostaRepo = repo.NewPropostaRepository(db).(*repo.PropostaRepository)
	c.NotificacaoRepo = repo.NewNotificacaoRepo(db).(*repo.NotificacaoRepo)

	//Init Services
	c.JwtService = middleware.NewJWTService()
	c.AuthService = services.NewAuthUser(c.UsuarioRepo, c.JwtService)
	//Init UseCases
	c.CategoriaUC = usecases.NewCategoriaUseCase(c.CategoriaRepo)
	c.UsuarioUC = usecases.NewUsuarioUseCase(c.UsuarioRepo, c.ClienteRepo, c.PrestadorRepo)
	c.CatalogoUC = usecases.NewCatalogoUC(c.CatalogoRepo)
	c.AgendamentoUC = usecases.NewAgendamentoUC(c.AgendamentoRepo, c.CatalogoRepo, c.ServicoRepo, c.NotificacaoRepo, c.UsuarioRepo)
	c.ServicoUC = usecases.NewServicoUseCase(c.ServicoRepo, c.AgendamentoRepo, c.VagaRepo, c.NotificacaoRepo, c.UsuarioRepo)
	c.PropostaUC = usecases.NewPropostaUseCase(c.PropostaRepo, c.VagaRepo, c.ServicoRepo, c.NotificacaoRepo, c.UsuarioRepo)
	c.VagaUC = usecases.NewVagaUseCase(c.VagaRepo)
	c.NotificacaoUC = usecases.NewNotificacaoUseCase(c.NotificacaoRepo, c.UsuarioRepo)

	//Init Handler
	c.CategoriaH = handler.NewCategoriaHandler(*c.CategoriaUC)
	c.UsuarioH = handler.NewUsuarioHandler(*c.UsuarioUC)
	c.AuthHandler = handler.NewAuthHandler(*c.AuthService)
	c.CatalogoH = handler.NewCatalogoHandler(*c.CatalogoUC)
	c.AgendamentoH = handler.NewAgendamentoHandler(*c.AgendamentoUC)
	c.ServicoH = handler.NewServicoHandler(*c.ServicoUC)
	c.PropostaH = handler.NewPropostaHandler(*c.PropostaUC)
	c.VagaH = handler.NewVagaHandler(*c.VagaUC)
	c.NotificacaoH = handler.NewNotificacaoHandler(*c.NotificacaoUC)
	return c
}
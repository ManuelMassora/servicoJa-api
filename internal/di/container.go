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
	//Handler
	CategoriaH	*handler.CategoriaHandler
	UsuarioH *handler.UsuarioHandler
	AuthHandler *handler.AuthHandler
	CatalogoH *handler.CatalogoHandler
	AgendamentoH *handler.AgendamentoHandler
	ServicoH *handler.ServicoHandler
	PropostaH *handler.PropostaHandler
	VagaH *handler.VagaHandler
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

	//Init Services
	c.JwtService = middleware.NewJWTService()
	c.AuthService = services.NewAuthUser(c.UsuarioRepo, c.JwtService)
	//Init UseCases
	c.CategoriaUC = usecases.NewCategoriaUseCase(c.CategoriaRepo)
	c.UsuarioUC = usecases.NewUsuarioUseCase(c.UsuarioRepo, c.ClienteRepo, c.PrestadorRepo)
	c.CatalogoUC = usecases.NewCatalogoUC(c.CatalogoRepo, c.PrestadorRepo)
	c.AgendamentoUC = usecases.NewAgendamentoUC(c.AgendamentoRepo, c.ClienteRepo, c.PrestadorRepo, c.CatalogoRepo, c.ServicoRepo)
	c.ServicoUC = usecases.NewServicoUseCase(c.ServicoRepo, c.PrestadorRepo, c.ClienteRepo, c.AgendamentoRepo, c.VagaRepo)
	c.PropostaUC = usecases.NewPropostaUseCase(c.PropostaRepo, c.PrestadorRepo, c.VagaRepo, c.ClienteRepo, c.ServicoRepo)
	c.VagaUC = usecases.NewVagaUseCase(c.VagaRepo, c.ClienteRepo)

	//Init Handler
	c.CategoriaH = handler.NewCategoriaHandler(*c.CategoriaUC)
	c.UsuarioH = handler.NewUsuarioHandler(*c.UsuarioUC)
	c.AuthHandler = handler.NewAuthHandler(*c.AuthService)
	c.CatalogoH = handler.NewCatalogoHandler(*c.CatalogoUC)
	c.AgendamentoH = handler.NewAgendamentoHandler(*c.AgendamentoUC)
	c.ServicoH = handler.NewServicoHandler(*c.ServicoUC)
	c.PropostaH = handler.NewPropostaHandler(*c.PropostaUC)
	c.VagaH = handler.NewVagaHandler(*c.VagaUC)
	return c
}
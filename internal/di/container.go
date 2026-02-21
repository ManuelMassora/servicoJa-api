package di

import (
	"github.com/ManuelMassora/servicoJa-api/internal/config"
	"github.com/ManuelMassora/servicoJa-api/internal/handler"
	gatewaympesa "github.com/ManuelMassora/servicoJa-api/internal/infra/gateway_mpesa"
	"github.com/ManuelMassora/servicoJa-api/internal/infra/repo"
	"github.com/ManuelMassora/servicoJa-api/internal/middleware"
	"github.com/ManuelMassora/servicoJa-api/internal/services"
	"github.com/ManuelMassora/servicoJa-api/internal/usecases"
	"gorm.io/gorm"
)

type Container struct {
	//Repositories
	CategoriaRepo          *repo.CategoriaRepository
	UsuarioRepo            *repo.UsuarioRepository
	ClienteRepo            *repo.ClienteRepository
	PrestadorRepo          *repo.PrestadorRepository
	CatalogoRepo           *repo.CatalogoRepository
	AgendamentoRepo        *repo.AgendamentoRepo
	ServicoRepo            *repo.ServicoRepository
	VagaRepo               *repo.VagaRepository
	PropostaRepo           *repo.PropostaRepository
	NotificacaoRepo        *repo.NotificacaoRepo
	AvaliacaoRepo          *repo.AvaliacaoRepo
	AnexoImagemRepo        *repo.AnexoImagemRepo
	GaleriaRepo            *repo.GaleriaRepo
	CategoriaPrestadorRepo *repo.CategoriaPrestadorRepo
	PagamentoRepo          *repo.PagamentoRepository

	//Services
	AuthService  *services.AuthUSer
	JwtService   *middleware.JwtService
	Uploader     *services.SupabaseUploader
	MpesaGateway *gatewaympesa.MpesaGateway

	//UseCases
	CategoriaUC          *usecases.CategoriaUseCase
	UsuarioUC            *usecases.UsuarioUseCase
	CatalogoUC           *usecases.CatalogoUseCase
	AgendamentoUC        *usecases.AgendamentoUC
	ServicoUC            *usecases.ServicoUseCase
	PropostaUC           *usecases.PropostaUseCase
	VagaUC               *usecases.VagaUseCase
	NotificacaoUC        *usecases.NotificacaoUseCase
	AvaliacaoUC          *usecases.AvaliacaoUseCase
	AnexoImagemUC        *usecases.AnexoImagemUseCase
	GaleriaUC            *usecases.GaleriaUseCase
	CategoriaPrestadorUC *usecases.CategoriaPrestadorUsecase
	PagamentoUC          *usecases.PagamentoUseCase

	//Handler
	CategoriaH   *handler.CategoriaHandler
	UsuarioH     *handler.UsuarioHandler
	AuthHandler  *handler.AuthHandler
	CatalogoH    *handler.CatalogoHandler
	AgendamentoH *handler.AgendamentoHandler
	ServicoH     *handler.ServicoHandler
	PropostaH    *handler.PropostaHandler
	VagaH        *handler.VagaHandler
	NotificacaoH *handler.NotificacaoHandler
	AvaliacaoH   *handler.AvaliacaoHandler
	// AnexoImagemH *handler.AnexoImagemHandler
	GaleriaH            *handler.GaleriaHandler
	CategoriaPrestadorH *handler.CategoriaPrestadorHandler
	PagamentoH          *handler.PagamentoHandler
}

func NewContainer(db *gorm.DB, cfg *config.Config, jwtService *middleware.JwtService) *Container {
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
	c.AvaliacaoRepo = repo.NewAvaliacaoRepository(db).(*repo.AvaliacaoRepo)
	c.AnexoImagemRepo = repo.NewAnexoImagemRepo(db).(*repo.AnexoImagemRepo)
	c.GaleriaRepo = repo.NewGaleriaRepo(db).(*repo.GaleriaRepo)
	c.CategoriaPrestadorRepo = repo.NewCategoriaPrestadorRepo(db).(*repo.CategoriaPrestadorRepo)
	c.PagamentoRepo = repo.NewPagamentoRepository(db).(*repo.PagamentoRepository)

	//Init Services
	c.AuthService = services.NewAuthUser(c.UsuarioRepo, jwtService)
	c.Uploader = services.NewSupabaseUploader(cfg)
	isSandbox := true
	httpClient := gatewaympesa.NewMpesaHttpClient(isSandbox)
	c.MpesaGateway = &gatewaympesa.MpesaGateway{
		Client: httpClient,
	}
	//Init UseCases
	c.CategoriaUC = usecases.NewCategoriaUseCase(c.CategoriaRepo)
	c.UsuarioUC = usecases.NewUsuarioUseCase(c.UsuarioRepo, c.ClienteRepo, c.PrestadorRepo, c.GaleriaRepo)
	c.CatalogoUC = usecases.NewCatalogoUC(c.CatalogoRepo, c.AnexoImagemRepo)

	// PagamentoUC needs a lot of things
	c.PagamentoUC = usecases.NewPagamentoUseCase(
		c.PagamentoRepo,
		c.ServicoRepo,
		c.VagaRepo,
		c.AgendamentoRepo,
		c.UsuarioRepo,
		c.NotificacaoRepo,
		c.MpesaGateway,
		cfg.MpesaAppKey, cfg.MpesaAppPub, cfg.MpesaShortCode,
	)

	c.AgendamentoUC = usecases.NewAgendamentoUC(c.AgendamentoRepo, c.CatalogoRepo, c.ServicoRepo, c.NotificacaoRepo, c.UsuarioRepo, c.AnexoImagemRepo, c.PagamentoRepo, c.PagamentoUC)
	c.ServicoUC = usecases.NewServicoUseCase(c.ServicoRepo, c.AgendamentoRepo, c.VagaRepo, c.NotificacaoRepo, c.UsuarioRepo, c.PagamentoUC)
	c.PropostaUC = usecases.NewPropostaUseCase(c.PropostaRepo, c.VagaRepo, c.ServicoRepo, c.NotificacaoRepo, c.UsuarioRepo, c.PagamentoRepo)
	c.VagaUC = usecases.NewVagaUseCase(c.VagaRepo, c.AnexoImagemRepo, c.UsuarioRepo, c.PagamentoRepo, c.PagamentoUC)
	c.NotificacaoUC = usecases.NewNotificacaoUseCase(c.NotificacaoRepo, c.UsuarioRepo)
	c.AvaliacaoUC = usecases.NewAvaliacaoUseCase(c.AvaliacaoRepo, c.ServicoRepo, c.NotificacaoRepo, c.UsuarioRepo)
	c.GaleriaUC = usecases.NewGaleriaUseCase(c.GaleriaRepo)
	c.CategoriaPrestadorUC = usecases.NewCategoriaPrestadorUsecase(c.CategoriaPrestadorRepo)

	//Init Handler
	c.CategoriaH = handler.NewCategoriaHandler(*c.CategoriaUC)
	c.UsuarioH = handler.NewUsuarioHandler(*c.UsuarioUC, c.Uploader)
	c.AuthHandler = handler.NewAuthHandler(*c.AuthService)
	c.CatalogoH = handler.NewCatalogoHandler(*c.CatalogoUC, c.Uploader)
	c.AgendamentoH = handler.NewAgendamentoHandler(*c.AgendamentoUC, c.Uploader)
	c.ServicoH = handler.NewServicoHandler(*c.ServicoUC)
	c.PropostaH = handler.NewPropostaHandler(*c.PropostaUC)
	c.VagaH = handler.NewVagaHandler(*c.VagaUC, c.Uploader)
	c.NotificacaoH = handler.NewNotificacaoHandler(*c.NotificacaoUC)
	c.AvaliacaoH = handler.NewAvaliacaoHandler(*c.AvaliacaoUC)
	c.GaleriaH = handler.NewGaleriaHandler(c.GaleriaUC, c.Uploader)
	c.CategoriaPrestadorH = handler.NewCategoriaPrestadorHandler(*c.CategoriaPrestadorUC)
	c.PagamentoH = handler.NewPagamentoHandler(*c.PagamentoUC)
	return c
}

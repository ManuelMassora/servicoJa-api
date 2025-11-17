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

	//Services
	AuthService *services.AuthUSer
	JwtService 	*middleware.JwtService

	//UseCases
	CategoriaUC *usecases.CategoriaUseCase
	UsuarioUC *usecases.UsuarioUseCase
	CatalogoUC *usecases.CatalogoUseCase

	//Handler
	CategoriaH	*handler.CategoriaHandler
	UsuarioH *handler.UsuarioHandler
	AuthHandler *handler.AuthHandler
	CatalogoH *handler.CatalogoHandler
}

func NewContainer(db *gorm.DB) *Container {
	c := &Container{}

	//Init Repositories
	c.CategoriaRepo = repo.NewCategoriaRepository(db).(*repo.CategoriaRepository)
	c.UsuarioRepo = repo.NewUsuarioRepository(db).(*repo.UsuarioRepository)
	c.ClienteRepo = repo.NewClienteRepository(db).(*repo.ClienteRepository)
	c.PrestadorRepo = repo.NewPrestadorRepository(db).(*repo.PrestadorRepository)
	c.CatalogoRepo = repo.NewCatalogoRepository(db).(*repo.CatalogoRepository)

	//Init Services
	c.JwtService = middleware.NewJWTService()
	c.AuthService = services.NewAuthUser(c.UsuarioRepo, c.JwtService)
	//Init UseCases
	c.CategoriaUC = usecases.NewCategoriaUseCase(c.CategoriaRepo)
	c.UsuarioUC = usecases.NewUsuarioUseCase(c.UsuarioRepo, c.ClienteRepo, c.PrestadorRepo)
	c.CatalogoUC = usecases.NewCatalogoUC(c.CatalogoRepo, c.PrestadorRepo)

	//Init Handler
	c.CategoriaH = handler.NewCategoriaHandler(*c.CategoriaUC)
	c.UsuarioH = handler.NewUsuarioHandler(*c.UsuarioUC)
	c.AuthHandler = handler.NewAuthHandler(*c.AuthService)
	c.CatalogoH = handler.NewCatalogoHandler(*c.CatalogoUC)
	return c
}
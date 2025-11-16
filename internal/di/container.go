package di

import (
	"github.com/ManuelMassora/servicoJa-api/internal/handler"
	"github.com/ManuelMassora/servicoJa-api/internal/infra/repo"
	"github.com/ManuelMassora/servicoJa-api/internal/usecases"
	"gorm.io/gorm"
)

type Container struct {
	//Repositories
	CategoriaRepo *repo.CategoriaRepository
	UsuarioRepo *repo.UsuarioRepository
	ClienteRepo *repo.ClienteRepository
	PrestadorRepo *repo.PrestadorRepository

	//Services

	//UseCases
	CategoriaUC *usecases.CategoriaUseCase
	UsuarioUC *usecases.UsuarioUseCase

	//Handler
	CategoriaH	*handler.CategoriaHandler
	UsuarioH *handler.UsuarioHandler
}

func NewContainer(db *gorm.DB) *Container {
	c := &Container{}

	//Init Repositories
	c.CategoriaRepo = repo.NewCategoriaRepository(db).(*repo.CategoriaRepository)
	c.UsuarioRepo = repo.NewUsuarioRepository(db).(*repo.UsuarioRepository)
	c.ClienteRepo = repo.NewClienteRepository(db).(*repo.ClienteRepository)
	c.PrestadorRepo = repo.NewPrestadorRepository(db).(*repo.PrestadorRepository)

	//Init Services

	//Init UseCases
	c.CategoriaUC = usecases.NewCategoriaUseCase(c.CategoriaRepo)
	c.UsuarioUC = usecases.NewUsuarioUseCase(c.UsuarioRepo, c.ClienteRepo, c.PrestadorRepo)

	//Init Handler
	c.CategoriaH = handler.NewCategoriaHandler(*c.CategoriaUC)
	c.UsuarioH = handler.NewUsuarioHandler(*c.UsuarioUC)
	return c
}
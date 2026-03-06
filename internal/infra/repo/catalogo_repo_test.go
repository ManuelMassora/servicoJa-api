package repo_test

import (
	"context"
	"testing"

	"github.com/ManuelMassora/servicoJa-api/internal/infra/repo"
	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CatalogoTestSuite struct {
	suite.Suite
	repo          model.CatalogoRepo
	usuarioRepo   model.UsuarioRepo
	prestadorRepo model.PrestadorRepo
	categoriaRepo model.CategoriaRepo
	ctx           context.Context
}

func (s *CatalogoTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.repo = repo.NewCatalogoRepository(TestDB)
	s.usuarioRepo = repo.NewUsuarioRepository(TestDB)
	s.prestadorRepo = repo.NewPrestadorRepository(TestDB)
	s.categoriaRepo = repo.NewCategoriaRepository(TestDB)
}

func (s *CatalogoTestSuite) SetupTest() {
	// Limpeza em ordem reversa de dependência
	TestDB.Exec("TRUNCATE TABLE catalogos RESTART IDENTITY CASCADE")
	TestDB.Exec("TRUNCATE TABLE prestadors RESTART IDENTITY CASCADE")
	TestDB.Exec("TRUNCATE TABLE categoria RESTART IDENTITY CASCADE")
	TestDB.Exec("TRUNCATE TABLE usuarios RESTART IDENTITY CASCADE")
}

func TestCatalogoTestSuite(t *testing.T) {
	suite.Run(t, new(CatalogoTestSuite))
}

func (s *CatalogoTestSuite) createTestDependencies() (uint, uint) {
	// 1. Criar Categoria
	cat := &model.Categoria{Nome: "Manutenção"}
	s.categoriaRepo.Criar(s.ctx, cat)

	// 2. Criar Usuário Prestador
	u := &model.Usuario{
		Nome:            "Prestador Teste",
		Telefone:        "840000001",
		Senha:           "123456",
		RolePermissaoID: 2, // PRESTADOR
	}
	s.usuarioRepo.Criar(s.ctx, u)

	// 3. Criar Perfil Prestador
	p := &model.Prestador{
		IDUsuario:   u.ID,
		Nome:        u.Nome,
		Telefone:    u.Telefone,
		Localizacao: "Maputo",
	}
	s.prestadorRepo.Criar(s.ctx, p)

	return cat.ID, p.IDUsuario
}

func (s *CatalogoTestSuite) TestCatalogo_Create() {
	catID, prestadorID := s.createTestDependencies()

	c := &model.Catalogo{
		Nome:        "Reparação de AC",
		Descricao:   "Serviço completo de reparação",
		IDCategoria: catID,
		IDPrestador: prestadorID,
		TipoPreco:   "fixo",
		ValorFixo:   1500.00,
		Latitude:    -25.96,
		Longitude:   32.58,
	}

	err := s.repo.Create(s.ctx, c)
	assert.NoError(s.T(), err)
	assert.NotZero(s.T(), c.ID)
}

func (s *CatalogoTestSuite) TestCatalogo_FindByID() {
	catID, prestadorID := s.createTestDependencies()
	c := &model.Catalogo{
		Nome:        "Pintura",
		IDCategoria: catID,
		IDPrestador: prestadorID,
	}
	s.repo.Create(s.ctx, c)

	found, err := s.repo.FindByID(s.ctx, c.ID)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), found)
	assert.Equal(s.T(), "Pintura", found.Nome)
}

func (s *CatalogoTestSuite) TestCatalogo_Update() {
	catID, prestadorID := s.createTestDependencies()
	c := &model.Catalogo{
		Nome:        "Jardinagem",
		IDCategoria: catID,
		IDPrestador: prestadorID,
	}
	s.repo.Create(s.ctx, c)

	campos := map[string]interface{}{
		"nome": "Jardinagem Pro",
	}
	err := s.repo.Update(s.ctx, c.ID, campos)
	assert.NoError(s.T(), err)

	found, _ := s.repo.FindByID(s.ctx, c.ID)
	assert.Equal(s.T(), "Jardinagem Pro", found.Nome)
}

func (s *CatalogoTestSuite) TestCatalogo_FindByLocation() {
	catID, prestadorID := s.createTestDependencies()

	// Criar um serviço em Maputo (-25.96, 32.58)
	c := &model.Catalogo{
		Nome:        "Serviço Perto",
		IDCategoria: catID,
		IDPrestador: prestadorID,
		Latitude:    -25.965,
		Longitude:   32.585,
	}
	s.repo.Create(s.ctx, c)

	// Buscar num raio de 5km de um ponto central em Maputo
	results, err := s.repo.FindByLocation(s.ctx, -25.96, 32.58, 5.0, nil, "", "", 10, 0)
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), results)
	assert.Equal(s.T(), "Serviço Perto", results[0].Nome)
}

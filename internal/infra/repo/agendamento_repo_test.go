package repo_test

import (
	"context"
	"testing"
	"time"

	"github.com/ManuelMassora/servicoJa-api/internal/infra/repo"
	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AgendamentoTestSuite struct {
	suite.Suite
	repo          model.AgendamentoRepo
	usuarioRepo   model.UsuarioRepo
	clienteRepo   model.ClienteRepo
	prestadorRepo model.PrestadorRepo
	categoriaRepo model.CategoriaRepo
	catalogoRepo  model.CatalogoRepo
	ctx           context.Context
}

func (s *AgendamentoTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.repo = repo.NewAgendamentoRepo(TestDB)
	s.usuarioRepo = repo.NewUsuarioRepository(TestDB)
	s.clienteRepo = repo.NewClienteRepository(TestDB)
	s.prestadorRepo = repo.NewPrestadorRepository(TestDB)
	s.categoriaRepo = repo.NewCategoriaRepository(TestDB)
	s.catalogoRepo = repo.NewCatalogoRepository(TestDB)
}

func (s *AgendamentoTestSuite) SetupTest() {
	TestDB.Exec("TRUNCATE TABLE agendamentos RESTART IDENTITY CASCADE")
	TestDB.Exec("TRUNCATE TABLE catalogos RESTART IDENTITY CASCADE")
	TestDB.Exec("TRUNCATE TABLE categoria RESTART IDENTITY CASCADE")
	TestDB.Exec("TRUNCATE TABLE clientes RESTART IDENTITY CASCADE")
	TestDB.Exec("TRUNCATE TABLE prestadors RESTART IDENTITY CASCADE")
	TestDB.Exec("TRUNCATE TABLE usuarios RESTART IDENTITY CASCADE")
}

func TestAgendamentoTestSuite(t *testing.T) {
	suite.Run(t, new(AgendamentoTestSuite))
}

func (s *AgendamentoTestSuite) createDependencies() (uint, uint) {
	// Categoria
	cat := &model.Categoria{Nome: "Cat"}
	s.categoriaRepo.Criar(s.ctx, cat)

	// Cliente
	u1 := &model.Usuario{Nome: "C", Telefone: "841", RolePermissaoID: 1}
	s.usuarioRepo.Criar(s.ctx, u1)
	c := &model.Cliente{IDUsuario: u1.ID, Nome: u1.Nome, Telefone: u1.Telefone}
	s.clienteRepo.Criar(s.ctx, c)

	// Prestador
	u2 := &model.Usuario{Nome: "P", Telefone: "842", RolePermissaoID: 2}
	s.usuarioRepo.Criar(s.ctx, u2)
	p := &model.Prestador{IDUsuario: u2.ID, Nome: u2.Nome, Telefone: u2.Telefone, Localizacao: "L"}
	s.prestadorRepo.Criar(s.ctx, p)

	// Catalogo
	ct := &model.Catalogo{Nome: "S", IDPrestador: p.IDUsuario, IDCategoria: cat.ID, Localizacao: "L", Latitude: -25.96, Longitude: 32.58}
	s.catalogoRepo.Create(s.ctx, ct)

	return c.IDUsuario, ct.ID
}

func (s *AgendamentoTestSuite) TestAgendamento_Criar() {
	cID, ctID := s.createDependencies()
	a := &model.Agendamento{
		IDCliente:  cID,
		IDCatalogo: ctID,
		DataHora:   time.Now(),
		Status:     "PENDENTE",
	}

	created, err := s.repo.Criar(s.ctx, a)
	assert.NoError(s.T(), err)
	assert.NotZero(s.T(), created.ID)
}

func (s *AgendamentoTestSuite) TestAgendamento_ListarPorPrestadorID() {
	cID, ctID := s.createDependencies()

	// Encontrar o prestador ID a partir do catálogo (ou usar o ID do prestador criado)
	ct, _ := s.catalogoRepo.FindByID(s.ctx, ctID)
	pID := ct.IDPrestador

	a := &model.Agendamento{IDCliente: cID, IDCatalogo: ctID, DataHora: time.Now(), Status: "PENDENTE"}
	s.repo.Criar(s.ctx, a)

	list, err := s.repo.ListarPorPrestadorID(s.ctx, pID, nil, "", "", 10, 0)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), list, 1)
}

func (s *AgendamentoTestSuite) TestAgendamento_FindByLocation() {
	cID, ctID := s.createDependencies()
	a := &model.Agendamento{IDCliente: cID, IDCatalogo: ctID, DataHora: time.Now(), Status: "PENDENTE", Latitude: -25.965, Longitude: 32.585}
	s.repo.Criar(s.ctx, a)

	results, err := s.repo.FindByLocation(s.ctx, cID, -25.96, 32.58, 5.0, nil, "", "", 10, 0)
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), results)
}

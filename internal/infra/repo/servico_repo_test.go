package repo_test

import (
	"context"
	"testing"

	"github.com/ManuelMassora/servicoJa-api/internal/infra/repo"
	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ServicoTestSuite struct {
	suite.Suite
	repo          model.ServicoRepo
	usuarioRepo   model.UsuarioRepo
	clienteRepo   model.ClienteRepo
	prestadorRepo model.PrestadorRepo
	ctx           context.Context
}

func (s *ServicoTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.repo = repo.NewServicoRepository(TestDB)
	s.usuarioRepo = repo.NewUsuarioRepository(TestDB)
	s.clienteRepo = repo.NewClienteRepository(TestDB)
	s.prestadorRepo = repo.NewPrestadorRepository(TestDB)
}

func (s *ServicoTestSuite) SetupTest() {
	TestDB.Exec("TRUNCATE TABLE servicos RESTART IDENTITY CASCADE")
	TestDB.Exec("TRUNCATE TABLE clientes RESTART IDENTITY CASCADE")
	TestDB.Exec("TRUNCATE TABLE prestadors RESTART IDENTITY CASCADE")
	TestDB.Exec("TRUNCATE TABLE usuarios RESTART IDENTITY CASCADE")
}

func TestServicoTestSuite(t *testing.T) {
	suite.Run(t, new(ServicoTestSuite))
}

func (s *ServicoTestSuite) createDependencies() (uint, uint) {
	// Cliente
	u1 := &model.Usuario{Nome: "Cliente", Telefone: "84111", RolePermissaoID: 1}
	s.usuarioRepo.Criar(s.ctx, u1)
	c := &model.Cliente{IDUsuario: u1.ID, Nome: u1.Nome, Telefone: u1.Telefone}
	s.clienteRepo.Criar(s.ctx, c)

	// Prestador
	u2 := &model.Usuario{Nome: "Prestador", Telefone: "84222", RolePermissaoID: 2}
	s.usuarioRepo.Criar(s.ctx, u2)
	p := &model.Prestador{IDUsuario: u2.ID, Nome: u2.Nome, Telefone: u2.Telefone}
	s.prestadorRepo.Criar(s.ctx, p)

	return c.IDUsuario, p.IDUsuario
}

func (s *ServicoTestSuite) TestServico_Criar() {
	cID, pID := s.createDependencies()

	serv := &model.Servico{
		IDCliente:   cID,
		IDPrestador: pID,
		Preco:       500.0,
		Status:      model.StatusPendente,
		Localizacao: "Maputo",
		Latitude:    -25.96,
		Longitude:   32.58,
	}

	created, err := s.repo.Criar(s.ctx, serv)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), created)
	assert.NotZero(s.T(), created.ID)
}

func (s *ServicoTestSuite) TestServico_AtualizarStatus() {
	cID, pID := s.createDependencies()
	serv := &model.Servico{
		IDCliente:   cID,
		IDPrestador: pID,
		Status:      model.StatusPendente,
		Localizacao: "Maputo",
	}
	s.repo.Criar(s.ctx, serv)

	err := s.repo.AtualizarStatus(s.ctx, serv.ID, string(model.StatusEmAndamento))
	assert.NoError(s.T(), err)

	found, _ := s.repo.BuscarPorID(s.ctx, serv.ID)
	assert.Equal(s.T(), model.StatusEmAndamento, found.Status)
}

func (s *ServicoTestSuite) TestServico_ListarPorCliente() {
	cID, pID := s.createDependencies()
	serv := &model.Servico{
		IDCliente:   cID,
		IDPrestador: pID,
		Status:      model.StatusPendente,
		Localizacao: "Maputo",
	}
	s.repo.Criar(s.ctx, serv)

	list, err := s.repo.ListarPorCliente(s.ctx, cID, nil, "created_at", "desc", 10, 0)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), list, 1)
	assert.Equal(s.T(), cID, list[0].IDCliente)
}

func (s *ServicoTestSuite) TestServico_FindByLocation() {
	cID, pID := s.createDependencies()
	serv := &model.Servico{
		IDCliente:   cID,
		IDPrestador: pID,
		Status:      model.StatusPendente,
		Localizacao: "Maputo",
		Latitude:    -25.965,
		Longitude:   32.585,
	}
	s.repo.Criar(s.ctx, serv)

	// Busca serviços vinculados ao usuário (cID) num raio de 5km
	results, err := s.repo.FindByLocation(s.ctx, cID, -25.96, 32.58, 5.0, nil, "", "", 10, 0)
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), results)
}

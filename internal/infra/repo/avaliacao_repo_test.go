package repo_test

import (
	"context"
	"testing"

	"github.com/ManuelMassora/servicoJa-api/internal/infra/repo"
	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AvaliacaoTestSuite struct {
	suite.Suite
	repo          model.AvaliacaoRepo
	usuarioRepo   model.UsuarioRepo
	clienteRepo   model.ClienteRepo
	prestadorRepo model.PrestadorRepo
	servicoRepo   model.ServicoRepo
	ctx           context.Context
}

func (s *AvaliacaoTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.repo = repo.NewAvaliacaoRepository(TestDB)
	s.usuarioRepo = repo.NewUsuarioRepository(TestDB)
	s.clienteRepo = repo.NewClienteRepository(TestDB)
	s.prestadorRepo = repo.NewPrestadorRepository(TestDB)
	s.servicoRepo = repo.NewServicoRepository(TestDB)
}

func (s *AvaliacaoTestSuite) SetupTest() {
	TestDB.Exec("TRUNCATE TABLE avaliacaos RESTART IDENTITY CASCADE")
	TestDB.Exec("TRUNCATE TABLE servicos RESTART IDENTITY CASCADE")
	TestDB.Exec("TRUNCATE TABLE clientes RESTART IDENTITY CASCADE")
	TestDB.Exec("TRUNCATE TABLE prestadors RESTART IDENTITY CASCADE")
	TestDB.Exec("TRUNCATE TABLE usuarios RESTART IDENTITY CASCADE")
}

func TestAvaliacaoTestSuite(t *testing.T) {
	suite.Run(t, new(AvaliacaoTestSuite))
}

func (s *AvaliacaoTestSuite) createDependencies() (uint, uint, uint) {
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

	// Servico
	serv := &model.Servico{
		IDCliente:   c.IDUsuario,
		IDPrestador: p.IDUsuario,
		Preco:       100.0,
		Status:      model.StatusConcluido,
		Localizacao: "Maputo",
	}
	s.servicoRepo.Criar(s.ctx, serv)

	return c.IDUsuario, p.IDUsuario, serv.ID
}

func (s *AvaliacaoTestSuite) TestAvaliacao_Criar() {
	cID, pID, sID := s.createDependencies()

	a := &model.Avaliacao{
		Nota:        5,
		Comentario:  "Excelente serviço!",
		IDCliente:   cID,
		IDPrestador: pID,
		IDServico:   sID,
	}

	err := s.repo.Criar(s.ctx, a)
	assert.NoError(s.T(), err)
	assert.NotZero(s.T(), a.ID)
}

func (s *AvaliacaoTestSuite) TestAvaliacao_MediaPorPrestador() {
	cID, pID, sID1 := s.createDependencies()

	// Nota 5
	a1 := &model.Avaliacao{Nota: 5, IDCliente: cID, IDPrestador: pID, IDServico: sID1}
	s.repo.Criar(s.ctx, a1)

	// Criar outro serviço para outra avaliação nota 3
	serv2 := &model.Servico{IDCliente: cID, IDPrestador: pID, Preco: 50.0, Status: model.StatusConcluido, Localizacao: "Xai-Xai"}
	s.servicoRepo.Criar(s.ctx, serv2)
	a2 := &model.Avaliacao{Nota: 3, IDCliente: cID, IDPrestador: pID, IDServico: serv2.ID}
	s.repo.Criar(s.ctx, a2)

	media, err := s.repo.MediaPorPrestador(s.ctx, pID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 4.0, media)
}

func (s *AvaliacaoTestSuite) TestAvaliacao_ListarPorPrestador() {
	cID, pID, sID := s.createDependencies()
	a := &model.Avaliacao{Nota: 4, IDCliente: cID, IDPrestador: pID, IDServico: sID}
	s.repo.Criar(s.ctx, a)

	list, err := s.repo.ListarPorPrestador(s.ctx, pID, nil, "nota", "desc", 10, 0)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), list, 1)
	assert.Equal(s.T(), pID, list[0].IDPrestador)
}

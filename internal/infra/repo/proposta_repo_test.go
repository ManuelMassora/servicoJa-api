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

type PropostaTestSuite struct {
	suite.Suite
	repo          model.PropostaRepo
	vagaRepo      model.VagaRepo
	usuarioRepo   model.UsuarioRepo
	prestadorRepo model.PrestadorRepo
	clienteRepo   model.ClienteRepo
	ctx           context.Context
}

func (s *PropostaTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.repo = repo.NewPropostaRepository(TestDB)
	s.vagaRepo = repo.NewVagaRepository(TestDB)
	s.usuarioRepo = repo.NewUsuarioRepository(TestDB)
	s.prestadorRepo = repo.NewPrestadorRepository(TestDB)
	s.clienteRepo = repo.NewClienteRepository(TestDB)
}

func (s *PropostaTestSuite) SetupTest() {
	TestDB.Exec("TRUNCATE TABLE propostas RESTART IDENTITY CASCADE")
	TestDB.Exec("TRUNCATE TABLE vagas RESTART IDENTITY CASCADE")
	TestDB.Exec("TRUNCATE TABLE clientes RESTART IDENTITY CASCADE")
	TestDB.Exec("TRUNCATE TABLE prestadors RESTART IDENTITY CASCADE")
	TestDB.Exec("TRUNCATE TABLE usuarios RESTART IDENTITY CASCADE")
}

func TestPropostaTestSuite(t *testing.T) {
	suite.Run(t, new(PropostaTestSuite))
}

func (s *PropostaTestSuite) createDependencies() (uint, uint) {
	// Cliente + Vaga
	u1 := &model.Usuario{Nome: "Cliente", Telefone: "84111", RolePermissaoID: 1}
	s.usuarioRepo.Criar(s.ctx, u1)
	c := &model.Cliente{IDUsuario: u1.ID, Nome: u1.Nome, Telefone: u1.Telefone}
	s.clienteRepo.Criar(s.ctx, c)

	v := &model.Vaga{Titulo: "Vaga Teste", Descricao: "D", IDCliente: c.IDUsuario, Localizacao: "L", Status: model.StatusDisponivel}
	s.vagaRepo.Criar(s.ctx, v)

	// Prestador
	u2 := &model.Usuario{Nome: "Prestador", Telefone: "84222", RolePermissaoID: 2}
	s.usuarioRepo.Criar(s.ctx, u2)
	p := &model.Prestador{IDUsuario: u2.ID, Nome: u2.Nome, Telefone: u2.Telefone}
	s.prestadorRepo.Criar(s.ctx, p)

	return v.ID, p.IDUsuario
}

func (s *PropostaTestSuite) TestProposta_Criar() {
	vID, pID := s.createDependencies()

	p := &model.Proposta{
		IDVaga:        vID,
		IDPrestador:   pID,
		ValorProposto: 1500.0,
		Mensagem:      "Eu faço este trabalho rápido",
		Status:        model.StatusPendente,
	}

	err := s.repo.Criar(s.ctx, p)
	assert.NoError(s.T(), err)
	assert.NotZero(s.T(), p.ID)
}

func (s *PropostaTestSuite) TestProposta_AtualizarStatus() {
	vID, pID := s.createDependencies()
	p := &model.Proposta{IDVaga: vID, IDPrestador: pID, Status: model.StatusPendente}
	s.repo.Criar(s.ctx, p)

	now := time.Now().Round(time.Second)
	err := s.repo.AtualizarStatus(s.ctx, p.ID, model.StatusAceito, now)
	assert.NoError(s.T(), err)

	found, _ := s.repo.BuscarPorID(s.ctx, p.ID)
	assert.Equal(s.T(), model.StatusAceito, found.Status)
	assert.True(s.T(), found.DataResposta.Equal(now))
}

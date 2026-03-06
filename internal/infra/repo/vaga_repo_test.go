package repo_test

import (
	"context"
	"testing"

	"github.com/ManuelMassora/servicoJa-api/internal/infra/repo"
	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type VagaTestSuite struct {
	suite.Suite
	repo        model.VagaRepo
	usuarioRepo model.UsuarioRepo
	clienteRepo model.ClienteRepo
	ctx         context.Context
}

func (s *VagaTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.repo = repo.NewVagaRepository(TestDB)
	s.usuarioRepo = repo.NewUsuarioRepository(TestDB)
	s.clienteRepo = repo.NewClienteRepository(TestDB)
}

func (s *VagaTestSuite) SetupTest() {
	TestDB.Exec("TRUNCATE TABLE vagas RESTART IDENTITY CASCADE")
	TestDB.Exec("TRUNCATE TABLE clientes RESTART IDENTITY CASCADE")
	TestDB.Exec("TRUNCATE TABLE usuarios RESTART IDENTITY CASCADE")
}

func TestVagaTestSuite(t *testing.T) {
	suite.Run(t, new(VagaTestSuite))
}

func (s *VagaTestSuite) createTestCliente() uint {
	u := &model.Usuario{Nome: "Cliente", Telefone: "84111", RolePermissaoID: 1}
	s.usuarioRepo.Criar(s.ctx, u)
	c := &model.Cliente{IDUsuario: u.ID, Nome: u.Nome, Telefone: u.Telefone}
	s.clienteRepo.Criar(s.ctx, c)
	return c.IDUsuario
}

func (s *VagaTestSuite) TestVaga_Criar() {
	cID := s.createTestCliente()
	v := &model.Vaga{
		Titulo:      "Preciso de Pedreiro",
		Descricao:   "Muro caiu",
		Localizacao: "Maputo",
		Preco:       2000.0,
		Status:      model.StatusDisponivel,
		IDCliente:   cID,
	}

	created, err := s.repo.Criar(s.ctx, v)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), created)
	assert.NotZero(s.T(), created.ID)
}

func (s *VagaTestSuite) TestVaga_ListarDisponiveis() {
	cID := s.createTestCliente()
	v1 := &model.Vaga{Titulo: "Vaga 1", Status: model.StatusDisponivel, IDCliente: cID, Localizacao: "A"}
	v2 := &model.Vaga{Titulo: "Vaga 2", Status: model.StatusOcupada, IDCliente: cID, Localizacao: "B"}
	s.repo.Criar(s.ctx, v1)
	s.repo.Criar(s.ctx, v2)

	list, err := s.repo.ListarDisponiveis(s.ctx, nil, "", "", 10, 0)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), list, 1)
	assert.Equal(s.T(), "Vaga 1", list[0].Titulo)
}

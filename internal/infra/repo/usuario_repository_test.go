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

// UsuarioRepoTestSuite usa TestDB definido no main_test.go — sem container próprio.
type UsuarioRepoTestSuite struct {
	suite.Suite
	repo          model.UsuarioRepo
	clienteRepo   model.ClienteRepo
	prestadorRepo model.PrestadorRepo
	ctx           context.Context
}

func (s *UsuarioRepoTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.repo = repo.NewUsuarioRepository(TestDB)
	s.clienteRepo = repo.NewClienteRepository(TestDB)
	s.prestadorRepo = repo.NewPrestadorRepository(TestDB)
}

// SetupTest limpa a tabela antes de cada teste para garantir isolamento.
func (s *UsuarioRepoTestSuite) SetupTest() {
	TestDB.Exec("TRUNCATE TABLE usuarios RESTART IDENTITY CASCADE")
}

func TestUsuarioRepoTestSuite(t *testing.T) {
	suite.Run(t, new(UsuarioRepoTestSuite))
}

func (s *UsuarioRepoTestSuite) TestUsuario_Criar() {
	u := &model.Usuario{
		Nome:            "Test User",
		Telefone:        "840000001",
		Senha:           "password",
		RolePermissaoID: 1,
	}

	err := s.repo.Criar(s.ctx, u)
	assert.NoError(s.T(), err)
	assert.NotZero(s.T(), u.ID)
}

func (s *UsuarioRepoTestSuite) TestUsuario_BuscarPorID() {
	u := &model.Usuario{
		Nome:            "Find Me",
		Telefone:        "840000002",
		Senha:           "password",
		RolePermissaoID: 1,
	}
	s.repo.Criar(s.ctx, u)

	found, err := s.repo.BuscarPorID(s.ctx, u.ID)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), found)
	assert.Equal(s.T(), u.Nome, found.Nome)
}

func (s *UsuarioRepoTestSuite) TestUsuario_BuscarPorTelefone() {
	u := &model.Usuario{
		Nome:            "Phone User",
		Telefone:        "840000003",
		Senha:           "password",
		RolePermissaoID: 1,
	}
	s.repo.Criar(s.ctx, u)

	found, err := s.repo.BuscarPorTelefone(s.ctx, u.Telefone)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), found)
	assert.Equal(s.T(), u.ID, found.ID)
}

func (s *UsuarioRepoTestSuite) TestUsuario_Atualizar() {
	u := &model.Usuario{
		Nome:            "Old Name",
		Telefone:        "840000004",
		Senha:           "password",
		RolePermissaoID: 1,
	}
	s.repo.Criar(s.ctx, u)

	u.Nome = "New Name"
	err := s.repo.Atualizar(s.ctx, u)
	assert.NoError(s.T(), err)

	found, _ := s.repo.BuscarPorID(s.ctx, u.ID)
	assert.Equal(s.T(), "New Name", found.Nome)
}

func (s *UsuarioRepoTestSuite) TestUsuario_Remover() {
	u := &model.Usuario{
		Nome:            "Delete Me",
		Telefone:        "840000005",
		Senha:           "password",
		RolePermissaoID: 1,
	}
	s.repo.Criar(s.ctx, u)

	err := s.repo.Remover(s.ctx, u.ID)
	assert.NoError(s.T(), err)

	found, err := s.repo.BuscarPorID(s.ctx, u.ID)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), found)
}

func (s *UsuarioRepoTestSuite) TestUsuario_Notificacoes() {
	u := &model.Usuario{
		Nome:            "Notify User",
		Telefone:        "840000006",
		Senha:           "password",
		RolePermissaoID: 1,
	}
	s.repo.Criar(s.ctx, u)

	err := s.repo.IncrementarNotificacoesNovas(s.ctx, u.ID)
	assert.NoError(s.T(), err)

	found, _ := s.repo.BuscarPorID(s.ctx, u.ID)
	assert.Equal(s.T(), uint(1), found.NotificacoesNovas)

	err = s.repo.ZerarNotificacoesNovas(s.ctx, u.ID)
	assert.NoError(s.T(), err)

	found, _ = s.repo.BuscarPorID(s.ctx, u.ID)
	assert.Equal(s.T(), uint(0), found.NotificacoesNovas)
}

func (s *UsuarioRepoTestSuite) TestCliente_Criar() {
	u := &model.Usuario{
		Nome:            "Cliente User",
		Telefone:        "840000007",
		Senha:           "password",
		RolePermissaoID: 1,
	}
	s.repo.Criar(s.ctx, u)

	c := &model.Cliente{
		IDUsuario: u.ID,
		Nome:      u.Nome,
		Telefone:  u.Telefone,
		ImagemURL: "http://image.com",
	}

	err := s.clienteRepo.Criar(s.ctx, c)
	assert.NoError(s.T(), err)
}

func (s *UsuarioRepoTestSuite) TestPrestador_Criar() {
	u := &model.Usuario{
		Nome:            "Prestador User",
		Telefone:        "840000008",
		Senha:           "password",
		RolePermissaoID: 2,
	}
	s.repo.Criar(s.ctx, u)

	p := &model.Prestador{
		IDUsuario:        u.ID,
		Nome:             u.Nome,
		Telefone:         u.Telefone,
		Localizacao:      "Maputo",
		StatusDisponivel: true,
	}

	err := s.prestadorRepo.Criar(s.ctx, p)
	assert.NoError(s.T(), err)
}

func (s *UsuarioRepoTestSuite) TestUsuario_Suspender() {
	u := &model.Usuario{
		Nome:            "Suspend User",
		Telefone:        "840000009",
		Senha:           "password",
		RolePermissaoID: 1,
	}
	s.repo.Criar(s.ctx, u)

	until := time.Now().Add(24 * time.Hour).Round(time.Second)
	err := s.repo.SuspenderUsuario(s.ctx, u.ID, until)
	assert.NoError(s.T(), err)

	found, _ := s.repo.BuscarPorID(s.ctx, u.ID)
	assert.NotNil(s.T(), found.SuspensoAte)
	assert.True(s.T(), found.SuspensoAte.After(time.Now()))
}

func (s *UsuarioRepoTestSuite) TestUsuario_IncrementarCancelamentos() {
	u := &model.Usuario{
		Nome:            "Cancel User",
		Telefone:        "849999999",
		Senha:           "password",
		RolePermissaoID: 1,
	}
	s.repo.Criar(s.ctx, u)

	count, err := s.repo.IncrementarCancelamentos(s.ctx, u.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), uint(1), count)

	found, _ := s.repo.BuscarPorID(s.ctx, u.ID)
	assert.Equal(s.T(), uint(1), found.CancelamentosCount)
}

func (s *UsuarioRepoTestSuite) TestUsuario_ListarTodos() {
	u1 := &model.Usuario{Nome: "Alice", Telefone: "84111", RolePermissaoID: 1}
	u2 := &model.Usuario{Nome: "Bob", Telefone: "84222", RolePermissaoID: 1}
	s.repo.Criar(s.ctx, u1)
	s.repo.Criar(s.ctx, u2)

	filters := map[string]interface{}{"nome": "Alice"}
	list, err := s.repo.ListarTodos(s.ctx, filters, "nome", "asc", 10, 0)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), list, 1)
	assert.Equal(s.T(), "Alice", list[0].Nome)
}

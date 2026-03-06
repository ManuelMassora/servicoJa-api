package repo_test

import (
	"context"
	"testing"

	"github.com/ManuelMassora/servicoJa-api/internal/infra/repo"
	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type NotificacaoTestSuite struct {
	suite.Suite
	repo        model.NotificacaoRepo
	usuarioRepo model.UsuarioRepo
	ctx         context.Context
}

func (s *NotificacaoTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.repo = repo.NewNotificacaoRepo(TestDB)
	s.usuarioRepo = repo.NewUsuarioRepository(TestDB)
}

func (s *NotificacaoTestSuite) SetupTest() {
	TestDB.Exec("TRUNCATE TABLE notificacaos RESTART IDENTITY CASCADE")
	TestDB.Exec("TRUNCATE TABLE usuarios RESTART IDENTITY CASCADE")
}

func TestNotificacaoTestSuite(t *testing.T) {
	suite.Run(t, new(NotificacaoTestSuite))
}

func (s *NotificacaoTestSuite) createTestUser() uint {
	u := &model.Usuario{Nome: "User", Telefone: "84111", RolePermissaoID: 1}
	s.usuarioRepo.Criar(s.ctx, u)
	return u.ID
}

func (s *NotificacaoTestSuite) TestNotificacao_Enviar() {
	uID := s.createTestUser()
	n := &model.Notificacao{
		IDUsuario: uID,
		Titulo:    "Alerta",
		Mensagem:  "Sua proposta foi aceita!",
		Lida:      false,
	}

	err := s.repo.Enviar(s.ctx, n)
	assert.NoError(s.T(), err)
	assert.NotZero(s.T(), n.ID)
}

func (s *NotificacaoTestSuite) TestNotificacao_MarcarComoLida() {
	uID := s.createTestUser()
	n := &model.Notificacao{IDUsuario: uID, Titulo: "T", Mensagem: "M", Lida: false}
	s.repo.Enviar(s.ctx, n)

	err := s.repo.MarcarComoLida(s.ctx, n.ID)
	assert.NoError(s.T(), err)

	found, _ := s.repo.BuscarPorID(s.ctx, n.ID)
	assert.True(s.T(), found.Lida)
}

func (s *NotificacaoTestSuite) TestNotificacao_ListarPorUsuario() {
	uID := s.createTestUser()
	n1 := &model.Notificacao{IDUsuario: uID, Titulo: "T1", Mensagem: "M1"}
	n2 := &model.Notificacao{IDUsuario: uID, Titulo: "T2", Mensagem: "M2"}
	s.repo.Enviar(s.ctx, n1)
	s.repo.Enviar(s.ctx, n2)

	list, err := s.repo.ListarPorUsuario(s.ctx, uID, nil, "", "", 10, 0)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), list, 2)
}

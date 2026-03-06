package repo_test

import (
	"context"
	"testing"

	"github.com/ManuelMassora/servicoJa-api/internal/infra/repo"
	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PagamentoTestSuite struct {
	suite.Suite
	repo        model.PagamentoRepo
	usuarioRepo model.UsuarioRepo
	ctx         context.Context
}

func (s *PagamentoTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.repo = repo.NewPagamentoRepository(TestDB)
	s.usuarioRepo = repo.NewUsuarioRepository(TestDB)
}

func (s *PagamentoTestSuite) SetupTest() {
	TestDB.Exec("TRUNCATE TABLE pagamentos RESTART IDENTITY CASCADE")
	TestDB.Exec("TRUNCATE TABLE usuarios RESTART IDENTITY CASCADE")
}

func TestPagamentoTestSuite(t *testing.T) {
	suite.Run(t, new(PagamentoTestSuite))
}

func (s *PagamentoTestSuite) createTestUsers() (uint, uint) {
	u1 := &model.Usuario{Nome: "Cliente", Telefone: "84111", RolePermissaoID: 1}
	s.usuarioRepo.Criar(s.ctx, u1)
	u2 := &model.Usuario{Nome: "Prestador", Telefone: "84222", RolePermissaoID: 2}
	s.usuarioRepo.Criar(s.ctx, u2)
	return u1.ID, u2.ID
}

func (s *PagamentoTestSuite) TestPagamento_Criar() {
	cID, pID := s.createTestUsers()
	pIDPtr := &pID

	p := &model.Pagamento{
		IDCliente:   cID,
		IDPrestador: pIDPtr,
		Valor:       1000.0,
		Status:      model.StatusPendente,
		Referencia:  "REF123",
	}

	err := s.repo.Criar(s.ctx, p)
	assert.NoError(s.T(), err)
	assert.NotZero(s.T(), p.ID)
}

func (s *PagamentoTestSuite) TestPagamento_BuscarPorReferencia() {
	cID, _ := s.createTestUsers()
	p := &model.Pagamento{
		IDCliente:  cID,
		Valor:      500.0,
		Status:     model.StatusPendente,
		Referencia: "REF456",
	}
	s.repo.Criar(s.ctx, p)

	found, err := s.repo.BuscarPorReferencia(s.ctx, "REF456")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), found)
	assert.Equal(s.T(), p.ID, found.ID)
}

func (s *PagamentoTestSuite) TestPagamento_AtualizarStatus() {
	cID, _ := s.createTestUsers()
	p := &model.Pagamento{
		IDCliente:  cID,
		Valor:      500.0,
		Status:     model.StatusPendente,
		Referencia: "REF789",
	}
	s.repo.Criar(s.ctx, p)

	err := s.repo.AtualizarStatus(s.ctx, p.ID, model.StatusConfirmado)
	assert.NoError(s.T(), err)

	found, _ := s.repo.BuscarPorID(s.ctx, p.ID)
	assert.Equal(s.T(), model.StatusConfirmado, found.Status)
}

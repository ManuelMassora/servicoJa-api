package repo_test

import (
	"context"
	"testing"

	"github.com/ManuelMassora/servicoJa-api/internal/infra/repo"
	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AnexoImagemTestSuite struct {
	suite.Suite
	repo        model.AnexoImagemRepo
	vagaRepo    model.VagaRepo
	usuarioRepo model.UsuarioRepo
	clienteRepo model.ClienteRepo
	ctx         context.Context
}

func (s *AnexoImagemTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.repo = repo.NewAnexoImagemRepo(TestDB)
	s.vagaRepo = repo.NewVagaRepository(TestDB)
	s.usuarioRepo = repo.NewUsuarioRepository(TestDB)
	s.clienteRepo = repo.NewClienteRepository(TestDB)
}

func (s *AnexoImagemTestSuite) SetupTest() {
	TestDB.Exec("TRUNCATE TABLE anexo_imagems RESTART IDENTITY CASCADE")
	TestDB.Exec("TRUNCATE TABLE vagas RESTART IDENTITY CASCADE")
	TestDB.Exec("TRUNCATE TABLE clientes RESTART IDENTITY CASCADE")
	TestDB.Exec("TRUNCATE TABLE usuarios RESTART IDENTITY CASCADE")
}

func TestAnexoImagemTestSuite(t *testing.T) {
	suite.Run(t, new(AnexoImagemTestSuite))
}

func (s *AnexoImagemTestSuite) createVaga() uint {
	u := &model.Usuario{Nome: "C", Telefone: "84111", RolePermissaoID: 1}
	s.usuarioRepo.Criar(s.ctx, u)
	c := &model.Cliente{IDUsuario: u.ID, Nome: u.Nome, Telefone: u.Telefone}
	s.clienteRepo.Criar(s.ctx, c)

	v := &model.Vaga{Titulo: "V", Descricao: "D", IDCliente: c.IDUsuario, Localizacao: "L", Status: model.StatusDisponivel}
	s.vagaRepo.Criar(s.ctx, v)
	return v.ID
}

func (s *AnexoImagemTestSuite) TestAnexoImagem_Create() {
	vID := s.createVaga()
	vIDPtr := &vID

	anexo := &model.AnexoImagem{
		URL:    "http://image.com/vaga.png",
		VagaID: vIDPtr,
	}

	err := s.repo.Create(s.ctx, anexo)
	assert.NoError(s.T(), err)
	assert.NotZero(s.T(), anexo.ID)
}

func (s *AnexoImagemTestSuite) TestAnexoImagem_FindByVagaID() {
	vID := s.createVaga()
	vIDPtr := &vID
	anexo := &model.AnexoImagem{URL: "http://image.com/a.png", VagaID: vIDPtr}
	s.repo.Create(s.ctx, anexo)

	list, err := s.repo.FindByVagaID(s.ctx, vID)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), list, 1)
	assert.Equal(s.T(), "http://image.com/a.png", list[0].URL)
}

package repo_test

import (
	"context"
	"testing"

	"github.com/ManuelMassora/servicoJa-api/internal/infra/repo"
	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CategoriaTestSuite struct {
	suite.Suite
	repo model.CategoriaRepo
	ctx  context.Context
}

func (s *CategoriaTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.repo = repo.NewCategoriaRepository(TestDB)
}

func (s *CategoriaTestSuite) SetupTest() {
	TestDB.Exec("TRUNCATE TABLE categoria RESTART IDENTITY CASCADE")
}

func TestCategoriaTestSuite(t *testing.T) {
	suite.Run(t, new(CategoriaTestSuite))
}

func (s *CategoriaTestSuite) TestCategoria_Criar() {
	c := &model.Categoria{
		Nome: "Limpeza",
	}

	err := s.repo.Criar(s.ctx, c)
	assert.NoError(s.T(), err)
	assert.NotZero(s.T(), c.ID)
}

func (s *CategoriaTestSuite) TestCategoria_Editar() {
	c := &model.Categoria{
		Nome: "Limpeza",
	}
	s.repo.Criar(s.ctx, c)

	campos := map[string]interface{}{
		"nome": "Limpeza Pesada",
	}

	err := s.repo.Editar(s.ctx, c.ID, campos)
	assert.NoError(s.T(), err)

	found, _ := s.repo.BuscarPorID(s.ctx, c.ID)
	assert.Equal(s.T(), "Limpeza Pesada", found.Nome)
}

func (s *CategoriaTestSuite) TestCategoria_Listar() {
	c1 := &model.Categoria{Nome: "Limpeza"}
	c2 := &model.Categoria{Nome: "Eletricista"}
	s.repo.Criar(s.ctx, c1)
	s.repo.Criar(s.ctx, c2)

	filters := map[string]interface{}{"nome": "Limp"}
	list, err := s.repo.Listar(s.ctx, filters, "nome", "asc", 10, 0)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), list, 1)
	assert.Equal(s.T(), "Limpeza", list[0].Nome)
}

func (s *CategoriaTestSuite) TestCategoria_BuscarPorID() {
	c := &model.Categoria{Nome: "Encanador"}
	s.repo.Criar(s.ctx, c)

	found, err := s.repo.BuscarPorID(s.ctx, c.ID)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), found)
	assert.Equal(s.T(), c.Nome, found.Nome)
}

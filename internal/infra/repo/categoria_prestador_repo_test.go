package repo_test

import (
	"context"
	"testing"

	"github.com/ManuelMassora/servicoJa-api/internal/infra/repo"
	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CategoriaPrestadorTestSuite struct {
	suite.Suite
	repo model.CategoriaPrestadorRepo
	ctx  context.Context
}

func (s *CategoriaPrestadorTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.repo = repo.NewCategoriaPrestadorRepo(TestDB)
}

func (s *CategoriaPrestadorTestSuite) SetupTest() {
	TestDB.Exec("TRUNCATE TABLE categoria_prestadors RESTART IDENTITY CASCADE")
}

func TestCategoriaPrestadorTestSuite(t *testing.T) {
	suite.Run(t, new(CategoriaPrestadorTestSuite))
}

func (s *CategoriaPrestadorTestSuite) TestCategoriaPrestador_Criar() {
	c := &model.CategoriaPrestador{
		Nome:  "Limpeza de Vidros",
		Icone: "window",
	}

	created, err := s.repo.Criar(s.ctx, c)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), created)
	assert.NotZero(s.T(), created.ID)
	assert.Equal(s.T(), "Limpeza de Vidros", created.Nome)
}

func (s *CategoriaPrestadorTestSuite) TestCategoriaPrestador_Editar() {
	c := &model.CategoriaPrestador{
		Nome: "Limpeza de Vidros",
	}
	s.repo.Criar(s.ctx, c)

	campos := map[string]interface{}{
		"nome":  "Limpeza Profissional",
		"icone": "star",
	}

	updated, err := s.repo.Editar(s.ctx, c.ID, campos)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), updated)
	assert.Equal(s.T(), "Limpeza Profissional", updated.Nome)
	assert.Equal(s.T(), "star", updated.Icone)
}

func (s *CategoriaPrestadorTestSuite) TestCategoriaPrestador_Listar() {
	c1 := &model.CategoriaPrestador{Nome: "Limpeza"}
	c2 := &model.CategoriaPrestador{Nome: "Eletricidade"}
	s.repo.Criar(s.ctx, c1)
	s.repo.Criar(s.ctx, c2)

	filters := map[string]interface{}{"nome": "Limp"}
	list, err := s.repo.Listar(s.ctx, filters, "nome", "asc", 10, 0)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), list, 1)
	assert.Equal(s.T(), "Limpeza", list[0].Nome)
}

func (s *CategoriaPrestadorTestSuite) TestCategoriaPrestador_BuscarPorID() {
	c := &model.CategoriaPrestador{Nome: "Canalizador"}
	s.repo.Criar(s.ctx, c)

	found, err := s.repo.BuscarPorID(s.ctx, c.ID)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), found)
	assert.Equal(s.T(), c.Nome, found.Nome)
}

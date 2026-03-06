package usecases_test

import (
	"context"
	"testing"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"github.com/ManuelMassora/servicoJa-api/internal/usecases"
	"github.com/ManuelMassora/servicoJa-api/internal/usecases/usecases_test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCategoriaUseCase_Criar_Success(t *testing.T) {
	repo := new(mocks.MockCategoriaRepo)
	uc := usecases.NewCategoriaUseCase(repo)

	ctx := context.Background()
	req := usecases.CategoriaRequest{Nome: "Limpeza", Descricao: "Limpeza de casas"}

	repo.On("Criar", ctx, &model.Categoria{Nome: "Limpeza", Descricao: "Limpeza de casas"}).Return(nil)

	id, err := uc.Criar(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, uint(0), id)
}

func TestCategoriaUseCase_Listar_Success(t *testing.T) {
	repo := new(mocks.MockCategoriaRepo)
	uc := usecases.NewCategoriaUseCase(repo)

	ctx := context.Background()
	categorias := []model.Categoria{{BaseModel: model.BaseModel{ID: 1}, Nome: "Cat 1"}}

	repo.On("Listar", ctx, mock.Anything, "", "", 10, 0).Return(categorias, nil)

	res, err := uc.Listar(ctx, nil, "", "", 10, 0)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "Cat 1", res[0].Nome)
}

func TestCategoriaUseCase_BuscarPorID_Success(t *testing.T) {
	repo := new(mocks.MockCategoriaRepo)
	uc := usecases.NewCategoriaUseCase(repo)

	ctx := context.Background()
	id := uint(1)
	repo.On("BuscarPorID", ctx, id).Return(&model.Categoria{BaseModel: model.BaseModel{ID: id}, Nome: "Cat 1"}, nil)

	res, err := uc.BuscarPorID(ctx, id)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "Cat 1", res.Nome)
}

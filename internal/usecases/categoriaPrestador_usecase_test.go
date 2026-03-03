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

func TestCategoriaPrestadorUsecase_Criar_Success(t *testing.T) {
	repo := new(mocks.MockCategoriaPrestadorRepo)
	uc := usecases.NewCategoriaPrestadorUsecase(repo)

	ctx := context.Background()
	req := usecases.CategoriaPrestadorRequest{Nome: "Encanador", Descricao: "Desc", Icone: "icon.png"}

	repo.On("Criar", ctx, mock.Anything).Return(&model.CategoriaPrestador{BaseModel: model.BaseModel{ID: 1}, Nome: "Encanador"}, nil)

	res, err := uc.Criar(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "Encanador", res.Nome)
}

func TestCategoriaPrestadorUsecase_Listar_Success(t *testing.T) {
	repo := new(mocks.MockCategoriaPrestadorRepo)
	uc := usecases.NewCategoriaPrestadorUsecase(repo)

	ctx := context.Background()
	categorias := []model.CategoriaPrestador{{BaseModel: model.BaseModel{ID: 1}, Nome: "Cat"}}

	repo.On("Listar", ctx, mock.Anything, "", "", 10, 0).Return(categorias, nil)

	res, err := uc.Listar(ctx, nil, "", "", 10, 0)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
}

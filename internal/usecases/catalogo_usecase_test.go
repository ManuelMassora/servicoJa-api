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

func TestCatalogoUseCase_Criar_Success(t *testing.T) {
	repo := new(mocks.MockCatalogoRepo)
	anexoRepo := new(mocks.MockAnexoImagemRepo)
	uc := usecases.NewCatalogoUC(repo, anexoRepo)

	ctx := context.Background()
	req := usecases.RequestCreateCatalogo{
		Nome:      "Servico A",
		Descricao: "Descricao A",
		TipoPreco: "fixo",
		ValorFixo: 500,
		Anexos:    []string{"url1.jpg"},
	}

	repo.On("Create", ctx, mock.Anything).Return(nil)
	anexoRepo.On("Create", ctx, mock.Anything).Return(nil)

	idCatalogo, err := uc.Criar(ctx, req, 1)
	assert.NoError(t, err)
	assert.Equal(t, uint(0), idCatalogo) // 0 because mock repo doesn't set ID
}

func TestCatalogoUseCase_Listar_Success(t *testing.T) {
	repo := new(mocks.MockCatalogoRepo)
	anexoRepo := new(mocks.MockAnexoImagemRepo)
	uc := usecases.NewCatalogoUC(repo, anexoRepo)

	ctx := context.Background()
	catalogos := []model.Catalogo{
		{
			BaseModel: model.BaseModel{ID: 1},
			Nome:      "Cat 1",
			Prestador: model.Prestador{Usuario: model.Usuario{Nome: "P1"}},
		},
	}

	repo.On("FindAll", ctx, mock.Anything, "", "", 10, 0).Return([]*model.Catalogo{&catalogos[0]}, nil)
	anexoRepo.On("FindByCatalogoIDs", ctx, []uint{1}).Return(nil, nil)

	res, err := uc.Listar(ctx, nil, "", "", 10, 0)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "Cat 1", res[0].Nome)
}

func TestCatalogoUseCase_Apagar_Success(t *testing.T) {
	repo := new(mocks.MockCatalogoRepo)
	uc := usecases.NewCatalogoUC(repo, nil)

	ctx := context.Background()
	id := uint(1)
	idPrestador := uint(10)

	catalogo := &model.Catalogo{BaseModel: model.BaseModel{ID: id}, IDPrestador: idPrestador}

	repo.On("FindByID", ctx, id).Return(catalogo, nil)
	repo.On("Delete", ctx, id).Return(nil)

	err := uc.Apagar(ctx, id, idPrestador)

	assert.NoError(t, err)
}

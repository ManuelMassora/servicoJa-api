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

func TestVagaUseCase_CriarVaga_Success(t *testing.T) {
	vRepo := new(mocks.MockVagaRepo)
	aRepo := new(mocks.MockAnexoImagemRepo)
	uRepo := new(mocks.MockUsuarioRepo)
	pRepo := new(mocks.MockPagamentoRepo)
	pUC := new(mocks.MockPagamentoUseCase)

	uc := usecases.NewVagaUseCase(vRepo, aRepo, uRepo, pRepo, pUC)

	ctx := context.Background()
	idUsuario := uint(1)
	req := usecases.VagaRequest{
		Titulo:            "Pintura",
		Descricao:         "Pintar muro",
		Preco:             1000,
		Localizacao:       "Maputo",
		Anexos:            []string{"muro.jpg"},
		TelefonePagamento: "841112233",
	}

	uRepo.On("BuscarPorID", ctx, idUsuario).Return(&model.Usuario{BaseModel: model.BaseModel{ID: idUsuario}}, nil)
	vRepo.On("Criar", ctx, mock.Anything).Return(&model.Vaga{BaseModel: model.BaseModel{ID: 10}, Preco: 1000}, nil)
	pRepo.On("Criar", ctx, mock.Anything).Return(nil)
	aRepo.On("Create", ctx, mock.Anything).Return(nil)

	pUC.On("IniciarPagamentoC2B", ctx, mock.Anything, req.TelefonePagamento).Return(nil)

	err := uc.CriarVaga(ctx, req, idUsuario)

	assert.NoError(t, err)
}

func TestVagaUseCase_ListarVagasDisponiveis_Success(t *testing.T) {
	vRepo := new(mocks.MockVagaRepo)
	aRepo := new(mocks.MockAnexoImagemRepo)
	pUC := new(mocks.MockPagamentoUseCase)
	uc := usecases.NewVagaUseCase(vRepo, aRepo, nil, nil, pUC)

	ctx := context.Background()
	vagas := []model.Vaga{
		{
			BaseModel: model.BaseModel{ID: 1},
			Titulo:    "Vaga 1",
			Cliente:   &model.Cliente{Usuario: model.Usuario{Nome: "C1"}},
		},
	}

	vRepo.On("ListarDisponiveis", ctx, mock.Anything, "", "", 10, 0).Return(vagas, nil)
	aRepo.On("FindByVagaIDs", ctx, []uint{1}).Return(nil, nil)

	res, err := uc.ListarVagasDisponiveis(ctx, nil, "", "", 10, 0)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
}

func TestVagaUseCase_CancelarVaga_Success(t *testing.T) {
	vRepo := new(mocks.MockVagaRepo)
	pUC := new(mocks.MockPagamentoUseCase)
	uc := usecases.NewVagaUseCase(vRepo, nil, nil, nil, pUC)

	ctx := context.Background()
	id := uint(1)
	idUsuario := uint(10)

	vaga := &model.Vaga{BaseModel: model.BaseModel{ID: id}, IDCliente: idUsuario}

	vRepo.On("BuscarPorID", ctx, id).Return(vaga, nil)
	vRepo.On("Salvar", ctx, mock.Anything).Return(nil)

	err := uc.CancelarVaga(ctx, id, idUsuario)

	assert.NoError(t, err)
}

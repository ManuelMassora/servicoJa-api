package usecases_test

import (
	"context"
	"testing"
	"time"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"github.com/ManuelMassora/servicoJa-api/internal/usecases"
	"github.com/ManuelMassora/servicoJa-api/internal/usecases/usecases_test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServicoUseCase_FinalizarServico_Success(t *testing.T) {
	sRepo := new(mocks.MockServicoRepo)
	agRepo := new(mocks.MockAgendamentoRepo)
	vRepo := new(mocks.MockVagaRepo)
	nRepo := new(mocks.MockNotificacaoRepo)
	uRepo := new(mocks.MockUsuarioRepo)
	pUC := new(mocks.MockPagamentoUseCase)

	uc := usecases.NewServicoUseCase(sRepo, agRepo, vRepo, nRepo, uRepo, pUC)

	ctx := context.Background()
	idServico := uint(1)
	idUsuario := uint(2) // Prestador

	servico := &model.Servico{
		BaseModel:   model.BaseModel{ID: idServico},
		IDPrestador: idUsuario,
		IDCliente:   1,
		Status:      model.StatusEmAndamento,
	}

	sRepo.On("BuscarPorID", ctx, idServico).Return(servico, nil)
	sRepo.On("Atualizar", ctx, mock.Anything).Return(nil)
	nRepo.On("Enviar", ctx, mock.Anything).Return(nil)
	uRepo.On("IncrementarNotificacoesNovas", ctx, uint(1)).Return(nil)

	err := uc.FinalizarServico(ctx, idServico, idUsuario)

	assert.NoError(t, err)
}

func TestServicoUseCase_ConfirmarServico_Success(t *testing.T) {
	sRepo := new(mocks.MockServicoRepo)
	agRepo := new(mocks.MockAgendamentoRepo)
	vRepo := new(mocks.MockVagaRepo)
	nRepo := new(mocks.MockNotificacaoRepo)
	uRepo := new(mocks.MockUsuarioRepo)
	pUC := new(mocks.MockPagamentoUseCase)

	uc := usecases.NewServicoUseCase(sRepo, agRepo, vRepo, nRepo, uRepo, pUC)

	ctx := context.Background()
	idServico := uint(1)
	idUsuario := uint(1) // Cliente

	servico := &model.Servico{
		BaseModel:   model.BaseModel{ID: idServico},
		IDCliente:   idUsuario,
		IDPrestador: 2,
		Status:      model.StatusConcluido,
	}

	sRepo.On("BuscarPorID", ctx, idServico).Return(servico, nil)
	sRepo.On("Atualizar", ctx, mock.Anything).Return(nil)
	nRepo.On("Enviar", ctx, mock.Anything).Return(nil)
	uRepo.On("IncrementarNotificacoesNovas", ctx, uint(2)).Return(nil)

	pUC.On("ProcessarPagamentoPrestador", ctx, idServico).Return(nil)

	err := uc.ConfirmarServico(ctx, idServico, idUsuario)

	assert.NoError(t, err)
}

func TestServicoUseCase_ListarPorCliente_Success(t *testing.T) {
	sRepo := new(mocks.MockServicoRepo)
	pUC := new(mocks.MockPagamentoUseCase)
	uc := usecases.NewServicoUseCase(sRepo, nil, nil, nil, nil, pUC)

	ctx := context.Background()
	idUsuario := uint(1)
	servicos := []model.Servico{
		{
			BaseModel:   model.BaseModel{ID: 1, CreatedAt: time.Now()},
			Localizacao: "Maputo",
		},
	}

	sRepo.On("ListarPorCliente", ctx, idUsuario, mock.Anything, "", "", 10, 0).Return(servicos, nil)

	res, err := uc.ListarPorCliente(ctx, idUsuario, nil, "", "", 10, 0)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
}

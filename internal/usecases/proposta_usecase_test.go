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

func TestPropostaUseCase_Criar_Success(t *testing.T) {
	pRepo := new(mocks.MockPropostaRepo)
	vRepo := new(mocks.MockVagaRepo)
	nRepo := new(mocks.MockNotificacaoRepo)
	uRepo := new(mocks.MockUsuarioRepo)
	uc := usecases.NewPropostaUseCase(pRepo, vRepo, nil, nRepo, uRepo, nil)

	ctx := context.Background()
	req := usecases.PropostaRequest{IDVaga: 1, ValorProposto: 100}

	vRepo.On("BuscarPorID", ctx, uint(1)).Return(&model.Vaga{BaseModel: model.BaseModel{ID: 1}, IDCliente: 10}, nil)
	nRepo.On("Enviar", ctx, mock.Anything).Return(nil)
	uRepo.On("IncrementarNotificacoesNovas", ctx, uint(10)).Return(nil)
	vRepo.On("IncrementarPropostasNovas", ctx, uint(1)).Return(nil)
	pRepo.On("Salvar", ctx, mock.Anything).Return(nil)

	err := uc.Criar(ctx, req, 2)

	assert.NoError(t, err)
}

func TestPropostaUseCase_Responder_Aceitar(t *testing.T) {
	pRepo := new(mocks.MockPropostaRepo)
	vRepo := new(mocks.MockVagaRepo)
	sRepo := new(mocks.MockServicoRepo)
	nRepo := new(mocks.MockNotificacaoRepo)
	uRepo := new(mocks.MockUsuarioRepo)
	payRepo := new(mocks.MockPagamentoRepo)

	uc := usecases.NewPropostaUseCase(pRepo, vRepo, sRepo, nRepo, uRepo, payRepo)

	ctx := context.Background()
	idProposta := uint(1)
	idUsuario := uint(10)

	proposta := &model.Proposta{BaseModel: model.BaseModel{ID: 1}, IDVaga: 5, IDPrestador: 2, Status: model.StatusPendente}
	vaga := &model.Vaga{BaseModel: model.BaseModel{ID: 5}, IDCliente: idUsuario, Titulo: "Vaga Teste"}

	pRepo.On("BuscarPorID", ctx, idProposta).Return(proposta, nil)
	vRepo.On("BuscarPorID", ctx, uint(5)).Return(vaga, nil)
	nRepo.On("Enviar", ctx, mock.Anything).Return(nil)
	uRepo.On("IncrementarNotificacoesNovas", ctx, uint(2)).Return(nil)
	vRepo.On("Salvar", ctx, mock.Anything).Return(nil)
	sRepo.On("Criar", ctx, mock.Anything).Return(&model.Servico{BaseModel: model.BaseModel{ID: 100}}, nil)
	payRepo.On("BuscarPorVaga", ctx, uint(5)).Return(&model.Pagamento{BaseModel: model.BaseModel{ID: 200}}, nil)
	payRepo.On("AtualizarIDServico", ctx, uint(200), uint(100)).Return(nil)
	pRepo.On("Salvar", ctx, mock.Anything).Return(nil)

	err := uc.Responder(ctx, idProposta, idUsuario, true)

	assert.NoError(t, err)
}

func TestPropostaUseCase_ListarPorVaga_Success(t *testing.T) {
	pRepo := new(mocks.MockPropostaRepo)
	vRepo := new(mocks.MockVagaRepo)
	uc := usecases.NewPropostaUseCase(pRepo, vRepo, nil, nil, nil, nil)

	ctx := context.Background()
	idVaga := uint(1)
	idUsuario := uint(10)

	vaga := &model.Vaga{BaseModel: model.BaseModel{ID: 1}, IDCliente: idUsuario}
	propostas := []model.Proposta{{BaseModel: model.BaseModel{ID: 1}, IDVaga: 1}}

	vRepo.On("BuscarPorID", ctx, idVaga).Return(vaga, nil)
	pRepo.On("ListarPorVaga", ctx, idVaga, mock.Anything, "", "", 10, 0).Return(propostas, nil)
	vRepo.On("ZerarPropostasNovas", ctx, idVaga).Return(nil)

	res, err := uc.ListarPorVaga(ctx, idUsuario, idVaga, nil, "", "", 10, 0)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
}

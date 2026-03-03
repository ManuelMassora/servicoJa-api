package mocks

import (
	"context"

	gatewaympesa "github.com/ManuelMassora/servicoJa-api/internal/infra/gateway_mpesa"
	"github.com/stretchr/testify/mock"
)

type MockPagamentoUseCase struct {
	mock.Mock
}

func (m *MockPagamentoUseCase) IniciarPagamentoC2B(ctx context.Context, idPagamento uint, telefone string) error {
	args := m.Called(ctx, idPagamento, telefone)
	return args.Error(0)
}

func (m *MockPagamentoUseCase) ConfirmarPagamentoC2B(ctx context.Context, referencia string) error {
	args := m.Called(ctx, referencia)
	return args.Error(0)
}

func (m *MockPagamentoUseCase) ProcessarCancelamentoComReembolso(ctx context.Context, idServico uint, idUsuarioCancelou uint) error {
	args := m.Called(ctx, idServico, idUsuarioCancelou)
	return args.Error(0)
}

func (m *MockPagamentoUseCase) ProcessarPagamentoPrestador(ctx context.Context, idServico uint) error {
	args := m.Called(ctx, idServico)
	return args.Error(0)
}

func (m *MockPagamentoUseCase) ProcessarCallbackMpesa(ctx context.Context, payload gatewaympesa.MpesaCallbackPayload) error {
	args := m.Called(ctx, payload)
	return args.Error(0)
}

func (m *MockPagamentoUseCase) ProcessarQuerySimulada(ctx context.Context, referencia string) error {
	args := m.Called(ctx, referencia)
	return args.Error(0)
}

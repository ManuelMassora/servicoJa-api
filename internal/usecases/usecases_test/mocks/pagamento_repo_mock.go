package mocks

import (
	"context"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockPagamentoRepo struct {
	mock.Mock
}

func (m *MockPagamentoRepo) Criar(ctx context.Context, pagamento *model.Pagamento) error {
	args := m.Called(ctx, pagamento)
	return args.Error(0)
}

func (m *MockPagamentoRepo) BuscarPorID(ctx context.Context, id uint) (*model.Pagamento, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Pagamento), args.Error(1)
}

func (m *MockPagamentoRepo) BuscarPorReferencia(ctx context.Context, referencia string) (*model.Pagamento, error) {
	args := m.Called(ctx, referencia)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Pagamento), args.Error(1)
}

func (m *MockPagamentoRepo) BuscarPorServico(ctx context.Context, idServico uint) (*model.Pagamento, error) {
	args := m.Called(ctx, idServico)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Pagamento), args.Error(1)
}

func (m *MockPagamentoRepo) BuscarPorVaga(ctx context.Context, idVaga uint) (*model.Pagamento, error) {
	args := m.Called(ctx, idVaga)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Pagamento), args.Error(1)
}

func (m *MockPagamentoRepo) BuscarPorAgendamento(ctx context.Context, idAgendamento uint) (*model.Pagamento, error) {
	args := m.Called(ctx, idAgendamento)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Pagamento), args.Error(1)
}

func (m *MockPagamentoRepo) AtualizarStatus(ctx context.Context, id uint, status model.Status) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockPagamentoRepo) AtualizarStatusPorReferencia(ctx context.Context, referencia string, status model.Status) error {
	args := m.Called(ctx, referencia, status)
	return args.Error(0)
}

func (m *MockPagamentoRepo) AtualizarIDServico(ctx context.Context, idPagamento uint, idServico uint) error {
	args := m.Called(ctx, idPagamento, idServico)
	return args.Error(0)
}

func (m *MockPagamentoRepo) ListarPorUsuario(ctx context.Context, idUsuario uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Pagamento, error) {
	args := m.Called(ctx, idUsuario, filters, orderBy, orderDir, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Pagamento), args.Error(1)
}

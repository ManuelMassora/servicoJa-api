package mocks

import (
	"context"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockNotificacaoRepo struct {
	mock.Mock
}

func (m *MockNotificacaoRepo) Enviar(ctx context.Context, notificacao *model.Notificacao) error {
	args := m.Called(ctx, notificacao)
	return args.Error(0)
}

func (m *MockNotificacaoRepo) BuscarPorID(ctx context.Context, id uint) (*model.Notificacao, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Notificacao), args.Error(1)
}

func (m *MockNotificacaoRepo) ListarPorUsuario(ctx context.Context, idUsuario uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Notificacao, error) {
	args := m.Called(ctx, idUsuario, filters, orderBy, orderDir, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Notificacao), args.Error(1)
}

func (m *MockNotificacaoRepo) MarcarComoLida(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockNotificacaoRepo) MarcarTodasComoLidas(ctx context.Context, idUsuario uint) error {
	args := m.Called(ctx, idUsuario)
	return args.Error(0)
}

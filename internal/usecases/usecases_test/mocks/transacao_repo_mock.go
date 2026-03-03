package mocks

import (
	"context"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockTransacaoRepo struct {
	mock.Mock
}

func (m *MockTransacaoRepo) Criar(ctx context.Context, transacao *model.Transacao) error {
	args := m.Called(ctx, transacao)
	return args.Error(0)
}

func (m *MockTransacaoRepo) ListarPorUsuario(ctx context.Context, idUsuario uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Transacao, error) {
	args := m.Called(ctx, idUsuario, filters, orderBy, orderDir, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Transacao), args.Error(1)
}

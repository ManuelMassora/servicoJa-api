package mocks

import (
	"context"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockClienteRepo struct {
	mock.Mock
}

func (m *MockClienteRepo) Criar(ctx context.Context, cliente *model.Cliente) error {
	args := m.Called(ctx, cliente)
	return args.Error(0)
}

func (m *MockClienteRepo) BuscarPorID(ctx context.Context, id uint) (*model.Cliente, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Cliente), args.Error(1)
}

func (m *MockClienteRepo) Listar(ctx context.Context, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Cliente, error) {
	args := m.Called(ctx, filters, orderBy, orderDir, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Cliente), args.Error(1)
}

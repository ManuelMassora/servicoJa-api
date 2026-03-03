package mocks

import (
	"context"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockCategoriaRepo struct {
	mock.Mock
}

func (m *MockCategoriaRepo) Criar(ctx context.Context, categoria *model.Categoria) error {
	args := m.Called(ctx, categoria)
	return args.Error(0)
}

func (m *MockCategoriaRepo) Editar(ctx context.Context, id uint, campos map[string]interface{}) error {
	args := m.Called(ctx, id, campos)
	return args.Error(0)
}

func (m *MockCategoriaRepo) Listar(ctx context.Context, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Categoria, error) {
	args := m.Called(ctx, filters, orderBy, orderDir, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Categoria), args.Error(1)
}

func (m *MockCategoriaRepo) BuscarPorID(ctx context.Context, id uint) (*model.Categoria, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Categoria), args.Error(1)
}

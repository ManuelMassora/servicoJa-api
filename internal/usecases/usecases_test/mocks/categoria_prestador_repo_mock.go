package mocks

import (
	"context"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockCategoriaPrestadorRepo struct {
	mock.Mock
}

func (m *MockCategoriaPrestadorRepo) Criar(ctx context.Context, categoria *model.CategoriaPrestador) (*model.CategoriaPrestador, error) {
	args := m.Called(ctx, categoria)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CategoriaPrestador), args.Error(1)
}

func (m *MockCategoriaPrestadorRepo) Editar(ctx context.Context, id uint, campos map[string]interface{}) (*model.CategoriaPrestador, error) {
	args := m.Called(ctx, id, campos)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CategoriaPrestador), args.Error(1)
}

func (m *MockCategoriaPrestadorRepo) Listar(ctx context.Context, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.CategoriaPrestador, error) {
	args := m.Called(ctx, filters, orderBy, orderDir, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.CategoriaPrestador), args.Error(1)
}

func (m *MockCategoriaPrestadorRepo) BuscarPorID(ctx context.Context, id uint) (*model.CategoriaPrestador, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CategoriaPrestador), args.Error(1)
}

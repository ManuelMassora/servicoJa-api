package mocks

import (
	"context"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockPrestadorRepo struct {
	mock.Mock
}

func (m *MockPrestadorRepo) Criar(ctx context.Context, prestador *model.Prestador) error {
	args := m.Called(ctx, prestador)
	return args.Error(0)
}

func (m *MockPrestadorRepo) AtualizarStatus(ctx context.Context, id uint, disponivel bool) error {
	args := m.Called(ctx, id, disponivel)
	return args.Error(0)
}

func (m *MockPrestadorRepo) BuscarPorID(ctx context.Context, id uint) (*model.Prestador, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Prestador), args.Error(1)
}

func (m *MockPrestadorRepo) BuscarPorUsuarioID(ctx context.Context, id uint) (*model.Prestador, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Prestador), args.Error(1)
}

func (m *MockPrestadorRepo) Listar(ctx context.Context, filters map[string]interface{}, statusDisponivel interface{}, orderBy, orderDir string, limit, offset int) ([]model.Prestador, error) {
	args := m.Called(ctx, filters, statusDisponivel, orderBy, orderDir, limit, offset)
	return args.Get(0).([]model.Prestador), args.Error(1)
}

func (m *MockPrestadorRepo) FindByLocation(ctx context.Context, latitude, longitude, radius float64, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Prestador, error) {
	args := m.Called(ctx, latitude, longitude, radius, filters, orderBy, orderDir, limit, offset)
	return args.Get(0).([]model.Prestador), args.Error(1)
}

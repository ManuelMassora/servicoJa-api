package mocks

import (
	"context"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockAvaliacaoRepo struct {
	mock.Mock
}

func (m *MockAvaliacaoRepo) Criar(ctx context.Context, avaliacao *model.Avaliacao) error {
	args := m.Called(ctx, avaliacao)
	return args.Error(0)
}

func (m *MockAvaliacaoRepo) ListarPorCliente(ctx context.Context, idCliente uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Avaliacao, error) {
	args := m.Called(ctx, idCliente, filters, orderBy, orderDir, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Avaliacao), args.Error(1)
}

func (m *MockAvaliacaoRepo) ListarPorPrestador(ctx context.Context, idPrestador uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Avaliacao, error) {
	args := m.Called(ctx, idPrestador, filters, orderBy, orderDir, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Avaliacao), args.Error(1)
}

func (m *MockAvaliacaoRepo) MediaPorPrestador(ctx context.Context, idPrestador uint) (float64, error) {
	args := m.Called(ctx, idPrestador)
	return args.Get(0).(float64), args.Error(1)
}

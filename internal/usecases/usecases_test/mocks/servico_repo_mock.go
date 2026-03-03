package mocks

import (
	"context"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockServicoRepo struct {
	mock.Mock
}

func (m *MockServicoRepo) Criar(ctx context.Context, servico *model.Servico) (*model.Servico, error) {
	args := m.Called(ctx, servico)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Servico), args.Error(1)
}

func (m *MockServicoRepo) Atualizar(ctx context.Context, servico *model.Servico) error {
	args := m.Called(ctx, servico)
	return args.Error(0)
}

func (m *MockServicoRepo) BuscarPorID(ctx context.Context, id uint) (*model.Servico, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Servico), args.Error(1)
}

func (m *MockServicoRepo) AtualizarStatus(ctx context.Context, id uint, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockServicoRepo) ListarPorCliente(ctx context.Context, idCliente uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Servico, error) {
	args := m.Called(ctx, idCliente, filters, orderBy, orderDir, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Servico), args.Error(1)
}

func (m *MockServicoRepo) ListarPorPrestador(ctx context.Context, IDPrestador uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Servico, error) {
	args := m.Called(ctx, IDPrestador, filters, orderBy, orderDir, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Servico), args.Error(1)
}

func (m *MockServicoRepo) FindByLocation(ctx context.Context, userID uint, latitude, longitude, radius float64, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Servico, error) {
	args := m.Called(ctx, userID, latitude, longitude, radius, filters, orderBy, orderDir, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Servico), args.Error(1)
}

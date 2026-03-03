package mocks

import (
	"context"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockAgendamentoRepo struct {
	mock.Mock
}

func (m *MockAgendamentoRepo) Criar(ctx context.Context, agendamento *model.Agendamento) (*model.Agendamento, error) {
	args := m.Called(ctx, agendamento)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Agendamento), args.Error(1)
}

func (m *MockAgendamentoRepo) BuscarPorID(ctx context.Context, id uint) (*model.Agendamento, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Agendamento), args.Error(1)
}

func (m *MockAgendamentoRepo) AtualizarStatus(ctx context.Context, id uint, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockAgendamentoRepo) Listar(ctx context.Context, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Agendamento, error) {
	args := m.Called(ctx, filters, orderBy, orderDir, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Agendamento), args.Error(1)
}

func (m *MockAgendamentoRepo) ListarPorClienteID(ctx context.Context, clienteID uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Agendamento, error) {
	args := m.Called(ctx, clienteID, filters, orderBy, orderDir, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Agendamento), args.Error(1)
}

func (m *MockAgendamentoRepo) ListarPorPrestadorID(ctx context.Context, prestadorID uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Agendamento, error) {
	args := m.Called(ctx, prestadorID, filters, orderBy, orderDir, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Agendamento), args.Error(1)
}

func (m *MockAgendamentoRepo) ListarPorCatalogID(ctx context.Context, catalogoID uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Agendamento, error) {
	args := m.Called(ctx, catalogoID, filters, orderBy, orderDir, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Agendamento), args.Error(1)
}

func (m *MockAgendamentoRepo) FindByLocation(ctx context.Context, userID uint, latitude, longitude, radius float64, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Agendamento, error) {
	args := m.Called(ctx, userID, latitude, longitude, radius, filters, orderBy, orderDir, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Agendamento), args.Error(1)
}

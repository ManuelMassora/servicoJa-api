package mocks

import (
	"context"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockVagaRepo struct {
	mock.Mock
}

func (m *MockVagaRepo) Criar(ctx context.Context, vaga *model.Vaga) (*model.Vaga, error) {
	args := m.Called(ctx, vaga)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Vaga), args.Error(1)
}

func (m *MockVagaRepo) Salvar(ctx context.Context, vaga *model.Vaga) error {
	args := m.Called(ctx, vaga)
	return args.Error(0)
}

func (m *MockVagaRepo) BuscarPorID(ctx context.Context, id uint) (*model.Vaga, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Vaga), args.Error(1)
}

func (m *MockVagaRepo) ListarDisponiveis(ctx context.Context, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Vaga, error) {
	args := m.Called(ctx, filters, orderBy, orderDir, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Vaga), args.Error(1)
}

func (m *MockVagaRepo) ListarPorCliente(ctx context.Context, idCliente uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Vaga, error) {
	args := m.Called(ctx, idCliente, filters, orderBy, orderDir, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Vaga), args.Error(1)
}

func (m *MockVagaRepo) AceitarVaga(ctx context.Context, idVaga, idPrestador uint) error {
	args := m.Called(ctx, idVaga, idPrestador)
	return args.Error(0)
}

func (m *MockVagaRepo) AtualizarStatus(ctx context.Context, idVaga uint, status model.Status) error {
	args := m.Called(ctx, idVaga, status)
	return args.Error(0)
}

func (m *MockVagaRepo) FindByLocation(ctx context.Context, latitude, longitude, radius float64, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Vaga, error) {
	args := m.Called(ctx, latitude, longitude, radius, filters, orderBy, orderDir, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Vaga), args.Error(1)
}

func (m *MockVagaRepo) IncrementarPropostasNovas(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockVagaRepo) ZerarPropostasNovas(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

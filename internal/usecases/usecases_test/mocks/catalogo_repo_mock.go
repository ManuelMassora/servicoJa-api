package mocks

import (
	"context"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockCatalogoRepo struct {
	mock.Mock
}

func (m *MockCatalogoRepo) Create(ctx context.Context, catalogo *model.Catalogo) error {
	args := m.Called(ctx, catalogo)
	return args.Error(0)
}

func (m *MockCatalogoRepo) Update(ctx context.Context, id uint, campos map[string]interface{}) error {
	args := m.Called(ctx, id, campos)
	return args.Error(0)
}

func (m *MockCatalogoRepo) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCatalogoRepo) FindByID(ctx context.Context, id uint) (*model.Catalogo, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Catalogo), args.Error(1)
}

func (m *MockCatalogoRepo) FindAll(ctx context.Context, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]*model.Catalogo, error) {
	args := m.Called(ctx, filters, orderBy, orderDir, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Catalogo), args.Error(1)
}

func (m *MockCatalogoRepo) FindByPrestadorID(ctx context.Context, prestadorID uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]*model.Catalogo, error) {
	args := m.Called(ctx, prestadorID, filters, orderBy, orderDir, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Catalogo), args.Error(1)
}

func (m *MockCatalogoRepo) FindByLocation(ctx context.Context, latitude, longitude, radius float64, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]*model.Catalogo, error) {
	args := m.Called(ctx, latitude, longitude, radius, filters, orderBy, orderDir, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Catalogo), args.Error(1)
}

func (m *MockCatalogoRepo) IncrementarAgendamentosNovos(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCatalogoRepo) ZerarAgendamentosNovos(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

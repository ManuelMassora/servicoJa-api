package mocks

import (
	"context"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"github.com/stretchr/testify/mock"
)

type GaleriaRepoMock struct {
	mock.Mock
}

func (m *GaleriaRepoMock) Create(ctx context.Context, galeria *model.Galeria) (*model.Galeria, error) {
	args := m.Called(ctx, galeria)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Galeria), args.Error(1)
}

func (m *GaleriaRepoMock) FindByID(ctx context.Context, id uint) (*model.Galeria, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Galeria), args.Error(1)
}

func (m *GaleriaRepoMock) FindByPrestadorID(ctx context.Context, prestadorID uint) (*model.Galeria, error) {
	args := m.Called(ctx, prestadorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Galeria), args.Error(1)
}

func (m *GaleriaRepoMock) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *GaleriaRepoMock) AddImage(ctx context.Context, imagem *model.Imagem) error {
	args := m.Called(ctx, imagem)
	return args.Error(0)
}

func (m *GaleriaRepoMock) CountImages(ctx context.Context, galeriaID uint) (int64, error) {
	args := m.Called(ctx, galeriaID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *GaleriaRepoMock) FindByGaleriaID(ctx context.Context, galeriaID uint) ([]model.Imagem, error) {
	args := m.Called(ctx, galeriaID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Imagem), args.Error(1)
}

func (m *GaleriaRepoMock) FindByPrestadorIDs(ctx context.Context, prestadorIDs []uint) ([]model.Galeria, error) {
	args := m.Called(ctx, prestadorIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Galeria), args.Error(1)
}

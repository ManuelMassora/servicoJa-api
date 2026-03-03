package mocks

import (
	"context"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockAnexoImagemRepo struct {
	mock.Mock
}

func (m *MockAnexoImagemRepo) Create(ctx context.Context, anexo *model.AnexoImagem) error {
	args := m.Called(ctx, anexo)
	return args.Error(0)
}

func (m *MockAnexoImagemRepo) FindByID(ctx context.Context, id uint) (*model.AnexoImagem, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.AnexoImagem), args.Error(1)
}

func (m *MockAnexoImagemRepo) FindByAgendamentoID(ctx context.Context, agendamentoID uint) ([]model.AnexoImagem, error) {
	args := m.Called(ctx, agendamentoID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.AnexoImagem), args.Error(1)
}

func (m *MockAnexoImagemRepo) FindByVagaID(ctx context.Context, vagaID uint) ([]model.AnexoImagem, error) {
	args := m.Called(ctx, vagaID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.AnexoImagem), args.Error(1)
}

func (m *MockAnexoImagemRepo) FindByCatalogoID(ctx context.Context, catalogoID uint) ([]model.AnexoImagem, error) {
	args := m.Called(ctx, catalogoID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.AnexoImagem), args.Error(1)
}

func (m *MockAnexoImagemRepo) FindByAgendamentoIDs(ctx context.Context, agendamentoIDs []uint) ([]model.AnexoImagem, error) {
	args := m.Called(ctx, agendamentoIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.AnexoImagem), args.Error(1)
}

func (m *MockAnexoImagemRepo) FindByVagaIDs(ctx context.Context, vagaIDs []uint) ([]model.AnexoImagem, error) {
	args := m.Called(ctx, vagaIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.AnexoImagem), args.Error(1)
}

func (m *MockAnexoImagemRepo) FindByCatalogoIDs(ctx context.Context, catalogoIDs []uint) ([]model.AnexoImagem, error) {
	args := m.Called(ctx, catalogoIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.AnexoImagem), args.Error(1)
}

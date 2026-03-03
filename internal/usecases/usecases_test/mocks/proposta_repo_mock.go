package mocks

import (
	"context"
	"time"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockPropostaRepo struct {
	mock.Mock
}

func (m *MockPropostaRepo) Criar(ctx context.Context, proposta *model.Proposta) error {
	args := m.Called(ctx, proposta)
	return args.Error(0)
}

func (m *MockPropostaRepo) Salvar(ctx context.Context, proposta *model.Proposta) error {
	args := m.Called(ctx, proposta)
	return args.Error(0)
}

func (m *MockPropostaRepo) BuscarPorID(ctx context.Context, idProposta uint) (*model.Proposta, error) {
	args := m.Called(ctx, idProposta)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Proposta), args.Error(1)
}

func (m *MockPropostaRepo) ListarPorVaga(ctx context.Context, idVaga uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Proposta, error) {
	args := m.Called(ctx, idVaga, filters, orderBy, orderDir, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Proposta), args.Error(1)
}

func (m *MockPropostaRepo) ListarPorPrestador(ctx context.Context, idPrestador uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Proposta, error) {
	args := m.Called(ctx, idPrestador, filters, orderBy, orderDir, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Proposta), args.Error(1)
}

func (m *MockPropostaRepo) AtualizarStatus(ctx context.Context, idProposta uint, status model.Status, dataResposta time.Time) error {
	args := m.Called(ctx, idProposta, status, dataResposta)
	return args.Error(0)
}

package mocks

import (
	"context"
	"time"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockUsuarioRepo struct {
	mock.Mock
}

func (m *MockUsuarioRepo) Criar(ctx context.Context, usuario *model.Usuario) error {
	args := m.Called(ctx, usuario)
	return args.Error(0)
}

func (m *MockUsuarioRepo) BuscarPorID(ctx context.Context, id uint) (*model.Usuario, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Usuario), args.Error(1)
}

func (m *MockUsuarioRepo) IncrementarNotificacoesNovas(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUsuarioRepo) ZerarNotificacoesNovas(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUsuarioRepo) IncrementarCancelamentos(ctx context.Context, id uint) (uint, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(uint), args.Error(1)
}

func (m *MockUsuarioRepo) SuspenderUsuario(ctx context.Context, id uint, ate time.Time) error {
	args := m.Called(ctx, id, ate)
	return args.Error(0)
}

func (m *MockUsuarioRepo) BuscarPorTelefone(ctx context.Context, numero string) (*model.Usuario, error) {
	args := m.Called(ctx, numero)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Usuario), args.Error(1)
}

func (m *MockUsuarioRepo) Atualizar(ctx context.Context, usuario *model.Usuario) error {
	args := m.Called(ctx, usuario)
	return args.Error(0)
}

func (m *MockUsuarioRepo) Remover(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUsuarioRepo) ListarTodos(ctx context.Context, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Usuario, error) {
	args := m.Called(ctx, filters, orderBy, orderDir, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Usuario), args.Error(1)
}

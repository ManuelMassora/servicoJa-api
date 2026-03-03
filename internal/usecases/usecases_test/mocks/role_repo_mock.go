package mocks

import (
	"context"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockRoleRepo struct {
	mock.Mock
}

func (m *MockRoleRepo) ListarPermissoes(ctx context.Context, role model.Role) ([]model.RolePermissao, error) {
	args := m.Called(ctx, role)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.RolePermissao), args.Error(1)
}

func (m *MockRoleRepo) AdicionarPermissao(ctx context.Context, role model.Role) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

func (m *MockRoleRepo) RemoverPermissao(ctx context.Context, role model.Role) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

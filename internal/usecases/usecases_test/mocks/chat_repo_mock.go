package mocks

import (
	"context"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockChatRepo struct {
	mock.Mock
}

func (m *MockChatRepo) CriarChat(ctx context.Context, chat *model.Chat) error {
	args := m.Called(ctx, chat)
	return args.Error(0)
}

func (m *MockChatRepo) ListarChatsPorUsuario(ctx context.Context, idUsuario uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Chat, error) {
	args := m.Called(ctx, idUsuario, filters, orderBy, orderDir, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Chat), args.Error(1)
}

type MockMensagemRepo struct {
	mock.Mock
}

func (m *MockMensagemRepo) EnviarMensagem(ctx context.Context, msg *model.Mensagem) error {
	args := m.Called(ctx, msg)
	return args.Error(0)
}

func (m *MockMensagemRepo) ListarMensagens(ctx context.Context, idChat uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Mensagem, error) {
	args := m.Called(ctx, idChat, filters, orderBy, orderDir, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Mensagem), args.Error(1)
}

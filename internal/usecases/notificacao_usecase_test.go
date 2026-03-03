package usecases_test

import (
	"context"
	"testing"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"github.com/ManuelMassora/servicoJa-api/internal/usecases"
	"github.com/ManuelMassora/servicoJa-api/internal/usecases/usecases_test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNotificacaoUseCase_ListarPorUsuario_Success(t *testing.T) {
	repo := new(mocks.MockNotificacaoRepo)
	uRepo := new(mocks.MockUsuarioRepo)
	uc := usecases.NewNotificacaoUseCase(repo, uRepo)

	ctx := context.Background()
	idUsuario := uint(1)
	notificacoes := []model.Notificacao{{BaseModel: model.BaseModel{ID: 1}, Titulo: "Olá"}}

	repo.On("ListarPorUsuario", ctx, idUsuario, mock.Anything, "", "", 10, 0).Return(notificacoes, nil)
	uRepo.On("ZerarNotificacoesNovas", ctx, idUsuario).Return(nil)

	res, err := uc.ListarPorUsuario(ctx, idUsuario, nil, "", "", 10, 0)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
}

func TestNotificacaoUseCase_MarcarComoLida_Success(t *testing.T) {
	repo := new(mocks.MockNotificacaoRepo)
	uc := usecases.NewNotificacaoUseCase(repo, nil)

	ctx := context.Background()
	id := uint(10)
	idUsuario := uint(1)

	repo.On("BuscarPorID", ctx, id).Return(&model.Notificacao{BaseModel: model.BaseModel{ID: id}, IDUsuario: idUsuario}, nil)
	repo.On("MarcarComoLida", ctx, id).Return(nil)

	err := uc.MarcarComoLida(ctx, id, idUsuario)

	assert.NoError(t, err)
}

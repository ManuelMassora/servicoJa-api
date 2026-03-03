package usecases_test

import (
	"context"
	"testing"
	"time"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"github.com/ManuelMassora/servicoJa-api/internal/usecases"
	"github.com/ManuelMassora/servicoJa-api/internal/usecases/usecases_test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAvaliacaoUseCase_Criar_Success(t *testing.T) {
	avaliacaoRepo := new(mocks.MockAvaliacaoRepo)
	servicoRepo := new(mocks.MockServicoRepo)
	notificacaoRepo := new(mocks.MockNotificacaoRepo)
	usuarioRepo := new(mocks.MockUsuarioRepo)

	uc := usecases.NewAvaliacaoUseCase(avaliacaoRepo, servicoRepo, notificacaoRepo, usuarioRepo)

	ctx := context.Background()
	idAvaliador := uint(1)
	req := usecases.AvaliacaoRequest{
		IDServico:  1,
		Pontuacao:  5,
		Comentario: "Excelente!",
	}

	servicoMock := &model.Servico{
		BaseModel:         model.BaseModel{ID: 1},
		Status:            model.StatusConcluido,
		IDCliente:         idAvaliador,
		IDPrestador:       2,
		IfAvaliadoCliente: false,
		Localizacao:       "Maputo",
		Cliente:           &model.Cliente{Usuario: model.Usuario{Nome: "Cliente"}},
		Prestador:         &model.Prestador{Usuario: model.Usuario{Nome: "Prestador"}},
	}

	servicoRepo.On("BuscarPorID", ctx, req.IDServico).Return(servicoMock, nil)
	servicoRepo.On("Atualizar", ctx, mock.Anything).Return(nil)
	avaliacaoRepo.On("Criar", ctx, mock.Anything).Return(nil)
	notificacaoRepo.On("Enviar", ctx, mock.Anything).Return(nil)
	usuarioRepo.On("IncrementarNotificacoesNovas", ctx, uint(2)).Return(nil)

	err := uc.Criar(ctx, req, idAvaliador)

	assert.NoError(t, err)
	servicoRepo.AssertExpectations(t)
	avaliacaoRepo.AssertExpectations(t)
	notificacaoRepo.AssertExpectations(t)
	usuarioRepo.AssertExpectations(t)
}

func TestAvaliacaoUseCase_Criar_ServicoNaoConcluido(t *testing.T) {
	avaliacaoRepo := new(mocks.MockAvaliacaoRepo)
	servicoRepo := new(mocks.MockServicoRepo)
	notificacaoRepo := new(mocks.MockNotificacaoRepo)
	usuarioRepo := new(mocks.MockUsuarioRepo)

	uc := usecases.NewAvaliacaoUseCase(avaliacaoRepo, servicoRepo, notificacaoRepo, usuarioRepo)

	ctx := context.Background()
	req := usecases.AvaliacaoRequest{IDServico: 1, Pontuacao: 5}

	servicoMock := &model.Servico{
		BaseModel: model.BaseModel{ID: 1},
		Status:    model.StatusPendente,
	}

	servicoRepo.On("BuscarPorID", ctx, req.IDServico).Return(servicoMock, nil)

	err := uc.Criar(ctx, req, 1)

	assert.Error(t, err)
	assert.Equal(t, "não é possível avaliar um serviço que não foi concluído", err.Error())
}

func TestAvaliacaoUseCase_Criar_JaAvaliado(t *testing.T) {
	avaliacaoRepo := new(mocks.MockAvaliacaoRepo)
	servicoRepo := new(mocks.MockServicoRepo)
	notificacaoRepo := new(mocks.MockNotificacaoRepo)
	usuarioRepo := new(mocks.MockUsuarioRepo)

	uc := usecases.NewAvaliacaoUseCase(avaliacaoRepo, servicoRepo, notificacaoRepo, usuarioRepo)

	ctx := context.Background()
	req := usecases.AvaliacaoRequest{IDServico: 1, Pontuacao: 5}

	servicoMock := &model.Servico{
		BaseModel:         model.BaseModel{ID: 1},
		Status:            model.StatusConcluido,
		IDCliente:         1,
		IfAvaliadoCliente: true,
	}

	servicoRepo.On("BuscarPorID", ctx, req.IDServico).Return(servicoMock, nil)

	err := uc.Criar(ctx, req, 1)

	assert.Error(t, err)
	assert.Equal(t, "usuário já avaliou este serviço", err.Error())
}

func TestAvaliacaoUseCase_ListarPorPrestador(t *testing.T) {
	avaliacaoRepo := new(mocks.MockAvaliacaoRepo)
	uc := usecases.NewAvaliacaoUseCase(avaliacaoRepo, nil, nil, nil)

	ctx := context.Background()
	idPrestador := uint(1)
	avaliacoes := []model.Avaliacao{
		{
			BaseModel:  model.BaseModel{ID: 1, CreatedAt: time.Now()},
			Nota:       5,
			Comentario: "Top",
			Servico:    &model.Servico{Localizacao: "Local 1"},
		},
	}

	avaliacaoRepo.On("ListarPorPrestador", ctx, idPrestador, mock.Anything, "", "", 10, 0).
		Return(avaliacoes, nil)

	res, err := uc.ListarPorPrestador(ctx, idPrestador, nil, "", "", 10, 0)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "Top", res[0].Comentario)
}

func TestAvaliacaoUseCase_MediaPorPrestador(t *testing.T) {
	avaliacaoRepo := new(mocks.MockAvaliacaoRepo)
	uc := usecases.NewAvaliacaoUseCase(avaliacaoRepo, nil, nil, nil)

	ctx := context.Background()
	idPrestador := uint(1)
	avaliacaoRepo.On("MediaPorPrestador", ctx, idPrestador).Return(4.5, nil)

	res, err := uc.MediaPorPrestador(ctx, idPrestador)

	assert.NoError(t, err)
	assert.Equal(t, 4.5, res)
}

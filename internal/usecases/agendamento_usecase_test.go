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

func TestAgendamentoUseCase_Criar_Success(t *testing.T) {
	repo := new(mocks.MockAgendamentoRepo)
	catalogoRepo := new(mocks.MockCatalogoRepo)
	servicoRepo := new(mocks.MockServicoRepo)
	notificacaoRepo := new(mocks.MockNotificacaoRepo)
	usuarioRepo := new(mocks.MockUsuarioRepo)
	anexoImagemRepo := new(mocks.MockAnexoImagemRepo)
	pagamentoRepo := new(mocks.MockPagamentoRepo)
	pagamentoUC := new(mocks.MockPagamentoUseCase)

	// Since PagamentoUseCase is a struct, we pass nil or a basic instance.
	// In the real code, IniciarPagamentoC2B is called and its error is ignored.
	uc := usecases.NewAgendamentoUC(repo, catalogoRepo, servicoRepo, notificacaoRepo, usuarioRepo, anexoImagemRepo, pagamentoRepo, pagamentoUC)

	ctx := context.Background()
	idCliente := uint(1)
	req := &usecases.AgendamentoRequest{
		IDCatalogo:        1,
		Detalhe:           "Teste",
		DataHora:          time.Now(),
		Localizacao:       "Maputo",
		Latitude:          -25.9,
		Longitude:         32.5,
		Anexos:            []string{"img1.jpg"},
		TelefonePagamento: "841112233",
	}

	usuarioRepo.On("BuscarPorID", ctx, idCliente).Return(&model.Usuario{BaseModel: model.BaseModel{ID: idCliente}}, nil)
	catalogoRepo.On("FindByID", ctx, req.IDCatalogo).Return(&model.Catalogo{
		BaseModel:   model.BaseModel{ID: 1},
		Nome:        "Cat 1",
		IDPrestador: 2,
		Prestador:   model.Prestador{IDUsuario: 2},
	}, nil)
	repo.On("Criar", ctx, mock.Anything).Return(&model.Agendamento{BaseModel: model.BaseModel{ID: 10}}, nil)
	pagamentoRepo.On("Criar", ctx, mock.Anything).Return(nil)
	pagamentoUC.On("IniciarPagamentoC2B", ctx, mock.Anything, req.TelefonePagamento).Return(nil)
	anexoImagemRepo.On("Create", ctx, mock.Anything).Return(nil)
	notificacaoRepo.On("Enviar", ctx, mock.Anything).Return(nil)
	usuarioRepo.On("IncrementarNotificacoesNovas", ctx, mock.Anything).Return(nil)
	catalogoRepo.On("IncrementarAgendamentosNovos", ctx, mock.Anything).Return(nil)

	err := uc.Criar(ctx, req, idCliente)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestAgendamentoUseCase_Buscar_Success(t *testing.T) {
	repo := new(mocks.MockAgendamentoRepo)
	anexoRepo := new(mocks.MockAnexoImagemRepo)
	pagamentoUC := new(mocks.MockPagamentoUseCase)
	uc := usecases.NewAgendamentoUC(repo, nil, nil, nil, nil, anexoRepo, nil, pagamentoUC)

	ctx := context.Background()
	id := uint(10)
	idUsuario := uint(1)

	agendamento := &model.Agendamento{
		BaseModel: model.BaseModel{ID: id},
		IDCliente: idUsuario,
		Cliente:   model.Cliente{IDUsuario: idUsuario, Nome: "Cliente"},
		Catalogo: model.Catalogo{
			Nome:      "Corta Cabelo",
			Prestador: model.Prestador{IDUsuario: 2, Nome: "Barbeiro"},
		},
	}

	repo.On("BuscarPorID", ctx, id).Return(agendamento, nil)
	anexoRepo.On("FindByAgendamentoID", ctx, id).Return([]model.AnexoImagem{{URL: "u1"}}, nil)

	res, err := uc.Buscar(ctx, id, idUsuario)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "Cliente", res.Cliente.ClienteNome)
}

func TestAgendamentoUseCase_Aceitar_Success(t *testing.T) {
	repo := new(mocks.MockAgendamentoRepo)
	notificacaoRepo := new(mocks.MockNotificacaoRepo)
	usuarioRepo := new(mocks.MockUsuarioRepo)
	servicoRepo := new(mocks.MockServicoRepo)
	pagamentoRepo := new(mocks.MockPagamentoRepo)

	pagamentoUC := new(mocks.MockPagamentoUseCase)
	uc := usecases.NewAgendamentoUC(repo, nil, servicoRepo, notificacaoRepo, usuarioRepo, nil, pagamentoRepo, pagamentoUC)

	ctx := context.Background()
	id := uint(10)
	idPrestadorUser := uint(2)

	agendamento := &model.Agendamento{
		BaseModel: model.BaseModel{ID: id},
		IDCliente: 1,
		Status:    "PENDENTE",
		Catalogo: model.Catalogo{
			Nome:        "Servico",
			IDPrestador: 2,
			TipoPreco:   "fixo",
			ValorFixo:   100,
			Prestador: model.Prestador{
				Usuario: model.Usuario{BaseModel: model.BaseModel{ID: idPrestadorUser}},
			},
		},
	}

	repo.On("BuscarPorID", ctx, id).Return(agendamento, nil)
	notificacaoRepo.On("Enviar", ctx, mock.Anything).Return(nil)
	usuarioRepo.On("IncrementarNotificacoesNovas", ctx, uint(1)).Return(nil)
	servicoRepo.On("Criar", ctx, mock.Anything).Return(&model.Servico{BaseModel: model.BaseModel{ID: 100}}, nil)
	pagamentoRepo.On("BuscarPorAgendamento", ctx, id).Return(&model.Pagamento{BaseModel: model.BaseModel{ID: 50}}, nil)
	pagamentoRepo.On("AtualizarIDServico", ctx, uint(50), uint(100)).Return(nil)
	repo.On("AtualizarStatus", ctx, id, "EM_ANDAMENTO").Return(nil)

	err := uc.Aceitar(ctx, id, idPrestadorUser)

	assert.NoError(t, err)
}

func TestAgendamentoUseCase_Listar_Success(t *testing.T) {
	repo := new(mocks.MockAgendamentoRepo)
	anexoRepo := new(mocks.MockAnexoImagemRepo)
	pagamentoUC := new(mocks.MockPagamentoUseCase)
	uc := usecases.NewAgendamentoUC(repo, nil, nil, nil, nil, anexoRepo, nil, pagamentoUC)

	ctx := context.Background()
	agendamentos := []model.Agendamento{
		{
			BaseModel: model.BaseModel{ID: 1},
			Catalogo:  model.Catalogo{Nome: "Cat"},
			Cliente:   model.Cliente{Nome: "Cli"},
		},
	}

	repo.On("Listar", ctx, mock.Anything, "", "", 10, 0).Return(agendamentos, nil)
	anexoRepo.On("FindByAgendamentoIDs", ctx, []uint{1}).Return(nil, nil)

	res, err := uc.Listar(ctx, nil, "", "", 10, 0)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
}

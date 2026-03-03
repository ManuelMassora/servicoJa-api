package usecases_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"github.com/ManuelMassora/servicoJa-api/internal/usecases"
	"github.com/ManuelMassora/servicoJa-api/internal/usecases/usecases_test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestUsuarioUseCase_CriarAdmin_Sucesso(t *testing.T) {
	usuarioRepo := new(mocks.MockUsuarioRepo)
	clienteRepo := new(mocks.MockClienteRepo)
	prestadorRepo := new(mocks.MockPrestadorRepo)
	galeriaRepo := new(mocks.MockGaleriaRepo)

	uc := usecases.NewUsuarioUseCase(usuarioRepo, clienteRepo, prestadorRepo, galeriaRepo)

	ctx := context.Background()
	req := usecases.UsuarioRequest{
		Nome:     "Admin User",
		Telefone: "841111111",
		Senha:    "senha123",
	}

	// Telefone não existe
	usuarioRepo.On("BuscarPorTelefone", ctx, req.Telefone).
		Return((*model.Usuario)(nil), gorm.ErrRecordNotFound)

	// Criar admin
	usuarioRepo.On("Criar", ctx, mock.Anything).Return(nil)

	err := uc.CriarAdmin(ctx, req)
	assert.NoError(t, err)

	usuarioRepo.AssertExpectations(t)
}

func TestUsuarioUseCase_CriarAdmin_TelefoneJaExiste(t *testing.T) {
	usuarioRepo := new(mocks.MockUsuarioRepo)
	clienteRepo := new(mocks.MockClienteRepo)
	prestadorRepo := new(mocks.MockPrestadorRepo)
	galeriaRepo := new(mocks.MockGaleriaRepo)

	uc := usecases.NewUsuarioUseCase(usuarioRepo, clienteRepo, prestadorRepo, galeriaRepo)

	ctx := context.Background()
	req := usecases.UsuarioRequest{
		Nome:     "Admin",
		Telefone: "841111111",
		Senha:    "password",
	}

	// Mock: telefone já existe
	usuarioRepo.On("BuscarPorTelefone", ctx, req.Telefone).Return(&model.Usuario{Telefone: req.Telefone}, nil)

	err := uc.CriarAdmin(ctx, req)
	assert.Error(t, err)
	assert.Equal(t, "ja existe usuario com mesmo contacto", err.Error())

	usuarioRepo.AssertExpectations(t)
}

func TestUsuarioUseCase_CriarCliente_Sucesso(t *testing.T) {
	usuarioRepo := new(mocks.MockUsuarioRepo)
	clienteRepo := new(mocks.MockClienteRepo)
	prestadorRepo := new(mocks.MockPrestadorRepo)
	galeriaRepo := new(mocks.MockGaleriaRepo)

	uc := usecases.NewUsuarioUseCase(usuarioRepo, clienteRepo, prestadorRepo, galeriaRepo)

	ctx := context.Background()
	req := usecases.UsuarioRequest{
		Nome:     "Cliente Teste",
		Telefone: "841111112",
		Senha:    "senha456",
	}

	usuarioRepo.On("BuscarPorTelefone", ctx, req.Telefone).Return((*model.Usuario)(nil), gorm.ErrRecordNotFound)
	clienteRepo.On("Criar", ctx, mock.AnythingOfType("*model.Cliente")).Return(nil)

	err := uc.CriarCliente(ctx, req)
	assert.NoError(t, err)

	usuarioRepo.AssertExpectations(t)
	clienteRepo.AssertExpectations(t)
}

func TestUsuarioUseCase_CriarPrestador_Sucesso(t *testing.T) {
	usuarioRepo := new(mocks.MockUsuarioRepo)
	clienteRepo := new(mocks.MockClienteRepo)
	prestadorRepo := new(mocks.MockPrestadorRepo)
	galeriaRepo := new(mocks.MockGaleriaRepo)

	uc := usecases.NewUsuarioUseCase(usuarioRepo, clienteRepo, prestadorRepo, galeriaRepo)

	ctx := context.Background()
	req := usecases.PrestadorRequest{
		Nome:        "Prestador X",
		Telefone:    "841111113",
		Senha:       "senha789",
		Localizacao: "Maputo",
	}

	usuarioRepo.On("BuscarPorTelefone", ctx, req.Telefone).Return((*model.Usuario)(nil), gorm.ErrRecordNotFound)
	prestadorRepo.On("Criar", ctx, mock.AnythingOfType("*model.Prestador")).Return(nil)

	err := uc.CriarPrestador(ctx, req)
	assert.NoError(t, err)

	usuarioRepo.AssertExpectations(t)
	prestadorRepo.AssertExpectations(t)
}

func TestUsuarioUseCase_ListarTodosUsuarios(t *testing.T) {
	usuarioRepo := new(mocks.MockUsuarioRepo)
	clienteRepo := new(mocks.MockClienteRepo)
	prestadorRepo := new(mocks.MockPrestadorRepo)
	galeriaRepo := new(mocks.MockGaleriaRepo)

	uc := usecases.NewUsuarioUseCase(usuarioRepo, clienteRepo, prestadorRepo, galeriaRepo)

	ctx := context.Background()
	usuariosMock := []model.Usuario{
		{Nome: "User 1", Telefone: "111", RolePermissao: model.RolePermissao{Role: "CLIENTE"}},
		{Nome: "User 2", Telefone: "222", RolePermissao: model.RolePermissao{Role: "PRESTADOR"}},
	}
	clienteMock := model.Cliente{ImagemURL: "cliente.jpg"}
	prestadorMock := model.Prestador{ImagemURL: "prestador.jpg"}

	usuarioRepo.On("ListarTodos", ctx, mock.Anything, "", "", 10, 0).
		Return(usuariosMock, nil)
	clienteRepo.On("BuscarPorID", ctx, mock.Anything).Return(&clienteMock, nil)
	prestadorRepo.On("BuscarPorID", ctx, mock.Anything).Return(&prestadorMock, nil)

	resp, err := uc.ListarTodosUsuarios(ctx, nil, "", "", 10, 0)
	assert.NoError(t, err)
	assert.Len(t, resp, 2)
	assert.Equal(t, "User 1", resp[0].Nome)
	assert.Equal(t, "222", resp[1].Telefone)

	usuarioRepo.AssertExpectations(t)
}

func TestUsuarioUseCase_ListarPrestadores(t *testing.T) {
	usuarioRepo := new(mocks.MockUsuarioRepo)
	clienteRepo := new(mocks.MockClienteRepo)
	prestadorRepo := new(mocks.MockPrestadorRepo)
	galeriaRepo := new(mocks.MockGaleriaRepo)

	uc := usecases.NewUsuarioUseCase(usuarioRepo, clienteRepo, prestadorRepo, galeriaRepo)

	ctx := context.Background()
	prestadoresMock := []model.Prestador{
		{
			IDUsuario:        1,
			Nome:             "Prest 1",
			Telefone:         "333",
			Localizacao:      "Matola",
			StatusDisponivel: true,
		},
	}
	galeriasMock := []model.Galeria{
		{
			PrestadorID: 1,
			Imagens: []model.Imagem{
				{URL: "image1.jpg"},
			},
		},
	}

	prestadorRepo.On("Listar", ctx, mock.Anything, mock.Anything, "", "", 10, 0).
		Return(prestadoresMock, nil)
	galeriaRepo.On("FindByPrestadorIDs", ctx, []uint{1}).Return(galeriasMock, nil)

	resp, err := uc.ListarPrestadores(ctx, nil, true, "", "", 10, 0)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "Prest 1", resp[0].Nome)
	assert.True(t, resp[0].Disponivel)

	prestadorRepo.AssertExpectations(t)
}

func TestUsuarioUseCase_seTelefoneExiste_NaoExiste(t *testing.T) {
	usuarioRepo := new(mocks.MockUsuarioRepo)
	uc := usecases.NewUsuarioUseCase(usuarioRepo, nil, nil, nil) // use o construtor!

	ctx := context.Background()
	usuarioRepo.On("BuscarPorTelefone", ctx, "999").Return((*model.Usuario)(nil), gorm.ErrRecordNotFound)

	err := uc.SeTelefoneExiste(ctx, "999")
	assert.NoError(t, err)

	usuarioRepo.AssertExpectations(t)
}

func TestUsuarioUseCase_seTelefoneExiste_JaExiste(t *testing.T) {
	usuarioRepo := new(mocks.MockUsuarioRepo)
	uc := usecases.NewUsuarioUseCase(usuarioRepo, nil, nil, nil)

	ctx := context.Background()
	usuarioRepo.On("BuscarPorTelefone", ctx, "888").Return(&model.Usuario{Telefone: "888"}, nil)

	err := uc.SeTelefoneExiste(ctx, "888")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ja existe usuario com mesmo contacto")

	usuarioRepo.AssertExpectations(t)
}

func TestUsuarioUseCase_seTelefoneExiste_ErroInterno(t *testing.T) {
	usuarioRepo := new(mocks.MockUsuarioRepo)
	uc := usecases.NewUsuarioUseCase(usuarioRepo, nil, nil, nil)

	ctx := context.Background()
	expectedErr := errors.New("db error")
	usuarioRepo.On("BuscarPorTelefone", ctx, "777").Return((*model.Usuario)(nil), expectedErr)

	err := uc.SeTelefoneExiste(ctx, "777")
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)

	usuarioRepo.AssertExpectations(t)
}

func TestUsuarioUseCase_EditarPrestador_Sucesso(t *testing.T) {
	usuarioRepo := new(mocks.MockUsuarioRepo)
	clienteRepo := new(mocks.MockClienteRepo)
	prestadorRepo := new(mocks.MockPrestadorRepo)
	galeriaRepo := new(mocks.MockGaleriaRepo)

	uc := usecases.NewUsuarioUseCase(usuarioRepo, clienteRepo, prestadorRepo, galeriaRepo)

	ctx := context.Background()
	userID := uint(1)
	campos := map[string]interface{}{
		"nome": "Novo Nome",
	}

	prestadorMock := &model.Prestador{
		IDUsuario: userID,
		Nome:      "Novo Nome",
		Telefone:  "123456789",
	}

	prestadorRepo.On("Editar", ctx, userID, campos).Return(prestadorMock, nil)

	resp, err := uc.EditarPrestador(ctx, userID, campos)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Novo Nome", resp.Nome)
	prestadorRepo.AssertExpectations(t)
}

func TestUsuarioUseCase_EditarPrestador_ErroRepositorio(t *testing.T) {
	usuarioRepo := new(mocks.MockUsuarioRepo)
	clienteRepo := new(mocks.MockClienteRepo)
	prestadorRepo := new(mocks.MockPrestadorRepo)
	galeriaRepo := new(mocks.MockGaleriaRepo)

	uc := usecases.NewUsuarioUseCase(usuarioRepo, clienteRepo, prestadorRepo, galeriaRepo)

	ctx := context.Background()
	userID := uint(1)
	campos := map[string]interface{}{
		"nome": "Novo Nome",
	}

	expectedErr := errors.New("erro ao editar")

	prestadorRepo.On("Editar", ctx, userID, campos).Return(nil, expectedErr)

	resp, err := uc.EditarPrestador(ctx, userID, campos)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, expectedErr, err)
	prestadorRepo.AssertExpectations(t)
}

func TestUsuarioUseCase_BuscarPrestador_Success(t *testing.T) {
	prestadorRepo := new(mocks.MockPrestadorRepo)
	galeriaRepo := new(mocks.MockGaleriaRepo)
	uc := usecases.NewUsuarioUseCase(nil, nil, prestadorRepo, galeriaRepo)

	ctx := context.Background()
	id := uint(1)

	prestador := &model.Prestador{
		IDUsuario:          id,
		Nome:               "P1",
		Telefone:           "111",
		CategoriaPrestador: &model.CategoriaPrestador{Nome: "Cat1"},
	}

	prestadorRepo.On("BuscarPorID", ctx, id).Return(prestador, nil)
	galeriaRepo.On("FindByPrestadorID", ctx, id).Return(&model.Galeria{Imagens: []model.Imagem{{URL: "img.jpg"}}}, nil)

	res, err := uc.BuscarPrestador(ctx, id)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "P1", res.Nome)
	assert.Len(t, res.Galeria, 1)
}

func TestUsuarioUseCase_ListarPrestadoresPorLocalizacao_Success(t *testing.T) {
	prestadorRepo := new(mocks.MockPrestadorRepo)
	galeriaRepo := new(mocks.MockGaleriaRepo)
	uc := usecases.NewUsuarioUseCase(nil, nil, prestadorRepo, galeriaRepo)

	ctx := context.Background()
	prestadores := []model.Prestador{{IDUsuario: 1, Nome: "P1"}}

	prestadorRepo.On("FindByLocation", ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, 10, 0).
		Return(prestadores, nil)
	galeriaRepo.On("FindByPrestadorIDs", ctx, []uint{1}).Return(nil, nil)

	res, err := uc.ListarPrestadoresPorLocalizacao(ctx, 0.0, 0.0, 10.0, nil, "", "", 10, 0)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
}

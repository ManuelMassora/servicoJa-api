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
	galeriaRepo := new(mocks.GaleriaRepoMock)

    uc := usecases.NewUsuarioUseCase(usuarioRepo, clienteRepo, prestadorRepo, galeriaRepo)

    ctx := context.Background()
    req := usecases.UsuarioRequest{
        Nome:     "Admin User",
        Telefone: "123456789",
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
	galeriaRepo := new(mocks.GaleriaRepoMock)

	uc := usecases.NewUsuarioUseCase(usuarioRepo, clienteRepo, prestadorRepo, galeriaRepo)

	ctx := context.Background()
	req := usecases.UsuarioRequest{
		Telefone: "123456789",
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
	galeriaRepo := new(mocks.GaleriaRepoMock)

	uc := usecases.NewUsuarioUseCase(usuarioRepo, clienteRepo, prestadorRepo, galeriaRepo)

	ctx := context.Background()
	req := usecases.UsuarioRequest{
		Nome:     "Cliente Teste",
		Telefone: "987654321",
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
	galeriaRepo := new(mocks.GaleriaRepoMock)

	uc := usecases.NewUsuarioUseCase(usuarioRepo, clienteRepo, prestadorRepo, galeriaRepo)

	ctx := context.Background()
	req := usecases.PrestadorRequest{
		Nome:     "Prestador X",
		Telefone: "111222333",
		Senha:    "senha789",
		Localizacao:   "Maputo",
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
	galeriaRepo := new(mocks.GaleriaRepoMock)

	uc := usecases.NewUsuarioUseCase(usuarioRepo, clienteRepo, prestadorRepo, galeriaRepo)

	ctx := context.Background()
	usuariosMock := []model.Usuario{
		{Nome: "User 1", Telefone: "111"},
		{Nome: "User 2", Telefone: "222"},
	}

	usuarioRepo.On("ListarTodos", ctx, mock.Anything, "", "", 10, 0).
	Return(usuariosMock, nil)

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
	galeriaRepo := new(mocks.GaleriaRepoMock)

	uc := usecases.NewUsuarioUseCase(usuarioRepo, clienteRepo, prestadorRepo, galeriaRepo)

	ctx := context.Background()
	prestadoresMock := []model.Prestador{
		{
			Usuario: model.Usuario{Nome: "Prest 1", Telefone: "333"},
			Localizacao: "Matola",
			StatusDisponivel: true,
		},
	}

	prestadorRepo.On("Listar", ctx, mock.Anything, mock.Anything, "", "", 10, 0).
		Return(prestadoresMock, nil)

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
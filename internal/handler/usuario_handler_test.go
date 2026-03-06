package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ManuelMassora/servicoJa-api/internal/handler"
	"github.com/ManuelMassora/servicoJa-api/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUsuarioUseCase struct {
	mock.Mock
}

func (m *mockUsuarioUseCase) CriarAdmin(ctx context.Context, request usecases.UsuarioRequest) error {
	args := m.Called(ctx, request)
	return args.Error(0)
}

func (m *mockUsuarioUseCase) CriarCliente(ctx context.Context, request usecases.UsuarioRequest) error {
	args := m.Called(ctx, request)
	return args.Error(0)
}

func (m *mockUsuarioUseCase) CriarPrestador(ctx context.Context, request usecases.PrestadorRequest) error {
	args := m.Called(ctx, request)
	return args.Error(0)
}

func (m *mockUsuarioUseCase) EditarPrestador(ctx context.Context, userId uint, campos map[string]interface{}) (*usecases.PrestadorResponse, error) {
	args := m.Called(ctx, userId, campos)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecases.PrestadorResponse), args.Error(1)
}

func (m *mockUsuarioUseCase) BuscarPrestador(ctx context.Context, id uint) (*usecases.PrestadorResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecases.PrestadorResponse), args.Error(1)
}

func (m *mockUsuarioUseCase) BuscarPorID(ctx context.Context, id uint) (*usecases.UsuarioResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecases.UsuarioResponse), args.Error(1)
}

func (m *mockUsuarioUseCase) ListarTodosUsuarios(ctx context.Context, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]usecases.UsuarioResponse, error) {
	args := m.Called(ctx, filters, orderBy, orderDir, limit, offset)
	return args.Get(0).([]usecases.UsuarioResponse), args.Error(1)
}

func (m *mockUsuarioUseCase) ListarPrestadores(ctx context.Context, filters map[string]interface{}, statusDisponivel interface{}, orderBy string, orderDir string, limit, offset int) ([]usecases.PrestadorResponse, error) {
	args := m.Called(ctx, filters, statusDisponivel, orderBy, orderDir, limit, offset)
	return args.Get(0).([]usecases.PrestadorResponse), args.Error(1)
}

func (m *mockUsuarioUseCase) ListarPrestadoresPorLocalizacao(ctx context.Context, latitude, longitude, radius float64, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]usecases.PrestadorResponse, error) {
	args := m.Called(ctx, latitude, longitude, radius, filters, orderBy, orderDir, limit, offset)
	return args.Get(0).([]usecases.PrestadorResponse), args.Error(1)
}

func TestUsuarioHandler_BuscarPrestadorPorID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		uc := new(mockUsuarioUseCase)
		h := handler.NewUsuarioHandler(uc, nil)

		router := gin.New()
		router.GET("/prestadores/:id", h.BuscarPrestadorPorID)

		resp := &usecases.PrestadorResponse{
			ID:   10,
			Nome: "Prestador Teste",
		}

		uc.On("BuscarPrestador", mock.Anything, uint(10)).Return(resp, nil)

		req := httptest.NewRequest("GET", "/prestadores/10", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var actual usecases.PrestadorResponse
		json.Unmarshal(w.Body.Bytes(), &actual)
		assert.Equal(t, resp.Nome, actual.Nome)
		uc.AssertExpectations(t)
	})

	t.Run("not_found", func(t *testing.T) {
		uc := new(mockUsuarioUseCase)
		h := handler.NewUsuarioHandler(uc, nil)

		router := gin.New()
		router.GET("/prestadores/:id", h.BuscarPrestadorPorID)

		uc.On("BuscarPrestador", mock.Anything, uint(10)).Return(nil, errors.New("não encontrado"))

		req := httptest.NewRequest("GET", "/prestadores/10", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		uc.AssertExpectations(t)
	})
}

func TestUsuarioHandler_CriarCliente(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		uc := new(mockUsuarioUseCase)
		h := handler.NewUsuarioHandler(uc, nil)

		router := gin.New()
		router.POST("/clientes", h.CriarCliente)

		uc.On("CriarCliente", mock.Anything, mock.Anything).Return(nil)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		_ = writer.WriteField("nome", "Cliente Teste")
		_ = writer.WriteField("telefone", "841112233")
		_ = writer.WriteField("senha", "senha123")
		writer.Close()

		req := httptest.NewRequest("POST", "/clientes", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		uc.AssertExpectations(t)
	})
}

func TestUsuarioHandler_CriarPrestador(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		uc := new(mockUsuarioUseCase)
		h := handler.NewUsuarioHandler(uc, nil)

		router := gin.New()
		router.POST("/prestadores", h.CriarPrestador)

		uc.On("CriarPrestador", mock.Anything, mock.Anything).Return(nil)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		_ = writer.WriteField("nome", "Prestador Teste")
		_ = writer.WriteField("telefone", "849998877")
		_ = writer.WriteField("senha", "senha123")
		_ = writer.WriteField("bi", "123456789")
		_ = writer.WriteField("nuit", "987654321")
		_ = writer.WriteField("latitude", "-25.0")
		_ = writer.WriteField("longitude", "32.0")
		_ = writer.WriteField("localizacao", "Maputo")
		writer.Close()

		req := httptest.NewRequest("POST", "/prestadores", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		uc.AssertExpectations(t)
	})
}

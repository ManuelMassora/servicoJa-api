package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ManuelMassora/servicoJa-api/internal/handler"
	"github.com/ManuelMassora/servicoJa-api/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockServicoUseCase struct {
	mock.Mock
}

func (m *mockServicoUseCase) FinalizarServico(ctx context.Context, idServico uint, idUsuario uint) error {
	args := m.Called(ctx, idServico, idUsuario)
	return args.Error(0)
}
func (m *mockServicoUseCase) ConfirmarServico(ctx context.Context, idServico uint, idUsuario uint) error {
	args := m.Called(ctx, idServico, idUsuario)
	return args.Error(0)
}
func (m *mockServicoUseCase) CancelarServico(ctx context.Context, idServico uint, idUsuario uint) error {
	args := m.Called(ctx, idServico, idUsuario)
	return args.Error(0)
}
func (m *mockServicoUseCase) ListarPorCliente(ctx context.Context, idUsuario uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]usecases.ServicoResponse, error) {
	args := m.Called(ctx, idUsuario, filters, orderBy, orderDir, limit, offset)
	return args.Get(0).([]usecases.ServicoResponse), args.Error(1)
}
func (m *mockServicoUseCase) ListarPorPrestador(ctx context.Context, idUsuario uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]usecases.ServicoResponse, error) {
	args := m.Called(ctx, idUsuario, filters, orderBy, orderDir, limit, offset)
	return args.Get(0).([]usecases.ServicoResponse), args.Error(1)
}
func (m *mockServicoUseCase) ListarPorLocalizacao(ctx context.Context, idUsuario uint, latitude, longitude, radius float64, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]usecases.ServicoResponse, error) {
	args := m.Called(ctx, idUsuario, latitude, longitude, radius, filters, orderBy, orderDir, limit, offset)
	return args.Get(0).([]usecases.ServicoResponse), args.Error(1)
}

func TestServicoHandler_FinalizarServico(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		uc := new(mockServicoUseCase)
		h := handler.NewServicoHandler(uc)

		router := gin.New()
		router.POST("/servicos/:id/finalizar", func(c *gin.Context) {
			c.Set("userID", uint(1)) // Mock setting user id in context
			h.FinalizarServico(c)
		})

		uc.On("FinalizarServico", mock.Anything, uint(2), uint(1)).Return(nil)

		req := httptest.NewRequest("POST", "/servicos/2/finalizar", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		uc.AssertExpectations(t)
	})
}

func TestServicoHandler_ConfirmarServico(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		uc := new(mockServicoUseCase)
		h := handler.NewServicoHandler(uc)

		router := gin.New()
		router.POST("/servicos/:id/confirmar", func(c *gin.Context) {
			c.Set("userID", uint(1))
			h.ConfirmarServico(c)
		})

		uc.On("ConfirmarServico", mock.Anything, uint(2), uint(1)).Return(nil)

		req := httptest.NewRequest("POST", "/servicos/2/confirmar", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		uc.AssertExpectations(t)
	})
}

func TestServicoHandler_CancelarServico(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		uc := new(mockServicoUseCase)
		h := handler.NewServicoHandler(uc)

		router := gin.New()
		router.POST("/servicos/:id/cancelar", func(c *gin.Context) {
			c.Set("userID", uint(1))
			h.CancelarServico(c)
		})

		uc.On("CancelarServico", mock.Anything, uint(2), uint(1)).Return(nil)

		req := httptest.NewRequest("POST", "/servicos/2/cancelar", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		uc.AssertExpectations(t)
	})
}

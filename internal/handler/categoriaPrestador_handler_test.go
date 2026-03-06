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

type mockCategoriaPrestadorUseCase struct {
	mock.Mock
}

func (m *mockCategoriaPrestadorUseCase) Criar(ctx context.Context, req usecases.CategoriaPrestadorRequest) (*usecases.CategoriaPrestadorResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecases.CategoriaPrestadorResponse), args.Error(1)
}
func (m *mockCategoriaPrestadorUseCase) Editar(ctx context.Context, id uint, campos map[string]interface{}) (*usecases.CategoriaPrestadorResponse, error) {
	args := m.Called(ctx, id, campos)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecases.CategoriaPrestadorResponse), args.Error(1)
}
func (m *mockCategoriaPrestadorUseCase) Listar(ctx context.Context, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]usecases.CategoriaPrestadorResponse, error) {
	args := m.Called(ctx, filters, orderBy, orderDir, limit, offset)
	return args.Get(0).([]usecases.CategoriaPrestadorResponse), args.Error(1)
}

func TestCategoriaPrestadorHandler_Listar(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		uc := new(mockCategoriaPrestadorUseCase)
		h := handler.NewCategoriaPrestadorHandler(uc)

		router := gin.New()
		router.GET("/categorias-prestadores", h.Listar)

		respData := []usecases.CategoriaPrestadorResponse{
			{ID: 1},
		}

		uc.On("Listar", mock.Anything, mock.Anything, "", "", 10, 0).Return(respData, nil)

		req := httptest.NewRequest("GET", "/categorias-prestadores", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		uc.AssertExpectations(t)
	})
}

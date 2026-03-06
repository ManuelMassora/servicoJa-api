package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ManuelMassora/servicoJa-api/internal/handler"
	"github.com/ManuelMassora/servicoJa-api/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockCategoriaUseCase struct {
	mock.Mock
}

func (m *mockCategoriaUseCase) Criar(ctx context.Context, request usecases.CategoriaRequest) (uint, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(uint), args.Error(1)
}
func (m *mockCategoriaUseCase) Editar(ctx context.Context, id uint, campos map[string]interface{}) error {
	args := m.Called(ctx, id, campos)
	return args.Error(0)
}
func (m *mockCategoriaUseCase) Listar(ctx context.Context, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]usecases.CategoriaResponse, error) {
	args := m.Called(ctx, filters, orderBy, orderDir, limit, offset)
	return args.Get(0).([]usecases.CategoriaResponse), args.Error(1)
}
func (m *mockCategoriaUseCase) BuscarPorID(ctx context.Context, id uint) (*usecases.CategoriaResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecases.CategoriaResponse), args.Error(1)
}

func TestCategoriaHandler_Criar(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		uc := new(mockCategoriaUseCase)
		h := handler.NewCategoriaHandler(uc)

		router := gin.New()
		router.POST("/categorias", h.Criar)

		reqBody := usecases.CategoriaRequest{Nome: "Beleza", Descricao: "Cuidados pessoais"}
		jsonReq, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/categorias", bytes.NewBuffer(jsonReq))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		uc.On("Criar", mock.Anything, reqBody).Return(uint(1), nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, float64(1), resp["id"])
		uc.AssertExpectations(t)
	})
}

func TestCategoriaHandler_Listar(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		uc := new(mockCategoriaUseCase)
		h := handler.NewCategoriaHandler(uc)

		router := gin.New()
		router.GET("/categorias", h.Listar)

		respData := []usecases.CategoriaResponse{
			{ID: 1, Nome: "Beleza"},
		}

		uc.On("Listar", mock.Anything, mock.Anything, "", "", 10, 0).Return(respData, nil)

		req := httptest.NewRequest("GET", "/categorias", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		uc.AssertExpectations(t)
	})
}

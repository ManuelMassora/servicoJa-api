package handler_test

import (
	"bytes"
	"context"
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

type mockCatalogoUseCase struct {
	mock.Mock
}

func (m *mockCatalogoUseCase) Criar(ctx context.Context, req usecases.RequestCreateCatalogo, prestadorID uint) (uint, error) {
	args := m.Called(ctx, req, prestadorID)
	return args.Get(0).(uint), args.Error(1)
}
func (m *mockCatalogoUseCase) Editar(ctx context.Context, id uint, prestadorID uint, campos map[string]interface{}) error {
	args := m.Called(ctx, id, prestadorID, campos)
	return args.Error(0)
}
func (m *mockCatalogoUseCase) Apagar(ctx context.Context, id uint, prestadorID uint) error {
	args := m.Called(ctx, id, prestadorID)
	return args.Error(0)
}
func (m *mockCatalogoUseCase) Listar(ctx context.Context, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]usecases.ResponseCatalogo, error) {
	args := m.Called(ctx, filters, orderBy, orderDir, limit, offset)
	return args.Get(0).([]usecases.ResponseCatalogo), args.Error(1)
}
func (m *mockCatalogoUseCase) ListarPorPrestador(ctx context.Context, prestadorID uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]usecases.ResponseCatalogo, error) {
	args := m.Called(ctx, prestadorID, filters, orderBy, orderDir, limit, offset)
	return args.Get(0).([]usecases.ResponseCatalogo), args.Error(1)
}
func (m *mockCatalogoUseCase) ListarPorLocalizacao(ctx context.Context, latitude, longitude, radius float64, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]usecases.ResponseCatalogo, error) {
	args := m.Called(ctx, latitude, longitude, radius, filters, orderBy, orderDir, limit, offset)
	return args.Get(0).([]usecases.ResponseCatalogo), args.Error(1)
}

func TestCatalogoHandler_Listar(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		uc := new(mockCatalogoUseCase)
		h := handler.NewCatalogoHandler(uc, nil)

		router := gin.New()
		router.GET("/catalogos", h.Listar)

		respData := []usecases.ResponseCatalogo{
			{ID: 1, Nome: "Cabelo"},
		}

		uc.On("Listar", mock.Anything, mock.Anything, "", "", 10, 0).Return(respData, nil)

		req := httptest.NewRequest("GET", "/catalogos", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		uc.AssertExpectations(t)
	})
}

func TestCatalogoHandler_Criar(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		uc := new(mockCatalogoUseCase)
		h := handler.NewCatalogoHandler(uc, nil)

		router := gin.New()
		router.POST("/catalogos", func(c *gin.Context) {
			c.Set("userID", uint(1))
			h.Criar(c)
		})

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		_ = writer.WriteField("nome", "Corte de Cabelo")
		_ = writer.WriteField("descricao", "Corte simples")
		_ = writer.WriteField("tipo_preco", "fixo")
		_ = writer.WriteField("valor_fixo", "500")
		_ = writer.WriteField("categoria_id", "1")
		_ = writer.WriteField("localizacao", "Maputo")
		_ = writer.WriteField("latitude", "-25.0")
		_ = writer.WriteField("longitude", "32.0")
		writer.Close()

		req := httptest.NewRequest("POST", "/catalogos", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		uc.On("Criar", mock.Anything, mock.Anything, uint(1)).Return(uint(1), nil)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
		}
		assert.Equal(t, http.StatusCreated, w.Code)
		uc.AssertExpectations(t)
	})
}

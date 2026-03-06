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

type mockAvaliacaoUseCase struct {
	mock.Mock
}

func (m *mockAvaliacaoUseCase) Criar(ctx context.Context, req usecases.AvaliacaoRequest, idAvaliador uint) error {
	args := m.Called(ctx, req, idAvaliador)
	return args.Error(0)
}
func (m *mockAvaliacaoUseCase) ListarPorCliente(ctx context.Context, idCliente uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]usecases.AvaliacaoResponse, error) {
	args := m.Called(ctx, idCliente, filters, orderBy, orderDir, limit, offset)
	return args.Get(0).([]usecases.AvaliacaoResponse), args.Error(1)
}
func (m *mockAvaliacaoUseCase) ListarPorPrestador(ctx context.Context, idPrestador uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]usecases.AvaliacaoResponse, error) {
	args := m.Called(ctx, idPrestador, filters, orderBy, orderDir, limit, offset)
	return args.Get(0).([]usecases.AvaliacaoResponse), args.Error(1)
}
func (m *mockAvaliacaoUseCase) MediaPorPrestador(ctx context.Context, idPrestador uint) (float64, error) {
	args := m.Called(ctx, idPrestador)
	return args.Get(0).(float64), args.Error(1)
}

func TestAvaliacaoHandler_MediaPorPrestador(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		uc := new(mockAvaliacaoUseCase)
		h := handler.NewAvaliacaoHandler(uc)

		router := gin.New()
		router.GET("/avaliacoes/prestador/:id/media", h.MediaPorPrestador)

		uc.On("MediaPorPrestador", mock.Anything, uint(10)).Return(4.5, nil)

		req := httptest.NewRequest("GET", "/avaliacoes/prestador/10/media", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		uc.AssertExpectations(t)
	})
}

func TestAvaliacaoHandler_Criar(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		uc := new(mockAvaliacaoUseCase)
		h := handler.NewAvaliacaoHandler(uc)

		router := gin.New()
		router.POST("/avaliacoes", func(c *gin.Context) {
			c.Set("userID", uint(1))
			h.CriarAvaliacao(c)
		})

		reqData := usecases.AvaliacaoRequest{
			IDServico:  2,
			Pontuacao:  5,
			Comentario: "Excelente servico!",
		}
		jsonReq, _ := json.Marshal(reqData)
		req := httptest.NewRequest("POST", "/avaliacoes", bytes.NewBuffer(jsonReq))
		req.Header.Set("Content-Type", "application/json")

		uc.On("Criar", mock.Anything, reqData, uint(1)).Return(nil)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		uc.AssertExpectations(t)
	})
}

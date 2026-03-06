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

type mockPropostaUseCase struct {
	mock.Mock
}

func (m *mockPropostaUseCase) Criar(ctx context.Context, req usecases.PropostaRequest, idUsuario uint) error {
	args := m.Called(ctx, req, idUsuario)
	return args.Error(0)
}
func (m *mockPropostaUseCase) Responder(ctx context.Context, idProposta uint, idUsuario uint, aceitar bool) error {
	args := m.Called(ctx, idProposta, idUsuario, aceitar)
	return args.Error(0)
}
func (m *mockPropostaUseCase) Cancelar(ctx context.Context, idProposta uint, idUsuario uint) error {
	args := m.Called(ctx, idProposta, idUsuario)
	return args.Error(0)
}
func (m *mockPropostaUseCase) ListarPorVaga(ctx context.Context, idUsuario uint, idVaga uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]usecases.PropostaResponse, error) {
	args := m.Called(ctx, idUsuario, idVaga, filters, orderBy, orderDir, limit, offset)
	return args.Get(0).([]usecases.PropostaResponse), args.Error(1)
}
func (m *mockPropostaUseCase) ListarPorPrestador(ctx context.Context, idUsuario uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]usecases.PropostaResponse, error) {
	args := m.Called(ctx, idUsuario, filters, orderBy, orderDir, limit, offset)
	return args.Get(0).([]usecases.PropostaResponse), args.Error(1)
}

func TestPropostaHandler_Responder(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success_aceitar", func(t *testing.T) {
		uc := new(mockPropostaUseCase)
		h := handler.NewPropostaHandler(uc)

		router := gin.New()
		router.POST("/propostas/:id/responder", func(c *gin.Context) {
			c.Set("userID", uint(1))
			h.Responder(c)
		})

		uc.On("Responder", mock.Anything, uint(10), uint(1), true).Return(nil)

		req := httptest.NewRequest("POST", "/propostas/10/responder?status=aceitar", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		uc.AssertExpectations(t)
	})
}

func TestPropostaHandler_Criar(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		uc := new(mockPropostaUseCase)
		h := handler.NewPropostaHandler(uc)

		router := gin.New()
		router.POST("/propostas", func(c *gin.Context) {
			c.Set("userID", uint(1))
			h.Criar(c)
		})

		reqData := usecases.PropostaRequest{
			IDVaga:        2,
			ValorProposto: 1000,
			Mensagem:      "Eu faco esse servico",
			PrazoEstimado: "2 dias",
		}
		jsonReq, _ := json.Marshal(reqData)
		req := httptest.NewRequest("POST", "/propostas", bytes.NewBuffer(jsonReq))
		req.Header.Set("Content-Type", "application/json")

		uc.On("Criar", mock.Anything, reqData, uint(1)).Return(nil)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		uc.AssertExpectations(t)
	})
}

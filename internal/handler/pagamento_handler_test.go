package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ManuelMassora/servicoJa-api/internal/handler"
	gatewaympesa "github.com/ManuelMassora/servicoJa-api/internal/infra/gateway_mpesa"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockPagamentoUseCase struct {
	mock.Mock
}

func (m *mockPagamentoUseCase) ProcessarCallbackMpesa(ctx context.Context, payload gatewaympesa.MpesaCallbackPayload) error {
	args := m.Called(ctx, payload)
	return args.Error(0)
}
func (m *mockPagamentoUseCase) ProcessarQuerySimulada(ctx context.Context, referencia string) error {
	args := m.Called(ctx, referencia)
	return args.Error(0)
}

func TestPagamentoHandler_SimularQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		uc := new(mockPagamentoUseCase)
		h := handler.NewPagamentoHandler(uc)

		router := gin.New()
		router.GET("/pagamentos/simular", h.SimularQuery)

		uc.On("ProcessarQuerySimulada", mock.Anything, "REF123").Return(nil)

		req := httptest.NewRequest("GET", "/pagamentos/simular?ref=REF123", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		uc.AssertExpectations(t)
	})
}

func TestPagamentoHandler_ReceiveCallback(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		uc := new(mockPagamentoUseCase)
		h := handler.NewPagamentoHandler(uc)

		router := gin.New()
		router.POST("/pagamentos/callback", h.ReceiveCallback)

		payload := gatewaympesa.MpesaCallbackPayload{
			ThirdPartyReference: "REF123",
			ResponseCode:        "INS-0",
		}
		jsonReq, _ := json.Marshal(payload)

		uc.On("ProcessarCallbackMpesa", mock.Anything, payload).Return(nil)

		req := httptest.NewRequest("POST", "/pagamentos/callback", bytes.NewBuffer(jsonReq))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		uc.AssertExpectations(t)
	})
}

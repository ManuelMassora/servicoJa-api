package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ManuelMassora/servicoJa-api/internal/handler"
	"github.com/ManuelMassora/servicoJa-api/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAuthUseCase struct {
	mock.Mock
}

func (m *mockAuthUseCase) Authenticate(ctx context.Context, request services.AuthRequest) (*services.AuthResponse, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.AuthResponse), args.Error(1)
}

func TestAuthHandler_Authenticate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		uc := new(mockAuthUseCase)
		h := handler.NewAuthHandler(uc)

		router := gin.New()
		router.POST("/auth", h.Authenticate)

		reqBody := services.AuthRequest{
			Telefone: "841112233",
			Senha:    "senha123",
		}
		jsonReq, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/auth", bytes.NewBuffer(jsonReq))

		resp := &services.AuthResponse{
			ID:    1,
			Role:  "CLIENTE",
			Token: "valid-token",
		}

		uc.On("Authenticate", mock.Anything, reqBody).Return(resp, nil)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusAccepted, w.Code)
		var actual services.AuthResponse
		json.Unmarshal(w.Body.Bytes(), &actual)
		assert.Equal(t, resp.Token, actual.Token)
		uc.AssertExpectations(t)
	})

	t.Run("invalid_credentials", func(t *testing.T) {
		uc := new(mockAuthUseCase)
		h := handler.NewAuthHandler(uc)

		router := gin.New()
		router.POST("/auth", h.Authenticate)

		reqBody := services.AuthRequest{
			Telefone: "841112233",
			Senha:    "errada",
		}
		jsonReq, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/auth", bytes.NewBuffer(jsonReq))

		uc.On("Authenticate", mock.Anything, reqBody).Return(nil, errors.New("credenciais inválidas"))

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		uc.AssertExpectations(t)
	})

	t.Run("bad_request_json", func(t *testing.T) {
		uc := new(mockAuthUseCase)
		h := handler.NewAuthHandler(uc)

		router := gin.New()
		router.POST("/auth", h.Authenticate)

		req := httptest.NewRequest("POST", "/auth", bytes.NewBufferString("invalid json"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

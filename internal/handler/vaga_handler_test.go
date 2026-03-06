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

type mockVagaUseCase struct {
	mock.Mock
}

func (m *mockVagaUseCase) CriarVaga(ctx context.Context, req usecases.VagaRequest, idUsuario uint) (uint, error) {
	args := m.Called(ctx, req, idUsuario)
	return args.Get(0).(uint), args.Error(1)
}
func (m *mockVagaUseCase) CancelarVaga(ctx context.Context, idVaga uint, idUsuario uint) error {
	args := m.Called(ctx, idVaga, idUsuario)
	return args.Error(0)
}
func (m *mockVagaUseCase) ListarVagasDisponiveis(ctx context.Context, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]usecases.VagaResponse, error) {
	args := m.Called(ctx, filters, orderBy, orderDir, limit, offset)
	return args.Get(0).([]usecases.VagaResponse), args.Error(1)
}
func (m *mockVagaUseCase) ListarPorCliente(ctx context.Context, idUsuario uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]usecases.VagaResponse, error) {
	args := m.Called(ctx, idUsuario, filters, orderBy, orderDir, limit, offset)
	return args.Get(0).([]usecases.VagaResponse), args.Error(1)
}
func (m *mockVagaUseCase) ListarPorLocalizacao(ctx context.Context, latitude, longitude, radius float64, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]usecases.VagaResponse, error) {
	args := m.Called(ctx, latitude, longitude, radius, filters, orderBy, orderDir, limit, offset)
	return args.Get(0).([]usecases.VagaResponse), args.Error(1)
}
func (m *mockVagaUseCase) BuscarPorIDIfCliente(ctx context.Context, id, idUsuario uint) (*usecases.VagaResponse, error) {
	args := m.Called(ctx, id, idUsuario)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecases.VagaResponse), args.Error(1)
}

func TestVagaHandler_ListarVagasDisponiveis(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		uc := new(mockVagaUseCase)
		h := handler.NewVagaHandler(uc, nil)

		router := gin.New()
		router.GET("/vagas", h.ListarVagasDisponiveis)

		respData := []usecases.VagaResponse{
			{ID: 1, Titulo: "Vaga 1"},
		}

		uc.On("ListarVagasDisponiveis", mock.Anything, mock.Anything, "", "", 10, 0).Return(respData, nil)

		req := httptest.NewRequest("GET", "/vagas", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		uc.AssertExpectations(t)
	})
}

func TestVagaHandler_CriarVaga(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		uc := new(mockVagaUseCase)
		h := handler.NewVagaHandler(uc, nil)

		router := gin.New()
		router.POST("/vagas", func(c *gin.Context) {
			c.Set("userID", uint(1))
			h.CriarVaga(c)
		})

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		_ = writer.WriteField("titulo", "Preciso de carpinteiro")
		_ = writer.WriteField("descricao", "Mesa partida")
		_ = writer.WriteField("preco", "2000")
		_ = writer.WriteField("localizacao", "Maputo")
		_ = writer.WriteField("latitude", "-25.0")
		_ = writer.WriteField("longitude", "32.0")
		_ = writer.WriteField("telefone_pagamento", "841112233")
		writer.Close()

		req := httptest.NewRequest("POST", "/vagas", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		uc.On("CriarVaga", mock.Anything, mock.Anything, uint(1)).Return(uint(1), nil)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
		}
		assert.Equal(t, http.StatusCreated, w.Code)
		uc.AssertExpectations(t)
	})
}

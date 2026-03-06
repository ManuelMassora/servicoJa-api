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
	"time"

	"github.com/ManuelMassora/servicoJa-api/internal/handler"
	"github.com/ManuelMassora/servicoJa-api/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAgendamentoUC struct {
	mock.Mock
}

func (m *mockAgendamentoUC) Criar(ctx context.Context, req *usecases.AgendamentoRequest, idCliente uint) (uint, error) {
	args := m.Called(ctx, req, idCliente)
	return args.Get(0).(uint), args.Error(1)
}

func (m *mockAgendamentoUC) Buscar(ctx context.Context, id uint, idUsuario uint) (*usecases.AgendamentoResponse, error) {
	args := m.Called(ctx, id, idUsuario)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecases.AgendamentoResponse), args.Error(1)
}

func (m *mockAgendamentoUC) Aceitar(ctx context.Context, id uint, idUsuario uint) (uint, error) {
	args := m.Called(ctx, id, idUsuario)
	return args.Get(0).(uint), args.Error(1)
}

func (m *mockAgendamentoUC) Recusar(ctx context.Context, id uint, idUsuario uint) error {
	args := m.Called(ctx, id, idUsuario)
	return args.Error(0)
}

func (m *mockAgendamentoUC) Cancelar(ctx context.Context, id uint, idUsuario uint) error {
	args := m.Called(ctx, id, idUsuario)
	return args.Error(0)
}

func (m *mockAgendamentoUC) Listar(ctx context.Context, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]usecases.AgendamentoResponse, error) {
	args := m.Called(ctx, filters, orderBy, orderDir, limit, offset)
	return args.Get(0).([]usecases.AgendamentoResponse), args.Error(1)
}

func (m *mockAgendamentoUC) ListarPorClienteID(ctx context.Context, idUsuario uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]usecases.AgendamentoGroupCategoriaResponse, error) {
	args := m.Called(ctx, idUsuario, filters, orderBy, orderDir, limit, offset)
	return args.Get(0).([]usecases.AgendamentoGroupCategoriaResponse), args.Error(1)
}

func (m *mockAgendamentoUC) ListarPorPrestadorIDAgrupado(ctx context.Context, idUsuario uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]usecases.AgendamentoGroupCategoriaResponse, error) {
	args := m.Called(ctx, idUsuario, filters, orderBy, orderDir, limit, offset)
	return args.Get(0).([]usecases.AgendamentoGroupCategoriaResponse), args.Error(1)
}

func (m *mockAgendamentoUC) ListarPorCatalogID(ctx context.Context, idUsuario, idCatalogo uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]usecases.AgendamentoResponse, error) {
	args := m.Called(ctx, idUsuario, idCatalogo, filters, orderBy, orderDir, limit, offset)
	return args.Get(0).([]usecases.AgendamentoResponse), args.Error(1)
}

func (m *mockAgendamentoUC) ListarPorLocalizacao(ctx context.Context, idUsuario uint, latitude, longitude, radius float64, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]usecases.AgendamentoResponse, error) {
	args := m.Called(ctx, idUsuario, latitude, longitude, radius, filters, orderBy, orderDir, limit, offset)
	return args.Get(0).([]usecases.AgendamentoResponse), args.Error(1)
}

func TestAgendamentoHandler_Buscar(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		uc := new(mockAgendamentoUC)
		h := handler.NewAgendamentoHandler(uc, nil)

		router := gin.New()
		router.GET("/agendamentos/:id", func(c *gin.Context) {
			c.Set("userID", uint(1))
			h.Buscar(c)
		})

		resp := &usecases.AgendamentoResponse{
			ID:      10,
			Detalhe: "Teste",
			Status:  "PENDENTE",
		}

		uc.On("Buscar", mock.Anything, uint(10), uint(1)).Return(resp, nil)

		req := httptest.NewRequest("GET", "/agendamentos/10", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var actual usecases.AgendamentoResponse
		json.Unmarshal(w.Body.Bytes(), &actual)
		assert.Equal(t, resp.ID, actual.ID)
		uc.AssertExpectations(t)
	})

	t.Run("invalid_id", func(t *testing.T) {
		uc := new(mockAgendamentoUC)
		h := handler.NewAgendamentoHandler(uc, nil)

		router := gin.New()
		router.GET("/agendamentos/:id", h.Buscar)

		req := httptest.NewRequest("GET", "/agendamentos/invalid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("not_found", func(t *testing.T) {
		uc := new(mockAgendamentoUC)
		h := handler.NewAgendamentoHandler(uc, nil)

		router := gin.New()
		router.GET("/agendamentos/:id", func(c *gin.Context) {
			c.Set("userID", uint(1))
			h.Buscar(c)
		})

		uc.On("Buscar", mock.Anything, uint(10), uint(1)).Return(nil, errors.New("agendamento não encontrado"))

		req := httptest.NewRequest("GET", "/agendamentos/10", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		uc.AssertExpectations(t)
	})
}

func TestAgendamentoHandler_Criar(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		uc := new(mockAgendamentoUC)
		h := handler.NewAgendamentoHandler(uc, nil)

		router := gin.New()
		router.POST("/agendamentos", func(c *gin.Context) {
			c.Set("userID", uint(1))
			h.Criar(c)
		})

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		_ = writer.WriteField("detalhe", "Teste")
		_ = writer.WriteField("id_catalogo", "2")
		_ = writer.WriteField("datahora", time.Now().Format(time.RFC3339))
		_ = writer.WriteField("localizacao", "Maputo")
		_ = writer.WriteField("latitude", "-25.1")
		_ = writer.WriteField("longitude", "32.1")
		_ = writer.WriteField("telefone_pagamento", "841112233")
		writer.Close()

		req := httptest.NewRequest("POST", "/agendamentos", body)
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

func TestAgendamentoHandler_Aceitar(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		uc := new(mockAgendamentoUC)
		h := handler.NewAgendamentoHandler(uc, nil)

		router := gin.New()
		router.POST("/agendamentos/:id/aceitar", func(c *gin.Context) {
			c.Set("userID", uint(1))
			h.Aceitar(c)
		})

		uc.On("Aceitar", mock.Anything, uint(10), uint(1)).Return(uint(100), nil)

		req := httptest.NewRequest("POST", "/agendamentos/10/aceitar", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var actual map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &actual)
		assert.Equal(t, float64(100), actual["id_servico"])
		uc.AssertExpectations(t)
	})
}

func TestAgendamentoHandler_Recusar(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		uc := new(mockAgendamentoUC)
		h := handler.NewAgendamentoHandler(uc, nil)

		router := gin.New()
		router.POST("/agendamentos/:id/recusar", func(c *gin.Context) {
			c.Set("userID", uint(1))
			h.Recusar(c)
		})

		uc.On("Recusar", mock.Anything, uint(10), uint(1)).Return(nil)

		req := httptest.NewRequest("POST", "/agendamentos/10/recusar", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		uc.AssertExpectations(t)
	})
}

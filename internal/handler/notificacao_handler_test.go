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

type mockNotificacaoUseCase struct {
	mock.Mock
}

func (m *mockNotificacaoUseCase) ListarPorUsuario(ctx context.Context, idUsuario uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]usecases.NotificacaoResponse, error) {
	args := m.Called(ctx, idUsuario, filters, orderBy, orderDir, limit, offset)
	return args.Get(0).([]usecases.NotificacaoResponse), args.Error(1)
}
func (m *mockNotificacaoUseCase) MarcarComoLida(ctx context.Context, idNotificacao uint, idUsuario uint) error {
	args := m.Called(ctx, idNotificacao, idUsuario)
	return args.Error(0)
}
func (m *mockNotificacaoUseCase) MarcarTodasComoLidas(ctx context.Context, idUsuario uint) error {
	args := m.Called(ctx, idUsuario)
	return args.Error(0)
}

func TestNotificacaoHandler_ListarPorUsuario(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		uc := new(mockNotificacaoUseCase)
		h := handler.NewNotificacaoHandler(uc)

		router := gin.New()
		router.GET("/notificacoes", func(c *gin.Context) {
			c.Set("userID", uint(1))
			h.ListarPorUsuario(c)
		})

		respData := []usecases.NotificacaoResponse{
			{ID: 1, Titulo: "Notif 1"},
		}

		uc.On("ListarPorUsuario", mock.Anything, uint(1), mock.Anything, "", "", 10, 0).Return(respData, nil)

		req := httptest.NewRequest("GET", "/notificacoes", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		uc.AssertExpectations(t)
	})
}

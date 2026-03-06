package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ManuelMassora/servicoJa-api/internal/dto"
	"github.com/ManuelMassora/servicoJa-api/internal/handler"
	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockGaleriaUseCase struct {
	mock.Mock
}

func (m *mockGaleriaUseCase) AddImagesToGaleria(ctx context.Context, prestadorID uint, input dto.GaleriaInput) (*model.Galeria, error) {
	args := m.Called(ctx, prestadorID, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Galeria), args.Error(1)
}

func TestGaleriaHandler_AddImagesToGaleria_NoImages(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("no_images_error", func(t *testing.T) {
		uc := new(mockGaleriaUseCase)
		h := handler.NewGaleriaHandler(uc, nil)

		router := gin.New()
		router.POST("/galeria", h.CriarGaleria)

		req := httptest.NewRequest("POST", "/galeria", nil)
		// No multipart form data
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

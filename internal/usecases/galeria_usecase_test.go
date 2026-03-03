package usecases_test

import (
	"context"
	"testing"

	"github.com/ManuelMassora/servicoJa-api/internal/dto"
	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"github.com/ManuelMassora/servicoJa-api/internal/usecases"
	"github.com/ManuelMassora/servicoJa-api/internal/usecases/usecases_test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGaleriaUseCase_AddImagesToGaleria_Success(t *testing.T) {
	repo := new(mocks.MockGaleriaRepo)
	uc := usecases.NewGaleriaUseCase(repo)

	ctx := context.Background()
	prestadorID := uint(1)
	input := dto.GaleriaInput{Imagens: []string{"img1.jpg", "img2.jpg"}}

	galeria := &model.Galeria{BaseModel: model.BaseModel{ID: 10}, PrestadorID: prestadorID}

	repo.On("FindByPrestadorID", ctx, prestadorID).Return(galeria, nil)
	repo.On("CountImages", ctx, uint(10)).Return(int64(0), nil)
	repo.On("AddImage", ctx, mock.Anything).Return(nil).Times(2)

	res, err := uc.AddImagesToGaleria(ctx, prestadorID, input)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	repo.AssertExpectations(t)
}

func TestGaleriaUseCase_AddImagesToGaleria_LimitExceeded(t *testing.T) {
	repo := new(mocks.MockGaleriaRepo)
	uc := usecases.NewGaleriaUseCase(repo)

	ctx := context.Background()
	prestadorID := uint(1)
	input := dto.GaleriaInput{Imagens: []string{"img1.jpg", "img2.jpg", "img3.jpg"}}

	galeria := &model.Galeria{BaseModel: model.BaseModel{ID: 10}, PrestadorID: prestadorID}

	repo.On("FindByPrestadorID", ctx, prestadorID).Return(galeria, nil)
	repo.On("CountImages", ctx, uint(10)).Return(int64(2), nil)

	res, err := uc.AddImagesToGaleria(ctx, prestadorID, input)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "a galeria pode ter no máximo 4 imagens", err.Error())
}

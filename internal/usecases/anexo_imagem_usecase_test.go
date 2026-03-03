package usecases_test

import (
	"context"
	"testing"

	"github.com/ManuelMassora/servicoJa-api/internal/dto"
	"github.com/ManuelMassora/servicoJa-api/internal/usecases"
	"github.com/ManuelMassora/servicoJa-api/internal/usecases/usecases_test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAnexoImagemUseCase_CreateAnexoImagem_Success(t *testing.T) {
	repo := new(mocks.MockAnexoImagemRepo)
	uc := usecases.NewAnexoImagemUseCase(repo)

	ctx := context.Background()
	input := dto.AnexoImagemInput{URL: "test.jpg"}

	repo.On("Create", ctx, mock.Anything).Return(nil)

	res, err := uc.CreateAnexoImagem(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "test.jpg", res.URL)
}

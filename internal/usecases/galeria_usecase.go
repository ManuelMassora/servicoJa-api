package usecases

import (
	"context"

	"github.com/ManuelMassora/servicoJa-api/internal/dto"
	"github.com/ManuelMassora/servicoJa-api/internal/infra/repo"
	"github.com/ManuelMassora/servicoJa-api/internal/model"
)

type GaleriaUseCase struct {
	repo *repo.GaleriaRepo
}

func NewGaleriaUseCase(repo *repo.GaleriaRepo) *GaleriaUseCase {
	return &GaleriaUseCase{repo: repo}
}

func (uc *GaleriaUseCase) CreateGaleria(ctx context.Context, input dto.GaleriaInput) (*model.Galeria, error) {
	galeria := &model.Galeria{
		PrestadorID: input.PrestadorID,
	}

	_, err := uc.repo.Create(ctx, galeria)
	if err != nil {
		return nil, err
	}

	for _, imagemInput := range input.Imagens {
		imagem := &model.Imagem{
			URL:       imagemInput.URL,
			GaleriaID: galeria.ID,
		}
		err := uc.repo.AddImage(ctx, imagem)
		if err != nil {
			return nil, err
		}
	}

	return galeria, nil
}

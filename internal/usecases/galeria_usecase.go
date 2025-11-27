package usecases

import (
	"context"
	"errors"

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

func (uc *GaleriaUseCase) AddImagesToGaleria(ctx context.Context, prestadorID uint, galeriaImagem dto.GaleriaInput) (*model.Galeria, error) {

    galeria, err := uc.repo.FindByPrestadorID(ctx, prestadorID)
    if err != nil {
        return nil, err
    }

    // Caso 1: Prestador ainda não tem galeria → criar automaticamente
    if galeria == nil {
        galeria = &model.Galeria{
            PrestadorID: prestadorID,
        }

        if _, err := uc.repo.Create(ctx, galeria); err != nil {
            return nil, err
        }
    }

    // Busca quantas imagens já existem na galeria
    totalExistentes, err := uc.repo.CountImages(ctx, galeria.ID)
    if err != nil {
        return nil, err
    }

    // Valida limite de imagens
    if totalExistentes+int64(len(galeriaImagem.Imagens)) > 4 {
        return nil, errors.New("a galeria pode ter no máximo 4 imagens")
    }

    // Adiciona as novas imagens
    for _, imgInput := range galeriaImagem.Imagens {
        img := &model.Imagem{
            URL:       imgInput.URL,
            GaleriaID: galeria.ID,
        }

        if err := uc.repo.AddImage(ctx, img); err != nil {
            return nil, err
        }
    }

    return galeria, nil
}
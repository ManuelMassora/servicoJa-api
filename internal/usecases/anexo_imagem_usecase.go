package usecases

import (
	"context"

	"github.com/ManuelMassora/servicoJa-api/internal/dto"
	"github.com/ManuelMassora/servicoJa-api/internal/infra/repo"
	"github.com/ManuelMassora/servicoJa-api/internal/model"
)

type AnexoImagemUseCase struct {
	repo *repo.AnexoImagemRepo
}

func NewAnexoImagemUseCase(repo *repo.AnexoImagemRepo) *AnexoImagemUseCase {
	return &AnexoImagemUseCase{repo: repo}
}

func (uc *AnexoImagemUseCase) CreateAnexoImagem(ctx context.Context, input dto.AnexoImagemInput) (*model.AnexoImagem, error) {
	anexo := &model.AnexoImagem{
		URL:           input.URL,
		AgendamentoID: input.AgendamentoID,
		VagaID:        input.VagaID,
		CatalogoID:    input.CatalogoID,
	}

	err := uc.repo.Create(ctx, anexo)
	if err != nil {
		return nil, err
	}

	return anexo, nil
}

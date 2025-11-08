package usecases

import (
	"context"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
)

type CategoriaUseCase struct {
	CategoriaRepo model.CategoriaRepo
}

func NewCategoriaUseCase(CategoriaRepo model.CategoriaRepo) *CategoriaUseCase {
	return &CategoriaUseCase{CategoriaRepo: CategoriaRepo}
}

type CategoriaRequest struct {
	Nome      string `json:"nome" binding:"required"`
	Descricao string `json:"descricao" binding:"required"`
}

type CategoriaResponse struct {
	ID			int64
	Nome      string `json:"nome"`
	Descricao string `json:"descricao"`
}

func (uc *CategoriaUseCase) Criar(ctx context.Context, request CategoriaRequest) error {
	categoria := &model.Categoria{
		Nome:      request.Nome,
		Descricao: request.Descricao,
	}
	if err := uc.CategoriaRepo.Criar(ctx, categoria); err != nil {
		return err
	}
	return nil
}

func (uc *CategoriaUseCase) Listar(ctx context.Context, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]CategoriaResponse, error) {
	categorias, err := uc.CategoriaRepo.Listar(ctx, filters, orderBy, orderDir, limit, offset)
	if err != nil {
		return nil, err
	}

	resp := make([]CategoriaResponse, 0, len(categorias))
	for _, c := range categorias {
		resp = append(resp, CategoriaResponse{
			ID:        c.ID,
			Nome:      c.Nome,
			Descricao: c.Descricao,
		})
	}
	return resp, nil
}

func (uc *CategoriaUseCase) BuscarPorID(ctx context.Context, id int64) (*CategoriaResponse, error) {
	categoria, err := uc.CategoriaRepo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, err
	}
	if categoria == nil {
		return nil, nil
	}
	return &CategoriaResponse{
		ID:        categoria.ID,
		Nome:      categoria.Nome,
		Descricao: categoria.Descricao,
	}, nil
}
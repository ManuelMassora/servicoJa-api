package usecases

import (
	"context"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
)

type CategoriaPrestadorUsecase struct {
	repo model.CategoriaPrestadorRepo
}

func NewCategoriaPrestadorUsecase(repo model.CategoriaPrestadorRepo) *CategoriaPrestadorUsecase {
	return &CategoriaPrestadorUsecase{repo: repo}
}

type CategoriaPrestadorRequest struct {
	Nome      string `json:"nome" binding:"required"`
	Descricao string `json:"descricao" binding:"required"`
	Icone     string `json:"icone" binding:"required"`
}

type CategoriaPrestadorResponse struct {
	ID          uint   `json:"id"`
	Nome        string `json:"nome"`
	Descricao   string `json:"descricao"`
	Icone       string `json:"icone"`
}

func (uc *CategoriaPrestadorUsecase) Criar(ctx context.Context, request CategoriaPrestadorRequest) (*CategoriaPrestadorResponse, error) {
	categoria := &model.CategoriaPrestador{
		Nome:      request.Nome,
		Descricao: request.Descricao,
		Icone:     request.Icone,
	}

	createdCategoria, err := uc.repo.Criar(ctx, categoria)
	if err != nil {
		return nil, err
	}

	return &CategoriaPrestadorResponse{
		ID:        createdCategoria.ID,
		Nome:      createdCategoria.Nome,
		Descricao: createdCategoria.Descricao,
		Icone:     createdCategoria.Icone,
	}, nil
}

func (uc *CategoriaPrestadorUsecase) Editar(ctx context.Context, id uint, campos map[string]interface{}) (*CategoriaPrestadorResponse, error) {
	updatedCategoria, err := uc.repo.Editar(ctx, id, campos)
	if err != nil {
		return nil, err
	}

	return &CategoriaPrestadorResponse{
		ID:        updatedCategoria.ID,
		Nome:      updatedCategoria.Nome,
		Descricao: updatedCategoria.Descricao,
		Icone:     updatedCategoria.Icone,
	}, nil
}

func (uc *CategoriaPrestadorUsecase) Listar(ctx context.Context, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]CategoriaPrestadorResponse, error) {
	categorias, err := uc.repo.Listar(ctx, filters, orderBy, orderDir, limit, offset)
	if err != nil {
		return nil, err
	}

	resp := make([]CategoriaPrestadorResponse, 0, len(categorias))
	for _, c := range categorias {
		resp = append(resp, CategoriaPrestadorResponse{
			ID:        c.ID,
			Nome:      c.Nome,
			Descricao: c.Descricao,
			Icone:     c.Icone,
		})
	}
	return resp, nil
}

func (uc *CategoriaPrestadorUsecase) BuscarPorID(ctx context.Context, id uint) (*CategoriaPrestadorResponse, error) {
	categoria, err := uc.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &CategoriaPrestadorResponse{ID: categoria.ID, Nome: categoria.Nome, Descricao: categoria.Descricao, Icone: categoria.Icone}, nil
}

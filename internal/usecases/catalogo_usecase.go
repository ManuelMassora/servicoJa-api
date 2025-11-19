package usecases

import (
	"context"
	"errors"
	"log"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
)

type CatalogoUseCase struct {
	r model.CatalogoRepo
	prestadorRepo model.PrestadorRepo
}

func NewCatalogoUC(
	r model.CatalogoRepo,
	prestadorRepo model.PrestadorRepo,
	) *CatalogoUseCase{
	return &CatalogoUseCase{r: r, prestadorRepo: prestadorRepo}
}

type RequestCreateCatalogo struct {
	Nome        	string   `json:"nome" binding:"required"`
	Descricao   	string   `json:"descricao" binding:"required"`
	PrecoBase   	float64  `json:"precobase" binding:"required"`
	IdCategoria  	uint   	`json:"categoriaid" binding:"required"`
	Localizacao 	string   `json:"localizacao" binding:"required"`
}
type ResponseCatalogo struct {
	ID		  	uint    `json:"id"`
	Nome        string  `json:"nome"`
	Descricao   string  `json:"descricao"`
	PrecoBase   float64 `json:"preco_base"`
	Categoria   string  `json:"categoria"`
	Disponivel  bool    `json:"disponivel"`
	Prestador   string 	`json:"prestador"`
	Localizacao string   `json:"localizacao"`
}

func(uc *CatalogoUseCase) Criar(ctx context.Context, request RequestCreateCatalogo, idPrestador uint) error {
	prestador, err := uc.prestadorRepo.BuscarPorUsuarioID(ctx, idPrestador)
	if err != nil {
		return err
	}
	catalogo := &model.Catalogo{
		Nome: request.Nome,
		Descricao: request.Descricao,
		PrecoBase: request.PrecoBase,
		IDCategoria: request.IdCategoria,
		IDPrestador: prestador.ID,
		Localizacao: request.Localizacao,
	}
	log.Printf("Id do Usuario: %d", idPrestador)
	log.Printf("Id do prestador: %d", prestador.ID)
	return uc.r.Create(ctx,catalogo)
}

func(uc *CatalogoUseCase) Editar(ctx context.Context,id, idPrestador uint, campos map[string]interface{}) error {
	catalogo, err := uc.r.FindByID(ctx, id)
	if err != nil {
		return err
	}
	prestador, err := uc.prestadorRepo.BuscarPorUsuarioID(ctx, idPrestador)
	if err != nil {
		return err
	}
	if catalogo.IDPrestador != prestador.ID {
		return errors.New("nao tem permissao para apagar esse catalogo")
	}
	return uc.r.Update(ctx, id, campos)
}

func(uc *CatalogoUseCase) Apagar(ctx context.Context, id, idPrestador uint) error {
	catalogo, err := uc.r.FindByID(ctx, id)
	if err != nil {
		return err
	}
	prestador, err := uc.prestadorRepo.BuscarPorUsuarioID(ctx, idPrestador)
	if err != nil {
		return err
	}
	if catalogo.IDPrestador != prestador.ID {
		return errors.New("nao tem permissao para apagar esse catalogo")
	}
	return uc.r.Delete(ctx, id)
}

func(uc *CatalogoUseCase) Listar(ctx context.Context, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]ResponseCatalogo, error) {
	catalogos, err := uc.r.FindAll(ctx, filters, orderBy, orderDir, limit, offset)
	if err != nil {
		return nil, err
	}
	catalogoResponse := make([]ResponseCatalogo,0, len(catalogos))
	for _, catalogo := range catalogos {
		catalogoResponse = append(catalogoResponse, ResponseCatalogo{
			ID: catalogo.ID,
			Nome: catalogo.Nome,
			Descricao: catalogo.Descricao,
			PrecoBase: catalogo.PrecoBase,
			Categoria: catalogo.Categoria.Nome,
			Disponivel: catalogo.Disponivel,
			Prestador: catalogo.Prestador.Usuario.Nome,
			Localizacao: catalogo.Localizacao,
		})
	}
	return catalogoResponse, nil
}

func(uc *CatalogoUseCase) ListarPorPrestador(ctx context.Context,prestadorID uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]ResponseCatalogo, error) {
	catalogos, err := uc.r.FindByPrestadorID(ctx, prestadorID, filters, orderBy, orderDir, limit, offset)
	if err != nil {
		return nil, err
	}
	catalogoResponse := make([]ResponseCatalogo,0, len(catalogos))
	for _, catalogo := range catalogos {
		catalogoResponse = append(catalogoResponse, ResponseCatalogo{
			ID: catalogo.ID,
			Nome: catalogo.Nome,
			Descricao: catalogo.Descricao,
			PrecoBase: catalogo.PrecoBase,
			Categoria: catalogo.Categoria.Nome,
			Disponivel: catalogo.Disponivel,
			Prestador: catalogo.Prestador.Usuario.Nome,
			Localizacao: catalogo.Localizacao,
		})
	}
	return catalogoResponse, nil
}
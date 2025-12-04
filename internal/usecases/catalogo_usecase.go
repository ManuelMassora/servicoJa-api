package usecases

import (
	"context"
	"errors"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
)

type CatalogoUseCase struct {
	r               model.CatalogoRepo
	anexoImagemRepo model.AnexoImagemRepo
}

func NewCatalogoUC(
	r model.CatalogoRepo,
	anexoImagemRepo model.AnexoImagemRepo,
) *CatalogoUseCase {
	return &CatalogoUseCase{r: r, anexoImagemRepo: anexoImagemRepo}
}

type RequestCreateCatalogo struct {
	Nome        	string   `json:"nome" form:"nome" binding:"required"`
	Descricao   	string   `json:"descricao" form:"descricao" binding:"required"`
	TipoPreco		string 	 `json:"tipo_preco" form:"tipo_preco" binding:"required,oneof=fixo por_hora"`
	ValorFixo   	float64  `json:"valor_fixo" form:"valor_fixo"`
	ValorPorHora	float64  `json:"valor_por_hora" form:"valor_por_hora"`
	IdCategoria  	uint   	`json:"categoria_id" form:"categoria_id" binding:"required"`
	Localizacao 	string   `json:"localizacao" form:"localizacao" binding:"required"`
	Latitude    	float64  `json:"latitude" form:"latitude" binding:"required"`
	Longitude   	float64  `json:"longitude" form:"longitude" binding:"required"`
	Anexos      	[]string `binding:"-"`
}
type ResponseCatalogo struct {
	ID		  	uint    `json:"id"`
	Nome        string  `json:"nome"`
	Descricao   string  `json:"descricao"`
	TipoPreco	string  `json:"tipo_preco"`
	ValorFixo   float64 `json:"valor_fixo"`
	ValorPorHora float64 `json:"valor_por_hora"`
	Categoria   string  `json:"categoria"`
	Disponivel  bool    `json:"disponivel"`
	Prestador   string 	`json:"prestador"`
	Localizacao string   `json:"localizacao"`
	Latitude    float64  `json:"latitude"`
	Longitude   float64  `json:"longitude"`
	Anexos      []string `json:"anexos"`
}

func(uc *CatalogoUseCase) Criar(ctx context.Context, request RequestCreateCatalogo, idPrestador uint) error {
	catalogo := &model.Catalogo{
		Nome:         request.Nome,
		Descricao:    request.Descricao,
		TipoPreco:    request.TipoPreco,
		ValorFixo:    request.ValorFixo,
		ValorPorHora: request.ValorPorHora,
		IDCategoria:  request.IdCategoria,
		IDPrestador:  idPrestador,
		Localizacao:  request.Localizacao,
		Latitude:     request.Latitude,
		Longitude:    request.Longitude,
	}
	// Validation for pricing based on TipoPreco
	if catalogo.TipoPreco == "fixo" && catalogo.ValorFixo <= 0 {
		return errors.New("para preco fixo, o valor fixo deve ser maior que zero")
	}
	if catalogo.TipoPreco == "por_hora" && catalogo.ValorPorHora <= 0 {
		return errors.New("para preco por hora, o valor por hora deve ser maior que zero")
	}
	if catalogo.TipoPreco == "por_hora" && catalogo.ValorFixo > 0 {
		return errors.New("preco por hora nao pode ter valor fixo")
	}
	if catalogo.TipoPreco == "fixo" && catalogo.ValorPorHora > 0 {
		return errors.New("preco fixo nao pode ter valor por hora")
	}

	if err := uc.r.Create(ctx, catalogo); err != nil {
		return err
	}

	for _, anexoURL := range request.Anexos {
		anexo := &model.AnexoImagem{
			URL:        anexoURL,
			CatalogoID: &catalogo.ID,
		}
		if err := uc.anexoImagemRepo.Create(ctx, anexo); err != nil {
			// In a real application, you might want to handle the rollback of the catalogo creation
			return err
		}
	}

	return nil
}

func(uc *CatalogoUseCase) Editar(ctx context.Context,id, idPrestador uint, campos map[string]interface{}) error {
	catalogo, err := uc.r.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if catalogo.IDPrestador != idPrestador {
		return errors.New("nao tem permissao para apagar esse catalogo")
	}
	return uc.r.Update(ctx, id, campos)
}

func(uc *CatalogoUseCase) Apagar(ctx context.Context, id, idPrestador uint) error {
	catalogo, err := uc.r.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if catalogo.IDPrestador != idPrestador {
		return errors.New("nao tem permissao para apagar esse catalogo")
	}
	return uc.r.Delete(ctx, id)
}

func(uc *CatalogoUseCase) Listar(ctx context.Context, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]ResponseCatalogo, error) {
	catalogos, err := uc.r.FindAll(ctx, filters, orderBy, orderDir, limit, offset)
	if err != nil {
		return nil, err
	}
	var catalogosID []uint
	for _, catalogo := range catalogos {
		catalogosID = append(catalogosID, catalogo.ID)
	}
	anexos, err := uc.anexoImagemRepo.FindByCatalogoIDs(ctx, catalogosID)
	if err != nil {
		return nil, err
	}
	anexosPorVagaMap := make(map[uint][]string)
	for _, anexo := range anexos {
		anexosPorVagaMap[*anexo.CatalogoID] = append(anexosPorVagaMap[*anexo.CatalogoID], anexo.URL)
	}
	catalogoResponse := make([]ResponseCatalogo,0, len(catalogos))
	for _, catalogo := range catalogos {
		urls := anexosPorVagaMap[catalogo.ID]
		catalogoResponse = append(catalogoResponse, ResponseCatalogo{
			ID: catalogo.ID,
			Nome: catalogo.Nome,
			Descricao: catalogo.Descricao,
			TipoPreco: catalogo.TipoPreco,
			ValorFixo: catalogo.ValorFixo,
			ValorPorHora: catalogo.ValorPorHora,
			Categoria: catalogo.Categoria.Nome,
			Disponivel: catalogo.Disponivel,
			Prestador: catalogo.Prestador.Nome,
			Localizacao: catalogo.Localizacao,
			Latitude: catalogo.Latitude,
			Longitude: catalogo.Longitude,
			Anexos: urls,
		})
	}
	return catalogoResponse, nil
}

func(uc *CatalogoUseCase) ListarPorPrestador(ctx context.Context,prestadorID uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]ResponseCatalogo, error) {
	catalogos, err := uc.r.FindByPrestadorID(ctx, prestadorID, filters, orderBy, orderDir, limit, offset)
	if err != nil {
		return nil, err
	}
	var catalogosID []uint
	for _, catalogo := range catalogos {
		catalogosID = append(catalogosID, catalogo.ID)
	}
	anexos, err := uc.anexoImagemRepo.FindByCatalogoIDs(ctx, catalogosID)
	if err != nil {
		return nil, err
	}
	anexosPorVagaMap := make(map[uint][]string)
	for _, anexo := range anexos {
		anexosPorVagaMap[*anexo.CatalogoID] = append(anexosPorVagaMap[*anexo.CatalogoID], anexo.URL)
	}
	catalogoResponse := make([]ResponseCatalogo,0, len(catalogos))
	for _, catalogo := range catalogos {
		urls := anexosPorVagaMap[catalogo.ID]
		catalogoResponse = append(catalogoResponse, ResponseCatalogo{
			ID: catalogo.ID,
			Nome: catalogo.Nome,
			Descricao: catalogo.Descricao,
			TipoPreco: catalogo.TipoPreco,
			ValorFixo: catalogo.ValorFixo,
			ValorPorHora: catalogo.ValorPorHora,
			Categoria: catalogo.Categoria.Nome,
			Disponivel: catalogo.Disponivel,
			Prestador: catalogo.Prestador.Nome,
			Localizacao: catalogo.Localizacao,
			Latitude: catalogo.Latitude,
			Longitude: catalogo.Longitude,
			Anexos: urls,
		})
	}
	return catalogoResponse, nil
}

func(uc *CatalogoUseCase) ListarPorLocalizacao(ctx context.Context, latitude, longitude, radius float64, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]ResponseCatalogo, error) {
	catalogos, err := uc.r.FindByLocation(ctx, latitude, longitude, radius, filters, orderBy, orderDir, limit, offset)
	if err != nil {
		return nil, err
	}
	var catalogosID []uint
	for _, catalogo := range catalogos {
		catalogosID = append(catalogosID, catalogo.ID)
	}
	anexos, err := uc.anexoImagemRepo.FindByCatalogoIDs(ctx, catalogosID)
	if err != nil {
		return nil, err
	}
	anexosPorVagaMap := make(map[uint][]string)
	for _, anexo := range anexos {
		anexosPorVagaMap[*anexo.CatalogoID] = append(anexosPorVagaMap[*anexo.CatalogoID], anexo.URL)
	}
	var catalogoResponse []ResponseCatalogo
	for _, catalogo := range catalogos {
		urls := anexosPorVagaMap[catalogo.ID]
		catalogoResponse = append(catalogoResponse, ResponseCatalogo{
			ID: catalogo.ID,
			Nome: catalogo.Nome,
			Descricao: catalogo.Descricao,
			TipoPreco: catalogo.TipoPreco,
			ValorFixo: catalogo.ValorFixo,
			ValorPorHora: catalogo.ValorPorHora,
			Categoria: catalogo.Categoria.Nome,
			Disponivel: catalogo.Disponivel,
			Prestador: catalogo.Prestador.Nome,
			Localizacao: catalogo.Localizacao,
			Latitude: catalogo.Latitude,
			Longitude: catalogo.Longitude,
			Anexos: urls,
		})
	}
	return catalogoResponse, nil
}
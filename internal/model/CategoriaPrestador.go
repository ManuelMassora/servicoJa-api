package model

import "context"

type CategoriaPrestador struct {
	BaseModel
	Nome      string `json:"nome"`
	Descricao string `json:"descricao"`
	Icone     string `json:"icone"`
}

type CategoriaPrestadorRepo interface {
	Criar(ctx context.Context, categoria *CategoriaPrestador) (*CategoriaPrestador,error)
	Editar(ctx context.Context, id uint, campos map[string]interface{}) (*CategoriaPrestador,error)
	Listar(ctx context.Context, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]CategoriaPrestador, error)
	BuscarPorID(ctx context.Context, id uint) (*CategoriaPrestador, error)
}
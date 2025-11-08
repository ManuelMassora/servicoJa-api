package model

import "context"

type Categoria struct {
	BaseModel
	Nome      string `json:"nome"`
	Descricao string `json:"descricao"`
}

type CategoriaRepo interface {
	Criar(ctx context.Context, categoria *Categoria) error
	Listar(ctx context.Context, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]Categoria, error)
	BuscarPorID(ctx context.Context, id int64) (*Categoria, error)
}

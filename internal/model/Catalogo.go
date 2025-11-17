package model

import "context"

type Catalogo struct {
	BaseModel
	Nome        string   `gorm:"column:nome;size:100;not null" json:"nome"`
	Descricao   string   `gorm:"column:descricao;size:2000;not null" json:"descricao"`
	PrecoBase   float64  `gorm:"column:preco_base;type:decimal(10,2);not null" json:"preco_base"`
	Categoria   string   `gorm:"column:categoria;size:100;not null" json:"categoria"` //Associar com categoria
	IDPrestador int64    `gorm:"column:id_prestador;type:bigint;not null" json:"prestador_id"`
	Disponivel  bool     `gorm:"column:disponivel;default:true" json:"disponivel"`
	Prestador   *Usuario `gorm:"foreignKey:IDPrestador;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"prestador,omitempty"` // Corrigir pra usar *Prestador nao *Usuario
}

type CatalogoRepo interface {
	Create(ctx context.Context, catalogo *Catalogo) error
	Update(ctx context.Context, id int64, campos map[string]interface{}) error
	Delete(ctx context.Context, id int64) error
	FindByID(ctx context.Context, id int64) (*Catalogo, error)
	FindAll(ctx context.Context, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]*Catalogo, error)
	FindByPrestadorID(ctx context.Context, prestadorID int64, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]*Catalogo, error)
}

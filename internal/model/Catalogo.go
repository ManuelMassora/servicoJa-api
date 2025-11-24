package model

import "context"

type Catalogo struct {
	BaseModel
	Nome        string   `gorm:"column:nome;size:100;not null" json:"nome"`
	Descricao   string   `gorm:"column:descricao;size:2000;not null" json:"descricao"`
	PrecoBase   float64  `gorm:"column:preco_base;type:decimal(10,2);not null" json:"preco_base"`
	IDCategoria  uint   `gorm:"column:id_categoria;size:100;not null" json:"categoria_id"`
	Categoria   Categoria   `gorm:"foreignKey:IDCategoria;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"categoria,omitempty"`
	Disponivel  bool     `gorm:"column:disponivel;default:true" json:"disponivel"`
	IDPrestador uint    `gorm:"column:id_prestador;type:bigint;not null" json:"prestador_id"`
	Prestador   Prestador `gorm:"foreignKey:IDPrestador;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"prestador,omitempty"` // Corrigir pra usar *Prestador nao *Usuario
	Localizacao string   `gorm:"column:localizacao;size:255;" json:"localizacao"`
	Latitude    float64  `gorm:"column:latitude;type:decimal(10,8);" json:"latitude"`
	Longitude   float64  `gorm:"column:longitude;type:decimal(11,8);" json:"longitude"`
}

type CatalogoRepo interface {
	Create(ctx context.Context, catalogo *Catalogo) error
	Update(ctx context.Context, id uint, campos map[string]interface{}) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*Catalogo, error)
	FindAll(ctx context.Context, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]*Catalogo, error)
	FindByPrestadorID(ctx context.Context, prestadorID uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]*Catalogo, error)
	FindByLocation(ctx context.Context, latitude, longitude, radius float64, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]*Catalogo, error)
}
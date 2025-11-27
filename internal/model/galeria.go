package model

import "context"

type Galeria struct {
	BaseModel
	PrestadorID uint    `gorm:"not null;index"`
	Imagens     []Imagem `gorm:"foreignKey:GaleriaID"`
}

type Imagem struct {
	BaseModel
	URL       string `gorm:"type:varchar(255);not null"`
	GaleriaID uint   `gorm:"not null;index"`
}

type GaleriaRepo interface {
	Create(ctx context.Context, galeria *Galeria) (*Galeria, error)
	FindByID(ctx context.Context, id uint) (*Galeria, error)
	FindByPrestadorID(ctx context.Context, prestadorID uint) (*Galeria, error)
	Delete(ctx context.Context, id uint) error
	AddImage(ctx context.Context, imagem *Imagem) error
	CountImages(ctx context.Context, galeriaID uint) (int64, error)
	FindByGaleriaID(ctx context.Context, galeriaID uint) ([]Imagem, error)
}

package model

import "context"

type AnexoImagem struct {
	BaseModel
	URL           string `gorm:"type:varchar(255);not null"`
	AgendamentoID *uint  `gorm:"index"`
	VagaID        *uint  `gorm:"index"`
	CatalogoID    *uint  `gorm:"index"`
}

type AnexoImagemRepo interface {
	Create(ctx context.Context, anexo *AnexoImagem) error
	FindByID(ctx context.Context, id uint) (*AnexoImagem, error)
	FindByAgendamentoID(ctx context.Context, agendamentoID uint) ([]AnexoImagem, error)
	FindByVagaID(ctx context.Context, vagaID uint) ([]AnexoImagem, error)
	FindByCatalogoID(ctx context.Context, catalogoID uint) ([]AnexoImagem, error)
	FindByAgendamentoIDs(ctx context.Context, agendamentoIDs []uint) ([]AnexoImagem, error)
	FindByVagaIDs(ctx context.Context, vagaIDs []uint) ([]AnexoImagem, error)
	FindByCatalogoIDs(ctx context.Context, catalogoIDs []uint) ([]AnexoImagem, error)
}
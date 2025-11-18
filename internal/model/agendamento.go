package model

import (
	"context"
	"time"
)

type Agendamento struct {
	BaseModel
	Detalhe 	string `json:"detalhe"`
	IDCatalogo  uint  `json:"id_catalogo"`
	Catalogo	Catalogo		`gorm:"foreignKey:IDCatalogo;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"catalogo,omitempty"`
	IDCliente  	uint  `json:"id_cliente"`
	Cliente		Cliente		`gorm:"foreignKey:IDCliente;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"cliente,omitempty"`
	DataHora 	time.Time      `json:"datahora"`
	Status		string         `json:"status"`
}

type AgendamentoRepo interface {
	Criar(ctx context.Context, agendamento *Agendamento) error
	BuscarPorID(ctx context.Context, id uint) (*Agendamento, error)
	Listar(
		ctx context.Context, 
		filters map[string]interface{}, 
		orderBy string, 
		orderDir string, 
		limit, 
		offset int,
	) ([]Agendamento, error)
}
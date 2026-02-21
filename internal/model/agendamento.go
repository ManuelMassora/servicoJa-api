package model

import (
	"context"
	"time"
)

type Agendamento struct {
	BaseModel
	Detalhe 	string `json:"detalhe"`
	IDCatalogo  uint  `json:"id_catalogo"`
	Catalogo	Catalogo		`gorm:"foreignKey:IDCatalogo;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"catalogo,omitempty"`
	IDCliente  	uint  `json:"id_cliente"`
	Cliente		Cliente		`gorm:"foreignKey:IDCliente;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"cliente,omitempty"`
	DataHora 	time.Time      `json:"datahora"`
	Status		string         `json:"status"`
	Localizacao string   `gorm:"column:localizacao;size:255;" json:"localizacao"`
	Latitude    float64  `gorm:"column:latitude;type:decimal(10,8);" json:"latitude"`
	Longitude   float64  `gorm:"column:longitude;type:decimal(11,8);" json:"longitude"`
	Anexos      []AnexoImagem  `gorm:"foreignKey:AgendamentoID"`
}

type AgendamentoRepo interface {
	Criar(ctx context.Context, agendamento *Agendamento) (*Agendamento, error)
	BuscarPorID(ctx context.Context, id uint) (*Agendamento, error)
	AtualizarStatus(ctx context.Context, id uint, status string) error
	Listar(
		ctx context.Context, 
		filters map[string]interface{}, 
		orderBy string, 
		orderDir string, 
		limit, 
		offset int,
	) ([]Agendamento, error)
	ListarPorClienteID(ctx context.Context, clienteID uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]Agendamento, error)
	ListarPorPrestadorID(ctx context.Context, prestadorID uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]Agendamento, error)
	ListarPorCatalogID(ctx context.Context, catalogoID uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]Agendamento, error)
	FindByLocation(ctx context.Context, userID uint, latitude, longitude, radius float64, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]Agendamento, error)
}
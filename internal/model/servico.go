package model

import (
	"context"
	"time"
)

type Servico struct {
	BaseModel
	Localizacao string   `gorm:"column:localizacao;size:255;not null" json:"localizacao"`
	Latitude    float64  `gorm:"column:latitude;type:decimal(10,8);" json:"latitude"`
	Longitude   float64  `gorm:"column:longitude;type:decimal(11,8);" json:"longitude"`
	Preco       float64  `gorm:"column:preco;type:decimal(10,2);not null" json:"preco"`
	Status      Status   `gorm:"column:status;type:varchar(20);not null" json:"status"`
	IDAgendamento   *uint    `gorm:"column:id_agendamento;type:bigint;" json:"id_agendamento,omitempty"`
	Agendamento     *Agendamento `gorm:"foreignKey:IDAgendamento;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"agendamento,omitempty"`
	IDVaga *uint `gorm:"column:id_vaga;type:bigint;default:null" json:"id_vaga,omitempty"`
	Vaga   *Vaga `gorm:"foreignKey:IDVaga;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"vaga,omitempty"`
	DataHoraInicio  time.Time `gorm:"column:data_inicio;type:timestamp;default:null" json:"data_inicio,omitempty"`
	DataHoraFim     time.Time  `gorm:"column:data_fim;type:timestamp;default:null" json:"data_fim,omitempty"`
	DataHoraConfirmado  time.Time  `gorm:"column:data_confirmado;type:timestamp;default:null" json:"data_confirmado,omitempty"`
	IDCliente    uint      `gorm:"column:id_cliente;not null"`
	Cliente      *Cliente  `gorm:"foreignKey:IDCliente"`
	IDPrestador  uint      `gorm:"column:id_prestador;not null"`
	Prestador    *Prestador `gorm:"foreignKey:IDPrestador"`
	IfAvaliadoCliente bool      `gorm:"column:if_avaliado_cliente;type:boolean;default:false" json:"if_avaliado_cliente"`
}

type ServicoRepo interface {
	Criar(ctx context.Context, servico *Servico) error
	Atualizar(ctx context.Context, servico *Servico) error
	BuscarPorID(ctx context.Context, id uint) (*Servico, error)
	AtualizarStatus(ctx context.Context, id uint, status string) error
	ListarPorCliente(ctx context.Context, idCliente uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]Servico, error)
	ListarPorPrestador(ctx context.Context, IDPrestador uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]Servico, error)
	FindByLocation(ctx context.Context, userID uint, latitude, longitude, radius float64, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]Servico, error)
}
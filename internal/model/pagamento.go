package model

import "context"

type Pagamento struct {
	BaseModel
	IDServico     *uint        `json:"servico_id" gorm:"default:null"`
	Servico       *Servico     `gorm:"foreignKey:IDServico;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"servico,omitempty"`
	IDVaga        *uint        `json:"vaga_id" gorm:"default:null"`
	Vaga          *Vaga        `gorm:"foreignKey:IDVaga;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"vaga,omitempty"`
	IDAgendamento *uint        `json:"agendamento_id" gorm:"default:null"`
	Agendamento   *Agendamento `gorm:"foreignKey:IDAgendamento;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"agendamento,omitempty"`
	IDCliente     uint         `json:"cliente_id" gorm:"not null"`
	Cliente       *Usuario     `gorm:"foreignKey:IDCliente;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"cliente,omitempty"`
	IDPrestador   *uint        `json:"prestador_id" gorm:"default:null"`
	Prestador     *Usuario     `gorm:"foreignKey:IDPrestador;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"prestador,omitempty"`
	Valor         float64      `json:"valor" gorm:"not null"`
	Status        Status       `json:"status" gorm:"column:status;type:varchar(20);not null"`
}

type PagamentoRepo interface {
	Criar(ctx context.Context, pagamento *Pagamento) error
	BuscarPorServico(ctx context.Context, idServico uint) (*Pagamento, error)
	BuscarPorVaga(ctx context.Context, idVaga uint) (*Pagamento, error)
	BuscarPorAgendamento(ctx context.Context, idAgendamento uint) (*Pagamento, error)
	AtualizarStatus(ctx context.Context, id uint, status Status) error
	ListarPorUsuario(ctx context.Context, idUsuario uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]Pagamento, error)
}

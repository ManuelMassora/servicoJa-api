package model

import "context"

type Pagamento struct {
	BaseModel
	IDServico   uint    `json:"servico_id" gorm:"not null"`
	Servico     *Servico `gorm:"foreignKey:IDServico;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"servico,omitempty"`
	IDCliente   uint    `json:"cliente_id" gorm:"not null"`
	Cliente     *Usuario `gorm:"foreignKey:IDCliente;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"cliente,omitempty"`
	IDPrestador uint    `json:"prestador_id" gorm:"not null"`
	Prestador   *Usuario `gorm:"foreignKey:IDPrestador;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"prestador,omitempty"`
	Valor       float64  `json:"valor" gorm:"not null"`
	Status      Status   `json:"status" gorm:"column:status;type:varchar(20);not null"`
}

type PagamentoRepo interface {
	Criar(ctx context.Context, pagamento *Pagamento) error
	BuscarPorServico(ctx context.Context, idServico uint) (*Pagamento, error)
	AtualizarStatus(ctx context.Context, id uint, status Status) error
	ListarPorUsuario(ctx context.Context, idUsuario uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]Pagamento, error)
}

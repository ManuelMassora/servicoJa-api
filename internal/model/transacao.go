package model

import "context"

type Transacao struct {
	BaseModel
	IDUsuario     uint        `json:"usuario_id" gorm:"not null"`
	Usuario       *Usuario     `gorm:"foreignKey:IDUsuario;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"usuario,omitempty"`
	IDServico     uint        `json:"servico_id" gorm:"not null"`
	Servico       *Servico     `gorm:"foreignKey:IDServico;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"servico,omitempty"`
	TipoMovimento TipoMovimento `json:"tipo_movimento" gorm:"column:tipo_movimento;not null"`
	Valor         float64      `json:"valor" gorm:"column:valor;not null"`
	Metodo        string       `json:"metodo" gorm:"column:metodo;not null"`
	Status        Status       `json:"status" gorm:"column:status;type:varchar(20);not null"`
}

type TransacaoRepo interface {
	Criar(ctx context.Context, transacao *Transacao) error
		ListarPorUsuario(ctx context.Context, idUsuario uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]Transacao, error)
}

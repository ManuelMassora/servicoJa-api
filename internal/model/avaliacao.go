package model

import "context"

type Avaliacao struct {
	BaseModel
	Nota       int    `json:"nota"`
	Comentario string `json:"comentario"`
	UsuarioID  int64  `json:"usuario_id"`
	Usuario		Usuario 	`gorm:"foreignKey:UsuarioID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"usuario,omitempty"`
	ServicoID  int64  `json:"servico_id"`
	Servico		Servico		`gorm:"foreignKey:ServicoID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"servico,omitempty"`
}

type AvaliacaoRepo interface {
	Criar(ctx context.Context, avaliacao *Avaliacao) error
	ListarPorServico(ctx context.Context, idServico int64, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]Avaliacao, error)
	MediaPorPrestador(ctx context.Context, idPrestador int64) (float64, error)
}

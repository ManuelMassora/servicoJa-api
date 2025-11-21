package model

import "context"

type Avaliacao struct {
	BaseModel
	Nota       int    `json:"nota"`
	Comentario string `json:"comentario"`
	UsuarioID  uint  `json:"usuario_id"`
	Usuario		Usuario 	`gorm:"foreignKey:UsuarioID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"usuario,omitempty"`
	ServicoID  uint  `json:"servico_id"`
	Servico		Servico		`gorm:"foreignKey:ServicoID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"servico,omitempty"`
}

type AvaliacaoRepo interface {
	Criar(ctx context.Context, avaliacao *Avaliacao) error
	ListarPorServico(ctx context.Context, idServico uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]Avaliacao, error)
	MediaPorPrestador(ctx context.Context, idPrestador uint) (float64, error)
}
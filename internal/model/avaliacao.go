package model

import "context"

type Avaliacao struct {
	BaseModel
	Nota        int       `gorm:"column:nota;not null" json:"nota"`
	Comentario  string    `gorm:"column:comentario;size:1000" json:"comentario"`
	IDCliente   uint      `gorm:"column:id_cliente;not null" json:"cliente_id"`
	Cliente     *Cliente  `gorm:"foreignKey:IDCliente;references:IDUsuario" json:"cliente,omitempty"`
	IDPrestador uint      `gorm:"column:id_prestador;not null" json:"prestador_id"`
	Prestador   *Prestador `gorm:"foreignKey:IDPrestador;references:IDUsuario" json:"prestador,omitempty"`
	IDServico   uint      `gorm:"column:id_servico;not null;uniqueIndex:idx_avaliacao_servico" json:"servico_id"`
	Servico     *Servico  `gorm:"foreignKey:IDServico;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"servico,omitempty"`
}

type AvaliacaoRepo interface {
	Criar(ctx context.Context, avaliacao *Avaliacao) error
	ListarPorCliente(ctx context.Context, idCliente uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]Avaliacao, error)
	ListarPorPrestador(ctx context.Context, idPrestador uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]Avaliacao, error)
	MediaPorPrestador(ctx context.Context, idPrestador uint) (float64, error)
}
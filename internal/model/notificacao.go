package model

import "context"

type Notificacao struct {
	BaseModel
	IDUsuario uint    `json:"usuario_id"`
	Usuario   *Usuario `gorm:"foreignKey:IDUsuario;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"usuario,omitempty"`
	Titulo    string   `json:"titulo"`
	Mensagem  string   `json:"mensagem"`
	Lida      bool     `json:"lida"`
}

type NotificacaoRepo interface {
	Enviar(ctx context.Context, notificacao *Notificacao) error
	ListarPorPrestador(ctx context.Context, idPrestador uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]Notificacao, error)
	MarcarComoLida(ctx context.Context, id uint) error
}

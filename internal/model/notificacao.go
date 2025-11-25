package model

import "context"

type Notificacao struct {
	BaseModel
	IDUsuario uint    `json:"usuario_id"`
	Usuario   *Usuario `gorm:"foreignKey:IDUsuario;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"usuario,omitempty"`
	Titulo    string   `json:"titulo"`
	Mensagem  string   `json:"mensagem"`
	Lida      bool     `json:"lida"`
}

type NotificacaoRepo interface {
	Enviar(ctx context.Context, notificacao *Notificacao) error
	BuscarPorID(ctx context.Context, id uint) (*Notificacao,error)
	ListarPorUsuario(ctx context.Context, idUsuario uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]Notificacao, error)
	MarcarComoLida(ctx context.Context, id uint) error
	MarcarTodasComoLidas(ctx context.Context, idUsuario uint) error
}

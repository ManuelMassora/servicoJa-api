package model

import "context"

type Chat struct {
	BaseModel
	ServicoID   int64    `json:"servico_id"`
	Servico     *Servico `gorm:"foreignKey:ServicoID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"servico,omitempty"`
	PrestadorID int64    `json:"prestador_id"`
	Prestador   *Prestador `gorm:"foreignKey:PrestadorID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"prestador,omitempty"`
	IDCliente   int64    `json:"cliente_id"`
	Cliente     *Prestador `gorm:"foreignKey:IDCliente;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"cliente,omitempty"`
}

type Mensagem struct {
	BaseModel
	IDChat        uint    `json:"chat_id"`
	Chat          *Chat    `gorm:"foreignKey:IDChat;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"chat,omitempty"`
	IDRemetente   uint    `json:"remetente_id"`
	Remetente     *Usuario `gorm:"foreignKey:IDRemetente;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"remetente,omitempty"`
	RemetenteTipo string   `json:"remetente_tipo"`
	Conteudo      string   `json:"conteudo"`
}

type ChatRepo interface {
	CriarChat(ctx context.Context, chat *Chat) error
		ListarChatsPorUsuario(ctx context.Context, idUsuario uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]Chat, error)
}

type MensagemRepo interface {
	EnviarMensagem(ctx context.Context, msg *Mensagem) error
		ListarMensagens(ctx context.Context, idChat uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]Mensagem, error)
}
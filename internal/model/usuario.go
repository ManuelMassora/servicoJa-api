package model

import "context"

type Usuario struct {
	BaseModel
	Nome     string `json:"nome"`
	Telefone string `json:"telefone"`
	Email    string `json:"email,omitempty"`
	Senha    string `json:"senha,omitempty"`
	RolePermissaoID   int64
	RolePermissao     RolePermissao  `gorm:"foreignKey:RolePermissaoID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"rolepermissao,omitempty"`
}

type UsuarioRepo interface {
	Criar(ctx context.Context, usuario *Usuario) error
	BuscarPorID(ctx context.Context, id int64) (*Usuario, error)
	Atualizar(ctx context.Context, usuario *Usuario) error
	Remover(ctx context.Context, id int64) error
	ListarTodos(ctx context.Context, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]Usuario, error)
}

// Cliente representa a ligação de um usuário com o perfil de cliente.
type Cliente struct {
	BaseModel
	UsuarioID int64   `json:"usuario_id" gorm:"not null"`
	Usuario   *Usuario `gorm:"foreignKey:UsuarioID;references:ID" json:"usuario,omitempty"`
}

type ClienteRepo interface {
	Criar(ctx context.Context, cliente *Cliente) error
	BuscarPorID(ctx context.Context, id int64) (*Cliente, error)
		Listar(ctx context.Context, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]Cliente, error)
}

// Prestador representa o perfil de prestador ligado a um usuário.
type Prestador struct {
	BaseModel
	UsuarioID       int64    `json:"usuario_id" gorm:"not null"`
	Usuario         *Usuario `gorm:"foreignKey:UsuarioID;references:ID" json:"usuario,omitempty"`
	Localizacao     string   `json:"localizacao"`
	StatusDisponivel bool    `json:"status_disponivel"`
	RaioAtuacao     float64  `json:"raio_atuacao"`
	Reputacao       float64  `json:"reputacao"`
}

type PrestadorRepo interface {
	Criar(ctx context.Context, prestador *Prestador) error
	AtualizarStatus(ctx context.Context, id int64, disponivel bool) error
	BuscarPorLocalizacao(ctx context.Context, local string) ([]Prestador, error)
	BuscarPorID(ctx context.Context, id int64) (*Prestador, error)
	Listar(ctx context.Context, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]Prestador, error)
}
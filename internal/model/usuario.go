package model

import (
	"context"
	"errors"
	"log"

	"github.com/alexedwards/argon2id"
)

type Usuario struct {
	BaseModel
	Nome     string `json:"nome"`
	Telefone string `json:"telefone"`
	Senha    string `json:"senha,omitempty"`
	RolePermissaoID   int64
	RolePermissao     RolePermissao  `gorm:"foreignKey:RolePermissaoID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"rolepermissao,omitempty"`
}

type UsuarioRepo interface {
	Criar(ctx context.Context, usuario *Usuario) error
	BuscarPorID(ctx context.Context, id int64) (*Usuario, error)
	BuscarPorTelefone(ctx context.Context, numero string) (*Usuario, error)
	Atualizar(ctx context.Context, usuario *Usuario) error
	Remover(ctx context.Context, id int64) error
	ListarTodos(ctx context.Context, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]Usuario, error)
}

// Cliente representa a ligação de um usuário com o perfil de cliente.
type Cliente struct {
	BaseModel
	UsuarioID int64   `json:"usuario_id" gorm:"not null"`
	Usuario   Usuario `gorm:"foreignKey:UsuarioID;references:ID" json:"usuario,omitempty"`
}

type ClienteRepo interface {
	Criar(ctx context.Context, cliente *Cliente) error
	BuscarPorID(ctx context.Context, id int64) (*Cliente, error)
	Listar(ctx context.Context, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]Cliente, error)
}

// Prestador representa o perfil de prestador ligado a um usuário.
type Prestador struct {
	BaseModel
	UsuarioID       	int64    	`json:"usuario_id" gorm:"not null"`
	Usuario         	Usuario 	`gorm:"foreignKey:UsuarioID;references:ID" json:"usuario,omitempty"`
	Localizacao     	string   	`json:"localizacao"`
	StatusDisponivel 	bool    	`json:"status_disponivel"`
	RaioAtuacao     	float64  	`json:"raio_atuacao"`
	Reputacao       	float64  	`json:"reputacao"`
}

type PrestadorRepo interface {
	Criar(ctx context.Context, prestador *Prestador) error
	AtualizarStatus(ctx context.Context, id int64, disponivel bool) error
	BuscarPorID(ctx context.Context, id int64) (*Prestador, error)
	Listar(ctx context.Context, filters map[string]interface{}, statusDisponivel interface{}, orderBy string, orderDir string, limit, offset int) ([]Prestador, error)
}

var params = &argon2id.Params{
    Memory: 16 * 1024,
    Iterations: 1,
    Parallelism: 1,
    SaltLength: 16,
    KeyLength: 16,
}

func NewAdmin(nome, telefone, senha string) (*Usuario,error) {
	if nome == "" {
        return nil, errors.New("nome não pode ser vazio")
    }

	if telefone == "" {
        return nil, errors.New("telefone não pode ser vazio")
    }
	if len([]byte(senha)) < 6 {
        return nil, errors.New("senha não pode ter menos de 6 letras")
    }

    if len([]byte(senha)) > 100 {
        return nil, errors.New("senha não pode ter mais de 60 letras")
    }
    senhaHash, err := argon2id.CreateHash(senha, params)
    if err != nil {
        log.Fatal(err)
        return nil, err
    }
	usuario := &Usuario{
		Nome: nome,
		Telefone: telefone,
		Senha: senhaHash,
		RolePermissaoID: 3,
	}
	return usuario, nil
}

func NewCliente(nome, telefone, senha string) (*Cliente,error) {
	if nome == "" {
        return nil, errors.New("nome não pode ser vazio")
    }

	if telefone == "" {
        return nil, errors.New("telefone não pode ser vazio")
    }
		
	if len([]byte(senha)) < 6 {
        return nil, errors.New("senha não pode ter menos de 6 letras")
    }

    if len([]byte(senha)) > 100 {
        return nil, errors.New("senha não pode ter mais de 60 letras")
    }
    senhaHash, err := argon2id.CreateHash(senha, params)
    if err != nil {
        log.Fatal(err)
        return nil, err
    }
	cliente := &Cliente{
		Usuario: Usuario{
			Nome: nome,
			Telefone: telefone,
			Senha: senhaHash,
			RolePermissaoID: 1,
		},
	}
	return cliente, nil
}

func NewPrestador(localizacao string, raio_atuacao float64, nome string, telefone, senha string) (*Prestador,error) {
	if nome == "" {
        return nil, errors.New("nome não pode ser vazio")
    }

	if telefone == "" {
        return nil, errors.New("telefone não pode ser vazio")
    }
	if len([]byte(senha)) < 6 {
        return nil, errors.New("senha não pode ter menos de 6 letras")
    }

    if len([]byte(senha)) > 100 {
        return nil, errors.New("senha não pode ter mais de 60 letras")
    }
    senhaHash, err := argon2id.CreateHash(senha, params)
    if err != nil {
        log.Fatal(err)
        return nil, err
    }
	prestador := &Prestador{
		Usuario: Usuario{
			Nome: nome,
			Telefone: telefone,
			Senha: senhaHash,
			RolePermissaoID: 1,
		},
		Localizacao: localizacao,
		RaioAtuacao: raio_atuacao,
		StatusDisponivel: true,
	}
	return prestador, nil
}
package model

import (
	"context"
	"errors"
	"log"
	"regexp"

	"github.com/alexedwards/argon2id"
)

type Usuario struct {
	BaseModel
	Nome     			string 	`json:"nome"`
	Telefone 			string 	`gorm:"unique" json:"telefone"`
	Senha    			string 	`json:"senha,omitempty"`
	NotificacoesNovas 	uint 	`json:"notificacoes_novas"`
	RolePermissaoID   	uint
	RolePermissao     	RolePermissao  `gorm:"foreignKey:RolePermissaoID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"rolepermissao,omitempty"`
}

type UsuarioRepo interface {
	Criar(ctx context.Context, usuario *Usuario) error
	BuscarPorID(ctx context.Context, id uint) (*Usuario, error)
	IncrementarNotificacoesNovas(ctx context.Context, id uint) error
	ZerarNotificacoesNovas(ctx context.Context, id uint) error
	BuscarPorTelefone(ctx context.Context, numero string) (*Usuario, error)
	Atualizar(ctx context.Context, usuario *Usuario) error
	Remover(ctx context.Context, id uint) error
	ListarTodos(ctx context.Context, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]Usuario, error)
}

// Cliente representa a ligação de um usuário com o perfil de cliente.
type Cliente struct {
	IDUsuario 	uint    `gorm:"primaryKey"`
    Usuario   	Usuario `gorm:"constraint:OnDelete:CASCADE;foreignKey:IDUsuario;references:ID"`
	Nome     	string 	`json:"nome"`
	Telefone 	string 	`json:"telefone"`
}

type ClienteRepo interface {
	Criar(ctx context.Context, cliente *Cliente) error
	BuscarPorID(ctx context.Context, id uint) (*Cliente, error)
	Listar(ctx context.Context, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]Cliente, error)
}

// Prestador representa o perfil de prestador ligado a um usuário.
type Prestador struct {
	IDUsuario 			uint    	`gorm:"primaryKey"`
    Usuario   			Usuario 	`gorm:"constraint:OnDelete:CASCADE;foreignKey:IDUsuario;references:ID"`
	Nome     			string 		`json:"nome"`
	Telefone 			string 		`json:"telefone"`
	Localizacao     	string   	`json:"localizacao"`
	StatusDisponivel 	bool    	`json:"status_disponivel"`
	Reputacao       	float64  	`json:"reputacao"`
}

type PrestadorRepo interface {
	Criar(ctx context.Context, prestador *Prestador) error
	AtualizarStatus(ctx context.Context, id uint, disponivel bool) error
	BuscarPorID(ctx context.Context, id uint) (*Prestador, error)
	Listar(ctx context.Context, filters map[string]interface{}, statusDisponivel interface{}, orderBy string, orderDir string, limit, offset int) ([]Prestador, error)
}

var params = &argon2id.Params{
    Memory: 16 * 1024,
    Iterations: 1,
    Parallelism: 1,
    SaltLength: 16,
    KeyLength: 16,
}

var strictE164Regex = regexp.MustCompile(`^\+[0-9]+$`)

func validateNumericPhone(telefone string) error {
	if telefone == "" {
		return errors.New("o telefone não pode ser vazio")
	}

	if !strictE164Regex.MatchString(telefone) {
		return errors.New("o telefone deve conter apenas números (dígitos de 0 a 9)")
	}

	if len(telefone) < 8 || len(telefone) > 15 {
	    return errors.New("o telefone tem um comprimento inválido")
	}

	return nil
}

func NewAdmin(nome, telefone, senha string) (*Usuario,error) {
	if nome == "" {
        return nil, errors.New("nome não pode ser vazio")
    }
	if err := validateNumericPhone(telefone); err != nil {
		return nil, err
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
	if err := validateNumericPhone(telefone); err != nil {
		return nil, err
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
		Nome: nome,
		Telefone: telefone,
	}
	return cliente, nil
}

func NewPrestador(localizacao string, nome string, telefone, senha string) (*Prestador,error) {
	if nome == "" {
        return nil, errors.New("nome não pode ser vazio")
    }
	if err := validateNumericPhone(telefone); err != nil {
		return nil, err
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
			RolePermissaoID: 2,
		},
		Nome: nome,
		Telefone: telefone,
		Localizacao: localizacao,
		StatusDisponivel: true,
	}
	return prestador, nil
}
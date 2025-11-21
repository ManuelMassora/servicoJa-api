package usecases

import (
	"context"
	"errors"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"gorm.io/gorm"
)

type UsuarioUseCase struct {
	usuarioRepo model.UsuarioRepo
	clienteRepo model.ClienteRepo
	prestadorRepo model.PrestadorRepo
}

func NewUsuarioUseCase(
	usuarioRepo model.UsuarioRepo,
	clienteRepo model.ClienteRepo,
	prestadorRepo model.PrestadorRepo,
) *UsuarioUseCase {
	return &UsuarioUseCase{
		usuarioRepo: usuarioRepo,
		clienteRepo: clienteRepo,
		prestadorRepo: prestadorRepo,
	}
}

type UsuarioRequest struct {
	Nome     string `json:"nome" binding:"required"`
	Telefone string `json:"telefone" binding:"required"`
	Senha    string `json:"senha,omitempty" binding:"required"`
}

type UsuarioResponse struct {
	Nome     string `json:"nome"`
	Telefone string `json:"telefone"`
}

type PrestadorRequest struct {
	Usuario UsuarioRequest	`json:"usuario" binding:"required"`
	Localizacao     string   `json:"localizacao" binding:"required"`
}

type PrestadorResponse struct {
	Nome     string `json:"nome"`
	Telefone string `json:"telefone"`
	Localizacao     string  `json:"localizacao"`
	Disponivel     	bool  	`json:"disponivel"`
}

func (uc *UsuarioUseCase) CriarAdmin(ctx context.Context, request UsuarioRequest) error{
	
	if err := uc.SeTelefoneExiste(ctx, request.Telefone); err != nil {
		return err
	}
	admin, err := model.NewAdmin(request.Nome, request.Telefone, request.Senha)
	if err != nil {
		return err
	}
	return uc.usuarioRepo.Criar(ctx, admin)
}

func (uc *UsuarioUseCase) CriarCliente(ctx context.Context, request UsuarioRequest) error{
	if err := uc.SeTelefoneExiste(ctx, request.Telefone); err != nil {
		return err
	}
	cliente, err := model.NewCliente(request.Nome, request.Telefone, request.Senha)
	if err != nil {
		return err
	} 
	return uc.clienteRepo.Criar(ctx, cliente)
}

func (uc *UsuarioUseCase) CriarPrestador(ctx context.Context, request PrestadorRequest) error{
	if err := uc.SeTelefoneExiste(ctx, request.Usuario.Telefone); err != nil {
		return err
	}
	prestador, err := model.NewPrestador(
		request.Localizacao,
		request.Usuario.Nome,
		request.Usuario.Telefone,
		request.Usuario.Senha)
	if err != nil {
		return err
	}
	return uc.prestadorRepo.Criar(ctx, prestador)
}

func(uc *UsuarioUseCase) ListarTodosUsuarios(ctx context.Context, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]UsuarioResponse, error) {
	usuarios, err := uc.usuarioRepo.ListarTodos(ctx, filters, orderBy, orderDir, limit, offset)
	if err != nil {
		return nil, err
	}

	response := make([]UsuarioResponse, 0, len(usuarios))
	for _, u := range usuarios {
		response = append(response, UsuarioResponse{
			Nome: u.Nome,
			Telefone: u.Telefone,
		})
	}
	return response, nil
}

func(uc *UsuarioUseCase) ListarPrestadores(ctx context.Context, filters map[string]interface{}, statusDisponivel interface{}, orderBy string, orderDir string, limit, offset int) ([]PrestadorResponse, error) {
	prestadores, err := uc.prestadorRepo.Listar(ctx, filters, statusDisponivel, orderBy, orderDir, limit, offset)
	if err != nil {
		return nil, err
	}

	response := make([]PrestadorResponse, 0, len(prestadores))
	for _, p := range prestadores {
		response = append(response, PrestadorResponse{
			Nome: p.Nome,
			Telefone: p.Telefone,
			Localizacao: p.Localizacao,
			Disponivel: p.StatusDisponivel,
		})
	}
	return response, nil
}

func(uc *UsuarioUseCase) SeTelefoneExiste(ctx context.Context, telefone string) error {
		if _, err := uc.usuarioRepo.BuscarPorTelefone(ctx, telefone); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	return errors.New("ja existe usuario com mesmo contacto")
}
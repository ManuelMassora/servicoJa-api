package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
)

type VagaUseCase struct {
	vagaRepo model.VagaRepo
	clienteRep model.ClienteRepo
}

func NewVagaUseCase(vagaRepo model.VagaRepo, clienteRep model.ClienteRepo) *VagaUseCase {
	return &VagaUseCase{vagaRepo: vagaRepo, clienteRep: clienteRep}
}

type VagaRequest struct {
	Titulo      string  `json:"titulo" binding:"required"`
	Descricao   string  `json:"descricao" binding:"required"`
	Localizacao string  `json:"localizacao" binding:"required"`
	Preco       float64 `json:"preco" binding:"required,gte=0"`
	Urgente     bool    `json:"urgente"`
}

type VagaResponse struct {
	ID          uint    `json:"id"`
	Titulo      string  `json:"titulo"`
	Descricao   string  `json:"descricao"`
	Localizacao string  `json:"localizacao"`
	Status      string 	`json:"status"`
	Preco       float64 `json:"preco"`
	Urgente     bool    `json:"urgente"`
	Cliente    	string  `json:"cliente"`
	DataCriacao string  `json:"data_criacao"`
}

func(uc *VagaUseCase) CriarVaga(ctx context.Context, req VagaRequest, idUsuario uint) error {
	cliente, err := uc.clienteRep.BuscarPorUsuarioID(ctx, idUsuario)
	if err != nil {
		return err
	}
	vaga := &model.Vaga{
		Titulo:      req.Titulo,
		Descricao:   req.Descricao,
		Localizacao: req.Localizacao,
		Preco:       req.Preco,
		Status:      model.StatusDisponivel,
		IDCliente:   cliente.ID,
		Urgente:     req.Urgente,
	}
	return uc.vagaRepo.Criar(ctx, vaga)
}

func(uc *VagaUseCase) CancelarVaga(ctx context.Context, id, idUsuario uint) error {
	cliente, err := uc.clienteRep.BuscarPorUsuarioID(ctx, idUsuario)
	if err != nil {
		return err
	}
	vaga, err := uc.vagaRepo.BuscarPorID(ctx, id)
	if err != nil {
		return err
	}
	if vaga.IDCliente != cliente.ID {
		return errors.New("vaga não pertence ao cliente")
	}
	vaga.DeletedAt.Time = time.Now()
	vaga.DeletedAt.Valid = true
	vaga.Status = model.StatusCancelado
	return uc.vagaRepo.Salvar(ctx, vaga)
}

func(uc *VagaUseCase) ListarVagasDisponiveis(ctx context.Context, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]VagaResponse, error) {
	vagas, err := uc.vagaRepo.ListarDisponiveis(ctx, filters, orderBy, orderDir, limit, offset)
	if err != nil {
		return nil, err
	}
	var resp []VagaResponse
	for _, vaga := range vagas {
		clienteNome := ""
		if vaga.Cliente != nil {
			clienteNome = vaga.Cliente.Usuario.Nome
		}
		resp = append(resp, VagaResponse{
			ID:          vaga.ID,
			Titulo:      vaga.Titulo,
			Descricao:   vaga.Descricao,
			Localizacao: vaga.Localizacao,
			Status:      string(vaga.Status),
			Preco:       vaga.Preco,
			Urgente:     vaga.Urgente,
			Cliente:     clienteNome,
			DataCriacao: vaga.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return resp, nil
}

func(uc *VagaUseCase) ListarPorCliente(ctx context.Context, idUsuario uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]VagaResponse, error) {
	cliente, err := uc.clienteRep.BuscarPorUsuarioID(ctx, idUsuario)
	if err != nil {
		return nil, err
	}
	vagas, err := uc.vagaRepo.ListarPorCliente(ctx, cliente.ID, filters, orderBy, orderDir, limit, offset)
	if err != nil {
		return nil, err
	}
	var resp []VagaResponse
	for _, vaga := range vagas {
		clienteNome := ""
		if vaga.Cliente != nil {
			clienteNome = vaga.Cliente.Usuario.Nome
		}
		resp = append(resp, VagaResponse{
			ID:          vaga.ID,
			Titulo:      vaga.Titulo,
			Descricao:   vaga.Descricao,
			Localizacao: vaga.Localizacao,
			Status:      string(vaga.Status),
			Preco:       vaga.Preco,
			Urgente:     vaga.Urgente,
			Cliente:     clienteNome,
			DataCriacao: vaga.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return resp, nil
}
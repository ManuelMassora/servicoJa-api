package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
)

type ServicoUseCase struct {
	r model.ServicoRepo
	prestadorRepo model.PrestadorRepo
	clienteRepo model.ClienteRepo
	agendamentoRepo model.AgendamentoRepo
	vagaRepo model.VagaRepo
}

type ServicoResponse struct {
	ID		  		uint      	`json:"id"`
	Localizacao 	string   	`json:"localizacao"`
	Preco       	float64  	`json:"preco"`
	Status      	string   	`json:"status"`
	IDAgendamento   *uint    	`json:"id_agendamento,omitempty"`
	IDVaga 			*uint 		`json:"id_vaga,omitempty"`
	DataHoraInicio  time.Time 	`json:"data_inicio,omitempty"`
	DataHoraFim     time.Time  	`json:"data_fim,omitempty"`
	Cliente    		uint      	`json:"cliente"`
	Prestador  		uint      	`json:"prestador"`
}

func NewServicoUseCase(r model.ServicoRepo, prestadorRepo model.PrestadorRepo, clienteRepo model.ClienteRepo, agendamentoRepo model.AgendamentoRepo, vagaRepo model.VagaRepo) *ServicoUseCase {
	return &ServicoUseCase{
		r: r,
		prestadorRepo: prestadorRepo,
		clienteRepo: clienteRepo,
		agendamentoRepo: agendamentoRepo,
		vagaRepo: vagaRepo,
	}
}

func (uc *ServicoUseCase) FinalizarServico(ctx context.Context, idServico, idUsuario uint) error {
	servico, err := uc.r.BuscarPorID(ctx, idServico)
	if err != nil {
		return err
	}

	if servico.Cliente.UsuarioID != idUsuario && servico.Prestador.UsuarioID != idUsuario {
		return errors.New("usuário não autorizado a finalizar este serviço")
	}

	if servico.Status == model.StatusConcluido || servico.Status == model.StatusCancelado {
		return nil
	}

	servico.Status = model.StatusConcluido
	servico.DataHoraFim = time.Now()

	return uc.r.Atualizar(ctx, servico)
}

func (uc *ServicoUseCase) CancelarServico(ctx context.Context, idServico, idUsuario uint) error {	
	servico, err := uc.r.BuscarPorID(ctx, idServico)
	if err != nil {
		return err
	}
	if servico.Cliente.UsuarioID != idUsuario && servico.Prestador.UsuarioID != idUsuario {
		return errors.New("usuário não autorizado a finalizar este serviço")
	}
	if servico.Status == model.StatusConcluido || servico.Status == model.StatusCancelado {
		return nil
	}
	servico.Status = model.StatusCancelado
	servico.DataHoraFim = time.Now()
	err = uc.r.Atualizar(ctx, servico)
	if err != nil {
		return err
	}
	return nil
}

func (uc *ServicoUseCase) ListarPorCliente(ctx context.Context, idUsuario uint, filters map[string]interface{}, orderBy, orderDir string, limit, offset int) ([]ServicoResponse, error) {
	cliente, err := uc.clienteRepo.BuscarPorUsuarioID(ctx, idUsuario)
	if err != nil {
		return nil, err
	}
	servicos, err := uc.r.ListarPorCliente(ctx, cliente.ID, filters, orderBy, orderDir, limit, offset)
	if err != nil {
		return nil, err
	}
	if len(servicos) == 0 {
		return []ServicoResponse{}, nil
	}
	var resp []ServicoResponse
	for _, s := range servicos {
		resp = append(resp, ServicoResponse{
			ID:            s.ID,
			Localizacao:   s.Localizacao,
			Preco:         s.Preco,
			Status:        string(s.Status),
			IDAgendamento: s.IDAgendamento,
			IDVaga:        s.IDVaga,
			DataHoraInicio: s.DataHoraInicio,
			DataHoraFim:    s.DataHoraFim,
			Cliente:        s.IDCliente,
			Prestador:      s.IDPrestador,
		})
	}
	return resp, nil
}

func (uc *ServicoUseCase) ListarPorPrestador(ctx context.Context, idUsuario uint, filters map[string]interface{}, orderBy, orderDir string, limit, offset int) ([]ServicoResponse, error) {
	prestador, err := uc.prestadorRepo.BuscarPorUsuarioID(ctx, idUsuario)
	if err != nil {
		return nil, err
	}
	servicos, err := uc.r.ListarPorPrestador(ctx, prestador.ID, filters, orderBy, orderDir, limit, offset)
	if err != nil {
		return nil, err
	}
	if len(servicos) == 0 {
		return []ServicoResponse{}, nil
	}
	var resp []ServicoResponse
	for _, s := range servicos {
		resp = append(resp, ServicoResponse{
			ID:            s.ID,
			Localizacao:   s.Localizacao,
			Preco:         s.Preco,
			Status:        string(s.Status),
			IDAgendamento: s.IDAgendamento,
			IDVaga:        s.IDVaga,
			DataHoraInicio: s.DataHoraInicio,
			DataHoraFim:    s.DataHoraFim,
			Cliente:        s.IDCliente,
			Prestador:      s.IDPrestador,
		})
	}
	return resp, nil
}

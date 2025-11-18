package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
)

type AgendamentoUC struct {
	r model.AgendamentoRepo
	clienteRepo model.ClienteRepo
	pretadorRepo model.PrestadorRepo
}

func NewAgendamentoUC(
	r model.AgendamentoRepo,
	clienteRepo model.ClienteRepo,
	pretadorRepo model.PrestadorRepo,
) *AgendamentoUC {
	return &AgendamentoUC{r: r,
		clienteRepo: clienteRepo,
		pretadorRepo: pretadorRepo,
	}
}

type AgendamentoRequest struct {
	Detalhe 	string 		`json:"detalhe" binding:"required"`
	IDCatalogo  uint  		`json:"id_catalogo" binding:"required"`
	DataHora 	time.Time   `json:"datahora" binding:"required"`
}

type AgendamentoResponse struct {
	Detalhe 	string 		`json:"detalhe"`
	Catalogo	string		`json:"catalogo"`
	Cliente		string		`json:"cliente"`
	DataHora 	time.Time   `json:"datahora"`
	Status 		string   	`json:"status"`
}

func(uc *AgendamentoUC) Criar(ctx context.Context, req *AgendamentoRequest, idCliente uint) error {
	cliente, err := uc.clienteRepo.BuscarPorUsuarioID(ctx, idCliente)
	if err != nil || cliente == nil {
		return errors.New("cliente não encontrado")
	}
	return uc.r.Criar(ctx, &model.Agendamento{
		Detalhe: req.Detalhe,
		IDCatalogo: req.IDCatalogo,
		IDCliente: cliente.ID,
		DataHora: req.DataHora,
		Status: "Pendente",
	})
}

func (uc *AgendamentoUC) Buscar(ctx context.Context, id uint, idUsuario uint) (*AgendamentoResponse, error) {
	
	agendamento, err := uc.r.BuscarPorID(ctx, id)
	if err != nil {
		return nil, err
	}
	clienteIDUsuario := agendamento.Cliente.Usuario.ID 	
	prestadorIDUsuario := agendamento.Catalogo.Prestador.Usuario.ID
	
	if idUsuario != clienteIDUsuario && idUsuario != prestadorIDUsuario {
		return nil, errors.New("acesso negado: você não é o cliente nem o prestador deste agendamento")
	}

	return &AgendamentoResponse{
		Detalhe: agendamento.Detalhe,
		Catalogo: agendamento.Catalogo.Nome,
		Cliente: agendamento.Cliente.Usuario.Nome,
		DataHora: agendamento.DataHora,
		Status: agendamento.Status,
	}, nil
}

func (uc *AgendamentoUC) Listar(ctx context.Context, idUsuario uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]AgendamentoResponse, error) {
	
	cliente, _ := uc.clienteRepo.BuscarPorUsuarioID(ctx, idUsuario) 

	prestador, _ := uc.pretadorRepo.BuscarPorUsuarioID(ctx, idUsuario) 
	if cliente != nil {
		
		filters["id_cliente"] = cliente.ID 
	} else if prestador != nil {   
        filters["id_prestador"] = prestador.ID 
	} else {
        
		return nil, errors.New("acesso negado: usuário não é cliente nem prestador")
	}
	agendamentos, err := uc.r.Listar(ctx, filters, orderBy, orderDir, limit, offset)
	if err != nil {
		return nil, err
	}
	var resp []AgendamentoResponse
	for _, agendamento := range agendamentos {
		resp = append(resp, AgendamentoResponse{
			Detalhe: agendamento.Detalhe,
			
			Catalogo: agendamento.Catalogo.Nome,
			Cliente:  agendamento.Cliente.Usuario.Nome,
			DataHora: agendamento.DataHora,
			Status:   agendamento.Status,
		})
	}
	return resp, nil
}
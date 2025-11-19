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
	catalogoRepo model.CatalogoRepo
	servico model.ServicoRepo
}

func NewAgendamentoUC(
	r model.AgendamentoRepo,
	clienteRepo model.ClienteRepo,
	pretadorRepo model.PrestadorRepo,
	catalogoRepo model.CatalogoRepo,
	servico model.ServicoRepo,
) *AgendamentoUC {
	return &AgendamentoUC{r: r,
		clienteRepo: clienteRepo,
		pretadorRepo: pretadorRepo,
		catalogoRepo: catalogoRepo,
		servico: servico,
	}
}

type AgendamentoRequest struct {
	Detalhe 	string 		`json:"detalhe" binding:"required"`
	IDCatalogo  uint  		`json:"id_catalogo" binding:"required"`
	DataHora 	time.Time   `json:"datahora" binding:"required"`
}

type AgendamentoResponse struct {
	ID			uint		`json:"id"`
	Detalhe 	string 		`json:"detalhe"`
	Catalogo	string		`json:"catalogo"`
	Cliente		string		`json:"cliente"`
	Prestador	string		`json:"prestador"`
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
		Status: "PENDENTE",
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
		ID: agendamento.ID,
		Detalhe: agendamento.Detalhe,
		Catalogo: agendamento.Catalogo.Nome,
		Cliente: agendamento.Cliente.Usuario.Nome,
		Prestador: agendamento.Catalogo.Prestador.Usuario.Nome,
		DataHora: agendamento.DataHora,
		Status: agendamento.Status,
	}, nil
}

func (uc *AgendamentoUC) Aceitar(ctx context.Context, id uint, idUsuario uint) error {
	agendamento, err := uc.r.BuscarPorID(ctx, id)
	if err != nil {
		return err
	}
	prestadorIDUsuario := agendamento.Catalogo.Prestador.Usuario.ID
	if idUsuario != prestadorIDUsuario {
		return errors.New("acesso negado: você não é o prestador deste agendamento")
	}
	if agendamento.Status == "EM_ANDAMENTO" {
		return nil
	}
	servico := &model.Servico{
		IDAgendamento: &id,
		Localizacao: agendamento.Catalogo.Localizacao,
		Preco: agendamento.Catalogo.PrecoBase,
		Status: model.StatusEmAndamento,
		IDCliente: agendamento.IDCliente,
		IDPrestador: agendamento.Catalogo.IDPrestador,
		DataHoraInicio: time.Now(),
	}
	err = uc.servico.Criar(ctx, servico)
	if err != nil {
		return err
	}
	return uc.r.AtualizarStatus(ctx, id, "EM_ANDAMENTO")
}

func (uc *AgendamentoUC) Recusar(ctx context.Context, id uint, idUsuario uint) error {
	agendamento, err := uc.r.BuscarPorID(ctx, id)
	if err != nil {
		return err
	}
	prestadorIDUsuario := agendamento.Catalogo.Prestador.Usuario.ID
	if idUsuario != prestadorIDUsuario {
		return errors.New("acesso negado: você não é o prestador deste agendamento")
	}
	return uc.r.AtualizarStatus(ctx, id, "RECUSADO")
}

func (uc *AgendamentoUC) Cancelar(ctx context.Context, id uint, idUsuario uint) error {
	agendamento, err := uc.r.BuscarPorID(ctx, id)
	if err != nil {
		return err
	}
	clienteIDUsuario := agendamento.Cliente.Usuario.ID
	if idUsuario != clienteIDUsuario {
		return errors.New("acesso negado: você não é o cliente deste agendamento")
	}
	if agendamento.Status == "EM_ANDAMENTO" {
		return nil
	}
	return uc.r.AtualizarStatus(ctx, id, "CANCELADO")
}

func (uc *AgendamentoUC) Listar(ctx context.Context, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]AgendamentoResponse, error) {
	agendamentos, err := uc.r.Listar(ctx, filters, orderBy, orderDir, limit, offset)
	if err != nil {
		return nil, err
	}
	var resp []AgendamentoResponse
	for _, agendamento := range agendamentos {
		resp = append(resp, AgendamentoResponse{
			ID: agendamento.ID,
			Detalhe: agendamento.Detalhe,
			Catalogo: agendamento.Catalogo.Nome,
			Cliente:  agendamento.Cliente.Usuario.Nome,
			Prestador: agendamento.Catalogo.Prestador.Usuario.Nome,
			DataHora: agendamento.DataHora,
			Status:   agendamento.Status,
		})
	}
	return resp, nil
}

func (uc *AgendamentoUC) ListarPorClienteID(ctx context.Context, idUsuario uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]AgendamentoResponse, error) {
	cliente, err := uc.clienteRepo.BuscarPorUsuarioID(ctx, idUsuario)
	if err != nil {
		return nil, err
	}
	agendamentos, err := uc.r.ListarPorClienteID(ctx, cliente.ID, filters, orderBy, orderDir, limit, offset)
	if err != nil {
		return nil, err
	}
	var resp []AgendamentoResponse
	for _, agendamento := range agendamentos {
		resp = append(resp, AgendamentoResponse{
			ID: agendamento.ID,
			Detalhe: agendamento.Detalhe,
			Catalogo: agendamento.Catalogo.Nome,
			Cliente:  agendamento.Cliente.Usuario.Nome,
			Prestador: agendamento.Catalogo.Prestador.Usuario.Nome,
			DataHora: agendamento.DataHora,
			Status:   agendamento.Status,
		})
	}
	return resp, nil
}

func (uc *AgendamentoUC) ListarPorCatalogID(ctx context.Context, idUsuario, idCatalogo uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]AgendamentoResponse, error) {
	prestador, err := uc.pretadorRepo.BuscarPorUsuarioID(ctx, idUsuario)
	if err != nil {
		return nil, err
	}
	catalogo, err := uc.catalogoRepo.FindByID(ctx, idCatalogo)
	if err != nil {
		return nil, err
	}
	if catalogo.IDPrestador != prestador.ID {
		return nil, errors.New("acesso negado: você não é o prestador deste catálogo")
	}
	agendamentos, err := uc.r.ListarPorCatalogID(ctx, catalogo.ID, filters, orderBy, orderDir, limit, offset)
	if err != nil {
		return nil, err
	}
	var resp []AgendamentoResponse
	for _, agendamento := range agendamentos {
		resp = append(resp, AgendamentoResponse{
			ID: agendamento.ID,
			Detalhe: agendamento.Detalhe,
			Catalogo: agendamento.Catalogo.Nome,
			Cliente:  agendamento.Cliente.Usuario.Nome,
			Prestador: agendamento.Catalogo.Prestador.Usuario.Nome,
			DataHora: agendamento.DataHora,
			Status:   agendamento.Status,
		})
	}
	return resp, nil
}
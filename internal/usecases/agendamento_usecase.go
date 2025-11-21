package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
)

type AgendamentoUC struct {
	r model.AgendamentoRepo
	catalogoRepo model.CatalogoRepo
	servico model.ServicoRepo
	notifacaoRepo model.NotificacaoRepo
	usuarioRepo model.UsuarioRepo
}

func NewAgendamentoUC(
	r model.AgendamentoRepo,
	catalogoRepo model.CatalogoRepo,
	servico model.ServicoRepo,
	notifacaoRepo model.NotificacaoRepo,
	usuarioRepo model.UsuarioRepo,
) *AgendamentoUC {
	return &AgendamentoUC{r: r,
		catalogoRepo: catalogoRepo,
		servico: servico,
		notifacaoRepo: notifacaoRepo,
		usuarioRepo: usuarioRepo,
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
	catalogo, err := uc.catalogoRepo.FindByID(ctx, req.IDCatalogo)
	if err != nil {
		return err
	}
	// // Verificar se o prestador está disponível na data/hora solicitada
	// if !catalogo.Prestador.DisponivelNaDataHora(ctx, req.DataHora) {
	// 	return errors.New("o prestador não está disponível na data/hora solicitada")
	// }
	// // Criar a notificação para o prestador
	err = uc.notifacaoRepo.Enviar(ctx, &model.Notificacao{
		IDUsuario: catalogo.Prestador.IDUsuario,
		Titulo: "Novo Agendamento",
		Mensagem: "Você tem um novo agendamento para o serviço: " + catalogo.Nome,
	})
	if err != nil {
		return err
	}
	err = uc.usuarioRepo.IncrementarNotificacoesNovas(ctx, catalogo.Prestador.IDUsuario)
	if err != nil {
		return err
	}
	return uc.r.Criar(ctx, &model.Agendamento{
		Detalhe: req.Detalhe,
		IDCatalogo: req.IDCatalogo,
		IDCliente: idCliente,
		DataHora: req.DataHora,
		Status: "PENDENTE",
	})
}

func (uc *AgendamentoUC) Buscar(ctx context.Context, id uint, idUsuario uint) (*AgendamentoResponse, error) {
	
	agendamento, err := uc.r.BuscarPorID(ctx, id)
	if err != nil {
		return nil, err
	}
	clienteIDUsuario := agendamento.Cliente.IDUsuario
	prestadorIDUsuario := agendamento.Catalogo.Prestador.IDUsuario
	
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
	err = uc.notifacaoRepo.Enviar(ctx, &model.Notificacao{
		IDUsuario: agendamento.IDCliente,
		Titulo: "Resposta ao Agendamento",
		Mensagem: "Seu agendamento foi aceito para o serviço: " + agendamento.Catalogo.Nome,
	})
	if err != nil {
		return err
	}
	err = uc.usuarioRepo.IncrementarNotificacoesNovas(ctx, agendamento.IDCliente)
	if err != nil {
		return err
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
	err = uc.notifacaoRepo.Enviar(ctx, &model.Notificacao{
		IDUsuario: agendamento.IDCliente,
		Titulo: "Resposta ao Agendamento",
		Mensagem: "Seu agendamento foi recusado para o serviço: " + agendamento.Catalogo.Nome,
	})
	if err != nil {
		return err
	}
	err = uc.usuarioRepo.IncrementarNotificacoesNovas(ctx, agendamento.IDCliente)
	if err != nil {
		return err
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
	agendamentos, err := uc.r.ListarPorClienteID(ctx, idUsuario, filters, orderBy, orderDir, limit, offset)
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
	catalogo, err := uc.catalogoRepo.FindByID(ctx, idCatalogo)
	if err != nil {
		return nil, err
	}
	if catalogo.IDPrestador != idUsuario {
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
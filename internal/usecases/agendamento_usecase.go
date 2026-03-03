package usecases

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
)

type AgendamentoUC struct {
	r               model.AgendamentoRepo
	catalogoRepo    model.CatalogoRepo
	servico         model.ServicoRepo
	notifacaoRepo   model.NotificacaoRepo
	usuarioRepo     model.UsuarioRepo
	anexoImagemRepo model.AnexoImagemRepo
	pagamentoRepo   model.PagamentoRepo
	pagamentoUC     PagamentoUseCase
}

func NewAgendamentoUC(
	r model.AgendamentoRepo,
	catalogoRepo model.CatalogoRepo,
	servico model.ServicoRepo,
	notifacaoRepo model.NotificacaoRepo,
	usuarioRepo model.UsuarioRepo,
	anexoImagemRepo model.AnexoImagemRepo,
	pagamentoRepo model.PagamentoRepo,
	pagamentoUC PagamentoUseCase,
) *AgendamentoUC {
	return &AgendamentoUC{r: r,
		catalogoRepo:    catalogoRepo,
		servico:         servico,
		notifacaoRepo:   notifacaoRepo,
		usuarioRepo:     usuarioRepo,
		anexoImagemRepo: anexoImagemRepo,
		pagamentoRepo:   pagamentoRepo,
		pagamentoUC:     pagamentoUC,
	}
}

type AgendamentoRequest struct {
	Detalhe           string    `json:"detalhe" form:"detalhe" binding:"required"`
	IDCatalogo        uint      `json:"id_catalogo" form:"id_catalogo" binding:"required"`
	DataHora          time.Time `json:"datahora" form:"datahora" binding:"required"`
	Localizacao       string    `json:"localizacao" form:"localizacao" binding:"required"`
	Latitude          float64   `json:"latitude" form:"latitude" binding:"required"`
	Longitude         float64   `json:"longitude" form:"longitude" binding:"required"`
	Anexos            []string  `binding:"-"`
	TelefonePagamento string    `json:"telefone_pagamento" form:"telefone_pagamento" binding:"required"`
}

type AgendamentoResponse struct {
	ID          uint                          `json:"id"`
	Detalhe     string                        `json:"detalhe"`
	Catalogo    string                        `json:"catalogo"`
	Cliente     ClientesAgendamentosResponse  `json:"cliente"`
	Prestador   PrestadorAgendamentosResponse `json:"prestador"`
	DataHora    time.Time                     `json:"datahora"`
	Status      string                        `json:"status"`
	Localizacao string                        `json:"localizacao"`
	Latitude    float64                       `json:"latitude"`
	Longitude   float64                       `json:"longitude"`
	Anexos      []string                      `json:"anexos"`
}

type AgendamentoGroupCategoriaResponse struct {
	IDCatalogo uint                          `json:"id_catalogo"`
	Catalogo   string                        `json:"catalogo"`
	SubGrupos  []AgendamentoSubGrupoResponse `json:"subgrupos"`
}

type AgendamentoSubGrupoResponse struct {
	DataGrupo    string                `json:"dataGrupo"`
	Agendamentos []AgendamentoResponse `json:"agendamentos"`
}

type ClientesAgendamentosResponse struct {
	ID          uint   `json:"id"`
	ClienteNome string `json:"nome"`
}
type PrestadorAgendamentosResponse struct {
	ID            uint   `json:"id"`
	PrestadorNome string `json:"nome"`
}

func (uc *AgendamentoUC) Criar(ctx context.Context, req *AgendamentoRequest, idCliente uint) error {
	usuario, err := uc.usuarioRepo.BuscarPorID(ctx, idCliente)
	if err != nil {
		return err
	}

	if usuario.SuspensoAte != nil && usuario.SuspensoAte.After(time.Now()) {
		return errors.New("sua conta está suspensa até " + usuario.SuspensoAte.Format("02/01/2006 15:04"))
	}

	catalogo, err := uc.catalogoRepo.FindByID(ctx, req.IDCatalogo)
	if err != nil {
		return err
	}

	agendamento := &model.Agendamento{
		Detalhe:     req.Detalhe,
		IDCatalogo:  req.IDCatalogo,
		IDCliente:   idCliente,
		DataHora:    req.DataHora,
		Status:      string(model.StatusAguardandoPagamento),
		Localizacao: req.Localizacao,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
	}

	agendamentoSave, err := uc.r.Criar(ctx, agendamento)
	if err != nil {
		return err
	}

	// Criar registro de pagamento pendente
	pagamento := &model.Pagamento{
		IDAgendamento: &agendamentoSave.ID,
		IDCliente:     idCliente,
		Valor:         catalogo.ValorFixo,
		Status:        model.StatusPendente, // Esperando C2B
		IDPrestador:   &catalogo.IDPrestador,
		Referencia:    "REF" + strconv.FormatUint(uint64(agendamentoSave.ID), 10),
	}
	if err := uc.pagamentoRepo.Criar(ctx, pagamento); err != nil {
		return err
	}

	// Iniciar o processo de pagamento via M-Pesa
	_ = uc.pagamentoUC.IniciarPagamentoC2B(ctx, pagamento.ID, req.TelefonePagamento)
	// if err != nil {
	// 	return err
	// }

	for _, anexoURL := range req.Anexos {
		anexo := &model.AnexoImagem{
			URL:           anexoURL,
			AgendamentoID: &agendamentoSave.ID,
		}
		if err := uc.anexoImagemRepo.Create(ctx, anexo); err != nil {
			// In a real application, you might want to handle the rollback of the agendamento creation
			return err
		}
	}

	err = uc.notifacaoRepo.Enviar(ctx, &model.Notificacao{
		IDUsuario: catalogo.Prestador.IDUsuario,
		Titulo:    "Novo Agendamento",
		Mensagem:  "Você tem um novo agendamento para o serviço: " + catalogo.Nome,
	})
	if err != nil {
		return err
	}
	err = uc.usuarioRepo.IncrementarNotificacoesNovas(ctx, catalogo.ID)
	if err != nil {
		return err
	}
	if err := uc.catalogoRepo.IncrementarAgendamentosNovos(ctx, catalogo.IDPrestador); err != nil {
		return err
	}
	return nil
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

	anexos, err := uc.anexoImagemRepo.FindByAgendamentoID(ctx, agendamento.ID)
	if err != nil {
		return nil, err
	}
	return &AgendamentoResponse{
		ID:       agendamento.ID,
		Detalhe:  agendamento.Detalhe,
		Catalogo: agendamento.Catalogo.Nome,
		Cliente: ClientesAgendamentosResponse{
			ID:          agendamento.IDCliente,
			ClienteNome: agendamento.Cliente.Nome,
		},
		Prestador: PrestadorAgendamentosResponse{
			ID:            agendamento.Catalogo.IDPrestador,
			PrestadorNome: agendamento.Catalogo.Prestador.Nome,
		},
		DataHora:    agendamento.DataHora,
		Status:      agendamento.Status,
		Localizacao: agendamento.Localizacao,
		Latitude:    agendamento.Latitude,
		Longitude:   agendamento.Longitude,
		Anexos: func() []string {
			var urls []string
			for _, anexo := range anexos {
				urls = append(urls, anexo.URL)
			}
			return urls
		}(),
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
		Titulo:    "Resposta ao Agendamento",
		Mensagem:  "Seu agendamento foi aceito para o serviço: " + agendamento.Catalogo.Nome,
	})
	if err != nil {
		return err
	}
	err = uc.usuarioRepo.IncrementarNotificacoesNovas(ctx, agendamento.IDCliente)
	if err != nil {
		return err
	}
	servico := &model.Servico{
		IDAgendamento:  &id,
		Localizacao:    agendamento.Localizacao,
		Latitude:       agendamento.Latitude,
		Longitude:      agendamento.Longitude,
		Status:         model.StatusEmAndamento,
		IDCliente:      agendamento.IDCliente,
		IDPrestador:    agendamento.Catalogo.IDPrestador,
		DataHoraInicio: time.Now().UTC(),
	}

	switch agendamento.Catalogo.TipoPreco {
	case "fixo":
		servico.Preco = agendamento.Catalogo.ValorFixo
	case "por_hora":
		servico.Preco = 0.0 // Initial price for hourly services, will be calculated at finalization
	default:
		return errors.New("tipo de preço de catálogo inválido")
	}

	servicoSave, err := uc.servico.Criar(ctx, servico)
	if err != nil {
		return err
	}

	// Associar ID do serviço ao pagamento
	p, err := uc.pagamentoRepo.BuscarPorAgendamento(ctx, id)
	if err == nil && p != nil {
		_ = uc.pagamentoRepo.AtualizarIDServico(ctx, p.ID, servicoSave.ID)
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
		Titulo:    "Resposta ao Agendamento",
		Mensagem:  "Seu agendamento foi recusado para o serviço: " + agendamento.Catalogo.Nome,
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
	var agendamentosIDs []uint
	for _, agendamento := range agendamentos {
		agendamentosIDs = append(agendamentosIDs, agendamento.ID)
	}
	anexos, err := uc.anexoImagemRepo.FindByAgendamentoIDs(ctx, agendamentosIDs)
	if err != nil {
		return nil, err
	}
	anexosPorAgendamentoMap := make(map[uint][]string)
	for _, anexo := range anexos {
		anexosPorAgendamentoMap[*anexo.AgendamentoID] = append(anexosPorAgendamentoMap[*anexo.AgendamentoID], anexo.URL)
	}
	var resp []AgendamentoResponse
	for _, agendamento := range agendamentos {
		urls := anexosPorAgendamentoMap[agendamento.ID]
		resp = append(resp, AgendamentoResponse{
			ID:       agendamento.ID,
			Detalhe:  agendamento.Detalhe,
			Catalogo: agendamento.Catalogo.Nome,
			Cliente: ClientesAgendamentosResponse{
				ID:          agendamento.IDCliente,
				ClienteNome: agendamento.Cliente.Nome,
			},
			Prestador: PrestadorAgendamentosResponse{
				ID:            agendamento.Catalogo.IDPrestador,
				PrestadorNome: agendamento.Catalogo.Prestador.Nome,
			},
			DataHora:    agendamento.DataHora,
			Status:      agendamento.Status,
			Localizacao: agendamento.Localizacao,
			Latitude:    agendamento.Latitude,
			Longitude:   agendamento.Longitude,
			Anexos:      urls,
		})
	}
	return resp, nil
}

func (uc *AgendamentoUC) ListarPorClienteID(
	ctx context.Context,
	idUsuario uint,
	filters map[string]interface{},
	orderBy string,
	orderDir string,
	limit, offset int,
) ([]AgendamentoGroupCategoriaResponse, error) {

	agendamentos, err := uc.r.ListarPorClienteID(ctx, idUsuario, filters, orderBy, orderDir, limit, offset)
	if err != nil {
		return nil, err
	}

	var agendamentosIDs []uint
	for _, ag := range agendamentos {
		agendamentosIDs = append(agendamentosIDs, ag.ID)
	}

	anexos, err := uc.anexoImagemRepo.FindByAgendamentoIDs(ctx, agendamentosIDs)
	if err != nil {
		return nil, err
	}

	anexosPorAgendamentoMap := make(map[uint][]string)
	for _, anexo := range anexos {
		anexosPorAgendamentoMap[*anexo.AgendamentoID] = append(anexosPorAgendamentoMap[*anexo.AgendamentoID], anexo.URL)
	}

	var resp []AgendamentoGroupCategoriaResponse

	mapCatalogo := make(map[string]*AgendamentoGroupCategoriaResponse)

	for _, ag := range agendamentos {
		dataGrupo := ag.CreatedAt.Format("2006-01-02")

		agResp := AgendamentoResponse{
			ID:       ag.ID,
			Detalhe:  ag.Detalhe,
			Catalogo: ag.Catalogo.Nome,
			Cliente: ClientesAgendamentosResponse{
				ID:          ag.IDCliente,
				ClienteNome: ag.Cliente.Nome,
			},
			Prestador: PrestadorAgendamentosResponse{
				ID:            ag.Catalogo.IDPrestador,
				PrestadorNome: ag.Catalogo.Prestador.Nome,
			},
			DataHora:    ag.CreatedAt,
			Status:      ag.Status,
			Localizacao: ag.Localizacao,
			Latitude:    ag.Latitude,
			Longitude:   ag.Longitude,
			Anexos:      anexosPorAgendamentoMap[ag.ID],
		}

		if _, ok := mapCatalogo[ag.Catalogo.Nome]; !ok {
			resp = append(resp, AgendamentoGroupCategoriaResponse{
				IDCatalogo: ag.Catalogo.ID,
				Catalogo:   ag.Catalogo.Nome,
				SubGrupos:  []AgendamentoSubGrupoResponse{},
			})
			mapCatalogo[ag.Catalogo.Nome] = &resp[len(resp)-1]
		}

		grupoCatalogo := mapCatalogo[ag.Catalogo.Nome]

		var subGrupo *AgendamentoSubGrupoResponse
		for i := range grupoCatalogo.SubGrupos {
			if grupoCatalogo.SubGrupos[i].DataGrupo == dataGrupo {
				subGrupo = &grupoCatalogo.SubGrupos[i]
				break
			}
		}

		if subGrupo == nil {
			grupoCatalogo.SubGrupos = append(grupoCatalogo.SubGrupos, AgendamentoSubGrupoResponse{
				DataGrupo:    dataGrupo,
				Agendamentos: []AgendamentoResponse{agResp},
			})

		} else {
			subGrupo.Agendamentos = append(subGrupo.Agendamentos, agResp)
		}
	}

	return resp, nil
}

func (uc *AgendamentoUC) ListarPorPrestadorIDAgrupado(
	ctx context.Context,
	idUsuario uint,
	filters map[string]interface{},
	orderBy string,
	orderDir string,
	limit, offset int,
) ([]AgendamentoGroupCategoriaResponse, error) {

	agendamentos, err := uc.r.ListarPorPrestadorID(ctx, idUsuario, filters, orderBy, orderDir, limit, offset)
	if err != nil {
		return nil, err
	}

	var agendamentosIDs []uint
	for _, ag := range agendamentos {
		agendamentosIDs = append(agendamentosIDs, ag.ID)
	}

	anexos, err := uc.anexoImagemRepo.FindByAgendamentoIDs(ctx, agendamentosIDs)
	if err != nil {
		return nil, err
	}

	anexosPorAgendamentoMap := make(map[uint][]string)
	for _, anexo := range anexos {
		anexosPorAgendamentoMap[*anexo.AgendamentoID] = append(anexosPorAgendamentoMap[*anexo.AgendamentoID], anexo.URL)
	}

	// Resultado final
	var resp []AgendamentoGroupCategoriaResponse

	// Map auxiliar para localizar catálogo já criado
	mapCatalogo := make(map[string]*AgendamentoGroupCategoriaResponse)

	for _, ag := range agendamentos {

		// define chave do subgrupo (somente dia)
		dataGrupo := ag.CreatedAt.Format("2006-01-02")

		// monta response do agendamento
		agResp := AgendamentoResponse{
			ID:       ag.ID,
			Detalhe:  ag.Detalhe,
			Catalogo: ag.Catalogo.Nome,
			Cliente: ClientesAgendamentosResponse{
				ID:          ag.IDCliente,
				ClienteNome: ag.Cliente.Nome,
			},
			Prestador: PrestadorAgendamentosResponse{
				ID:            ag.Catalogo.IDPrestador,
				PrestadorNome: ag.Catalogo.Prestador.Nome,
			},
			DataHora:    ag.CreatedAt,
			Status:      ag.Status,
			Localizacao: ag.Localizacao,
			Latitude:    ag.Latitude,
			Longitude:   ag.Longitude,
			Anexos:      anexosPorAgendamentoMap[ag.ID],
		}

		// se ainda não existe o grupo do catálogo, cria
		if _, ok := mapCatalogo[ag.Catalogo.Nome]; !ok {
			resp = append(resp, AgendamentoGroupCategoriaResponse{
				IDCatalogo: ag.Catalogo.ID,
				Catalogo:   ag.Catalogo.Nome,
				SubGrupos:  []AgendamentoSubGrupoResponse{},
			})
			mapCatalogo[ag.Catalogo.Nome] = &resp[len(resp)-1]
		}

		grupoCatalogo := mapCatalogo[ag.Catalogo.Nome]

		// verifica se já existe subgrupo com esse dia
		var subGrupo *AgendamentoSubGrupoResponse
		for i := range grupoCatalogo.SubGrupos {
			if grupoCatalogo.SubGrupos[i].DataGrupo == dataGrupo {
				subGrupo = &grupoCatalogo.SubGrupos[i]
				break
			}
		}

		// se não existir, cria
		if subGrupo == nil {
			grupoCatalogo.SubGrupos = append(grupoCatalogo.SubGrupos, AgendamentoSubGrupoResponse{
				DataGrupo:    dataGrupo,
				Agendamentos: []AgendamentoResponse{agResp},
			})

		} else {
			subGrupo.Agendamentos = append(subGrupo.Agendamentos, agResp)
		}
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
	var agendamentosIDs []uint
	for _, agendamento := range agendamentos {
		agendamentosIDs = append(agendamentosIDs, agendamento.ID)
	}
	anexos, err := uc.anexoImagemRepo.FindByAgendamentoIDs(ctx, agendamentosIDs)
	if err != nil {
		return nil, err
	}
	anexosPorAgendamentoMap := make(map[uint][]string)
	for _, anexo := range anexos {
		anexosPorAgendamentoMap[*anexo.AgendamentoID] = append(anexosPorAgendamentoMap[*anexo.AgendamentoID], anexo.URL)
	}
	var resp []AgendamentoResponse
	for _, agendamento := range agendamentos {
		urls := anexosPorAgendamentoMap[agendamento.ID]
		resp = append(resp, AgendamentoResponse{
			ID:       agendamento.ID,
			Detalhe:  agendamento.Detalhe,
			Catalogo: agendamento.Catalogo.Nome,
			Cliente: ClientesAgendamentosResponse{
				ID:          agendamento.IDCliente,
				ClienteNome: agendamento.Cliente.Nome,
			},
			Prestador: PrestadorAgendamentosResponse{
				ID:            agendamento.Catalogo.IDPrestador,
				PrestadorNome: agendamento.Catalogo.Prestador.Nome,
			},
			DataHora:    agendamento.DataHora,
			Status:      agendamento.Status,
			Localizacao: agendamento.Localizacao,
			Latitude:    agendamento.Latitude,
			Longitude:   agendamento.Longitude,
			Anexos:      urls,
		})
	}
	return resp, nil
}

func (uc *AgendamentoUC) ListarPorLocalizacao(ctx context.Context, userID uint, latitude, longitude, radius float64, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]AgendamentoResponse, error) {
	agendamentos, err := uc.r.FindByLocation(ctx, userID, latitude, longitude, radius, filters, orderBy, orderDir, limit, offset)
	if err != nil {
		return nil, err
	}
	var agendamentosIDs []uint
	for _, agendamento := range agendamentos {
		agendamentosIDs = append(agendamentosIDs, agendamento.ID)
	}
	anexos, err := uc.anexoImagemRepo.FindByAgendamentoIDs(ctx, agendamentosIDs)
	if err != nil {
		return nil, err
	}
	anexosPorAgendamentoMap := make(map[uint][]string)
	for _, anexo := range anexos {
		anexosPorAgendamentoMap[*anexo.AgendamentoID] = append(anexosPorAgendamentoMap[*anexo.AgendamentoID], anexo.URL)
	}
	var resp []AgendamentoResponse
	for _, agendamento := range agendamentos {
		urls := anexosPorAgendamentoMap[agendamento.ID]
		resp = append(resp, AgendamentoResponse{
			ID:       agendamento.ID,
			Detalhe:  agendamento.Detalhe,
			Catalogo: agendamento.Catalogo.Nome,
			Cliente: ClientesAgendamentosResponse{
				ID:          agendamento.IDCliente,
				ClienteNome: agendamento.Cliente.Nome,
			},
			Prestador: PrestadorAgendamentosResponse{
				ID:            agendamento.Catalogo.IDPrestador,
				PrestadorNome: agendamento.Catalogo.Prestador.Nome,
			},
			DataHora:    agendamento.DataHora,
			Status:      agendamento.Status,
			Localizacao: agendamento.Localizacao,
			Latitude:    agendamento.Latitude,
			Longitude:   agendamento.Longitude,
			Anexos:      urls,
		})
	}
	return resp, nil
}

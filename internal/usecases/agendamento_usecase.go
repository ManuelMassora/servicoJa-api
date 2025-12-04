package usecases

import (
	"context"
	"errors"
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
}

func NewAgendamentoUC(
	r model.AgendamentoRepo,
	catalogoRepo model.CatalogoRepo,
	servico model.ServicoRepo,
	notifacaoRepo model.NotificacaoRepo,
	usuarioRepo model.UsuarioRepo,
	anexoImagemRepo model.AnexoImagemRepo,
) *AgendamentoUC {
	return &AgendamentoUC{r: r,
		catalogoRepo:    catalogoRepo,
		servico:         servico,
		notifacaoRepo:   notifacaoRepo,
		usuarioRepo:     usuarioRepo,
		anexoImagemRepo: anexoImagemRepo,
	}
}

type AgendamentoRequest struct {
	Detalhe 	string 		`json:"detalhe" form:"detalhe" binding:"required"`
	IDCatalogo  uint  		`json:"id_catalogo" form:"id_catalogo" binding:"required"`
	DataHora 	time.Time   `json:"datahora" form:"datahora" binding:"required"`
	Localizacao string   	`json:"localizacao" form:"localizacao" binding:"required"`
	Latitude    float64  	`json:"latitude" form:"latitude" binding:"required"`
	Longitude   float64  	`json:"longitude" form:"longitude" binding:"required"`
	Anexos      []string 	`binding:"-"`
}

type AgendamentoResponse struct {
	ID			uint		`json:"id"`
	Detalhe 	string 		`json:"detalhe"`
	Catalogo	string		`json:"catalogo"`
	Cliente		ClientesAgendamentosResponse		`json:"cliente"`
	Prestador	PrestadorAgendamentosResponse		`json:"prestador"`
	DataHora 	time.Time   `json:"datahora"`
	Status 		string   	`json:"status"`
	Localizacao string   	`json:"localizacao"`
	Latitude    float64  	`json:"latitude"`
	Longitude   float64  	`json:"longitude"`
	Anexos      []string 	`json:"anexos"`
}

type ClientesAgendamentosResponse struct {
	ID			uint		`json:"id"`
	ClienteNome     string `json:"nome"`
}
type PrestadorAgendamentosResponse struct {
	ID				uint	`json:"id"`
	PrestadorNome   string `json:"nome"`
}

func(uc *AgendamentoUC) Criar(ctx context.Context, req *AgendamentoRequest, idCliente uint) error {
	catalogo, err := uc.catalogoRepo.FindByID(ctx, req.IDCatalogo)
	if err != nil {
		return err
	}

	agendamento := &model.Agendamento{
		Detalhe:     req.Detalhe,
		IDCatalogo:  req.IDCatalogo,
		IDCliente:   idCliente,
		DataHora:    req.DataHora,
		Status:      "PENDENTE",
		Localizacao: req.Localizacao,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
	}

	err = uc.r.Criar(ctx, agendamento)
	if err != nil {
		return err
	}

	for _, anexoURL := range req.Anexos {
		anexo := &model.AnexoImagem{
			URL:           anexoURL,
			AgendamentoID: &agendamento.ID,
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
	err = uc.usuarioRepo.IncrementarNotificacoesNovas(ctx, catalogo.Prestador.IDUsuario)
	if err != nil {
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
		ID: agendamento.ID,
		Detalhe: agendamento.Detalhe,
		Catalogo: agendamento.Catalogo.Nome,
		Cliente: ClientesAgendamentosResponse{
			ID: agendamento.IDCliente,
			ClienteNome: agendamento.Cliente.Nome,
		},
		Prestador: PrestadorAgendamentosResponse{
			ID: agendamento.Catalogo.IDPrestador,
			PrestadorNome: agendamento.Catalogo.Prestador.Nome,
		},
		DataHora: agendamento.DataHora,
		Status: agendamento.Status,
		Localizacao: agendamento.Localizacao,
		Latitude: agendamento.Latitude,
		Longitude: agendamento.Longitude,
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
		Localizacao: agendamento.Localizacao,
		Latitude: agendamento.Latitude,
		Longitude: agendamento.Longitude,
		Status: model.StatusEmAndamento,
		IDCliente: agendamento.IDCliente,
		IDPrestador: agendamento.Catalogo.IDPrestador,
		DataHoraInicio: time.Now(),
	}

	switch agendamento.Catalogo.TipoPreco {
		case "fixo":
			servico.Preco = agendamento.Catalogo.ValorFixo
		case "por_hora":
			servico.Preco = 0.0 // Initial price for hourly services, will be calculated at finalization
		default:
			return errors.New("tipo de preço de catálogo inválido")
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
			ID: agendamento.ID,
			Detalhe: agendamento.Detalhe,
			Catalogo: agendamento.Catalogo.Nome,
			Cliente: ClientesAgendamentosResponse{
				ID: agendamento.IDCliente,
				ClienteNome: agendamento.Cliente.Nome,
			},
			Prestador: PrestadorAgendamentosResponse{
				ID: agendamento.Catalogo.IDPrestador,
				PrestadorNome: agendamento.Catalogo.Prestador.Nome,
			},
			DataHora: agendamento.DataHora,
			Status:   agendamento.Status,
			Localizacao: agendamento.Localizacao,
			Latitude: agendamento.Latitude,
			Longitude: agendamento.Longitude,
			Anexos: urls,
		})
	}
	return resp, nil
}

func (uc *AgendamentoUC) ListarPorClienteID(ctx context.Context, idUsuario uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]AgendamentoResponse, error) {
	agendamentos, err := uc.r.ListarPorClienteID(ctx, idUsuario, filters, orderBy, orderDir, limit, offset)
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
			ID: agendamento.ID,
			Detalhe: agendamento.Detalhe,
			Catalogo: agendamento.Catalogo.Nome,
			Cliente: ClientesAgendamentosResponse{
				ID: agendamento.IDCliente,
				ClienteNome: agendamento.Cliente.Nome,
			},
			Prestador: PrestadorAgendamentosResponse{
				ID: agendamento.Catalogo.IDPrestador,
				PrestadorNome: agendamento.Catalogo.Prestador.Nome,
			},
			DataHora: agendamento.DataHora,
			Status:   agendamento.Status,
			Localizacao: agendamento.Localizacao,
			Latitude: agendamento.Latitude,
			Longitude: agendamento.Longitude,
			Anexos: urls,
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
			ID: agendamento.ID,
			Detalhe: agendamento.Detalhe,
			Catalogo: agendamento.Catalogo.Nome,
			Cliente: ClientesAgendamentosResponse{
				ID: agendamento.IDCliente,
				ClienteNome: agendamento.Cliente.Nome,
			},
			Prestador: PrestadorAgendamentosResponse{
				ID: agendamento.Catalogo.IDPrestador,
				PrestadorNome: agendamento.Catalogo.Prestador.Nome,
			},
			DataHora: agendamento.DataHora,
			Status:   agendamento.Status,
			Localizacao: agendamento.Localizacao,
			Latitude: agendamento.Latitude,
			Longitude: agendamento.Longitude,
			Anexos: urls,
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
			ID: agendamento.ID,
			Detalhe: agendamento.Detalhe,
			Catalogo: agendamento.Catalogo.Nome,
			Cliente: ClientesAgendamentosResponse{
				ID: agendamento.IDCliente,
				ClienteNome: agendamento.Cliente.Nome,
			},
			Prestador: PrestadorAgendamentosResponse{
				ID: agendamento.Catalogo.IDPrestador,
				PrestadorNome: agendamento.Catalogo.Prestador.Nome,
			},
			DataHora: agendamento.DataHora,
			Status:   agendamento.Status,
			Localizacao: agendamento.Localizacao,
			Latitude: agendamento.Latitude,
			Longitude: agendamento.Longitude,
			Anexos: urls,
		})
	}
	return resp, nil
}
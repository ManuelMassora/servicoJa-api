package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
)

type ServicoUseCase struct {
	r model.ServicoRepo
	agendamentoRepo model.AgendamentoRepo
	vagaRepo model.VagaRepo
	notificacaoRepo model.NotificacaoRepo
    usuarioRepo model.UsuarioRepo
}

type ServicoResponse struct {
	ID		  		uint      	`json:"id"`
	Localizacao 	string   	`json:"localizacao"`
	Latitude    	float64  	`json:"latitude"`
	Longitude   	float64  	`json:"longitude"`
	Preco       	float64  	`json:"preco"`
	Status      	string   	`json:"status"`
	IDAgendamento   *uint    	`json:"id_agendamento,omitempty"`
	IDVaga 			*uint 		`json:"id_vaga,omitempty"`
	DataHoraInicio  time.Time 	`json:"data_inicio,omitempty"`
	DataHoraFim     time.Time  	`json:"data_fim,omitempty"`
	Cliente    		uint      	`json:"cliente"`
	Prestador  		uint      	`json:"prestador"`
	Catalogo 		string		`json:"catalogo,omitempty"`
	Descricao		string		`json:"descricao"`
}

func NewServicoUseCase(r model.ServicoRepo, agendamentoRepo model.AgendamentoRepo, vagaRepo model.VagaRepo, notificacaoRepo model.NotificacaoRepo, usuarioRepo model.UsuarioRepo) *ServicoUseCase {
	return &ServicoUseCase{
		r: r,
		agendamentoRepo: agendamentoRepo,
		vagaRepo: vagaRepo,
		notificacaoRepo: notificacaoRepo,
        usuarioRepo: usuarioRepo,
	}
}

func (uc *ServicoUseCase) FinalizarServico(ctx context.Context, idServico, idUsuario uint) error {
	servico, err := uc.r.BuscarPorID(ctx, idServico)
	if err != nil {
		return err
	}
	if servico.Cliente.IDUsuario != idUsuario && servico.Prestador.IDUsuario != idUsuario {
		return errors.New("usuário não autorizado a finalizar este serviço")
	}
	if servico.Status == model.StatusConcluido || servico.Status == model.StatusCancelado {
		return nil
	}
	err = uc.notificacaoRepo.Enviar(ctx, &model.Notificacao{
		IDUsuario: servico.IDCliente,
		Titulo: "Serviço Concluído",
		Mensagem: "O serviço foi concluído com sucesso. Obrigado por usar nossos serviços!",
	})
	if err != nil {
		return err
	}
	err = uc.usuarioRepo.IncrementarNotificacoesNovas(ctx, servico.IDCliente)
	if err != nil {
		return err
	}
	servico.Status = model.StatusConcluido
	servico.DataHoraFim = time.Now()

	// Calculate final price for hourly services
	if servico.Agendamento.Catalogo.TipoPreco == "por_hora" {
		if servico.DataHoraInicio.IsZero() {
			return errors.New("data de início do serviço não definida para cálculo por hora")
		}
		duration := servico.DataHoraFim.Sub(servico.DataHoraInicio)
		servico.Preco = CalculateFinalServicePrice(&servico.Agendamento.Catalogo, duration)
	}
	return uc.r.Atualizar(ctx, servico)
}

func (uc *ServicoUseCase) CancelarServico(ctx context.Context, idServico, idUsuario uint) error {	
	servico, err := uc.r.BuscarPorID(ctx, idServico)
	if err != nil {
		return err
	}
	err = uc.notificacaoRepo.Enviar(ctx, &model.Notificacao{
		IDUsuario: servico.IDCliente,
		Titulo: "Serviço Cancelado",
		Mensagem: "O serviço foi cancelado com sucesso. Esperamos que tenha gostado de nossos serviços!",
	})
	if err != nil {
		return err
	}
	err = uc.usuarioRepo.IncrementarNotificacoesNovas(ctx, servico.IDCliente)
	if err != nil {
		return err
	}
	if servico.Cliente.IDUsuario != idUsuario && servico.Prestador.IDUsuario != idUsuario {
		return errors.New("usuário não autorizado a cancelar este serviço")
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
	servicos, err := uc.r.ListarPorCliente(ctx, idUsuario, filters, orderBy, orderDir, limit, offset)
	if err != nil {
		return nil, err
	}
	if len(servicos) == 0 {
		return []ServicoResponse{}, nil
	}
	var resp []ServicoResponse
	for _, s := range servicos {
		var descricao string
		if s.IDAgendamento != nil && s.Agendamento != nil {
			descricao = s.Agendamento.Detalhe
		} else if s.IDVaga != nil && s.Vaga != nil {
			descricao = s.Vaga.Descricao
		}

		resp = append(resp, ServicoResponse{
			ID:             s.ID,
			Localizacao:    s.Localizacao,
			Latitude:       s.Latitude,
			Longitude:      s.Longitude,
			Preco:          s.Preco,
			Status:         string(s.Status),
			IDAgendamento:  s.IDAgendamento,
			IDVaga:         s.IDVaga,
			DataHoraInicio: s.DataHoraInicio,
			DataHoraFim:    s.DataHoraFim,
			Cliente:         s.IDCliente,
			Prestador:       s.IDPrestador,
			Descricao:       descricao,
		})
	}
	return resp, nil
}

func (uc *ServicoUseCase) ListarPorPrestador(ctx context.Context, idUsuario uint, filters map[string]interface{}, orderBy, orderDir string, limit, offset int) ([]ServicoResponse, error) {
	servicos, err := uc.r.ListarPorPrestador(ctx, idUsuario, filters, orderBy, orderDir, limit, offset)
	if err != nil {
		return nil, err
	}
	if len(servicos) == 0 {
		return []ServicoResponse{}, nil
	}
	var resp []ServicoResponse
	for _, s := range servicos {
		var descricao string
		var catalogo string
		if s.IDAgendamento != nil && s.Agendamento != nil {
			descricao = s.Agendamento.Detalhe
			catalogo = s.Agendamento.Catalogo.Nome
		} else if s.IDVaga != nil && s.Vaga != nil {
			descricao = s.Vaga.Descricao
		}

		resp = append(resp, ServicoResponse{
			ID:             s.ID,
			Localizacao:    s.Localizacao,
			Latitude:       s.Latitude,
			Longitude:      s.Longitude,
			Preco:          s.Preco,
			Status:         string(s.Status),
			IDAgendamento:  s.IDAgendamento,
			IDVaga:         s.IDVaga,
			DataHoraInicio: s.DataHoraInicio,
			DataHoraFim:    s.DataHoraFim,
			Cliente:        s.IDCliente,
			Prestador:      s.IDPrestador,
			Descricao:      descricao,
			Catalogo: 		catalogo,
		})
	}
	return resp, nil
}

func (uc *ServicoUseCase) ListarPorLocalizacao(ctx context.Context, userID uint, latitude, longitude, radius float64, filters map[string]interface{}, orderBy, orderDir string, limit, offset int) ([]ServicoResponse, error) {
	servicos, err := uc.r.FindByLocation(ctx, userID, latitude, longitude, radius, filters, orderBy, orderDir, limit, offset)
	if err != nil {
		return nil, err
	}
	if len(servicos) == 0 {
		return []ServicoResponse{}, nil
	}
	var resp []ServicoResponse
	for _, s := range servicos {
		var descricao string
		if s.IDAgendamento != nil && s.Agendamento != nil {
			descricao = s.Agendamento.Detalhe
		} else if s.IDVaga != nil && s.Vaga != nil {
			descricao = s.Vaga.Descricao
		}

		resp = append(resp, ServicoResponse{
			ID:             s.ID,
			Localizacao:    s.Localizacao,
			Latitude:       s.Latitude,
			Longitude:      s.Longitude,
			Preco:          s.Preco,
			Status:         string(s.Status),
			IDAgendamento:  s.IDAgendamento,
			IDVaga:         s.IDVaga,
			DataHoraInicio: s.DataHoraInicio,
			DataHoraFim:    s.DataHoraFim,
			Cliente:         s.IDCliente,
			Prestador:       s.IDPrestador,
			Descricao:       descricao,
		})
	}
	return resp, nil
}

package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
)

type PropostaUseCase struct {
	propostaRepo    model.PropostaRepo
	vagaRepo        model.VagaRepo
	servicoRepo     model.ServicoRepo
	notificacaoRepo model.NotificacaoRepo
	usuarioRepo     model.UsuarioRepo
	pagamentoRepo   model.PagamentoRepo
}

func NewPropostaUseCase(
	propostaRepo model.PropostaRepo,
	vagaRepo model.VagaRepo,
	servicoRepo model.ServicoRepo,
	notificacaoRepo model.NotificacaoRepo,
	usuarioRepo model.UsuarioRepo,
	pagamentoRepo model.PagamentoRepo,
) *PropostaUseCase {
	return &PropostaUseCase{
		propostaRepo:    propostaRepo,
		vagaRepo:        vagaRepo,
		servicoRepo:     servicoRepo,
		notificacaoRepo: notificacaoRepo,
		usuarioRepo:     usuarioRepo,
		pagamentoRepo:   pagamentoRepo,
	}
}

type PropostaRequest struct {
	IDVaga        uint    `json:"id_vaga" binding:"required"`
	ValorProposto float64 `json:"valor_proposto" binding:"required"`
	Mensagem      string  `json:"mensagem" binding:"required"`
	PrazoEstimado string  `json:"prazo_estimado" binding:"required"`
}

type PropostaResponse struct {
	ID            uint      `json:"id"`
	IDVaga        uint      `json:"id_vaga"`
	Vaga          string    `json:"vaga"`
	IDPrestador   uint      `json:"id_prestador"`
	Prestador     string    `json:"prestador"`
	ValorProposto float64   `json:"valor_proposto"`
	Mensagem      string    `json:"mensagem"`
	PrazoEstimado string    `json:"prazo_estimado"`
	Status        string    `json:"status"`
	DataResposta  time.Time `json:"data_resposta"`
}

func (uc *PropostaUseCase) Criar(ctx context.Context, request PropostaRequest, idUsuario uint) error {
	vaga, err := uc.vagaRepo.BuscarPorID(ctx, request.IDVaga)
	if err != nil {
		return err
	}
	err = uc.notificacaoRepo.Enviar(ctx, &model.Notificacao{
		IDUsuario: vaga.IDCliente,
		Titulo:    "Nova Proposta",
		Mensagem:  "Você tem uma nova proposta na vaga: " + vaga.Titulo,
	})
	if err != nil {
		return err
	}
	err = uc.usuarioRepo.IncrementarNotificacoesNovas(ctx, vaga.IDCliente)
	if err != nil {
		return err
	}
	if err := uc.vagaRepo.IncrementarPropostasNovas(ctx, vaga.ID); err != nil {
		return err
	}
	proposta := &model.Proposta{
		IDVaga:        request.IDVaga,
		IDPrestador:   idUsuario,
		ValorProposto: request.ValorProposto,
		Mensagem:      request.Mensagem,
		PrazoEstimado: request.PrazoEstimado,
		Status:        model.StatusPendente,
	}
	return uc.propostaRepo.Salvar(ctx, proposta)
}

func (uc *PropostaUseCase) Responder(ctx context.Context, idProposta, idUsuario uint, aceitar bool) error {
	proposta, err := uc.propostaRepo.BuscarPorID(ctx, idProposta)
	if err != nil {
		return err
	}
	vaga, err := uc.vagaRepo.BuscarPorID(ctx, proposta.IDVaga)
	if err != nil {
		return err
	}
	if vaga.IDCliente != idUsuario {
		return errors.New("acesso negado: apenas o cliente dono da vaga pode responder a proposta")
	}
	if proposta.Status != model.StatusPendente {
		return errors.New("acesso negado: apenas propostas pendentes podem ser respondidas")
	}
	proposta.DataResposta = time.Now()
	if aceitar {
		err = uc.notificacaoRepo.Enviar(ctx, &model.Notificacao{
			IDUsuario: proposta.IDPrestador,
			Titulo:    "Proposta Aceita",
			Mensagem:  "Parabéns! A sua proposta foi aceita na vaga: " + vaga.Titulo,
		})
		if err != nil {
			return err
		}
		err = uc.usuarioRepo.IncrementarNotificacoesNovas(ctx, proposta.IDPrestador)
		if err != nil {
			return err
		}
		vaga.IDPrestador = &proposta.IDPrestador
		vaga.Status = model.StatusProposta
		if err := uc.vagaRepo.Salvar(ctx, vaga); err != nil {
			return err
		}
		proposta.Status = model.StatusAceito
		servico := &model.Servico{
			IDVaga:         &proposta.IDVaga,
			Localizacao:    vaga.Localizacao,
			Preco:          proposta.ValorProposto,
			Status:         model.StatusEmAndamento,
			IDCliente:      vaga.IDCliente,
			IDPrestador:    proposta.IDPrestador,
			DataHoraInicio: time.Now(),
		}
		if servicoSave, err := uc.servicoRepo.Criar(ctx, servico); err != nil {
			return err
		} else {
			// Associar ID do serviço ao pagamento
			p, err := uc.pagamentoRepo.BuscarPorVaga(ctx, proposta.IDVaga)
			if err == nil && p != nil {
				_ = uc.pagamentoRepo.AtualizarIDServico(ctx, p.ID, servicoSave.ID)
			}
		}
	} else {
		err = uc.notificacaoRepo.Enviar(ctx, &model.Notificacao{
			IDUsuario: proposta.IDPrestador,
			Titulo:    "Proposta Rejeitada",
			Mensagem:  "Infelismente a sua proposta foi rejeitada na vaga: " + vaga.Titulo,
		})
		if err != nil {
			return err
		}
		err = uc.usuarioRepo.IncrementarNotificacoesNovas(ctx, proposta.IDPrestador)
		if err != nil {
			return err
		}
		proposta.Status = model.StatusRejeitado
	}
	return uc.propostaRepo.Salvar(ctx, proposta)
}

func (uc *PropostaUseCase) Cancelar(ctx context.Context, idProposta, idUsuario uint) error {
	proposta, err := uc.propostaRepo.BuscarPorID(ctx, idProposta)
	if err != nil {
		return err
	}
	if proposta.IDPrestador != idUsuario {
		return errors.New("acesso negado: apenas o prestador que fez a proposta pode cancelá-la")
	}
	proposta.Status = model.StatusCancelado
	proposta.DeletedAt.Time = time.Now()
	proposta.DeletedAt.Valid = true
	return uc.propostaRepo.Salvar(ctx, proposta)
}

func mapPropostasToResponse(propostas []model.Proposta) []PropostaResponse {
	if len(propostas) == 0 {
		return []PropostaResponse{}
	}
	respostas := make([]PropostaResponse, 0, len(propostas))
	for _, proposta := range propostas {
		// Garantir que Vaga e Prestador estão carregados (preloaded) pelo Repositório.
		// Se proposta.Vaga ou proposta.Prestador for nil, pode causar pânico (panic).
		// Aqui assumimos que eles estão carregados (o que é usual em listagens).

		prestadorNome := ""
		if proposta.Prestador != nil {
			prestadorNome = proposta.Prestador.Usuario.Nome
		}

		vagaTitulo := ""
		if proposta.Vaga.Titulo != "" {
			vagaTitulo = proposta.Vaga.Titulo
		}

		respostas = append(respostas, PropostaResponse{
			ID:            proposta.ID,
			IDVaga:        proposta.IDVaga,
			Vaga:          vagaTitulo,
			IDPrestador:   proposta.IDPrestador,
			Prestador:     prestadorNome,
			ValorProposto: proposta.ValorProposto,
			Mensagem:      proposta.Mensagem,
			PrazoEstimado: proposta.PrazoEstimado,
			Status:        string(proposta.Status),
			DataResposta:  proposta.DataResposta,
		})
	}
	return respostas
}

func (uc *PropostaUseCase) ListarPorVaga(ctx context.Context, idUsuario, idVaga uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]PropostaResponse, error) {
	vaga, err := uc.vagaRepo.BuscarPorID(ctx, idVaga)
	if err != nil {
		return nil, err
	}

	// Verificação de Autorização: O usuário deve ser o Cliente dono da Vaga
	if idUsuario != vaga.IDCliente {
		return nil, errors.New("acesso negado: apenas o cliente que criou a vaga pode ver as propostas")
	}

	// Busca no Repositório
	propostas, err := uc.propostaRepo.ListarPorVaga(ctx, idVaga, filters, orderBy, orderDir, limit, offset)
	if err != nil {
		return nil, err
	}

	if err := uc.vagaRepo.ZerarPropostasNovas(ctx, vaga.ID); err != nil {
		return nil, err
	}

	// Uso da função auxiliar
	return mapPropostasToResponse(propostas), nil
}

func (uc *PropostaUseCase) ListarPorPrestador(ctx context.Context, idUsuario uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]PropostaResponse, error) {
	// Busca no Repositório
	propostas, err := uc.propostaRepo.ListarPorPrestador(ctx, idUsuario, filters, orderBy, orderDir, limit, offset)
	if err != nil {
		return nil, err
	}

	// Uso da função auxiliar
	return mapPropostasToResponse(propostas), nil
}

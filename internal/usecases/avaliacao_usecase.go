package usecases

import (
	"context"
	"errors"
	"time"
	
	"github.com/ManuelMassora/servicoJa-api/internal/model"
)

type AvaliacaoUseCase struct {
	avaliacaoRepo model.AvaliacaoRepo
	servicoRepo model.ServicoRepo
	notificacaoRepo model.NotificacaoRepo
	usuarioRepo model.UsuarioRepo
}
func NewAvaliacaoUseCase(
	avaliacaoRepo model.AvaliacaoRepo,
	servicoRepo model.ServicoRepo,
	notificacaoRepo model.NotificacaoRepo,
	usuarioRepo model.UsuarioRepo,
) *AvaliacaoUseCase {
	return &AvaliacaoUseCase{
		avaliacaoRepo:   avaliacaoRepo,
		servicoRepo:     servicoRepo,
		notificacaoRepo: notificacaoRepo,
		usuarioRepo:     usuarioRepo,
	}
}

type AvaliacaoRequest struct {
	IDServico uint    `json:"id_servico" binding:"required"`
	Pontuacao int     `json:"pontuacao" binding:"required,min=1,max=5"`
	Comentario string `json:"comentario"`
}

type AvaliacaoResponse struct {
	ID          uint      `json:"id"`
	IDServico   uint      `json:"id_servico"`
	Servico     string    `json:"servico"`
	Avaliador   string    `json:"avaliador"`
	Avaliado    string    `json:"avaliado"`
	Pontuacao   int       `json:"pontuacao"`
	Comentario  string    `json:"comentario"`
	DataCriacao time.Time `json:"data_criacao"`
}

func (uc *AvaliacaoUseCase) Criar(ctx context.Context, req AvaliacaoRequest, idAvaliador uint) error {
	servico, err := uc.servicoRepo.BuscarPorID(ctx, req.IDServico)
	if err != nil {
		return err
	}

	if servico.Status != model.StatusConcluido {
		return errors.New("não é possível avaliar um serviço que não foi concluído")
	}

	var idCliente, idPrestador uint
	if servico.IDCliente == idAvaliador {
		idCliente = idAvaliador
		idPrestador = servico.IDPrestador
	} else if servico.IDPrestador == idAvaliador {
		idCliente = servico.IDCliente
		idPrestador = idAvaliador
	} else {
		return errors.New("usuário não autorizado a avaliar este serviço")
	}

	avaliacao := &model.Avaliacao{
		Nota: req.Pontuacao,
		Comentario: req.Comentario,
		IDCliente: idCliente,
		IDPrestador: idPrestador,
		IDServico: req.IDServico,
	}

	if err := uc.avaliacaoRepo.Criar(ctx, avaliacao); err != nil {
		return err
	}

	// Enviar notificação para o usuário avaliado
	var idUsuarioAvaliado uint
	if idAvaliador == idCliente { // Cliente avaliou o prestador
		idUsuarioAvaliado = idPrestador
	} else { // Prestador avaliou o cliente
		idUsuarioAvaliado = idCliente
	}

	err = uc.notificacaoRepo.Enviar(ctx, &model.Notificacao{
		IDUsuario: idUsuarioAvaliado,
		Titulo: "Nova Avaliação Recebida",
		Mensagem: "Você recebeu uma nova avaliação para o serviço " + servico.Localizacao + " com " + servico.Cliente.Usuario.Nome + " e " + servico.Prestador.Usuario.Nome + ".",
	})
	if err != nil {
		return err
	}
	err = uc.usuarioRepo.IncrementarNotificacoesNovas(ctx, idUsuarioAvaliado)
	if err != nil {
		return err
	}
	return nil
}

func mapAvaliacoesToResponse(avaliacoes []model.Avaliacao) []AvaliacaoResponse {
	if len(avaliacoes) == 0 {
		return []AvaliacaoResponse{}
	}
	respostas := make([]AvaliacaoResponse, 0, len(avaliacoes))
	for _, a := range avaliacoes {
		avaliadorNome := ""
		if a.Cliente != nil {
			avaliadorNome = a.Cliente.Usuario.Nome
		}

		avaliadoNome := ""
		if a.Prestador != nil {
			avaliadoNome = a.Prestador.Usuario.Nome
		}

		respostas = append(respostas, AvaliacaoResponse{
			ID:          a.ID,
			IDServico:   a.IDServico,
			Servico:     a.Servico.Localizacao, // Assumindo que queremos a localização como identificador do serviço
			Avaliador:   avaliadorNome,
			Avaliado:    avaliadoNome,
			Pontuacao:   a.Nota,
			Comentario:  a.Comentario,
			DataCriacao: a.CreatedAt,
		})
	}
	return respostas
}

func (uc *AvaliacaoUseCase) ListarPorCliente(ctx context.Context, idUsuario uint, filters map[string]interface{}, orderBy, orderDir string, limit, offset int) ([]AvaliacaoResponse, error) {

	avaliacoes, err := uc.avaliacaoRepo.ListarPorCliente(ctx, idUsuario, filters, orderBy, orderDir, limit, offset)
	if err != nil {
		return nil, err
	}

	return mapAvaliacoesToResponse(avaliacoes), nil
}

func (uc *AvaliacaoUseCase) ListarPorPrestador(ctx context.Context, idPrestador uint, filters map[string]interface{}, orderBy, orderDir string, limit, offset int) ([]AvaliacaoResponse, error) {
	// Sem validação de usuário, qualquer um pode ver as avaliações de um prestador
	avaliacoes, err := uc.avaliacaoRepo.ListarPorPrestador(ctx, idPrestador, filters, orderBy, orderDir, limit, offset)
	if err != nil {
		return nil, err
	}

	return mapAvaliacoesToResponse(avaliacoes), nil
}

func (uc *AvaliacaoUseCase) MediaPorPrestador(ctx context.Context, idPrestador uint) (float64, error) {
	return uc.avaliacaoRepo.MediaPorPrestador(ctx, idPrestador)
}

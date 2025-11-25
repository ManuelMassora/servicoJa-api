package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
)

type NotificacaoUseCase struct {
	NotificacaoRepo model.NotificacaoRepo
	UsuarioRepo model.UsuarioRepo
}

type NotificacaoResponse struct {
	ID			uint		`json:"id"`
	IDUsuario 	uint    	`json:"usuario_id"`
	Titulo    	string   	`json:"titulo"`
	Mensagem  	string   	`json:"mensagem"`
	Lida      	bool     	`json:"lida"`
	DataCriacao time.Time 	`json:"data_criacao"`
}

func NewNotificacaoUseCase(notificacaoRepo model.NotificacaoRepo, usuarioRepo model.UsuarioRepo) *NotificacaoUseCase {
	return &NotificacaoUseCase{
		NotificacaoRepo: notificacaoRepo,
		UsuarioRepo: usuarioRepo,
	}
}

func(uc *NotificacaoUseCase) ListarPorUsuario(
    ctx context.Context, 
    idUsuario uint, 
    filters map[string]interface{}, 
    orderBy string, 
    orderDir string, 
    limit, offset int,
) ([]NotificacaoResponse, error) {
	notificacoes, err := uc.NotificacaoRepo.ListarPorUsuario(ctx, idUsuario, filters, orderBy, orderDir, limit, offset)
	if err != nil {
		return nil, err
	}
	var respostas []NotificacaoResponse
	for _, notificacao := range notificacoes {
		resposta := NotificacaoResponse{
			ID:           notificacao.ID,
			IDUsuario:    notificacao.IDUsuario,
			Titulo:       notificacao.Titulo,
			Mensagem:     notificacao.Mensagem,
			Lida:         notificacao.Lida,
			DataCriacao:  notificacao.CreatedAt,
		}
		respostas = append(respostas, resposta)
	}
	err = uc.UsuarioRepo.ZerarNotificacoesNovas(ctx, idUsuario)
	if err != nil {
		return nil, err
	}
	return respostas, nil
}

func(uc *NotificacaoUseCase) MarcarComoLida(ctx context.Context, id, idUsuario uint) error {
	notificacao, err := uc.NotificacaoRepo.BuscarPorID(ctx, id)
	if err != nil {
		return err
	}
	if notificacao.IDUsuario != idUsuario {
		return errors.New("acesso negado: nao pode marcar essa notificacao como lida")
	}
	return uc.NotificacaoRepo.MarcarComoLida(ctx, id)
}

func(uc *NotificacaoUseCase) MarcarTodasComoLidas(ctx context.Context, idUsuario uint) error {
	err := uc.NotificacaoRepo.MarcarTodasComoLidas(ctx, idUsuario)
	if err != nil {
		return err
	}
	// After marking all as read, set the new notifications count to zero
	return uc.UsuarioRepo.ZerarNotificacoesNovas(ctx, idUsuario)
}
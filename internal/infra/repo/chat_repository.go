package repo

import (
	"context"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"gorm.io/gorm"
)

type ChatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) model.ChatRepo {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) CriarChat(ctx context.Context, chat *model.Chat) error {
	return r.db.WithContext(ctx).Create(chat).Error
}

func (r *ChatRepository) ListarChatsPorUsuario(ctx context.Context, idUsuario uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Chat, error) {
	var chats []model.Chat
	query := r.db.WithContext(ctx).
		Preload("Servico").
		Preload("Prestador").
		Preload("Cliente").
		Where("prestador_id = ? OR id_cliente = ?", idUsuario, idUsuario)

	// Apply filters
	for field, value := range filters {
		query = query.Where(field+" = ?", value)
	}

	// Apply ordering
	if orderBy != "" {
		if orderDir == "" {
			orderDir = "asc"
		}
		query = query.Order(orderBy + " " + orderDir)
	}

	// Apply pagination
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&chats).Error
	if err != nil {
		return nil, err
	}
	return chats, nil
}

type MensagemRepository struct {
	db *gorm.DB
}

func NewMensagemRepository(db *gorm.DB) *MensagemRepository {
	return &MensagemRepository{db: db}
}

func (r *MensagemRepository) EnviarMensagem(ctx context.Context, msg *model.Mensagem) error {
	return r.db.WithContext(ctx).Create(msg).Error
}

func (r *MensagemRepository) ListarMensagens(ctx context.Context, idChat uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Mensagem, error) {
	var mensagens []model.Mensagem
	query := r.db.WithContext(ctx).
		Preload("Chat").
		Preload("Remetente").
		Where("id_chat = ?", idChat)

	// Apply filters
	for field, value := range filters {
		query = query.Where(field+" = ?", value)
	}

	// Apply ordering
	if orderBy != "" {
		if orderDir == "" {
			orderDir = "asc"
		}
		query = query.Order(orderBy + " " + orderDir)
	} else {
		// Default ordering by created_at if no order specified
		query = query.Order("created_at ASC")
	}

	// Apply pagination
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&mensagens).Error
	if err != nil {
		return nil, err
	}
	return mensagens, nil
}
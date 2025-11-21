package repo

import (
	"context"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"gorm.io/gorm"
)

type NotificacaoRepo struct {
	db *gorm.DB
}

func NewNotificacaoRepo(db *gorm.DB) model.NotificacaoRepo {
	return &NotificacaoRepo{db: db}
}

func(r *NotificacaoRepo)Enviar(ctx context.Context, notificacao *model.Notificacao) error {
	return r.db.WithContext(ctx).Create(notificacao).Error
}

func(r *NotificacaoRepo) BuscarPorID(ctx context.Context, id uint) (*model.Notificacao,error){
	var notificacao model.Notificacao
	err := r.db.WithContext(ctx).Preload("Usuario").First(&notificacao, id).Error
	if err != nil {
		return nil, err
	}
	return &notificacao, nil
}

func (r *NotificacaoRepo) ListarPorUsuario(
    ctx context.Context, 
    idUsuario uint, 
    filters map[string]interface{}, 
    orderBy string, 
    orderDir string, 
    limit, offset int,
) ([]model.Notificacao, error) { 
    var notificacoes []model.Notificacao
    query := r.db.Preload("Usuario").Model(&model.Notificacao{}).
        Where("id_usuario = ?", idUsuario)
   
    for key, value := range filters {
        query = query.Where(key+" LIKE ?", "%"+value.(string)+"%")
    } 
    if orderBy != "" {        
        if orderDir == "" {
            orderDir = "asc"
        }
        query = query.Order(orderBy + " " + orderDir)
    } else {
        query = query.Order("created_at DESC")
    }
    if limit > 0 {
        query = query.Limit(limit)
    }
    if offset > 0 {
        query = query.Offset(offset)
    }
    if err := query.Find(&notificacoes).Error; err != nil {
        return nil, err
    }
    return notificacoes, nil
}

func(r *NotificacaoRepo)MarcarComoLida(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&model.Notificacao{}).Where("id = ?", id).Update("lida", true).Error
}
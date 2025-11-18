package repo

import (
	"context"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"gorm.io/gorm"
)

type AgendamentoRepo struct {
	db *gorm.DB
}

func NewAgendamentoRepo(db *gorm.DB) model.AgendamentoRepo {
	return &AgendamentoRepo{db: db}
}

func(r *AgendamentoRepo) Criar(ctx context.Context, agendamento *model.Agendamento) error{
	return r.db.WithContext(ctx).Create(agendamento).Error
}
func(r *AgendamentoRepo) BuscarPorID(ctx context.Context, id uint) (*model.Agendamento, error){
	var agendamento model.Agendamento
	err := r.db.WithContext(ctx).
		Where("id=?", id).
		First(&agendamento).Error
	if err != nil {
		return nil, err
	}
	return &agendamento, nil
}
func(r *AgendamentoRepo) Listar(
		ctx context.Context, 
		filters map[string]interface{}, 
		orderBy string, 
		orderDir string, 
		limit, 
		offset int,
	) ([]model.Agendamento, error){
	var agendamentos []model.Agendamento
	query := r.db.WithContext(ctx).Model(&model.Agendamento{})
    query = query.Preload("Catalogo").Preload("Catalogo.Prestador").Preload("Catalogo.Prestador.Usuario").Preload("Cliente").Preload("Cliente.Usuario")
    for key, value := range filters {
        switch v := value.(type) {
        case string:           
            query = query.Where(key+" LIKE ?", "%"+v+"%")
        case uint, int:
            query = query.Where(key+" = ?", v)
        default:         
        }
    }
	
	if orderBy != "" {
		if orderDir == "" {
			orderDir = "asc"
		}
		query = query.Order(orderBy + " " + orderDir)
	}
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	if err := query.Find(&agendamentos).Error;err != nil {
		return nil, err
	}
	return agendamentos, nil		
}
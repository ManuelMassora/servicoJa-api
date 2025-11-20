package repo

import (
	"context"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"gorm.io/gorm"
)

type ServicoRepository struct {
	db *gorm.DB
}

func NewServicoRepository(db *gorm.DB) model.ServicoRepo {
	return &ServicoRepository{db: db}
}

func (r *ServicoRepository) Criar(ctx context.Context, servico *model.Servico) error {
	return r.db.WithContext(ctx).Create(servico).Error
}

func (r *ServicoRepository) Atualizar(ctx context.Context, servico *model.Servico) error {
	return r.db.WithContext(ctx).Save(servico).Error
}

func (r *ServicoRepository) BuscarPorID(ctx context.Context, id uint) (*model.Servico, error) {
	var servico model.Servico
	err := r.db.WithContext(ctx).
		Preload("Cliente").
		Preload("Prestador").
		Preload("Agendamento").
		Preload("Agendamento.Catalogo").
		Preload("Vaga").
		First(&servico, id).Error
	if err != nil {
		return nil, err
	}
	return &servico, nil
}

func (r *ServicoRepository) AtualizarStatus(ctx context.Context, id uint, status string) error {
	return r.db.WithContext(ctx).
		Model(&model.Servico{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *ServicoRepository) ListarPorCliente(ctx context.Context, idCliente uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Servico, error) {
	var servicos []model.Servico
	query := r.db.WithContext(ctx).
		Preload("Cliente").
		Preload("Prestador").
		Preload("Agendamento").
		Preload("Agendamento.Catalogo").
		Preload("Vaga").
		Where("id_cliente = ?", idCliente)

	// Apply filters
    for key, value := range filters {
        switch v := value.(type) {
        case string:           
            query = query.Where(key+" LIKE ?", "%"+v+"%")
        case uint, int:
            query = query.Where(key+" = ?", v)
        default:         
        }
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

	err := query.Find(&servicos).Error
	if err != nil {
		return nil, err
	}
	return servicos, nil
}

func (r *ServicoRepository) ListarPorPrestador(ctx context.Context, IDPrestador uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Servico, error) {
	var servicos []model.Servico
	query := r.db.WithContext(ctx).
		Preload("Cliente").
		Preload("Prestador").
		Preload("Agendamento").
		Preload("Agendamento.Catalogo").
		Preload("Vaga").
		Where("id_prestador =?", IDPrestador)

	// Apply filters
    for key, value := range filters {
        switch v := value.(type) {
        case string:           
            query = query.Where(key+" LIKE ?", "%"+v+"%")
        case uint, int:
            query = query.Where(key+" = ?", v)
        default:         
        }
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

	err := query.Find(&servicos).Error
	if err != nil {
		return nil, err
	}
	return servicos, nil
}
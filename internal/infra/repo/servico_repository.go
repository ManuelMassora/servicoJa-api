package repo

import (
	"context"
	"fmt"

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

func (r *ServicoRepository) FindByLocation(ctx context.Context, latitude, longitude, radius float64, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Servico, error) {
    var servicos []model.Servico

    haversine := fmt.Sprintf(
        "6371 * acos(cos(radians(%f)) * cos(radians(latitude)) * cos(radians(longitude) - radians(%f)) + sin(radians(%f)) * sin(radians(latitude)))",
        latitude, longitude, latitude,
    )

    query := r.db.Select(fmt.Sprintf("*, (%s) AS distance", haversine)).
        Where(fmt.Sprintf("(%s) < ?", haversine), radius)

    for key, value := range filters {
        query = query.Where(fmt.Sprintf("%s LIKE ?", key), fmt.Sprintf("%%%v%%", value))
    }

    if orderBy != "" {
        if orderDir == "" {
            orderDir = "asc"
        }
        query = query.Order(fmt.Sprintf("%s %s", orderBy, orderDir))
    } else {
        query = query.Order("distance asc")
    }

    if limit > 0 {
        query = query.Limit(limit)
    }
    if offset > 0 {
        query = query.Offset(offset)
    }

    if err := query.Find(&servicos).Error; err != nil {
        return nil, err
    }

    return servicos, nil
}

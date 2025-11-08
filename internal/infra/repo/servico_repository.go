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

func (r *ServicoRepository) BuscarPorID(ctx context.Context, id int64) (*model.Servico, error) {
	var servico model.Servico
	err := r.db.WithContext(ctx).
		Preload("Cliente").
		Preload("Prestador").
		Preload("Categoria").
		First(&servico, id).Error
	if err != nil {
		return nil, err
	}
	return &servico, nil
}

func (r *ServicoRepository) AtualizarStatus(ctx context.Context, id int64, status string) error {
	return r.db.WithContext(ctx).
		Model(&model.Servico{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *ServicoRepository) ListarPorCliente(ctx context.Context, idCliente int64, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Servico, error) {
	var servicos []model.Servico
	query := r.db.WithContext(ctx).
		Preload("Cliente").
		Preload("Prestador").
		Preload("Categoria").
		Where("id_cliente = ?", idCliente)

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

	err := query.Find(&servicos).Error
	if err != nil {
		return nil, err
	}
	return servicos, nil
}

func (r *ServicoRepository) ListarDisponiveis(ctx context.Context, localizacao string, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Servico, error) {
	var servicos []model.Servico
	query := r.db.WithContext(ctx).
		Preload("Cliente").
		Preload("Prestador").
		Preload("Categoria").
		Where("id_prestador IS NULL")

	if localizacao != "" {
		query = query.Where("localizacao LIKE ?", "%"+localizacao+"%")
	}

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

	err := query.Find(&servicos).Error
	if err != nil {
		return nil, err
	}
	return servicos, nil
}
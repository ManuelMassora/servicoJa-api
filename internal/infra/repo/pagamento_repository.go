package repo

import (
	"context"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"gorm.io/gorm"
)

type PagamentoRepository struct {
	db *gorm.DB
}

func NewPagamentoRepository(db *gorm.DB) model.PagamentoRepo {
	return &PagamentoRepository{db: db}
}

func (r *PagamentoRepository) Criar(ctx context.Context, pagamento *model.Pagamento) error {
	return r.db.WithContext(ctx).Create(pagamento).Error
}

func (r *PagamentoRepository) BuscarPorServico(ctx context.Context, idServico int64) (*model.Pagamento, error) {
	var pagamento model.Pagamento
	err := r.db.WithContext(ctx).
		Preload("Servico").
		Preload("Cliente").
		Preload("Prestador").
		Where("id_servico = ?", idServico).
		First(&pagamento).Error
	if err != nil {
		return nil, err
	}
	return &pagamento, nil
}

func (r *PagamentoRepository) AtualizarStatus(ctx context.Context, id int64, status model.Status) error {
	return r.db.WithContext(ctx).
		Model(&model.Pagamento{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *PagamentoRepository) ListarPorUsuario(ctx context.Context, idUsuario int64, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Pagamento, error) {
	var pagamentos []model.Pagamento
	query := r.db.WithContext(ctx).
		Preload("Servico").
		Preload("Cliente").
		Preload("Prestador").
		Where("id_cliente = ? OR id_prestador = ?", idUsuario, idUsuario)

	
	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}

	
	if orderBy != "" {
		query = query.Order(orderBy + " " + orderDir)
	}

	
	err := query.Limit(limit).Offset(offset).Find(&pagamentos).Error
	if err != nil {
		return nil, err
	}
	return pagamentos, nil
}
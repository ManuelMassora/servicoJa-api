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

func (r *PagamentoRepository) BuscarPorID(ctx context.Context, id uint) (*model.Pagamento, error) {
	var pagamento model.Pagamento
	err := r.db.WithContext(ctx).
		Preload("Servico").
		Preload("Cliente").
		Preload("Prestador").
		Where("id = ?", id).
		First(&pagamento).Error
	if err != nil {
		return nil, err
	}
	return &pagamento, nil
}

func (r *PagamentoRepository) BuscarPorReferencia(ctx context.Context, referencia string) (*model.Pagamento, error) {
	var pagamento model.Pagamento
	err := r.db.WithContext(ctx).
		Preload("Servico").
		Preload("Cliente").
		Preload("Prestador").
		Where("referencia = ?", referencia).
		First(&pagamento).Error
	if err != nil {
		return nil, err
	}
	return &pagamento, nil
}

func (r *PagamentoRepository) BuscarPorServico(ctx context.Context, idServico uint) (*model.Pagamento, error) {
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

func (r *PagamentoRepository) BuscarPorVaga(ctx context.Context, idVaga uint) (*model.Pagamento, error) {
	var pagamento model.Pagamento
	err := r.db.WithContext(ctx).
		Preload("Vaga").
		Preload("Cliente").
		Preload("Prestador").
		Where("id_vaga = ?", idVaga).
		First(&pagamento).Error
	if err != nil {
		return nil, err
	}
	return &pagamento, nil
}

func (r *PagamentoRepository) BuscarPorAgendamento(ctx context.Context, idAgendamento uint) (*model.Pagamento, error) {
	var pagamento model.Pagamento
	err := r.db.WithContext(ctx).
		Preload("Agendamento").
		Preload("Cliente").
		Preload("Prestador").
		Where("id_agendamento = ?", idAgendamento).
		First(&pagamento).Error
	if err != nil {
		return nil, err
	}
	return &pagamento, nil
}

func (r *PagamentoRepository) AtualizarStatus(ctx context.Context, id uint, status model.Status) error {
	return r.db.WithContext(ctx).
		Model(&model.Pagamento{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *PagamentoRepository) AtualizarStatusPorReferencia(ctx context.Context, referencia string, status model.Status) error {
	return r.db.WithContext(ctx).
		Model(&model.Pagamento{}).
		Where("referencia = ?", referencia).
		Update("status", status).Error
}

func (r *PagamentoRepository) AtualizarIDServico(ctx context.Context, idPagamento uint, idServico uint) error {
	return r.db.WithContext(ctx).
		Model(&model.Pagamento{}).
		Where("id = ?", idPagamento).
		Update("id_servico", idServico).Error
}

func (r *PagamentoRepository) ListarPorUsuario(ctx context.Context, idUsuario uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Pagamento, error) {
	var pagamentos []model.Pagamento
	query := r.db.WithContext(ctx).
		Preload("Servico").
		Preload("Cliente").
		Preload("Prestador")

	if idUsuario > 0 {
		query = query.Where("id_cliente = ? OR id_prestador = ?", idUsuario, idUsuario)
	}

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

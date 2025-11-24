package repo

import (
	"context"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"gorm.io/gorm"
)

type AvaliacaoRepo struct {
	db *gorm.DB
}

func NewAvaliacaoRepository(db *gorm.DB) model.AvaliacaoRepo {
	return &AvaliacaoRepo{db: db}
}

func (r *AvaliacaoRepo) Criar(ctx context.Context, avaliacao *model.Avaliacao) error {
	return r.db.WithContext(ctx).Create(avaliacao).Error
}

func (r *AvaliacaoRepo) ListarPorCliente(ctx context.Context, idCliente uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Avaliacao, error) {
	var avaliacoes []model.Avaliacao
	query := r.db.WithContext(ctx).
		Preload("Cliente").
		Preload("Prestador").
		Preload("Servico").
		Where("id_cliente = ?", idCliente)

	for field, value := range filters {
		query = query.Where(field+" = ?", value)
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
	if err := query.Find(&avaliacoes).Error; err != nil {
		return nil, err
	}
	return avaliacoes, nil
}

func (r *AvaliacaoRepo) ListarPorPrestador(ctx context.Context, idPrestador uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Avaliacao, error) {
	var avaliacoes []model.Avaliacao
	query := r.db.WithContext(ctx).
		Preload("Cliente").
		Preload("Prestador").
		Preload("Servico").
		Where("id_prestador = ?", idPrestador)

	for field, value := range filters {
		query = query.Where(field+" = ?", value)
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
	if err := query.Find(&avaliacoes).Error; err != nil {
		return nil, err
	}
	return avaliacoes, nil
}

func (r *AvaliacaoRepo) MediaPorPrestador(ctx context.Context, idPrestador uint) (float64, error) {
	var media float64
	err := r.db.WithContext(ctx).Model(&model.Avaliacao{}).
		Where("id_prestador = ?", idPrestador).
		Select("COALESCE(AVG(nota), 0)").
		Scan(&media).Error
	if err != nil {
		return 0, err
	}
	return media, nil
}
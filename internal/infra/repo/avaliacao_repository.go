package repo

import (
	"context"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"gorm.io/gorm"
)

type AvaliacaoRepository struct {
	db *gorm.DB
}

func NewAvaliacaoRepository(db *gorm.DB) model.AvaliacaoRepo {
	return &AvaliacaoRepository{db: db}
}

func (r *AvaliacaoRepository) Criar(ctx context.Context, avaliacao *model.Avaliacao) error {
	return r.db.WithContext(ctx).Create(avaliacao).Error
}

func (r *AvaliacaoRepository) ListarPorServico(ctx context.Context, idServico uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Avaliacao, error) {
	var avaliacoes []model.Avaliacao
	query := r.db.WithContext(ctx).
		Preload("Usuario").
		Preload("Servico").
		Where("servico_id = ?", idServico)

	// Apply additional filters
	for field, value := range filters {
		query = query.Where(field+" = ?", value)
	}

	// Apply ordering
	if orderBy != "" {
		if orderDir != "" {
			query = query.Order(orderBy + " " + orderDir)
		} else {
			query = query.Order(orderBy)
		}
	}

	// Apply pagination
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&avaliacoes).Error
	if err != nil {
		return nil, err
	}
	return avaliacoes, nil
}

func (r *AvaliacaoRepository) MediaPorPrestador(ctx context.Context, idPrestador uint) (float64, error) {
	var media float64
	err := r.db.WithContext(ctx).
		Table("avaliacaos"). // GORM pluralizes the table name
		Joins("JOIN servicos ON avaliacaos.servico_id = servicos.id").
		Where("servicos.id_prestador = ?", idPrestador).
		Select("COALESCE(AVG(avaliacaos.nota), 0)").
		Scan(&media).Error
	if err != nil {
		return 0, err
	}
	return media, nil
}

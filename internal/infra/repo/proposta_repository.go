package repo

import (
	"context"
	"time"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"gorm.io/gorm"
)

type PropostaRepository struct {
	db *gorm.DB
}

func NewPropostaRepository(db *gorm.DB) model.PropostaRepo {
	return &PropostaRepository{db: db}
}

func (r *PropostaRepository) Criar(ctx context.Context, proposta *model.Proposta) error {
	return r.db.WithContext(ctx).Create(proposta).Error
}

func (r *PropostaRepository) ListarPorVaga(ctx context.Context, idVaga uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Proposta, error) {
	var propostas []model.Proposta
	query := r.db.WithContext(ctx).Preload("Vaga").Preload("Prestador").Where("id_vaga = ?", idVaga)

	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}

	if orderDir == "desc" {
		query = query.Order(orderBy + " desc")
	} else {
		query = query.Order(orderBy + " asc")
	}

	err := query.Limit(limit).Offset(offset).Find(&propostas).Error
	if err != nil {
		return nil, err
	}
	return propostas, nil
}

func (r *PropostaRepository) ListarPorPrestador(ctx context.Context, idPrestador uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Proposta, error) {
	var propostas []model.Proposta
	query := r.db.WithContext(ctx).Preload("Vaga").Preload("Prestador").Where("id_prestador = ?", idPrestador)

	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}

	if orderDir == "desc" {
		query = query.Order(orderBy + " desc")
	} else {
		query = query.Order(orderBy + " asc")
	}

	err := query.Limit(limit).Offset(offset).Find(&propostas).Error
	if err != nil {
		return nil, err
	}
	return propostas, nil
}

func (r *PropostaRepository) AtualizarStatus(ctx context.Context, idProposta uint, status model.Status, dataResposta time.Time) error {
	return r.db.WithContext(ctx).
		Model(&model.Proposta{}).
		Where("id = ?", idProposta).
		Updates(map[string]interface{}{
			"status":         status,
			"data_resposta": dataResposta,
		}).Error
}
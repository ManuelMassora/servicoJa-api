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

func (r *PropostaRepository) Salvar(ctx context.Context, proposta *model.Proposta) error {
	return r.db.WithContext(ctx).Save(proposta).Error
}

func (r *PropostaRepository) BuscarPorID(ctx context.Context, idProposta uint) (*model.Proposta, error) {
	var proposta model.Proposta
	err := r.db.WithContext(ctx).Preload("Vaga").Preload("Prestador").First(&proposta, idProposta).Error
	if err != nil {
		return nil, err
	}
	return &proposta, nil
}

func (r *PropostaRepository) ListarPorVaga(ctx context.Context, idVaga uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Proposta, error) {
	var propostas []model.Proposta
	query := r.db.WithContext(ctx).Preload("Vaga").Preload("Prestador").Where("id_vaga = ?", idVaga)

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

	err := query.Find(&propostas).Error
	if err != nil {
		return nil, err
	}
	return propostas, nil
}

func (r *PropostaRepository) ListarPorPrestador(ctx context.Context, idPrestador uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Proposta, error) {
	var propostas []model.Proposta
	query := r.db.WithContext(ctx).Preload("Vaga").Preload("Prestador").Where("id_prestador = ?", idPrestador)

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

	err := query.Find(&propostas).Error
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
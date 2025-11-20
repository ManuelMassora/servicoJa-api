package repo

import (
	"context"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"gorm.io/gorm"
)

type VagaRepository struct {
	db *gorm.DB
}

func NewVagaRepository(db *gorm.DB) model.VagaRepo {
	return &VagaRepository{db: db}
}

func (r *VagaRepository) Criar(ctx context.Context, vaga *model.Vaga) error {
	return r.db.WithContext(ctx).Create(&vaga).Error
}

func (r *VagaRepository) Salvar(ctx context.Context, vaga *model.Vaga) error {
	return r.db.WithContext(ctx).Save(&vaga).Error
}

func (r *VagaRepository) BuscarPorID(ctx context.Context, id uint) (*model.Vaga, error) {
	var vaga model.Vaga
	err := r.db.WithContext(ctx).
		Preload("Cliente").
		Preload("Prestador").
		First(&vaga, id).Error
	if err != nil {
		return nil, err
	}
	return &vaga, nil
}

func (r *VagaRepository) ListarDisponiveis(ctx context.Context, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Vaga, error) {
	var vagas []model.Vaga
	query := r.db.WithContext(ctx).
		Preload("Cliente").
		Where("id_prestador IS NULL")
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
	err := query.Find(&vagas).Error
	if err != nil {
		return nil, err
	}
	return vagas, nil
}

func (r *VagaRepository)	ListarPorCliente(ctx context.Context, idCliente uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Vaga, error) {
	var vagas []model.Vaga
	query := r.db.WithContext(ctx).
		Preload("Cliente").
		Where("id_cliente = ?", idCliente)
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

	err := query.Find(&vagas).Error
	if err != nil {
		return nil, err
	}
	return vagas, nil
}

func (r *VagaRepository) AceitarVaga(ctx context.Context, idVaga, idPrestador uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Atualizar o prestador da vaga
		err := tx.Model(&model.Vaga{}).
			Where("id = ? AND id_prestador IS NULL", idVaga).
			Update("id_prestador", idPrestador).Error
		if err != nil {
			return err
		}

		// Atualizar o status da vaga para "em_andamento"
		return tx.Model(&model.Vaga{}).
			Where("id = ?", idVaga).
			Update("status", model.Status("em_andamento")).Error
	})
}

func (r *VagaRepository) AtualizarStatus(ctx context.Context, idVaga uint, status model.Status) error {
	return r.db.WithContext(ctx).
		Model(&model.Vaga{}).
		Where("id = ?", idVaga).
		Update("status", status).Error
}
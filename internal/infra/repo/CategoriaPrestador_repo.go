package repo

import (
	"context"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"gorm.io/gorm"
)

type CategoriaPrestadorRepo struct {
	db *gorm.DB
}

func NewCategoriaPrestadorRepo(db *gorm.DB) model.CategoriaPrestadorRepo {
	return &CategoriaPrestadorRepo{db: db}
}

func (r *CategoriaPrestadorRepo) Criar(ctx context.Context, categoria *model.CategoriaPrestador) (*model.CategoriaPrestador, error) {
	if err := r.db.WithContext(ctx).Create(categoria).Error; err != nil {
		return nil, err
	}
	return categoria, nil
}

func (r *CategoriaPrestadorRepo) Editar(ctx context.Context, id uint, campos map[string]interface{}) (*model.CategoriaPrestador, error) {
	var categoria model.CategoriaPrestador
	if err := r.db.WithContext(ctx).Model(&categoria).Where("id = ?", id).Updates(campos).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).First(&categoria, id).Error; err != nil {
		return nil, err
	}
	return &categoria, nil
}

func (r *CategoriaPrestadorRepo) Listar(ctx context.Context, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.CategoriaPrestador, error) {
	var categorias []model.CategoriaPrestador
	query := r.db.WithContext(ctx).Model(&model.CategoriaPrestador{})

	for field, value := range filters {
		if strVal, ok := value.(string); ok {
			query = query.Where(field+" LIKE ?", "%"+strVal+"%")
		} else {
			query = query.Where(field+" = ?", value)
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

	if err := query.Find(&categorias).Error; err != nil {
		return nil, err
	}
	return categorias, nil
}

func (r *CategoriaPrestadorRepo) BuscarPorID(ctx context.Context, id uint) (*model.CategoriaPrestador, error) {
	var categoria model.CategoriaPrestador
	err := r.db.WithContext(ctx).First(&categoria, id).Error
	if err != nil {
		return nil, err
	}
	return &categoria, nil
}

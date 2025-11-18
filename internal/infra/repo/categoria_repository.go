package repo

import (
	"context"
	"errors"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"gorm.io/gorm"
)

type CategoriaRepository struct {
	db *gorm.DB
}

func NewCategoriaRepository(db *gorm.DB) model.CategoriaRepo {
	return &CategoriaRepository{db: db}
}

func (r *CategoriaRepository) Criar(ctx context.Context, categoria *model.Categoria) error {
	return r.db.WithContext(ctx).Create(categoria).Error
}

func (r *CategoriaRepository) Editar(ctx context.Context, id uint, campos map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&model.Categoria{}).Where("id = ?", id).Updates(campos).Error
}

func (r *CategoriaRepository) Listar(ctx context.Context, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Categoria, error) {
	var categorias []model.Categoria
	query := r.db.WithContext(ctx)

	// Apply filters with LIKE for contains
	for field, value := range filters {
		query = query.Where(field+" LIKE ?", "%"+value.(string)+"%")
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

	err := query.Find(&categorias).Error
	if err != nil {
		return nil, err
	}
	return categorias, nil
}

func (r *CategoriaRepository) BuscarPorID(ctx context.Context, id uint) (*model.Categoria, error) {
	var categoria model.Categoria
	err := r.db.WithContext(ctx).First(&categoria, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &categoria, nil
}
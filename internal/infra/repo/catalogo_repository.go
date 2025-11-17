package repo

import (
	"context"
	"time"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"gorm.io/gorm"
)

type CatalogoRepository struct {
	db *gorm.DB
}

func NewCatalogoRepository(db *gorm.DB) model.CatalogoRepo {
	return &CatalogoRepository{db: db}
}

func (r *CatalogoRepository) Create(ctx context.Context,catalogo *model.Catalogo) error {
	return r.db.Create(catalogo).Error
}

func (r *CatalogoRepository) Update(ctx context.Context, id int64, campos map[string]interface{}) error {
	return r.db.Model(&model.Catalogo{}).Where("id=?", id).Updates(campos).Error
}

func (r *CatalogoRepository) Delete(ctx context.Context, id int64) error {
    return r.db.Model(&model.Catalogo{}).
        Where("id = ?", id).
        Update("deleted_at", time.Now()).Error 
}

func (r *CatalogoRepository) FindByID(ctx context.Context,id int64) (*model.Catalogo, error) {
	var catalogo model.Catalogo
	err := r.db.Preload("Prestador").First(&catalogo, id).Error
	if err != nil {
		return nil, err
	}
	return &catalogo, nil
}

func (r *CatalogoRepository) FindAll(ctx context.Context,filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]*model.Catalogo, error) {
	var catalogos []*model.Catalogo
	query := r.db.Preload("Prestador").Model(&model.Catalogo{})

	
	for key, value := range filters {
		query = query.Where(key+" LIKE ?", "%"+value.(string)+"%")
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
	if err := query.Find(&catalogos).Error;err != nil {
		return nil, err
	}
	return catalogos, nil
}

func (r *CatalogoRepository) FindByPrestadorID(ctx context.Context,prestadorID int64, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]*model.Catalogo, error) {
	var catalogos []*model.Catalogo
	query := r.db.Preload("Prestador").Model(&model.Catalogo{}).Where("id_prestador = ?", prestadorID)

	
	for key, value := range filters {
		query = query.Where(key+" LIKE ?", "%"+value.(string)+"%")
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
	if err := query.Find(&catalogos).Error;err != nil {
		return nil, err
	}
	return catalogos, nil
}
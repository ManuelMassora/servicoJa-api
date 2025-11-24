package repo

import (
	"context"
	"fmt"
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

func (r *CatalogoRepository) Update(ctx context.Context, id uint, campos map[string]interface{}) error {
	return r.db.Model(&model.Catalogo{}).Where("id=?", id).Updates(campos).Error
}

func (r *CatalogoRepository) Delete(ctx context.Context, id uint) error {
    return r.db.Model(&model.Catalogo{}).
        Where("id = ?", id).
        Update("deleted_at", time.Now()).Error
}

func (r *CatalogoRepository) FindByID(ctx context.Context,id uint) (*model.Catalogo, error) {
	var catalogo model.Catalogo
	err := r.db.Preload("Prestador").First(&catalogo, id).Error
	if err != nil {
		return nil, err
	}
	return &catalogo, nil
}

func (r *CatalogoRepository) FindAll(ctx context.Context,filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]*model.Catalogo, error) {
	var catalogos []*model.Catalogo
	query := r.db.Preload("Prestador").
	Preload("Prestador.Usuario").
	Preload("Categoria").
	Model(&model.Catalogo{})

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

func (r *CatalogoRepository) FindByPrestadorID(ctx context.Context,prestadorID uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]*model.Catalogo, error) {
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

func (r *CatalogoRepository) FindByLocation(ctx context.Context, latitude, longitude, radius float64, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]*model.Catalogo, error) {
    var catalogos []*model.Catalogo

    haversine := fmt.Sprintf(
        "6371 * acos(cos(radians(%f)) * cos(radians(latitude)) * cos(radians(longitude) - radians(%f)) + sin(radians(%f)) * sin(radians(latitude)))",
        latitude, longitude, latitude,
    )

    query := r.db.Select(fmt.Sprintf("*, (%s) AS distance", haversine)).
        Where(fmt.Sprintf("(%s) < ?", haversine), radius)

    for key, value := range filters {
        query = query.Where(fmt.Sprintf("%s LIKE ?", key), fmt.Sprintf("%%%v%%", value))
    }

    if orderBy != "" {
        if orderDir == "" {
            orderDir = "asc"
        }
        query = query.Order(fmt.Sprintf("%s %s", orderBy, orderDir))
    } else {
        query = query.Order("distance asc")
    }

    if limit > 0 {
        query = query.Limit(limit)
    }
    if offset > 0 {
        query = query.Offset(offset)
    }

    if err := query.Find(&catalogos).Error; err != nil {
        return nil, err
    }

    return catalogos, nil
}

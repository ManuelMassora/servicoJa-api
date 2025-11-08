package repo

import (
	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"gorm.io/gorm"
)

type CatalogoRepository struct {
	db *gorm.DB
}

func NewCatalogoRepository(db *gorm.DB) model.CatalogoRepo {
	return &CatalogoRepository{db: db}
}

func (r *CatalogoRepository) Create(catalogo *model.Catalogo) error {
	return r.db.Create(catalogo).Error
}

func (r *CatalogoRepository) Update(catalogo *model.Catalogo) error {
	return r.db.Save(catalogo).Error
}

func (r *CatalogoRepository) Delete(id int64) error {
	return r.db.Delete(&model.Catalogo{}, id).Error
}

func (r *CatalogoRepository) FindByID(id int64) (*model.Catalogo, error) {
	var catalogo model.Catalogo
	err := r.db.Preload("Prestador").First(&catalogo, id).Error
	if err != nil {
		return nil, err
	}
	return &catalogo, nil
}

func (r *CatalogoRepository) FindAll(filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]*model.Catalogo, error) {
	var catalogos []*model.Catalogo
	query := r.db.Preload("Prestador").Model(&model.Catalogo{})

	
	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}

	
	if orderBy != "" {
		query = query.Order(orderBy + " " + orderDir)
	}

	
	err := query.Limit(limit).Offset(offset).Find(&catalogos).Error
	if err != nil {
		return nil, err
	}
	return catalogos, nil
}

func (r *CatalogoRepository) FindByPrestadorID(prestadorID int64, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]*model.Catalogo, error) {
	var catalogos []*model.Catalogo
	query := r.db.Preload("Prestador").Model(&model.Catalogo{}).Where("id_prestador = ?", prestadorID)

	
	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}

	
	if orderBy != "" {
		query = query.Order(orderBy + " " + orderDir)
	}

	
	err := query.Limit(limit).Offset(offset).Find(&catalogos).Error
	if err != nil {
		return nil, err
	}
	return catalogos, nil
}
package repo

import (
	"context"
	"errors"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"gorm.io/gorm"
)

type GaleriaRepo struct {
	db *gorm.DB
}

func NewGaleriaRepo(db *gorm.DB) model.GaleriaRepo {
	return &GaleriaRepo{db: db}
}

func (r *GaleriaRepo) Create(ctx context.Context, galeria *model.Galeria) (*model.Galeria, error) {
	if err := r.db.WithContext(ctx).Create(galeria).Error; err != nil {
		return nil, err
	}
	return galeria, nil
}

func (r *GaleriaRepo) AddImage(ctx context.Context, imagem *model.Imagem) error {
	return r.db.WithContext(ctx).Create(&imagem).Error
}

func (r *GaleriaRepo) FindByID(ctx context.Context, id uint) (*model.Galeria, error) {
	var galeria model.Galeria
	err := r.db.WithContext(ctx).Preload("Imagens").First(&galeria, id).Error
	if err != nil {
		return nil, err
	}
	return &galeria, nil
}

func (r *GaleriaRepo) FindByPrestadorID(ctx context.Context, prestadorID uint) (*model.Galeria, error) {
    var galeria model.Galeria
    if err := r.db.WithContext(ctx).
	Preload("Imagens").
	Where("prestador_id = ?", prestadorID).First(&galeria).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &galeria, nil
}

func (r *GaleriaRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Select("Imagens").Delete(&model.Galeria{}, id).Error
}

func (r *GaleriaRepo) CountImages(ctx context.Context, galeriaID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Imagem{}).Where("galeria_id = ?", galeriaID).Count(&count).Error
	return count, err
}

func (r *GaleriaRepo) FindByGaleriaID(ctx context.Context, galeriaID uint) ([]model.Imagem, error) {
	var imagens []model.Imagem
	err := r.db.WithContext(ctx).Where("galeria_id = ?", galeriaID).Find(&imagens).Error
	return imagens, err
}

func(r *GaleriaRepo) FindByPrestadorIDs(ctx context.Context, prestadorIDs []uint) ([]model.Galeria, error) {
	var galerias []model.Galeria
	err := r.db.WithContext(ctx).
		Preload("Imagens").
		Where("prestador_id IN ?", prestadorIDs).
		Find(&galerias).Error
	return galerias, err
}
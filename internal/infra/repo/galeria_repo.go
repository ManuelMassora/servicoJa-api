package repo

import (
	"context"

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
	return r.db.WithContext(ctx).Create(imagem).Error
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
	err := r.db.WithContext(ctx).
		Where("prestador_id = ?", prestadorID).
		Find(&galeria).Error
	return &galeria, err
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
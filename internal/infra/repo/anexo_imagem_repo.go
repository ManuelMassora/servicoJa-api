package repo

import (
	"context"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"gorm.io/gorm"
)

type AnexoImagemRepo struct {
	db *gorm.DB
}

func NewAnexoImagemRepo(db *gorm.DB) model.AnexoImagemRepo {
	return &AnexoImagemRepo{db: db}
}

func (r *AnexoImagemRepo) Create(ctx context.Context, anexo *model.AnexoImagem) error {
	return r.db.WithContext(ctx).Create(anexo).Error
}

func (r *AnexoImagemRepo) FindByID(ctx context.Context, id uint) (*model.AnexoImagem, error) {
	var anexo model.AnexoImagem
	err := r.db.WithContext(ctx).First(&anexo, id).Error
	if err != nil {
		return nil, err
	}
	return &anexo, nil
}

func (r *AnexoImagemRepo) FindByAgendamentoID(ctx context.Context, agendamentoID uint) ([]model.AnexoImagem, error) {
	var anexos []model.AnexoImagem
	err := r.db.WithContext(ctx).Where("agendamento_id = ?", agendamentoID).Find(&anexos).Error
	if err != nil {
		return nil, err
	}
	return anexos, nil
}

func (r *AnexoImagemRepo) FindByVagaID(ctx context.Context, vagaID uint) ([]model.AnexoImagem, error) {
	var anexos []model.AnexoImagem
	err := r.db.WithContext(ctx).Where("vaga_id = ?", vagaID).Find(&anexos).Error
	if err != nil {
		return nil, err
	}
	return anexos, nil
}

func (r *AnexoImagemRepo) FindByCatalogoID(ctx context.Context, catalogoID uint) ([]model.AnexoImagem, error) {
	var anexos []model.AnexoImagem
	err := r.db.WithContext(ctx).Where("catalogo_id = ?", catalogoID).Find(&anexos).Error
	return anexos, err
}

func (r *AnexoImagemRepo) FindByVagaIDs(ctx context.Context, vagaIDs []uint) ([]model.AnexoImagem, error) {
	var anexos []model.AnexoImagem
	err := r.db.WithContext(ctx).Where("vaga_id IN ?", vagaIDs).Find(&anexos).Error
	if err != nil {
		return nil, err
	}
	return anexos, nil
}

func (r *AnexoImagemRepo) FindByCatalogoIDs(ctx context.Context, catalogoIDs []uint) ([]model.AnexoImagem, error) {
	var anexos []model.AnexoImagem
	err := r.db.WithContext(ctx).Where("catalogo_id IN ?", catalogoIDs).Find(&anexos).Error
	if err != nil {
		return nil, err
	}
	return anexos, nil
}

func (r *AnexoImagemRepo) FindByAgendamentoIDs(ctx context.Context, agendamentoIDs []uint) ([]model.AnexoImagem, error) {
	var anexos []model.AnexoImagem
	err := r.db.WithContext(ctx).Where("agendamento_id IN ?", agendamentoIDs).Find(&anexos).Error
	if err != nil {
		return nil, err
	}
	return anexos, nil
}
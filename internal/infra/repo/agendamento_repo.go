package repo

import (
	"context"
	"fmt"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"gorm.io/gorm"
)

type AgendamentoRepo struct {
	db *gorm.DB
}

func NewAgendamentoRepo(db *gorm.DB) model.AgendamentoRepo {
	return &AgendamentoRepo{db: db}
}

func (r *AgendamentoRepo) Criar(ctx context.Context, agendamento *model.Agendamento) (*model.Agendamento, error) {
	err := r.db.WithContext(ctx).Create(agendamento).Error
	if err != nil {
		return nil, err
	}
	return agendamento, nil
}
func (r *AgendamentoRepo) BuscarPorID(ctx context.Context, id uint) (*model.Agendamento, error) {
	var agendamento model.Agendamento
	err := r.db.WithContext(ctx).
		Preload("Cliente").
		Preload("Cliente.Usuario").
		Preload("Catalogo").
		Preload("Catalogo.Prestador.Usuario").
		Where("id=?", id).
		First(&agendamento).Error
	if err != nil {
		return nil, err
	}
	return &agendamento, nil
}

func (r *AgendamentoRepo) AtualizarStatus(ctx context.Context, id uint, status string) error {
	return r.db.WithContext(ctx).
		Model(&model.Agendamento{}).
		Where("id = ?", id).
		Update("status", status).Error
}
func (r *AgendamentoRepo) Listar(
	ctx context.Context,
	filters map[string]interface{},
	orderBy string,
	orderDir string,
	limit,
	offset int,
) ([]model.Agendamento, error) {
	var agendamentos []model.Agendamento
	query := r.db.WithContext(ctx).Model(&model.Agendamento{})
	query = query.
		Preload("Cliente").
		Preload("Catalogo").
		Preload("Catalogo.Prestador")
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
	if err := query.Find(&agendamentos).Error; err != nil {
		return nil, err
	}
	return agendamentos, nil
}

func (r *AgendamentoRepo) ListarPorClienteID(ctx context.Context, clienteID uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Agendamento, error) {
	var agendamentos []model.Agendamento
	query := r.db.Preload("Cliente").
		Preload("Catalogo").
		Preload("Catalogo.Prestador").
		Model(&model.Agendamento{}).
		Where("id_cliente = ?", clienteID)

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
	if err := query.Find(&agendamentos).Error; err != nil {
		return nil, err
	}
	return agendamentos, nil
}

func (r *AgendamentoRepo) ListarPorPrestadorID(ctx context.Context, prestadorID uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Agendamento, error) {
	var agendamentos []model.Agendamento
	query := r.db.Preload("Cliente").
		Preload("Catalogo").
		Preload("Catalogo.Prestador").
		Model(&model.Agendamento{}).
		Joins("LEFT JOIN catalogos ON agendamentos.id_catalogo = catalogos.id").
		Where("catalogos.id_prestador = ? AND agendamentos.status = ?", prestadorID, model.StatusPendente)

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
	if err := query.Find(&agendamentos).Error; err != nil {
		return nil, err
	}
	return agendamentos, nil
}

func (r *AgendamentoRepo) ListarPorCatalogID(ctx context.Context, catalogoID uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Agendamento, error) {
	var agendamentos []model.Agendamento
	query := r.db.Preload("Cliente").
		Preload("Catalogo").
		Preload("Catalogo.Prestador").
		Model(&model.Agendamento{}).
		Where("id_catalogo = ?", catalogoID)

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
	if err := query.Find(&agendamentos).Error; err != nil {
		return nil, err
	}
	return agendamentos, nil
}

func (r *AgendamentoRepo) FindByLocation(ctx context.Context, userID uint, latitude, longitude, radius float64, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Agendamento, error) {
	var agendamentos []model.Agendamento

	haversine := fmt.Sprintf(
		"6371 * acos(cos(radians(%f)) * cos(radians(agendamentos.latitude)) * cos(radians(agendamentos.longitude) - radians(%f)) + sin(radians(%f)) * sin(radians(agendamentos.latitude)))",
		latitude, longitude, latitude,
	)

	query := r.db.WithContext(ctx).Select(fmt.Sprintf("agendamentos.*, (%s) AS distance", haversine)).
		Joins("JOIN catalogos ON agendamentos.id_catalogo = catalogos.id").
		Where(fmt.Sprintf("(%s) < ?", haversine), radius).
		Where("agendamentos.id_cliente = ? OR catalogos.id_prestador = ?", userID, userID).
		Preload("Cliente").
		Preload("Catalogo").
		Preload("Catalogo.Prestador")

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

	if err := query.Find(&agendamentos).Error; err != nil {
		return nil, err
	}

	return agendamentos, nil
}

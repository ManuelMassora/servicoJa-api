package repo

import (
	"context"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
	"gorm.io/gorm"
)

type UsuarioRepository struct {
	db *gorm.DB
}

func NewUsuarioRepository(db *gorm.DB) model.UsuarioRepo {
	return &UsuarioRepository{db: db}
}

func (r *UsuarioRepository) Criar(ctx context.Context, usuario *model.Usuario) error {
	return r.db.WithContext(ctx).Create(usuario).Error
}

func (r *UsuarioRepository) BuscarPorID(ctx context.Context, id int64) (*model.Usuario, error) {
	var usuario model.Usuario
	err := r.db.WithContext(ctx).
		Preload("RolePermissao").
		First(&usuario, id).Error
	if err != nil {
		return nil, err
	}
	return &usuario, nil
}

func (r *UsuarioRepository) BuscarPorTelefone(ctx context.Context, telefone string) (*model.Usuario, error) {
	var usuario model.Usuario
	err := r.db.WithContext(ctx).
		Preload("RolePermissao").
		Where("telefone=?", telefone).
		First(&usuario).Error
	if err != nil {
		return nil, err
	}
	return &usuario, nil
}

func (r *UsuarioRepository) Atualizar(ctx context.Context, usuario *model.Usuario) error {
	return r.db.WithContext(ctx).Save(usuario).Error
}

func (r *UsuarioRepository) Remover(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Usuario{}, id).Error
}

func (r *UsuarioRepository) ListarTodos(ctx context.Context, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Usuario, error) {
	var usuarios []model.Usuario
	query := r.db.WithContext(ctx).Preload("RolePermissao")

	
	for field, value := range filters {
		query = query.Where(field+" LIKE ?", "%"+value.(string)+"%")
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

	err := query.Find(&usuarios).Error
	if err != nil {
		return nil, err
	}
	return usuarios, nil
}

type ClienteRepository struct {
	db *gorm.DB
}

func NewClienteRepository(db *gorm.DB) model.ClienteRepo {
	return &ClienteRepository{db: db}
}

func (r *ClienteRepository) Criar(ctx context.Context, cliente *model.Cliente) error {
	return r.db.WithContext(ctx).Create(cliente).Error
}

func (r *ClienteRepository) BuscarPorID(ctx context.Context, id int64) (*model.Cliente, error) {
	var cliente model.Cliente
	err := r.db.WithContext(ctx).
		Preload("Usuario").
		First(&cliente, id).Error
	if err != nil {
		return nil, err
	}
	return &cliente, nil
}

func (r *ClienteRepository) Listar(ctx context.Context, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Cliente, error) {
	var clientes []model.Cliente
	query := r.db.WithContext(ctx).Preload("Usuario")

	
	for field, value := range filters {
		query = query.Where(field+" = ?", value)
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

	err := query.Find(&clientes).Error
	if err != nil {
		return nil, err
	}
	return clientes, nil
}

type PrestadorRepository struct {
	db *gorm.DB
}

func NewPrestadorRepository(db *gorm.DB) model.PrestadorRepo {
	return &PrestadorRepository{db: db}
}

func (r *PrestadorRepository) Criar(ctx context.Context, prestador *model.Prestador) error {
	return r.db.WithContext(ctx).Create(prestador).Error
}

func (r *PrestadorRepository) AtualizarStatus(ctx context.Context, id int64, disponivel bool) error {
	return r.db.WithContext(ctx).
		Model(&model.Prestador{}).
		Where("id = ?", id).
		Update("status_disponivel", disponivel).Error
}

func (r *PrestadorRepository) BuscarPorID(ctx context.Context, id int64) (*model.Prestador, error) {
	var prestador model.Prestador
	err := r.db.WithContext(ctx).
		Preload("Usuario").
		First(&prestador, id).Error
	if err != nil {
		return nil, err
	}
	return &prestador, nil
}

func (r *PrestadorRepository) BuscarPorUsuarioID(ctx context.Context, id int64) (*model.Prestador, error) {
	var prestador model.Prestador
	err := r.db.WithContext(ctx).
		Preload("Usuario").
		Where("usuario_id", id).
		First(&prestador).Error
	if err != nil {
		return nil, err
	}
	return &prestador, nil
}

func (r *PrestadorRepository) Listar(ctx context.Context, filters map[string]interface{}, statusDisponivel interface{}, orderBy string, orderDir string, limit, offset int) ([]model.Prestador, error) {
	var prestadores []model.Prestador
	query := r.db.WithContext(ctx).Preload("Usuario")

	// 1. Aplica filtros de string (usando LIKE)
	for field, value := range filters {
		// Assume que todos os filtros no mapa 'filters' são strings para a busca LIKE
		// Você pode precisar de uma verificação de tipo mais robusta aqui,
		// mas seguindo sua implementação original, faremos o Type Assertion
		if strVal, ok := value.(string); ok {
			query = query.Where(field+" LIKE ?", "%"+strVal+"%")
		}
	}

	// 2. Aplica filtro StatusDisponivel (booleano, busca exata)
	// O parâmetro foi alterado para 'interface{}' para permitir nil ou um bool.
	if statusDisponivel != nil {
		if boolVal, ok := statusDisponivel.(bool); ok {
			query = query.Where("status_disponivel = ?", boolVal)
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

	err := query.Find(&prestadores).Error
	if err != nil {
		return nil, err
	}
	return prestadores, nil
}
package model

import "context"

type RolePermissao struct {
	BaseModel
	Role Role `gorm:"type:varchar(20);not null;unique;index" json:"role"`
}

type RoleRepo interface {
	ListarPermissoes(ctx context.Context, role Role) ([]RolePermissao, error)
	AdicionarPermissao(ctx context.Context, role Role) error
	RemoverPermissao(ctx context.Context, role Role) error
}
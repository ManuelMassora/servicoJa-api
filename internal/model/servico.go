package model

import "context"

type Servico struct {
	BaseModel
	Titulo      string   `gorm:"column:titulo;size:100;not null" json:"titulo"`
	Descricao   string   `gorm:"column:descricao;size:2000;not null" json:"descricao"`
	Localizacao string   `gorm:"column:localizacao;size:255;not null" json:"localizacao"`
	Preco       float64  `gorm:"column:preco;type:decimal(10,2);not null" json:"preco"`
	Status      Status   `gorm:"column:status;type:varchar(20);not null" json:"status"`
	IDCliente   int64    `gorm:"column:id_cliente;type:bigint;not null" json:"cliente_id"`
	Cliente     *Usuario `gorm:"foreignKey:IDCliente;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"cliente,omitempty"`
	IDPrestador *int64   `gorm:"column:id_prestador;type:bigint;default:null" json:"prestador_id,omitempty"`
	Prestador   *Usuario `gorm:"foreignKey:IDPrestador;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"prestador,omitempty"`
	IDCategoria *int64   `json:"categoria_id,omitempty"`
	Categoria   *Categoria `gorm:"foreignKey:IDCategoria;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"categoria,omitempty"` //Mudar pra catalogoID
}

type ServicoRepo interface {
	Criar(ctx context.Context, servico *Servico) error
	BuscarPorID(ctx context.Context, id int64) (*Servico, error)
	AtualizarStatus(ctx context.Context, id int64, status string) error
		ListarPorCliente(ctx context.Context, idCliente int64, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]Servico, error)
		ListarDisponiveis(ctx context.Context, localizacao string, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]Servico, error)
}

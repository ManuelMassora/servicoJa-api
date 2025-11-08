package model

import "context"

type Vaga struct {
	BaseModel
	Titulo      string   `gorm:"column:titulo;not null;size:100" json:"titulo"`
	Descricao   string   `gorm:"column:descricao;not null" json:"descricao"`
	Localizacao string   `gorm:"column:localizacao;not null" json:"localizacao"`
	Preco       float64  `gorm:"column:preco;not null;check:preco >= 0" json:"preco"`
	Status      Status   `gorm:"column:status;not null" json:"status"`
	IDCliente   int64    `gorm:"not null" json:"cliente_id"`
	Cliente     *Usuario `gorm:"foreignKey:IDCliente" json:"cliente,omitempty"`
	IDPrestador *int64   `json:"prestador_id,omitempty"` // nulo até alguém aceitar
	Prestador   *Usuario `gorm:"foreignKey:IDPrestador" json:"prestador,omitempty"`
	Urgente     bool     `gorm:"column:urgente;not null;default:false" json:"urgente"`
}

type VagaRepo interface {
	Criar(ctx context.Context, vaga Vaga) error
	BuscarPorID(ctx context.Context, id int64) (*Vaga, error)
	ListarDisponiveis(ctx context.Context, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]Vaga, error)
	AceitarVaga(ctx context.Context, idVaga, idPrestador int64) error
	AtualizarStatus(ctx context.Context, idVaga int64, status Status) error
}
package model

import "context"

type Vaga struct {
	BaseModel
	Titulo      string   `gorm:"column:titulo;not null;size:100" json:"titulo"`
	Descricao   string   `gorm:"column:descricao;not null" json:"descricao"`
	Localizacao string   `gorm:"column:localizacao;not null" json:"localizacao"`
	Latitude    float64  `gorm:"column:latitude;type:decimal(10,8);" json:"latitude"`
	Longitude   float64  `gorm:"column:longitude;type:decimal(11,8);" json:"longitude"`
	Preco       float64  `gorm:"column:preco;not null;check:preco >= 0" json:"preco"`
	Status      Status   `gorm:"column:status;not null" json:"status"`
	IDCliente   uint     `gorm:"not null" json:"cliente_id"`
	Cliente     *Cliente `gorm:"foreignKey:IDCliente" json:"cliente,omitempty"`
	IDPrestador *uint    `json:"prestador_id,omitempty"` // nulo até alguém aceitar
	Prestador   *Prestador `gorm:"foreignKey:IDPrestador" json:"prestador,omitempty"`
	Urgente     bool     `gorm:"column:urgente;not null;default:false" json:"urgente"`
	Anexos      []AnexoImagem `gorm:"foreignKey:VagaID"`
}

type VagaRepo interface {
	Criar(ctx context.Context, vaga *Vaga) error
	Salvar(ctx context.Context, vaga *Vaga) error
	BuscarPorID(ctx context.Context, id uint) (*Vaga, error)
	ListarDisponiveis(ctx context.Context, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]Vaga, error)
	ListarPorCliente(ctx context.Context, idCliente uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]Vaga, error)
	AceitarVaga(ctx context.Context, idVaga, idPrestador uint) error
	AtualizarStatus(ctx context.Context, idVaga uint, status Status) error
	FindByLocation(ctx context.Context, latitude, longitude, radius float64, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]Vaga, error)
}
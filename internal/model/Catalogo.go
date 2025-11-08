package model

type Catalogo struct {
	BaseModel
	Nome        string   `gorm:"column:nome;size:100;not null" json:"nome"`
	Descricao   string   `gorm:"column:descricao;size:2000;not null" json:"descricao"`
	PrecoBase   float64  `gorm:"column:preco_base;type:decimal(10,2);not null" json:"preco_base"`
	Categoria   string   `gorm:"column:categoria;size:100;not null" json:"categoria"`
	IDPrestador int64    `gorm:"column:id_prestador;type:bigint;not null" json:"prestador_id"`
	Disponivel  bool     `gorm:"column:disponivel;default:true" json:"disponivel"`
	Prestador   *Usuario `gorm:"foreignKey:IDPrestador;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"prestador,omitempty"`
}

type CatalogoRepo interface {
	Create(catalogo *Catalogo) error
	Update(catalogo *Catalogo) error
	Delete(id int64) error
	FindByID(id int64) (*Catalogo, error)
	FindAll(filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]*Catalogo, error)
	FindByPrestadorID(prestadorID int64, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]*Catalogo, error)
}

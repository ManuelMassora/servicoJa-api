package model

type Plano struct {
	BaseModel
	Nome        	string  `gorm:"column:nome;size:255;not null" json:"nome"`
	Descricao   	string  `gorm:"column:descricao;size:1000" json:"descricao"`
	Preco       	float64 `gorm:"column:preco;type:decimal(10,2);not null" json:"preco"`
	DuracaoDias 	int     `gorm:"column:duracao_dias;not null" json:"duracao_dias"`
	MaxServicos 	int     `gorm:"column:max_servicos;not null" json:"max_servicos"`
	MaxVagas    	int     `gorm:"column:max_vagas;not null" json:"max_vagas"`
	MaxAgendamentos int 	`gorm:"column:max_agendamentos;not null" json:"max_agendamentos"`
	MaxCatalogos 	int    	`gorm:"column:max_catalogos;not null" json:"max_catalogos"`
	MaxAnexos   	int    	`gorm:"column:max_anexos;not null" json:"max_anexos"`
}
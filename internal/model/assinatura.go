package model

type Assinatura struct {
	BaseModel
	UsuarioID     uint   `gorm:"column:usuario_id;not null" json:"usuario_id"`
	PlanoID       uint   `gorm:"column:plano_id;not null" json:"plano_id"`
	Status        string `gorm:"column:status;type:varchar(20);not null" json:"status"`
	DataExpiracao string `gorm:"column:data_expiracao;type:timestamp;not null" json:"data_expiracao"`
}

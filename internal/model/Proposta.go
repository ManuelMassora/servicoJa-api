package model

import (
	"context"
	"time"
)


type Proposta struct {
	BaseModel
	IDVaga        	uint      `json:"id_vaga"`
	Vaga			Vaga		`gorm:"foreignKey:IDVaga;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"vaga,omitempty"`
	IDPrestador   	uint      `json:"id_prestador"`
	Prestador		Usuario		`gorm:"foreignKey:IDPrestador;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"prestador,omitempty"`
	ValorProposto 	float64    `json:"valor_proposto"`
	Mensagem      	string     `json:"mensagem"`
	PrazoEstimado 	string     `json:"prazo_estimado"`
	Status        	Status     `json:"status"`
	DataResposta  	*time.Time `json:"data_resposta"`
}

type PropostaRepo interface {
	Criar(ctx context.Context, proposta *Proposta) error
	ListarPorVaga(ctx context.Context, idVaga uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]Proposta, error)
	ListarPorPrestador(ctx context.Context, idPrestador uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]Proposta, error)
	AtualizarStatus(ctx context.Context, idProposta uint, status Status, dataResposta time.Time) error
}
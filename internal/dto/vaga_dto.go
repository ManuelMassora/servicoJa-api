package dto

import "github.com/ManuelMassora/servicoJa-api/internal/model"

type VagaInput struct {
	Titulo      string             `json:"titulo" form:"titulo" binding:"required"`
	Descricao   string             `json:"descricao" form:"descricao" binding:"required"`
	Localizacao string             `json:"localizacao" form:"localizacao" binding:"required"`
	Latitude    float64            `json:"latitude" form:"latitude" binding:"required"`
	Longitude   float64            `json:"longitude" form:"longitude" binding:"required"`
	Preco       float64            `json:"preco" form:"preco" binding:"required"`
	Status      model.Status       `json:"status" form:"status" binding:"required"`
	IDCliente   uint               `json:"cliente_id" form:"cliente_id" binding:"required"`
	Urgente     bool               `json:"urgente" form:"urgente"`
	Anexos      []AnexoImagemInput `json:"anexos" form:"anexos"`
}
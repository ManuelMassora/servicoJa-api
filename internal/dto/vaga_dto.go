package dto

import "github.com/ManuelMassora/servicoJa-api/internal/model"

type VagaInput struct {
	Titulo      string    `json:"titulo" binding:"required"`
	Descricao   string    `json:"descricao" binding:"required"`
	Localizacao string    `json:"localizacao" binding:"required"`
	Latitude    float64   `json:"latitude" binding:"required"`
	Longitude   float64   `json:"longitude" binding:"required"`
	Preco       float64   `json:"preco" binding:"required"`
	Status      model.Status `json:"status" binding:"required"`
	IDCliente   uint      `json:"cliente_id" binding:"required"`
	Urgente     bool      `json:"urgente"`
	Anexos      []AnexoImagemInput `json:"anexos"`
}

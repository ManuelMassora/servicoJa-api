package dto

import "time"

type AgendamentoInput struct {
	Detalhe     string           `json:"detalhe" form:"detalhe" binding:"required"`
	IDCatalogo  uint             `json:"id_catalogo" form:"id_catalogo" binding:"required"`
	IDCliente   uint             `json:"id_cliente" form:"id_cliente" binding:"required"`
	DataHora    time.Time        `json:"datahora" form:"datahora" binding:"required"`
	Status      string           `json:"status" form:"status" binding:"required"`
	Localizacao string           `json:"localizacao" form:"localizacao"`
	Latitude    float64          `json:"latitude" form:"latitude"`
	Longitude   float64          `json:"longitude" form:"longitude"`
	Anexos      []AnexoImagemInput `json:"anexos" form:"anexos"`
}

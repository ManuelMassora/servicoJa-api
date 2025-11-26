package dto

import "time"

type AgendamentoInput struct {
	Detalhe     string    `json:"detalhe" binding:"required"`
	IDCatalogo  uint      `json:"id_catalogo" binding:"required"`
	IDCliente   uint      `json:"id_cliente" binding:"required"`
	DataHora    time.Time `json:"datahora" binding:"required"`
	Status      string    `json:"status" binding:"required"`
	Localizacao string    `json:"localizacao"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Anexos      []AnexoImagemInput `json:"anexos"`
}

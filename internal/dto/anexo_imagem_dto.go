package dto

type AnexoImagemInput struct {
	URL           string `json:"url" binding:"required"`
	AgendamentoID *uint  `json:"agendamento_id"`
	VagaID        *uint  `json:"vaga_id"`
	CatalogoID    *uint  `json:"catalogo_id"`
}

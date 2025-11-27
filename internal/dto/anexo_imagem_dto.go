package dto

type AnexoImagemInput struct {
	URL           string `json:"url" form:"url"`
	AgendamentoID *uint  `json:"agendamento_id" form:"agendamento_id"`
	VagaID        *uint  `json:"vaga_id" form:"vaga_id"`
	CatalogoID    *uint  `json:"catalogo_id" form:"catalogo_id"`
}

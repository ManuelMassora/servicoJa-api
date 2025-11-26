package dto

type CatalogoInput struct {
	Nome         string    `json:"nome" binding:"required"`
	Descricao    string    `json:"descricao" binding:"required"`
	TipoPreco    string    `json:"tipo_preco" binding:"required"`
	ValorFixo    float64   `json:"valor_fixo"`
	ValorPorHora float64   `json:"valor_por_hora"`
	IDCategoria  uint      `json:"categoria_id" binding:"required"`
	Disponivel   bool      `json:"disponivel"`
	IDPrestador  uint      `json:"prestador_id" binding:"required"`
	Localizacao  string    `json:"localizacao"`
	Latitude     float64   `json:"latitude"`
	Longitude    float64   `json:"longitude"`
	Anexos       []AnexoImagemInput `json:"anexos"`
}

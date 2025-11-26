package dto

type GaleriaInput struct {
	PrestadorID uint          `json:"prestador_id" binding:"required"`
	Imagens     []ImagemInput `json:"imagens" binding:"required"`
}

type ImagemInput struct {
	URL string `json:"url" binding:"required"`
}

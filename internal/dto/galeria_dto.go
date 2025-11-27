package dto

type GaleriaInput struct {
	Imagens []ImagemInput `json:"imagens" form:"imagens" binding:"required"`
}

type ImagemInput struct {
	URL string `json:"url" form:"url" binding:"required"`
}

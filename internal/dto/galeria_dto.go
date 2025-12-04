package dto

type GaleriaInput struct {
	Imagens []string `binding:"-"`
}

type ImagemInput struct {
	URL string `binding:"-"`
}
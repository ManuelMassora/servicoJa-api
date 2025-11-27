package dto

type ClienteInput struct {
	Nome      string `json:"nome" form:"nome" binding:"required"`
	Telefone  string `json:"telefone" form:"telefone" binding:"required"`
	Senha     string `json:"senha" form:"senha" binding:"required"`
	ImagemURL string
}

type PrestadorInput struct {
    Nome        string  `json:"nome" form:"nome" binding:"required"`
    Telefone    string  `json:"telefone" form:"telefone" binding:"required"`
    Senha       string  `json:"senha" form:"senha" binding:"required"`
    ImagemURL   string  `json:"-" form:"-"` // não vem do form
    Localizacao string  `json:"localizacao" form:"localizacao" binding:"required"`
    Latitude    float64 `json:"latitude" form:"latitude" binding:"required"`
    Longitude   float64 `json:"longitude" form:"longitude" binding:"required"`
}

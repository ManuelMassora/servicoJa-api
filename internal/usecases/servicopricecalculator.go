package usecases

import (
	"log"
	"time"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
)

func CalculateFinalServicePrice(catalogo *model.Catalogo, duration time.Duration) float64 {
	if catalogo.TipoPreco == "por_hora" {
		minutes := duration.Minutes()
		log.Println("minutos: ", minutes)
		hours := minutes/60
		log.Println("horas: ", hours)

		total := catalogo.ValorPorHora * hours
		log.Println("total: ", total)
		return total
	}
	return catalogo.ValorFixo
}

package usecases

import (
	"time"

	"github.com/ManuelMassora/servicoJa-api/internal/model"
)

// CalculateFinalServicePrice calculates the final price of a service based on its catalog and duration.
func CalculateFinalServicePrice(catalogo *model.Catalogo, duration time.Duration) float64 {
	if catalogo.TipoPreco == "por_hora" {
		hours := duration.Hours()
		// Round up to the nearest hour if there's any partial hour
		if hours > float64(int(hours)) {
			hours = float64(int(hours) + 1)
		}
		return catalogo.ValorPorHora * hours
	}
	// For fixed price, return the fixed value. This should ideally be set during service creation.
	return catalogo.ValorFixo
}

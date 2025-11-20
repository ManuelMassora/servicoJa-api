package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// ExtractFilters extrai parâmetros da query string e os converte em um mapa de filtros.
// A função é pública pois começa com letra maiúscula (E).
func ExtractFilters(c *gin.Context) map[string]interface{} {
	filters := make(map[string]interface{})
	
	for key, vals := range c.Request.URL.Query() {
		// Ignora parâmetros de paginação/ordenação
		if key == "orderBy" || key == "orderDir" || key == "page" || key == "pageSize" {
			continue
		}
		if len(vals) == 0 || vals[0] == "" {
			continue
		}

		v := vals[0]

		// Tentativa de conversão de tipos comuns
		if key == "id" {
			if id, err := strconv.ParseUint(v, 10, 64); err == nil {
				filters[key] = id
				continue
			}
		}
		// Exemplo de conversão booleana, se necessário:
		// if key == "status_bool" {
		// 	if b, err := strconv.ParseBool(v); err == nil {
		// 		filters[key] = b
		// 		continue
		// 	}
		// }
		
		// Valor padrão: string
		filters[key] = v
	}
	return filters
}

// ExtractPagination extrai os parâmetros de paginação e retorna limit, offset, page e pageSize.
// A função é pública pois começa com letra maiúscula (E).
func ExtractPagination(c *gin.Context) (limit, offset, page, pageSize int) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err = strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	offset = (page - 1) * pageSize
	limit = pageSize
	return limit, offset, page, pageSize
}
package handler

import (
	"net/http"
	"strconv"

	"github.com/ManuelMassora/servicoJa-api/internal/usecases"
	"github.com/gin-gonic/gin"
)

type ServicoHandler struct {
	uc usecases.ServicoUseCase
}

func NewServicoHandler(uc usecases.ServicoUseCase) *ServicoHandler {
	return &ServicoHandler{uc: uc}
}

func (h *ServicoHandler) FinalizarServico(c *gin.Context) {
	idServico, err := getServicoID(c) // Assumindo que 'getServicoID' extrai o ID do parâmetro da URL
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do serviço inválido"})
		return
	}

	idUsuario, err := getUsuarioID(c) // Assumindo que 'getUsuarioID' extrai o ID do usuário do contexto (por exemplo, do token JWT)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	err = h.uc.FinalizarServico(c.Request.Context(), idServico, uint(idUsuario))
	if err != nil {
		if err.Error() == "usuário não autorizado a finalizar este serviço" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Serviço finalizado com sucesso"})
}

func (h *ServicoHandler) ConfirmarServico(c *gin.Context) {
	idServico, err := getServicoID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do serviço inválido"})
		return
	}

	idUsuario, err := getUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	err = h.uc.ConfirmarServico(c.Request.Context(), idServico, uint(idUsuario))
	if err != nil {
		if err.Error() == "usuário não autorizado confirmar este serviço" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Serviço finalizado com sucesso"})
}

func (h *ServicoHandler) CancelarServico(c *gin.Context) {
	idServico, err := getServicoID(c) // Assumindo que 'getServicoID' extrai o ID do parâmetro da URL
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do serviço inválido"})
		return
	}

	idUsuario, err := getUsuarioID(c) // Assumindo que 'getUsuarioID' extrai o ID do usuário do contexto
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	err = h.uc.CancelarServico(c.Request.Context(), idServico, uint(idUsuario))
	if err != nil {
		if err.Error() == "usuário não autorizado a finalizar este serviço" {
			// Nota: O erro da UC é "usuário não autorizado a finalizar este serviço", mesmo para Cancelar.
			// É recomendável que a UC retorne uma mensagem de erro mais específica para Cancelar.
			c.JSON(http.StatusForbidden, gin.H{"error": "usuário não autorizado a cancelar este serviço"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Serviço cancelado com sucesso"})
}

func (h *ServicoHandler) ListarPorCliente(c *gin.Context) {
	idUsuario, err := getUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	filters := make(map[string]interface{})

	// Lógica de extração de filtros, ordenação e paginação, idêntica ao exemplo fornecido
	for key, vals := range c.Request.URL.Query() {
		if key == "orderBy" || key == "orderDir" || key == "page" || key == "pageSize" {
			continue
		}
		if len(vals) == 0 || vals[0] == "" {
			continue
		}

		v := vals[0]

		if key == "id" {
			if id, err := strconv.ParseUint(v, 10, 64); err == nil {
				filters[key] = id
				continue
			}
		}

		if key == "status" {
			// A verificação de status no AgendamentoHandler era para 'bool',
			// mas um status de Serviço geralmente é uma string.
			// Mantendo a lógica de usar o valor como string, a menos que seja 'id' ou 'status' (com parse bool no exemplo original)
			filters[key] = v
			continue
		}

		filters[key] = v
	}

	orderBy := c.Query("orderBy")
	orderDir := c.Query("orderDir")

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	servicos, err := h.uc.ListarPorCliente(
		c.Request.Context(),
		uint(idUsuario),
		filters,
		orderBy,
		orderDir,
		pageSize,
		offset,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      servicos,
		"page":      page,
		"pageSize":  pageSize,
		"orderBy":   orderBy,
		"direction": orderDir,
		"filters":   filters,
	})
}

func (h *ServicoHandler) ListarPorPrestador(c *gin.Context) {
	idUsuario, err := getUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	filters := make(map[string]interface{})

	// Lógica de extração de filtros, ordenação e paginação, idêntica ao exemplo fornecido
	for key, vals := range c.Request.URL.Query() {
		if key == "orderBy" || key == "orderDir" || key == "page" || key == "pageSize" {
			continue
		}
		if len(vals) == 0 || vals[0] == "" {
			continue
		}

		v := vals[0]

		if key == "id" {
			if id, err := strconv.ParseUint(v, 10, 64); err == nil {
				filters[key] = id
				continue
			}
		}

		if key == "status" {
			filters[key] = v
			continue
		}

		filters[key] = v
	}

	orderBy := c.Query("orderBy")
	orderDir := c.Query("orderDir")

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	servicos, err := h.uc.ListarPorPrestador(
		c.Request.Context(),
		uint(idUsuario),
		filters,
		orderBy,
		orderDir,
		pageSize,
		offset,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      servicos,
		"page":      page,
		"pageSize":  pageSize,
		"orderBy":   orderBy,
		"direction": orderDir,
		"filters":   filters,
	})
}

func getServicoID(c *gin.Context) (uint, error) {
	idStr := c.Param("id") // Assume que o ID é passado como um parâmetro de rota chamado "id"
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

func (h *ServicoHandler) ListarPorLocalizacao(c *gin.Context) {
	idUsuario, err := getUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	latitude, err := strconv.ParseFloat(c.Query("latitude"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Latitude inválida"})
		return
	}
	longitude, err := strconv.ParseFloat(c.Query("longitude"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Longitude inválida"})
		return
		}
	radius, err := strconv.ParseFloat(c.DefaultQuery("radius", "10"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Raio inválido"})
		return
	}

	filters := make(map[string]interface{})
	for key, vals := range c.Request.URL.Query() {
		if key == "latitude" || key == "longitude" || key == "radius" || key == "orderBy" || key == "orderDir" || key == "page" || key == "pageSize" {
			continue
		}
		if len(vals) > 0 && vals[0] != "" {
			filters[key] = vals[0]
		}
	}

	orderBy := c.Query("orderBy")
	orderDir := c.Query("orderDir")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	offset := (page - 1) * pageSize

	servicos, err := h.uc.ListarPorLocalizacao(c.Request.Context(), uint(idUsuario), latitude, longitude, radius, filters, orderBy, orderDir, pageSize, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      servicos,
		"page":      page,
		"pageSize":  pageSize,
		"orderBy":   orderBy,
		"direction": orderDir,
		"filters":   filters,
	})
}
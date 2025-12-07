package handler

import (
	"bytes"
	"fmt"
	"mime"
	"mime/multipart"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/ManuelMassora/servicoJa-api/internal/services"
	"github.com/ManuelMassora/servicoJa-api/internal/usecases"
	"github.com/ManuelMassora/servicoJa-api/pkg"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type AgendamentoHandler struct {
	uc       usecases.AgendamentoUC
	uploader *services.SupabaseUploader
}

func NewAgendamentoHandler(uc usecases.AgendamentoUC, uploader *services.SupabaseUploader) *AgendamentoHandler {
	return &AgendamentoHandler{uc: uc, uploader: uploader}
}

func (h *AgendamentoHandler) Criar(c *gin.Context) {
	var req usecases.AgendamentoRequest

	if err := c.ShouldBindWith(&req, binding.FormMultipart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	form, err := c.MultipartForm()
	
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao processar formulário multipart: " + err.Error()})
		return
	}

	files := form.File["anexos"]
	if len(files) > 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Limite de 3 imagens por agendamento excedido."})
		return
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(files))
	urlsCh := make(chan string, len(files))

	for _, file := range files {
		wg.Add(1)
		go func(fileHeader *multipart.FileHeader) {
			defer wg.Done()

			compressedBuf, format, err := pkg.CompressImage(fileHeader, 150)
			if err != nil {
				errCh <- fmt.Errorf("falha ao processar imagem para upload: %w", err)
				return
			}

			fileName := fmt.Sprintf("%d.%s", time.Now().UnixNano(), format)
			contentType := mime.TypeByExtension("." + format)
			if contentType == "" {
				contentType = "application/octet-stream"
			}

			_, uploadedFileName, err := h.uploader.UploadFromReader(c.Request.Context(), bytes.NewReader(compressedBuf.Bytes()), fileName, contentType)
			if err != nil {
				errCh <- fmt.Errorf("falha ao fazer upload do anexo: %w", err)
				return
			}

			publicURL := h.uploader.GetPublicURL("serviceja-image", uploadedFileName)
			urlsCh <- publicURL
		}(file)
	}

	wg.Wait()
	close(errCh)
	close(urlsCh)

	// Verifica se ocorreu algum erro durante o upload
	for err := range errCh {
		if err != nil {
			// Retorna o primeiro erro encontrado
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Coleta todas as URLs dos uploads bem-sucedidos
	var urls []string
	for url := range urlsCh {
		urls = append(urls, url)
	}
	req.Anexos = urls

	idCliente, err := getUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := h.uc.Criar(c.Request.Context(), &req, idCliente); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao criar agendamento: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Agendamento criado com sucesso!"})
}

func (h *AgendamentoHandler) Buscar(c *gin.Context) {

	idParam := c.Param("id")
	agendamentoID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de agendamento inválido"})
		return
	}

	idUsuario, err := getUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.uc.Buscar(c.Request.Context(), uint(agendamentoID), idUsuario)
	if err != nil {
		if err.Error() == "acesso negado: você não é o cliente nem o prestador deste agendamento" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Agendamento não encontrado ou falha na busca: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AgendamentoHandler) Aceitar(c *gin.Context) {

	idParam := c.Param("id")
	agendamentoID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de agendamento inválido"})
		return
	}

	idUsuario, err := getUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	err = h.uc.Aceitar(c.Request.Context(), uint(agendamentoID), idUsuario)
	if err != nil {
		if err.Error() == "acesso negado: você não é o cliente nem o prestador deste agendamento" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Agendamento não encontrado ou falha na busca: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (h *AgendamentoHandler) Recusar(c *gin.Context) {

	idParam := c.Param("id")
	agendamentoID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de agendamento inválido"})
		return
	}

	idUsuario, err := getUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	err = h.uc.Recusar(c.Request.Context(), uint(agendamentoID), idUsuario)
	if err != nil {
		if err.Error() == "acesso negado: você não é o cliente nem o prestador deste agendamento" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Agendamento não encontrado ou falha na busca: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (h *AgendamentoHandler) Cancelar(c *gin.Context) {

	idParam := c.Param("id")
	agendamentoID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de agendamento inválido"})
		return
	}

	idUsuario, err := getUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	err = h.uc.Cancelar(c.Request.Context(), uint(agendamentoID), idUsuario)
	if err != nil {
		if err.Error() == "acesso negado: você não é o cliente nem o prestador deste agendamento" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Agendamento não encontrado ou falha na busca: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (h *AgendamentoHandler) Listar(c *gin.Context) {
	filters := make(map[string]interface{})

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
			if b, err := strconv.ParseBool(v); err == nil {
				filters[key] = b
				continue
			}
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

	agendamentos, err := h.uc.Listar(
		c.Request.Context(),
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
		"data":      agendamentos,
		"page":      page,
		"pageSize":  pageSize,
		"orderBy":   orderBy,
		"direction": orderDir,
		"filters":   filters,
	})
}

func (h *AgendamentoHandler) ListarPorClienteID(c *gin.Context) {
	idUsuario, err := getUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	filters := make(map[string]interface{})

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
			if b, err := strconv.ParseBool(v); err == nil {
				filters[key] = b
				continue
			}
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

	agendamentos, err := h.uc.ListarPorClienteID(
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
		"data":      agendamentos,
		"page":      page,
		"pageSize":  pageSize,
		"orderBy":   orderBy,
		"direction": orderDir,
		"filters":   filters,
	})
}

func (h *AgendamentoHandler) ListarPorPrestadorID(c *gin.Context) {
	idUsuario, err := getUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	filters := make(map[string]interface{})

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
			if b, err := strconv.ParseBool(v); err == nil {
				filters[key] = b
				continue
			}
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

	agendamentos, err := h.uc.ListarPorPrestadorIDAgrupado(
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
		"data":      agendamentos,
		"page":      page,
		"pageSize":  pageSize,
		"orderBy":   orderBy,
		"direction": orderDir,
		"filters":   filters,
	})
}

func (h *AgendamentoHandler) ListarPorCatalogID(c *gin.Context) {
	idUsuario, err := getUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	idParam := c.Param("catalogoID")
	idCatalogo, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de catálogo inválido"})
		return
	}

	filters := make(map[string]interface{})

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
			if b, err := strconv.ParseBool(v); err == nil {
				filters[key] = b
				continue
			}
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

	agendamentos, err := h.uc.ListarPorCatalogID(
		c.Request.Context(),
		uint(idUsuario),
		uint(idCatalogo),
		filters,
		orderBy,
		orderDir,
		pageSize,
		offset,
	)
	if err != nil {
		if err.Error() == "acesso negado: você não é o prestador deste catálogo" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      agendamentos,
		"page":      page,
		"pageSize":  pageSize,
		"orderBy":   orderBy,
		"direction": orderDir,
		"filters":   filters,
	})
}

func (h *AgendamentoHandler) ListarPorLocalizacao(c *gin.Context) {
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

	agendamentos, err := h.uc.ListarPorLocalizacao(c.Request.Context(), idUsuario, latitude, longitude, radius, filters, orderBy, orderDir, pageSize, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      agendamentos,
		"page":      page,
		"pageSize":  pageSize,
		"orderBy":   orderBy,
		"direction": orderDir,
		"filters":   filters,
	})
}
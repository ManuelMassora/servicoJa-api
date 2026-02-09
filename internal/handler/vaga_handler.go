package handler

import (
	"bytes"
	"fmt"
	"mime"
	"net/http"
	"strconv"
	"time"

	"github.com/ManuelMassora/servicoJa-api/internal/services"
	"github.com/ManuelMassora/servicoJa-api/internal/usecases"
	"github.com/ManuelMassora/servicoJa-api/pkg"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type VagaHandler struct {
	uc       usecases.VagaUseCase
	uploader *services.SupabaseUploader
}

func NewVagaHandler(uc usecases.VagaUseCase, uploader *services.SupabaseUploader) *VagaHandler {
	return &VagaHandler{uc: uc, uploader: uploader}
}

func (h *VagaHandler) CriarVaga(c *gin.Context) {
	var req usecases.VagaRequest
	if err := c.ShouldBindWith(&req, binding.FormMultipart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados de requisição inválidos", "details": err.Error()})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao processar formulário multipart: " + err.Error()})
		return
	}

	files := form.File["anexos"]
	if len(files) > 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Limite de 3 imagens por vaga excedido."})
		return
	}

	for _, file := range files {
		// Validação rigorosa da imagem
		if err := pkg.ValidateImage(file); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		compressedBuf, format, err := pkg.CompressImage(file, 150)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao processar imagem para upload: " + err.Error()})
			return
		}

		// Gera um nome de arquivo único
		fileName := fmt.Sprintf("%d.%s", time.Now().UnixNano(), format)

		// Determina o Content-Type
		contentType := mime.TypeByExtension("." + format)
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		_, uploadedFileName, err := h.uploader.UploadFromReader(c.Request.Context(), bytes.NewReader(compressedBuf.Bytes()), fileName, contentType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao fazer upload do anexo: " + err.Error()})
			return
		}
		publicURL := h.uploader.GetPublicURL("serviceja-image", uploadedFileName)
		req.Anexos = append(req.Anexos, publicURL)
	}

	idUsuario, err := getUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	err = h.uc.CriarVaga(c.Request.Context(), req, uint(idUsuario))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Vaga criada com sucesso"})
}

func (h *VagaHandler) CancelarVaga(c *gin.Context) {
	idVagaStr := c.Param("id")
	idVaga, err := strconv.ParseUint(idVagaStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID da vaga inválido"})
		return
	}

	idUsuario, err := getUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	err = h.uc.CancelarVaga(c.Request.Context(), uint(idVaga), uint(idUsuario))
	if err != nil {
		if err.Error() == "vaga não pertence ao cliente" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Vaga cancelada com sucesso"})
}

func (h *VagaHandler) ListarVagasDisponiveis(c *gin.Context) {
	filters := ExtractFilters(c)
	orderBy := c.Query("orderBy")
	orderDir := c.Query("orderDir")
	limit, offset, page, pageSize := ExtractPagination(c)

	vagas, err := h.uc.ListarVagasDisponiveis(
		c.Request.Context(),
		filters,
		orderBy,
		orderDir,
		limit,
		offset,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":      vagas,
		"page":      page,
		"pageSize":  pageSize,
		"orderBy":   orderBy,
		"direction": orderDir,
		"filters":   filters,
	})
}

func (h *VagaHandler) ListarPorCliente(c *gin.Context) {
	idUsuario, err := getUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	filters := ExtractFilters(c)
	orderBy := c.Query("orderBy")
	orderDir := c.Query("orderDir")
	limit, offset, page, pageSize := ExtractPagination(c)

	vagas, err := h.uc.ListarPorCliente(
		c.Request.Context(),
		uint(idUsuario),
		filters,
		orderBy,
		orderDir,
		limit,
		offset,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      vagas,
		"page":      page,
		"pageSize":  pageSize,
		"orderBy":   orderBy,
		"direction": orderDir,
		"filters":   filters,
	})
}

func (h *VagaHandler) ListarPorLocalizacao(c *gin.Context) {
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

	vagas, err := h.uc.ListarPorLocalizacao(c.Request.Context(), latitude, longitude, radius, filters, orderBy, orderDir, pageSize, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      vagas,
		"page":      page,
		"pageSize":  pageSize,
		"orderBy":   orderBy,
		"direction": orderDir,
		"filters":   filters,
	})
}

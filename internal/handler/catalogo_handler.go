package handler

import (
	"bytes"
	"context"
	"fmt"
	"mime"

	// "mime/multipart"
	"net/http"
	"strconv"

	// "sync"
	"time"

	"github.com/ManuelMassora/servicoJa-api/internal/services"
	"github.com/ManuelMassora/servicoJa-api/internal/usecases"
	"github.com/ManuelMassora/servicoJa-api/pkg"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"golang.org/x/sync/errgroup"
)

type CatalogoUseCase interface {
	Criar(ctx context.Context, req usecases.RequestCreateCatalogo, prestadorID uint) (uint, error)
	Editar(ctx context.Context, id uint, prestadorID uint, campos map[string]interface{}) error
	Apagar(ctx context.Context, id uint, prestadorID uint) error
	Listar(ctx context.Context, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]usecases.ResponseCatalogo, error)
	ListarPorPrestador(ctx context.Context, prestadorID uint, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]usecases.ResponseCatalogo, error)
	ListarPorLocalizacao(ctx context.Context, latitude, longitude, radius float64, filters map[string]interface{}, orderBy string, orderDir string, limit, offset int) ([]usecases.ResponseCatalogo, error)
}

type CatalogoHandler struct {
	uc       CatalogoUseCase
	uploader *services.SupabaseUploader
}

func NewCatalogoHandler(uc CatalogoUseCase, uploader *services.SupabaseUploader) *CatalogoHandler {
	return &CatalogoHandler{uc: uc, uploader: uploader}
}

func (h *CatalogoHandler) Criar(c *gin.Context) {
	var request usecases.RequestCreateCatalogo
	if err := c.ShouldBindWith(&request, binding.FormMultipart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao processar formulário multipart: " + err.Error()})
		return
	}

	files := form.File["anexos"]
	if len(files) > 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Limite de 3 imagens por catálogo excedido."})
		return
	}

	g, ctx := errgroup.WithContext(c.Request.Context())
	type result struct {
		idx int
		url string
	}
	resCh := make(chan result, len(files))

	for i, file := range files {
		i := i
		file := file
		g.Go(func() error {
			// Validação rigorosa da imagem
			if err := pkg.ValidateImage(file); err != nil {
				return err
			}

			comp, format, err := pkg.CompressImage(file, 150)
			if err != nil {
				return fmt.Errorf("compress: %w", err)
			}

			fileName := fmt.Sprintf("%d.%s", time.Now().UnixNano(), format)
			contentType := mime.TypeByExtension("." + format)
			if contentType == "" {
				contentType = "application/octet-stream"
			}

			_, uploadedFileName, err := h.uploader.UploadFromReader(ctx, bytes.NewReader(comp.Bytes()), fileName, contentType)
			if err != nil {
				return fmt.Errorf("upload: %w", err)
			}

			resCh <- result{idx: i, url: h.uploader.GetPublicURL("serviceja-image", uploadedFileName)}
			return nil
		})
	}

	go func() {
		if err := g.Wait(); err != nil {
			// g will cancel ctx automatically, handle error later
		}
		close(resCh)
	}()

	// coletar resultados e ordenar por idx
	var tmp = make([]string, len(files))
	for r := range resCh {
		tmp[r.idx] = r.url
	}
	request.Anexos = tmp
	if err := g.Wait(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	prestadorID, err := getUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	idCatalogo, err := h.uc.Criar(c.Request.Context(), request, prestadorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Catálogo criado com sucesso", "id": idCatalogo})
}

func (h *CatalogoHandler) Editar(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de catálogo inválido"})
		return
	}

	prestadorID, err := getUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var campos map[string]interface{}
	if err := c.ShouldBindJSON(&campos); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(campos) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nenhum campo para atualizar"})
		return
	}

	delete(campos, "id")
	delete(campos, "idprestador")
	if err := h.uc.Editar(c.Request.Context(), uint(id), prestadorID, campos); err != nil {

		if err.Error() == "nao tem permissao para apagar esse catalogo" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Você não tem permissão para editar este catálogo."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *CatalogoHandler) Apagar(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de catálogo inválido"})
		return
	}

	prestadorID, err := getUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.uc.Apagar(c.Request.Context(), uint(id), prestadorID); err != nil {

		if err.Error() == "nao tem permissao para apagar esse catalogo" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Você não tem permissão para apagar este catálogo."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *CatalogoHandler) Listar(c *gin.Context) {
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

		if key == "disponivel" {
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

	catalogos, err := h.uc.Listar(c.Request.Context(), filters, orderBy, orderDir, pageSize, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      catalogos,
		"page":      page,
		"pageSize":  pageSize,
		"orderBy":   orderBy,
		"direction": orderDir,
		"filters":   filters,
	})
}

func (h *CatalogoHandler) ListarPorPrestador(c *gin.Context) {

	prestadorIDStr := c.Param("prestadorID")
	prestadorID, err := strconv.ParseUint(prestadorIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de prestador inválido"})
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
		if key == "disponivel" {
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

	catalogos, err := h.uc.ListarPorPrestador(c.Request.Context(), uint(prestadorID), filters, orderBy, orderDir, pageSize, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      catalogos,
		"page":      page,
		"pageSize":  pageSize,
		"orderBy":   orderBy,
		"direction": orderDir,
		"filters":   filters,
	})
}

func (h *CatalogoHandler) ListarPorLocalizacao(c *gin.Context) {
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

	catalogos, err := h.uc.ListarPorLocalizacao(c.Request.Context(), latitude, longitude, radius, filters, orderBy, orderDir, pageSize, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      catalogos,
		"page":      page,
		"pageSize":  pageSize,
		"orderBy":   orderBy,
		"direction": orderDir,
		"filters":   filters,
	})
}

package handler

import (
	"bytes"
	"fmt"
	"mime"
	"net/http"
	"time"

	"github.com/ManuelMassora/servicoJa-api/internal/dto"
	"github.com/ManuelMassora/servicoJa-api/internal/services"
	"github.com/ManuelMassora/servicoJa-api/internal/usecases"
	"github.com/ManuelMassora/servicoJa-api/pkg"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"golang.org/x/sync/errgroup"
)

type GaleriaHandler struct {
	uc       *usecases.GaleriaUseCase
	uploader *services.SupabaseUploader
}

func NewGaleriaHandler(uc *usecases.GaleriaUseCase, uploader *services.SupabaseUploader) *GaleriaHandler {
	return &GaleriaHandler{uc: uc, uploader: uploader}
}

func (h *GaleriaHandler) CriarGaleria(c *gin.Context) {
	var input dto.GaleriaInput
	if err := c.ShouldBindWith(&input, binding.FormMultipart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao processar formulário multipart: " + err.Error()})
		return
	}

	files := form.File["imagens"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nenhuma imagem foi enviada"})
		return
	}

	if len(files) > 4 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Limite de 4 imagens por galeria excedido."})
		return
	}

	prestadorID, err := getUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g, ctx := errgroup.WithContext(c.Request.Context())
	type result struct{ idx int; url string }
	resCh := make(chan result, len(files))

	for i, file := range files {
		i := i
		file := file
		g.Go(func() error {
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

	if err := g.Wait(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	close(resCh)

	imagens := make([]string, len(files))
	for r := range resCh {
		imagens[r.idx] = r.url
	}

	input.Imagens = imagens

	galeria, err := h.uc.AddImagesToGaleria(c.Request.Context(), prestadorID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, galeria)
}

package handler

import (
	"net/http"

	"github.com/ManuelMassora/servicoJa-api/internal/dto"
	"github.com/ManuelMassora/servicoJa-api/internal/services"
	"github.com/ManuelMassora/servicoJa-api/internal/usecases"
	"github.com/gin-gonic/gin"
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
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao processar formulário multipart: " + err.Error()})
		return
	}

	files := form.File["imagens"]
	for _, file := range files {
		_, fileName, err := h.uploader.Upload(c.Request.Context(), file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao fazer upload da imagem: " + err.Error()})
			return
		}
		publicURL := h.uploader.GetPublicURL("serviceja-image", fileName)
		input.Imagens = append(input.Imagens, dto.ImagemInput{URL: publicURL})
	}

	prestadorID, err := getUsuarioID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	input.PrestadorID = prestadorID

	galeria, err := h.uc.CreateGaleria(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, galeria)
}

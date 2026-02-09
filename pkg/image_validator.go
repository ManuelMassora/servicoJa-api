package pkg

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"mime/multipart"
	"net/http"
	"strings"
)

const (
	MaxImageSize = 5 * 1024 * 1024 // 5MB
	MaxDimension = 5000            // 5000px
)

// ValidateImage performs rigorous validation on an image file.
// It checks the file size, MIME type, and dimensions.
func ValidateImage(file *multipart.FileHeader) error {
	// 1. Check file size
	if file.Size > MaxImageSize {
		return fmt.Errorf("arquivo muito grande: o limite é de 5MB")
	}

	// 2. Open the file to check its content
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo: %w", err)
	}
	defer src.Close()

	// 3. Detect MIME type correctly (checking the first 512 bytes)
	buffer := make([]byte, 512)
	n, err := src.Read(buffer)
	if err != nil && n == 0 {
		return fmt.Errorf("erro ao ler arquivo: %w", err)
	}

	contentType := http.DetectContentType(buffer[:n])
	if !strings.HasPrefix(contentType, "image/jpeg") && !strings.HasPrefix(contentType, "image/png") {
		return fmt.Errorf("formato de arquivo inválido: apenas JPEG e PNG são permitidos")
	}

	// 4. Seek back to the beginning to decode metadata
	if _, err := src.Seek(0, 0); err != nil {
		return fmt.Errorf("erro ao processar arquivo: %w", err)
	}

	// 5. Decode image configuration to check dimensions
	config, format, err := image.DecodeConfig(src)
	if err != nil {
		return fmt.Errorf("arquivo de imagem inválido ou corrompido: %w", err)
	}

	// Double check format just in case
	if format != "jpeg" && format != "png" {
		return fmt.Errorf("formato de imagem não suportado: %s", format)
	}

	if config.Width > MaxDimension || config.Height > MaxDimension {
		return fmt.Errorf("dimensões da imagem muito grandes: o limite é %dx%d pixels", MaxDimension, MaxDimension)
	}

	return nil
}

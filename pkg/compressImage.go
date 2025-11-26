package pkg

import (
	"bytes"
	"image"
	_ "image/gif"
	"image/jpeg"
	"mime/multipart"

	// "image/png"
	"github.com/disintegration/imaging"
)

// Função que comprime a imagem para no máximo maxSizeKB (ex: 300)
func CompressImage(file *multipart.FileHeader, maxSizeKB int) (*bytes.Buffer, string, error) {
    // Abre o arquivo
    srcFile, err := file.Open()
    if err != nil {
        return nil, "", err
    }
    defer srcFile.Close()

    // Decodifica a imagem (suporta jpeg, png, gif, webp)
    srcImg, format, err := image.Decode(srcFile)
    if err != nil {
        return nil, "", err
    }

    // Buffer de saída
    var buf bytes.Buffer

    // Começa com qualidade 90 e vai reduzindo até caber no limite
    quality := 90
    for quality >= 30 {
        buf.Reset()

        // Redimensiona mantendo proporção (opcional, mas recomendado)
        // Aqui eu limito a largura máxima para 1200px (ótimo para perfil)
        resizedImg := imaging.Resize(srcImg, 1200, 0, imaging.Lanczos)

        // Codifica com a qualidade atual
        switch format {
			case "jpeg", "jpg":
				err = jpeg.Encode(&buf, resizedImg, &jpeg.Options{Quality: quality})
			case "png":
				// PNG não tem qualidade, mas você pode converter pra JPEG pra reduzir mais
				err = jpeg.Encode(&buf, resizedImg, &jpeg.Options{Quality: quality})
				format = "jpeg" // força conversão pra jpeg
			default:
				// fallback: salva como jpeg
				err = jpeg.Encode(&buf, resizedImg, &jpeg.Options{Quality: quality})
				format = "jpeg"
        }
        if err != nil {
            return nil, "", err
        }
        // Se já está abaixo do limite → ótimo!
        if buf.Len()/1024 <= maxSizeKB {
            break
        }
        // Se ainda está grande, reduz a qualidade
        quality -= 10
        if quality < 30 {
            quality = 30
        }
    }
    return &buf, format, nil
}
package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/ManuelMassora/servicoJa-api/internal/config"
	storage "github.com/supabase-community/storage-go"
	"github.com/supabase-community/supabase-go"
)

type SupabaseUploader struct {
	client *supabase.Client
	projectID string
}

func NewSupabaseUploader(cfg *config.Config) *SupabaseUploader {
	client, err := supabase.NewClient(cfg.SupabaseURL, cfg.SupabaseKey, nil)
	if err != nil {
		panic(err)
	}
	return &SupabaseUploader{client: client, projectID: "ypeauobysuqxcqagxowt"} // Adicione seu Project ID aqui
}

func (s *SupabaseUploader) Upload(ctx context.Context, file *multipart.FileHeader) (string, string, error) {
	fileContent, err := file.Open()
	if err != nil {
		return "", "", err
	}
	defer fileContent.Close()

	ext := filepath.Ext(file.Filename)
	if ext == "" {
		ext = ".jpg"
	}
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	upsert := true

	res, err := s.client.Storage.
		UploadFile("serviceja-image", fileName, fileContent, storage.FileOptions{
			ContentType: &contentType,
			Upsert: &upsert,
		})
	if err != nil {
		return "", "", fmt.Errorf("upload failed: %w", err)
	}

	return res.Key, fileName, nil
}

func (s *SupabaseUploader) UploadFromReader(ctx context.Context, fileContent io.Reader, fileName string, contentType string) (string, string, error) {
	upsert := true
	res, err := s.client.Storage.
		UploadFile("serviceja-image", fileName, fileContent, storage.FileOptions{
			ContentType: &contentType,
			Upsert: &upsert,
		})
	if err != nil {
		return "", "", fmt.Errorf("upload failed: %w", err)
	}

	return res.Key, fileName, nil
}

func (s *SupabaseUploader) GetPublicURL(bucketName, fileName string) string {
	return fmt.Sprintf("https://%s.supabase.co/storage/v1/object/public/%s/%s", s.projectID, bucketName, fileName)
}
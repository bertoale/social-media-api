package cloudinary

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type Service struct {
	cld    *cloudinary.Cloudinary
	folder string
}

func NewService(cloudName, apiKey, apiSecret, folder string) (*Service, error) {
	url := fmt.Sprintf("cloudinary://%s:%s@%s", apiKey, apiSecret, cloudName)
	cld, err := cloudinary.NewFromURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed to create cloudinary client: %w", err)
	}
	return &Service{
		cld:    cld,
		folder: folder,
	}, nil
}

func (s *Service) Upload(ctx context.Context, file multipart.File, header *multipart.FileHeader, subfolder string) (string, error) {
	folderPath := s.folder
	if subfolder != "" {
		folderPath = s.folder + "/" + subfolder
	}

	result, err := s.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: folderPath,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to cloudinary: %w", err)
	}

	return result.SecureURL, nil
}

func (s *Service) DeleteByURL(ctx context.Context, url string) error {
	if url == "" {
		return nil
	}

	if !strings.Contains(url, "res.cloudinary.com") {
		return nil
	}

	publicID := s.extractPublicID(url)
	if publicID == "" {
		return nil
	}

	_, err := s.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete from cloudinary: %w", err)
	}

	return nil
}

func (s *Service) extractPublicID(url string) string {
	idx := strings.Index(url, s.folder+"/")
	if idx == -1 {
		return ""
	}

	afterFolder := url[idx+len(s.folder)+1:]
	publicID := s.folder + "/" + strings.TrimSuffix(afterFolder, filepath.Ext(afterFolder))
	return publicID
}

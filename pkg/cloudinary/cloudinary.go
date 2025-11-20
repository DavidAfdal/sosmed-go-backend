package cloudinary

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"sync"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/davidafdal/post-app/config"
	"golang.org/x/sync/errgroup"
)

type fileUpload struct {
	Url  string
	Type string
}

type CloudinaryUseCase interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader) (string, error)
	UploadMultipelFiles(ctx context.Context, files []*multipart.FileHeader) ([]fileUpload, error)
	DeleteFile(ctx context.Context, secureURL string) error
}

type cloudinaryUseCaseImpl struct {
	cld *cloudinary.Cloudinary
}

func NewCloudinaryUseCase(cfg *config.CloudinaryConfig) (CloudinaryUseCase, error) {
	cld, err := cloudinary.NewFromURL(cfg.Url)

	if err != nil {
		return nil, err
	}

	return &cloudinaryUseCaseImpl{
		cld: cld,
	}, nil
}

func (c *cloudinaryUseCaseImpl) UploadFile(ctx context.Context, file *multipart.FileHeader) (string, error) {
	uploadResult, err := c.cld.Upload.Upload(ctx, file, uploader.UploadParams{})

	if err != nil {
		return "", nil
	}

	return uploadResult.SecureURL, nil
}

func (c *cloudinaryUseCaseImpl) UploadMultipelFiles(ctx context.Context, files []*multipart.FileHeader) ([]fileUpload, error) {
	var uploadedFiles []fileUpload
	var mu sync.Mutex
	g := new(errgroup.Group)
	g.SetLimit(5)

	for _, file := range files {
		g.Go(func() error {
			url, err := c.UploadFile(ctx, file)

			if err != nil {
				return err
			}

			fileType := strings.Split(file.Header.Get("Content-Type"), "/")[0]

			mu.Lock()
			uploadedFiles = append(uploadedFiles, fileUpload{
				Url:  url,
				Type: fileType,
			})
			mu.Unlock()

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return uploadedFiles, nil
}

func (c *cloudinaryUseCaseImpl) DeleteFile(ctx context.Context, secureURL string) error {
	publicID, err := c.extractPublicIDFromURL(secureURL)

	if err != nil {
		return err
	}

	_, err = c.cld.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: publicID})
	return err
}

func (c *cloudinaryUseCaseImpl) extractPublicIDFromURL(secureURL string) (string, error) {
	parts := strings.Split(secureURL, "/upload/")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid cloudinary URL format")
	}

	afterUpload := parts[1]

	pathParts := strings.SplitN(afterUpload, "/", 2)
	if len(pathParts) < 2 {
		return "", fmt.Errorf("invalid versioned path in URL")
	}

	pathWithExt := pathParts[1]

	ext := filepath.Ext(pathWithExt)
	publicID := strings.TrimSuffix(pathWithExt, ext)

	return publicID, nil
}

package cloudinary

import (
	"context"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/davidafdal/post-app/config"
)

type CloudinaryUseCase interface {
}

type cloudinaryUseCaseImpl struct {
	cld *cloudinary.Cloudinary
}

func NewCloudinaryUseCase(cfg config.CloudinaryConfig) (CloudinaryUseCase, error) {
	cld, err := cloudinary.NewFromURL("")

	if err != nil {
		return nil, err
	}

	return &cloudinaryUseCaseImpl{
		cld: cld,
	}, nil
}

func (c *cloudinaryUseCaseImpl) UploadImage(ctx context.Context, file multipart.File) (string, error) {
	uploadResult, err := c.cld.Upload.Upload(ctx, file, uploader.UploadParams{})

	if err != nil {
		return "", nil
	}

	return uploadResult.SecureURL, nil
}

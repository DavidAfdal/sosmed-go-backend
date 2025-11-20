package upload

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type UploadUseCase interface {
	SaveTempFile(fileHeader *multipart.FileHeader, tempDir string) (string, error)
}

type uploadUseCaseImpl struct {
}

func NewUploadUseCase() UploadUseCase {
	return &uploadUseCaseImpl{}
}

func (u *uploadUseCaseImpl) SaveTempFile(fileHeader *multipart.FileHeader, tempDir string) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	ext := filepath.Ext(fileHeader.Filename)
	filename := uuid.New().String() + ext
	tempPath := filepath.Join(tempDir, filename)

	dst, err := os.Create(tempPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	return tempPath, nil
}

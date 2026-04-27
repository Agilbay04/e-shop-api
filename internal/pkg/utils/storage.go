package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
	"github.com/google/uuid"
)

type FileUpload struct {
	AllowedExtensions []string
	MaxFileSize       int64
	UploadDir         string
}

type FileOption func(*FileUpload)

func WithExtensions(exts []string) FileOption {
	return func(f *FileUpload) {
		f.AllowedExtensions = exts
	}
}

func WithMaxSize(megabytes int64) FileOption {
	return func(f *FileUpload) {
		f.MaxFileSize = megabytes * 1024 * 1024
	}
}

func WithDirectory(dir string) FileOption {
	return func(f *FileUpload) {
		f.UploadDir = dir
	}
}

func NewFileUploader(opts ...FileOption) *FileUpload {
	uploader := &FileUpload{
		AllowedExtensions: []string{".jpg", ".jpeg", ".png"},
		MaxFileSize:       2 * 1024 * 1024, // 2MB
		UploadDir:         "uploads/others",
	}

	for _, opt := range opts {
		opt(uploader)
	}

	return uploader
}

func MakeDir(dir string) (string, error) {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return dir, err
	}

	return dir, nil
}

func (u *FileUpload) UploadFile(file *multipart.FileHeader) (string, error) {
	if _, err := MakeDir(u.UploadDir); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %v", err)
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	newFileName := fmt.Sprintf("%s-%d%s", uuid.New().String(), time.Now().Unix(), filepath.Ext(file.Filename))
	dst := filepath.Join(u.UploadDir, newFileName)

	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return dst, err
}

func (u *FileUpload) DeleteFile(filePath string) (string, error) {
	err := os.Remove(filePath)
	if err != nil {
		return filePath, err
	}

	return filePath, nil
}


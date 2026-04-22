package util

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

func (u *FileUpload) UploadFile(file *multipart.FileHeader) (string, error) {
	if err := os.MkdirAll(u.UploadDir, 0755); err != nil {
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

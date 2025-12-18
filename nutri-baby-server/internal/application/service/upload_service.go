package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/wxlbd/nutri-baby-server/internal/infrastructure/config"
	"github.com/wxlbd/nutri-baby-server/pkg/errors"
)

// UploadService 文件上传服务
type UploadService struct {
	cfg *config.Config
}

// UploadType 上传类型
type UploadType string

const (
	UploadTypeUserAvatar UploadType = "user_avatar"
	UploadTypeBabyAvatar UploadType = "baby_avatar"
)

// UploadResult 上传结果
type UploadResult struct {
	URL      string `json:"url"`      // 完整访问URL
	Path     string `json:"path"`     // 相对路径
	Filename string `json:"filename"` // 文件名
	Size     int64  `json:"size"`     // 文件大小
}

// NewUploadService 创建上传服务
func NewUploadService(cfg *config.Config) *UploadService {
	return &UploadService{
		cfg: cfg,
	}
}

// UploadFile 上传文件
func (s *UploadService) UploadFile(ctx context.Context, fileHeader *multipart.FileHeader, uploadType UploadType, relatedID string) (*UploadResult, error) {
	if fileHeader == nil {
		return nil, errors.New(errors.ParamError, "File cannot be empty")
	}

	// MIME type to extension mapping
	mimeToExt := map[string]string{
		"image/jpeg": ".jpg",
		"image/jpg":  ".jpg",
		"image/png":  ".png",
		"image/gif":  ".gif",
	}

	// Validate file type
	allowedTypes := s.cfg.Upload.AllowedTypes
	fileExt := strings.ToLower(filepath.Ext(fileHeader.Filename))
	isAllowed := false
	allowedExts := []string{}

	// Collect all allowed extensions
	for _, mimeType := range allowedTypes {
		if ext, ok := mimeToExt[strings.ToLower(mimeType)]; ok {
			allowedExts = append(allowedExts, ext)
		}
	}

	// Validate file extension
	for _, ext := range allowedExts {
		if fileExt == strings.ToLower(ext) {
			isAllowed = true
			break
		}
	}

	if !isAllowed {
		return nil, errors.New(errors.ParamError, "Unsupported file type. Allowed types: "+strings.Join(allowedExts, ", "))
	}

	// Validate file size
	if fileHeader.Size > s.cfg.Upload.MaxSize {
		return nil, errors.New(errors.ParamError, "File size exceeds limit")
	}

	// Get subdirectory
	subDir := s.getSubDir(uploadType)
	if subDir == "" {
		return nil, errors.New(errors.ParamError, "Unsupported upload type")
	}

	// Generate filename
	filename := s.generateFilename(fileHeader.Filename, uploadType, relatedID)

	// Create save path
	storagePath := s.cfg.Upload.StoragePath
	dirPath := filepath.Join(storagePath, "images", subDir)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return nil, errors.Wrap(errors.InternalError, "Failed to create directory", err)
	}

	// Save file
	filePath := filepath.Join(dirPath, filename)
	if err := s.saveFile(fileHeader, filePath); err != nil {
		return nil, err
	}

	// Generate access URL
	relPath := filepath.Join("uploads", "images", subDir, filename)
	url := s.cfg.Server.BaseURL + "/" + filepath.ToSlash(relPath)

	return &UploadResult{
		URL:      url,
		Path:     "/" + relPath,
		Filename: filename,
		Size:     fileHeader.Size,
	}, nil
}

// getSubDir Get subdirectory by upload type
func (s *UploadService) getSubDir(uploadType UploadType) string {
	switch uploadType {
	case UploadTypeUserAvatar:
		return "users"
	case UploadTypeBabyAvatar:
		return "babies"
	default:
		return ""
	}
}

// generateFilename Generate unique filename
func (s *UploadService) generateFilename(originalFilename string, uploadType UploadType, relatedID string) string {
	// Get file extension
	ext := filepath.Ext(originalFilename)
	if ext == "" {
		ext = ".jpg" // Default extension
	}

	// Build prefix
	prefix := string(uploadType)
	if relatedID != "" {
		prefix = prefix + "_" + relatedID
	}

	// Add timestamp
	timestamp := time.Now().Format("20060102_150405")

	// Generate unique filename: {prefix}_{timestamp}{ext}
	return fmt.Sprintf("%s_%s%s", prefix, timestamp, ext)
}

// saveFile Save file to disk
func (s *UploadService) saveFile(fileHeader *multipart.FileHeader, filePath string) error {
	srcFile, err := fileHeader.Open()
	if err != nil {
		return errors.Wrap(errors.InternalError, "Failed to open file", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(filePath)
	if err != nil {
		return errors.Wrap(errors.InternalError, "Failed to create file", err)
	}
	defer dstFile.Close()

	// Copy file content
	if _, err := dstFile.ReadFrom(srcFile); err != nil {
		return errors.Wrap(errors.InternalError, "Failed to save file", err)
	}

	return nil
}

package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/wxlbd/nutri-baby-server/internal/infrastructure/config"
	"github.com/wxlbd/nutri-baby-server/pkg/errors"
)

// UploadService 文件上传服务
type UploadService struct {
	cfg       *config.Config
	cosClient *cos.Client
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
	Path     string `json:"path"`     // relative path or key
	Filename string `json:"filename"` // 文件名
	Size     int64  `json:"size"`     // 文件大小
}

// NewUploadService 创建上传服务
func NewUploadService(cfg *config.Config) *UploadService {
	s := &UploadService{
		cfg: cfg,
	}

	// 初始化 COS Client
	if cfg.COS.BucketURL != "" && cfg.COS.SecretID != "" && cfg.COS.SecretKey != "" {
		u, _ := url.Parse(cfg.COS.BucketURL)
		b := &cos.BaseURL{BucketURL: u}
		// 1. Permanent key
		s.cosClient = cos.NewClient(b, &http.Client{
			Transport: &cos.AuthorizationTransport{
				SecretID:  cfg.COS.SecretID,
				SecretKey: cfg.COS.SecretKey,
			},
		})
		fmt.Printf("Info: Upload Service initialized with COS driver (Bucket: %s)\n", cfg.COS.BucketURL)
	} else {
		fmt.Printf("Info: Upload Service initialized with Local driver (Path: %s)\n", cfg.Upload.StoragePath)
	}

	return s
}

// UploadFile 上传文件
func (s *UploadService) UploadFile(ctx context.Context, fileHeader *multipart.FileHeader, uploadType UploadType, relatedID string) (*UploadResult, error) {
	if fileHeader == nil {
		return nil, errors.New(errors.ParamError, "File cannot be empty")
	}

	// 1. 验证文件类型和大小 (保留原有逻辑)
	if err := s.validateFile(fileHeader); err != nil {
		return nil, err
	}

	// 2. 生成文件名和路径
	subDir := s.getSubDir(uploadType)
	if subDir == "" {
		return nil, errors.New(errors.ParamError, "Unsupported upload type")
	}
	filename := s.generateFilename(fileHeader.Filename, uploadType, relatedID)
	// COS 对象键 (Key): images/users/xxx.jpg
	objectKey := fmt.Sprintf("images/%s/%s", subDir, filename)

	// 3. 执行上传
	if s.cosClient != nil {
		return s.uploadToCOS(ctx, fileHeader, objectKey, filename)
	}

	// 4. 降级到本地上传 (只要未配置COS)
	return s.uploadToLocal(fileHeader, subDir, filename)
}

// uploadToCOS 上传到腾讯云COS
func (s *UploadService) uploadToCOS(ctx context.Context, fileHeader *multipart.FileHeader, objectKey string, filename string) (*UploadResult, error) {
	srcFile, err := fileHeader.Open()
	if err != nil {
		return nil, errors.Wrap(errors.InternalError, "Failed to open file", err)
	}
	defer srcFile.Close()

	// 上传
	_, err = s.cosClient.Object.Put(ctx, objectKey, srcFile, nil)
	if err != nil {
		return nil, errors.Wrap(errors.InternalError, "Failed to upload to COS", err)
	}

	// 构建返回 URL
	// 如果 BucketURL 不带 scheme，需要补全? 一般配置里应该带 https://
	fullURL := s.cfg.COS.BucketURL
	if !strings.HasSuffix(fullURL, "/") {
		fullURL += "/"
	}
	fullURL += objectKey

	return &UploadResult{
		URL:      fullURL,
		Path:     objectKey,
		Filename: filename,
		Size:     fileHeader.Size,
	}, nil
}

// uploadToLocal 上传到本地磁盘
func (s *UploadService) uploadToLocal(fileHeader *multipart.FileHeader, subDir string, filename string) (*UploadResult, error) {
	storagePath := s.cfg.Upload.StoragePath // e.g., "uploads/"

	// 确保目录存在
	dirPath := filepath.Join(storagePath, "images", subDir)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return nil, errors.Wrap(errors.InternalError, "Failed to create directory", err)
	}

	filePath := filepath.Join(dirPath, filename)
	srcFile, err := fileHeader.Open()
	if err != nil {
		return nil, errors.Wrap(errors.InternalError, "Failed to open source file", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(filePath)
	if err != nil {
		return nil, errors.Wrap(errors.InternalError, "Failed to create destination file", err)
	}
	defer dstFile.Close()

	if _, err := dstFile.ReadFrom(srcFile); err != nil {
		return nil, errors.Wrap(errors.InternalError, "Failed to save file content", err)
	}

	// 关键修复：URL 应当基于 /uploads/ 开头，而 Path 作为本地存储路径
	// 假设 r.Static("/uploads", "./uploads")，且 storagePath="uploads/"
	// 那么 relPathFromWebRoot 是 "images/users/xxx.jpg"
	relPath := filepath.ToSlash(filepath.Join("images", subDir, filename))
	url := fmt.Sprintf("%s/uploads/%s", strings.TrimSuffix(s.cfg.Server.BaseURL, "/"), relPath)

	return &UploadResult{
		URL:      url,
		Path:     filepath.ToSlash(filepath.Join(storagePath, "images", subDir, filename)),
		Filename: filename,
		Size:     fileHeader.Size,
	}, nil
}

func (s *UploadService) validateFile(fileHeader *multipart.FileHeader) error {
	mimeToExt := map[string]string{
		"image/jpeg": ".jpg",
		"image/jpg":  ".jpg",
		"image/png":  ".png",
		"image/gif":  ".gif",
	}

	allowedTypes := s.cfg.Upload.AllowedTypes
	fileExt := strings.ToLower(filepath.Ext(fileHeader.Filename))
	isAllowed := false
	allowedExts := []string{}

	for _, mimeType := range allowedTypes {
		if ext, ok := mimeToExt[strings.ToLower(mimeType)]; ok {
			allowedExts = append(allowedExts, ext)
		}
	}

	for _, ext := range allowedExts {
		if fileExt == strings.ToLower(ext) {
			isAllowed = true
			break
		}
	}

	if !isAllowed {
		return errors.New(errors.ParamError, "Unsupported file type")
	}

	if fileHeader.Size > s.cfg.Upload.MaxSize {
		return errors.New(errors.ParamError, "File size exceeds limit")
	}
	return nil
}

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

func (s *UploadService) generateFilename(originalFilename string, uploadType UploadType, relatedID string) string {
	ext := filepath.Ext(originalFilename)
	if ext == "" {
		ext = ".jpg"
	}
	prefix := string(uploadType)
	if relatedID != "" {
		prefix = prefix + "_" + relatedID
	}
	timestamp := time.Now().Format("20060102_150405")
	return fmt.Sprintf("%s_%s%s", prefix, timestamp, ext)
}

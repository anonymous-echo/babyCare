package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/wxlbd/nutri-baby-server/internal/application/service"
	"github.com/wxlbd/nutri-baby-server/pkg/errors"
	"github.com/wxlbd/nutri-baby-server/pkg/response"
)

// UploadHandler 文件上传处理器
type UploadHandler struct {
	uploadService *service.UploadService
}

// NewUploadHandler 创建上传处理器
func NewUploadHandler(uploadService *service.UploadService) *UploadHandler {
	return &UploadHandler{uploadService: uploadService}
}

// Upload 上传文件
// @Router /upload [post]
func (h *UploadHandler) Upload(c *gin.Context) {
	// Get upload type
	uploadType := service.UploadType(c.PostForm("type"))
	if uploadType == "" {
		response.ErrorWithMessage(c, errors.ParamError, "Missing upload type parameter")
		return
	}

	// Get related ID (optional)
	relatedID := c.PostForm("related_id")

	// Get file
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.ErrorWithMessage(c, errors.ParamError, "Please select a file to upload")
		return
	}

	// Validate upload type
	if uploadType != service.UploadTypeUserAvatar && uploadType != service.UploadTypeBabyAvatar {
		response.ErrorWithMessage(c, errors.ParamError, "Unsupported upload type")
		return
	}

	// Upload file
	result, err := h.uploadService.UploadFile(c.Request.Context(), fileHeader, uploadType, relatedID)
	if err != nil {
		response.Error(c, err)
		return
	}

	// Return success response
	response.Success(c, gin.H{
		"url":      result.URL,
		"path":     result.Path,
		"filename": result.Filename,
		"size":     result.Size,
	})
}

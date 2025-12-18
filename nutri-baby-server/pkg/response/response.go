package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wxlbd/nutri-baby-server/pkg/errors"
)

// Response 统一响应结构
type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:      int(errors.Success),
		Message:   "success",
		Data:      data,
		Timestamp: time.Now().Unix(),
	})
}

// Error 错误响应
func Error(c *gin.Context, err error) {
	appErr, ok := err.(*errors.AppError)
	if !ok {
		// 未知错误,返回内部错误
		c.JSON(http.StatusInternalServerError, Response{
			Code:      int(errors.InternalError),
			Message:   "服务器内部错误",
			Timestamp: time.Now().Unix(),
		})
		return
	}

	// 根据错误码确定HTTP状态码
	httpStatus := getHTTPStatus(appErr.Code)

	c.JSON(httpStatus, Response{
		Code:      int(appErr.Code),
		Message:   appErr.Message,
		Timestamp: time.Now().Unix(),
	})
}

// ErrorWithMessage 自定义错误消息
func ErrorWithMessage(c *gin.Context, code errors.ErrorCode, message string) {
	httpStatus := getHTTPStatus(code)

	c.JSON(httpStatus, Response{
		Code:      int(code),
		Message:   message,
		Timestamp: time.Now().Unix(),
	})
}

// getHTTPStatus 根据错误码获取HTTP状态码
func getHTTPStatus(code errors.ErrorCode) int {
	switch code {
	case errors.Success:
		return http.StatusOK
	case errors.ParamError:
		return http.StatusBadRequest
	case errors.Unauthorized, errors.InvalidToken, errors.TokenExpired:
		return http.StatusUnauthorized
	case errors.NotFound, errors.UserNotFound, errors.BabyNotFound,
		errors.FamilyNotFound, errors.RecordNotFound:
		return http.StatusNotFound
	case errors.Conflict:
		return http.StatusConflict
	case errors.PermissionDenied:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

// Paginated 分页响应
type Paginated struct {
	Records  interface{} `json:"records"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

// SuccessPaginated 分页成功响应
func SuccessPaginated(c *gin.Context, records interface{}, total int64, page, pageSize int) {
	Success(c, Paginated{
		Records:  records,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

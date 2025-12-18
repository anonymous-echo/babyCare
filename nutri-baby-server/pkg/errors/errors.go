package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

var (
	WithStack = errors.WithStack
	Wrapf     = errors.Wrapf
	Is        = errors.Is
	Errorf    = errors.Errorf
)

// ErrorCode 错误码
type ErrorCode int

const (
	// 成功
	Success ErrorCode = 0

	// 通用错误 1000-1999
	ParamError       ErrorCode = 1001
	Unauthorized     ErrorCode = 1002
	NotFound         ErrorCode = 1003
	Conflict         ErrorCode = 1004
	PermissionDenied ErrorCode = 1005

	// 服务器错误 2000-2999
	InternalError ErrorCode = 2001
	DatabaseError ErrorCode = 2002
	CacheError    ErrorCode = 2003

	// 业务错误 3000-3999
	UserNotFound      ErrorCode = 3001
	InvalidToken      ErrorCode = 3002
	TokenExpired      ErrorCode = 3003
	BabyNotFound      ErrorCode = 3004
	FamilyNotFound    ErrorCode = 3005
	InvalidInvitation ErrorCode = 3006
	RecordNotFound    ErrorCode = 3007
)

// AppError 应用错误
type AppError struct {
	Code    ErrorCode
	Message string
	Err     error
}

// Error 实现error接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// New 创建新错误
func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Wrap 包装错误
func Wrap(code ErrorCode, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// 预定义错误
var (
	ErrParamInvalid      = New(ParamError, "参数错误")
	ErrUnauthorized      = New(Unauthorized, "未授权")
	ErrNotFound          = New(NotFound, "资源不存在")
	ErrConflict          = New(Conflict, "数据冲突")
	ErrPermissionDenied  = New(PermissionDenied, "权限不足")
	ErrInternal          = New(InternalError, "服务器内部错误")
	ErrDatabase          = New(DatabaseError, "数据库错误")
	ErrUserNotFound      = New(UserNotFound, "用户不存在")
	ErrInvalidToken      = New(InvalidToken, "无效的令牌")
	ErrTokenExpired      = New(TokenExpired, "令牌已过期")
	ErrBabyNotFound      = New(BabyNotFound, "宝宝不存在")
	ErrFamilyNotFound    = New(FamilyNotFound, "家庭不存在")
	ErrInvalidInvitation = New(InvalidInvitation, "邀请码无效或已过期")
	ErrRecordNotFound    = New(RecordNotFound, "记录不存在")
)

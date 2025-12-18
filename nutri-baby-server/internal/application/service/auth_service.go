package service

import (
	"context"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/wxlbd/nutri-baby-server/internal/application/dto"
	"github.com/wxlbd/nutri-baby-server/internal/domain/entity"
	"github.com/wxlbd/nutri-baby-server/internal/domain/repository"
	"github.com/wxlbd/nutri-baby-server/internal/infrastructure/config"
	"github.com/wxlbd/nutri-baby-server/internal/infrastructure/wechat"
	"github.com/wxlbd/nutri-baby-server/pkg/errors"
)

// AuthService 认证服务 (去家庭化架构)
type AuthService struct {
	userRepo     repository.UserRepository
	cfg          *config.Config
	wechatClient *wechat.Client
}

// NewAuthService 创建认证服务
func NewAuthService(
	userRepo repository.UserRepository,
	cfg *config.Config,
	wechatClient *wechat.Client,
) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		cfg:          cfg,
		wechatClient: wechatClient,
	}
}

// WechatLogin 微信小程序登录 (去家庭化架构)
func (s *AuthService) WechatLogin(ctx context.Context, req *dto.WechatLoginRequest) (*dto.LoginResponse, error) {
	// 使用 SDK 调用微信API获取openid
	miniProgram := s.wechatClient.GetMiniProgram()
	auth := miniProgram.GetAuth()

	session, err := auth.Code2SessionContext(ctx, req.Code)
	if err != nil {
		return nil, errors.Wrap(errors.InternalError, "调用微信API失败", err)
	}

	if session.ErrCode != 0 {
		return nil, errors.New(errors.Unauthorized, "微信登录失败: "+session.ErrMsg)
	}

	// 查找或创建用户
	user, err := s.userRepo.FindByOpenID(ctx, session.OpenID)
	if err != nil && !errors.Is(err, errors.ErrUserNotFound) {
		return nil, err
	}

	now := time.Now().UnixMilli()
	isNewUser := false

	if user == nil {
		// 创建新用户
		user = &entity.User{
			OpenID:        session.OpenID,
			NickName:      req.NickName,
			AvatarURL:     req.AvatarURL,
			LastLoginTime: now,
		}

		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, err
		}

		isNewUser = true
	} else {
		// 更新用户信息
		user.NickName = req.NickName
		user.AvatarURL = req.AvatarURL
		user.LastLoginTime = now

		if err := s.userRepo.Update(ctx, user); err != nil {
			return nil, err
		}
	}

	// 生成Token
	token, err := s.generateToken(user.OpenID)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
		UserInfo: dto.UserInfoDTO{
			OpenID:        user.OpenID,
			NickName:      user.NickName,
			AvatarURL:     user.AvatarURL,
			DefaultBabyID: strconv.FormatInt(user.DefaultBabyID, 10),
		},
		IsNewUser: isNewUser, // 前端根据此字段判断是否需要引导创建宝宝
	}, nil
}

// RefreshToken 刷新Token
func (s *AuthService) RefreshToken(ctx context.Context, openID string) (*dto.RefreshTokenResponse, error) {
	// 验证用户存在
	user, err := s.userRepo.FindByOpenID(ctx, openID)
	if err != nil {
		return nil, err
	}

	// 生成新Token
	token, err := s.generateToken(user.OpenID)
	if err != nil {
		return nil, err
	}

	return &dto.RefreshTokenResponse{
		Token: token,
	}, nil
}

// GetUserInfo 获取用户信息
func (s *AuthService) GetUserInfo(ctx context.Context, openID string) (*dto.UserInfoDTO, error) {
	user, err := s.userRepo.FindByOpenID(ctx, openID)
	if err != nil {
		return nil, err
	}

	return &dto.UserInfoDTO{
		OpenID:        user.OpenID,
		NickName:      user.NickName,
		AvatarURL:     user.AvatarURL,
		DefaultBabyID: strconv.FormatInt(user.DefaultBabyID, 10),
		CreateTime:    user.CreatedAt,
		LastLoginTime: user.LastLoginTime,
	}, nil
}

// SetDefaultBaby 设置默认宝宝
func (s *AuthService) SetDefaultBaby(ctx context.Context, openID string, req *dto.SetDefaultBabyRequest) error {
	// 验证用户存在
	_, err := s.userRepo.FindByOpenID(ctx, openID)
	if err != nil {
		return err
	}

	// 转换BabyID from string to int64
	babyID, err := strconv.ParseInt(req.BabyID, 10, 64)
	if err != nil {
		return errors.New(errors.ParamError, "invalid baby id format")
	}

	// 更新默认宝宝ID
	if err := s.userRepo.UpdateDefaultBabyID(ctx, openID, babyID); err != nil {
		return err
	}

	return nil
}

// UpdateUserInfo 更新用户信息
func (s *AuthService) UpdateUserInfo(ctx context.Context, openID string, req *dto.UpdateUserInfoRequest) (*dto.UserInfoDTO, error) {
	// 查找用户
	user, err := s.userRepo.FindByOpenID(ctx, openID)
	if err != nil {
		return nil, err
	}

	// 更新用户信息
	user.NickName = req.NickName
	user.AvatarURL = req.AvatarURL

	// 保存到数据库
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	// 返回更新后的用户信息
	return &dto.UserInfoDTO{
		OpenID:        user.OpenID,
		NickName:      user.NickName,
		AvatarURL:     user.AvatarURL,
		DefaultBabyID: strconv.FormatInt(user.DefaultBabyID, 10),
		CreateTime:    user.CreatedAt,
		LastLoginTime: user.LastLoginTime,
	}, nil
}

// generateToken 生成JWT Token
func (s *AuthService) generateToken(openID string) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   openID,
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * time.Duration(s.cfg.JWT.ExpireHours))),
		IssuedAt:  jwt.NewNumericDate(now),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.cfg.JWT.Secret))
	if err != nil {
		return "", errors.Wrap(errors.InternalError, "生成Token失败", err)
	}

	return tokenString, nil
}

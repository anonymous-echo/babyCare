package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"log"`
	Upload   UploadConfig   `mapstructure:"upload"`
	Wechat   WechatConfig   `mapstructure:"wechat"`
	AI       AIConfig       `mapstructure:"ai"` // AI配置
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         int    `mapstructure:"port"`
	Mode         string `mapstructure:"mode"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
	BaseURL      string `mapstructure:"base_url"` // 服务器基础 URL，用于生成完整资源访问地址
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host            string `mapstructure:"host"` // 主库地址
	Port            int    `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	DBName          string `mapstructure:"dbname"`
	SSLMode         string `mapstructure:"sslmode"`
	Timezone        string `mapstructure:"timezone"` // 时区配置，如 Asia/Shanghai
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
	// 读副本配置（可选）
	ReadReplicaHosts  []string `mapstructure:"read_replica_hosts"`  // 只读副本地址列表
	ReadReplicaPort   int      `mapstructure:"read_replica_port"`   // 只读副本端口
	EnableReadReplica bool     `mapstructure:"enable_read_replica"` // 是否启用读副本
}

// DSN 返回PostgreSQL连接字符串
func (d DatabaseConfig) DSN() string {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.DBName, d.SSLMode)
	if d.Timezone != "" {
		dsn += " timezone=" + d.Timezone
	}
	return dsn
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// Addr 返回Redis地址
func (r RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

// UploadConfig 上传配置
type UploadConfig struct {
	MaxSize      int64    `mapstructure:"max_size"`
	AllowedTypes []string `mapstructure:"allowed_types"`
	StoragePath  string   `mapstructure:"storage_path"`
}

// WechatConfig 微信配置
type WechatConfig struct {
	AppID              string            `mapstructure:"app_id"`
	AppSecret          string            `mapstructure:"app_secret"`
	SubscribeTemplates map[string]string `mapstructure:"subscribe_templates"` // 订阅消息模板映射: templateType -> templateID
}

// AIConfig AI配置
type AIConfig struct {
	Provider string         `mapstructure:"provider"`
	OpenAI   OpenAIConfig   `mapstructure:"openai"`
	Claude   ClaudeConfig   `mapstructure:"claude"`
	ERNIE    ERNIEConfig    `mapstructure:"ernie"`
	DeepSeek DeepSeekConfig `mapstructure:"deepSeek"`
	Analysis AnalysisConfig `mapstructure:"analysis"`
	Gemini   GeminiConfig   `mapstructure:"gemini"`
}

type GeminiConfig struct {
	APIKey  string `mapstructure:"api_key"`
	BaseURL string `mapstructure:"base_url"`
	Model   string `mapstructure:"model"`
}

// OpenAIConfig OpenAI配置
type OpenAIConfig struct {
	APIKey      string  `mapstructure:"api_key"`
	BaseURL     string  `mapstructure:"base_url"`
	Model       string  `mapstructure:"model"`
	MaxTokens   int     `mapstructure:"max_tokens"`
	Temperature float64 `mapstructure:"temperature"`
}

// ClaudeConfig Claude配置
type ClaudeConfig struct {
	APIKey      string  `mapstructure:"api_key"`
	BaseURL     string  `mapstructure:"base_url"`
	Model       string  `mapstructure:"model"`
	MaxTokens   int     `mapstructure:"max_tokens"`
	Temperature float64 `mapstructure:"temperature"`
}

// ERNIEConfig ERNIE配置
type ERNIEConfig struct {
	APIKey    string `mapstructure:"api_key"`
	SecretKey string `mapstructure:"secret_key"`
	BaseURL   string `mapstructure:"base_url"`
	Model     string `mapstructure:"model"`
}

// DeepSeekConfig Gemini配置
type DeepSeekConfig struct {
	APIKey  string `mapstructure:"api_key"`
	BaseURL string `mapstructure:"base_url"`
	Model   string `mapstructure:"model"`
}

// AnalysisConfig 分析配置
type AnalysisConfig struct {
	Timeout    int               `mapstructure:"timeout"`
	RetryCount int               `mapstructure:"retry_count"`
	BatchSize  int               `mapstructure:"batch_size"`
	CacheTTL   int               `mapstructure:"cache_ttl"`
	Prompts    map[string]string `mapstructure:"prompts"`
}

// Load 加载配置
func Load(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// 读取环境变量
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// GetDefaultConfig 获取默认配置
func GetDefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         8080,
			Mode:         "debug",
			ReadTimeout:  30,
			WriteTimeout: 30,
			BaseURL:      "http://localhost:8080",
		},
		Database: DatabaseConfig{
			Host:              "localhost",
			Port:              5432,
			User:              "postgres",
			Password:          "",
			DBName:            "nutri_baby",
			SSLMode:           "disable",
			Timezone:          "Asia/Shanghai", // 默认时区为中国北京时间
			MaxOpenConns:      100,
			MaxIdleConns:      10,
			ConnMaxLifetime:   3600,
			ReadReplicaHosts:  []string{},
			ReadReplicaPort:   5432,
			EnableReadReplica: false,
		},
		Redis: RedisConfig{
			Host:     "localhost",
			Port:     6379,
			Password: "",
			DB:       0,
			PoolSize: 100,
		},
		JWT: JWTConfig{
			Secret:      "your-secret-key",
			ExpireHours: 72,
		},
		Log: LogConfig{
			Level:      "info",
			Filename:   "logs/app.log",
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     7,
			Compress:   true,
		},
		Upload: UploadConfig{
			MaxSize:      10 * 1024 * 1024, // 10MB
			AllowedTypes: []string{"image/jpeg", "image/png", "image/gif"},
			StoragePath:  "uploads/",
		},
		Wechat: WechatConfig{
			AppID:              "",
			AppSecret:          "",
			SubscribeTemplates: map[string]string{},
		},
		AI: GetDefaultAIConfig(),
	}
}

// GetDefaultAIConfig 获取默认AI配置
func GetDefaultAIConfig() AIConfig {
	return AIConfig{
		Provider: "mock", // 默认使用mock模式，便于开发测试
		OpenAI: OpenAIConfig{
			APIKey:      "",
			BaseURL:     "https://api.openai.com/v1",
			Model:       "gpt-4",
			MaxTokens:   2000,
			Temperature: 0.7,
		},
		Claude: ClaudeConfig{
			APIKey:      "",
			BaseURL:     "https://api.anthropic.com",
			Model:       "claude-3-sonnet-20240229",
			MaxTokens:   2000,
			Temperature: 0.7,
		},
		ERNIE: ERNIEConfig{
			APIKey:    "",
			SecretKey: "",
			BaseURL:   "https://aip.baidubce.com",
			Model:     "ernie-3.5",
		},
		Analysis: AnalysisConfig{
			Timeout:    30,
			RetryCount: 3,
			BatchSize:  10,
			CacheTTL:   3600,
			Prompts: map[string]string{
				"feeding":  "分析以下宝宝的喂养数据，提供专业的营养建议：",
				"sleep":    "分析以下宝宝的睡眠数据，提供改善建议：",
				"growth":   "分析以下宝宝的成长数据，评估发育状况：",
				"health":   "综合分析以下宝宝的健康数据：",
				"behavior": "分析以下宝宝的行为模式：",
			},
		},
	}
}

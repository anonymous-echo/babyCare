package config

import (
	"fmt"
	"net/url"
	"os"
	"reflect"
	"strings"

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
	COS      COSConfig      `mapstructure:"cos"` // COS配置
	AI       AIConfig       `mapstructure:"ai"`  // AI配置
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

// DSN 返回PostgreSQL连接字符串 (使用 URL 格式以更好地支持连接池和特殊字符)
func (d DatabaseConfig) DSN() string {
	// 针对 Supabase Pooler，使用 URL 格式且必须对用户名和密码进行转义
	// 同时增加 prepareThreshold=0 禁用预处理语句，因为 Transaction Pooler 不支持它
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(d.User, d.Password),
		Host:   fmt.Sprintf("%s:%d", d.Host, d.Port),
		Path:   "/" + d.DBName,
	}

	q := u.Query()
	q.Set("sslmode", d.SSLMode)
	q.Set("prepareThreshold", "0")
	// 强制使用简单查询协议，解决 Supabase Transaction Pooler 不支持扩展协议的问题 (Code 2002)
	q.Set("default_query_exec_mode", "simple_protocol")
	if d.Timezone != "" {
		q.Set("timezone", d.Timezone)
	}
	u.RawQuery = q.Encode()

	return u.String()
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
	UseTLS   bool   `mapstructure:"use_tls"`
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

// COSConfig 腾讯云COS配置
type COSConfig struct {
	BucketURL string `mapstructure:"bucket_url"`
	SecretID  string `mapstructure:"secret_id"`
	SecretKey string `mapstructure:"secret_key"`
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
	Doubao   DoubaoConfig   `mapstructure:"doubao"`
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

// DoubaoConfig 豆包配置
type DoubaoConfig struct {
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
	// 1. 获取默认配置作为基础结构和默认值
	defaultCfg := GetDefaultConfig()

	// 2. 设置 Viper 基础规则 (使用独立实例避免全局干扰)
	v := viper.New()
	v.SetEnvPrefix("") // 移除 NB_ 前缀，回归标准环境变量名以匹配用户现有云端配置
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	v.SetConfigType("yaml")

	// 3. 核心修复：深度递归绑定所有已知的配置项
	// 这确保了即便不提供 YAML，Viper 也能识别出 DATABASE_HOST 等嵌套环境变量
	bindAllEnvVars(v, defaultCfg, "")

	// 4. 加载配置文件（优先级：指定路径 > .example 模板）
	if configPath != "" {
		v.SetConfigFile(configPath)
		if err := v.ReadInConfig(); err != nil {
			// 如果主配置加载失败，尝试加载 .example 作为结构参考
			examplePath := configPath + ".example"
			fmt.Printf("Warning: Primary config %s failed: %v. Trying example %s\n", configPath, err, examplePath)

			v.SetConfigFile(examplePath)
			if err := v.ReadInConfig(); err != nil {
				// 如果两个都失败，且没有环境变量，那也没办法，只能返回错误或者依赖 safe defaults
			} else {
				//fmt.Printf("Info: Successfully loaded structure from %s\n", examplePath)
			}
		} else {
			//fmt.Printf("Info: Successfully loaded config from %s\n", configPath)
		}
	}

	// 调试：在 Unmarshal 前打印识别到的值
	fmt.Printf("Viper Debug (Pre-Unmarshal) - Database Host: %s\n", v.GetString("database.host"))
	fmt.Printf("Viper Debug (Pre-Unmarshal) - Database User: %s\n", v.GetString("database.user"))
	fmt.Printf("Viper Debug (Pre-Unmarshal) - Database DBNAME: %s\n", v.GetString("database.dbname"))
	fmt.Printf("Viper Debug (Pre-Unmarshal) - Database Port: %s\n", v.GetString("database.port"))
	fmt.Printf("Viper Debug (Pre-Unmarshal) - Port Env (Direct): %s\n", os.Getenv("DATABASE_PORT"))
	fmt.Printf("Viper Debug (Pre-Unmarshal) - SSLMode: %s\n", v.GetString("database.sslmode"))
	fmt.Printf("Viper Debug (Pre-Unmarshal) - Timezone: %s\n", v.GetString("database.timezone"))
	fmt.Printf("Viper Debug (Pre-Unmarshal) - Redis UseTLS: %v\n", v.GetBool("redis.use_tls"))

	// 5. 将所有来源合并到 defaultCfg 结构体中
	if err := v.Unmarshal(defaultCfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 调试：在 Unmarshal 后打印结构体中的实际值
	fmt.Printf("Viper Debug (After Unmarshal) - Struct Host: %s\n", defaultCfg.Database.Host)
	fmt.Printf("Viper Debug (After Unmarshal) - Struct User: %s\n", defaultCfg.Database.User)
	fmt.Printf("Viper Debug (After Unmarshal) - Struct DBName: %s\n", defaultCfg.Database.DBName)
	fmt.Printf("Viper Debug (After Unmarshal) - Struct SSLMode: %s\n", defaultCfg.Database.SSLMode)

	if defaultCfg.Database.Password == "" {
		fmt.Printf("Warning: DATABASE_PASSWORD is empty!\n")
	} else {
		fmt.Printf("Info: DATABASE_PASSWORD is set (length: %d)\n", len(defaultCfg.Database.Password))
	}

	// 打印脱敏后的 DSN 供核对格式 (掩盖密码中段)
	dsn := defaultCfg.Database.DSN()
	maskedDSN := maskDSN(dsn)
	fmt.Printf("Viper Debug - Final Masked DSN: %s\n", maskedDSN)

	// 6. 最终检查：如果数据库主机为空，说明映射完全失败
	if defaultCfg.Database.Host == "" {
		return nil, fmt.Errorf("critical error: database host is empty. please check your config file or environment variables")
	}

	return defaultCfg, nil
}

// bindAllEnvVars 深度递归绑定结构体中的所有字段到 Viper
// 这样即便没有 YAML 文件，Viper 也能知道该去读取哪些环境变量
func bindAllEnvVars(v *viper.Viper, i interface{}, prefix string) {
	val := reflect.ValueOf(i)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return // Only process structs
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		// Get the mapstructure tag for the field name
		tag := field.Tag.Get("mapstructure")
		if tag == "" {
			tag = strings.ToLower(field.Name) // Default to lowercase field name if no tag
		}

		currentKey := tag
		if prefix != "" {
			currentKey = prefix + "." + tag
		}

		// Bind the environment variable for the current key
		// Viper automatically converts "parent.child" to "PARENT_CHILD" for env vars
		_ = v.BindEnv(currentKey)

		// 核心修复：如果环境变量中有值，显式设置到 Viper 中
		// 这样 Unmarshal 过程才能确保命中这些来自环境的值（即使 YAML 中没有对应项）
		if val := v.Get(currentKey); val != nil {
			v.Set(currentKey, val)
		}

		// Recursively call for nested structs
		if fieldVal.Kind() == reflect.Struct {
			bindAllEnvVars(v, fieldVal.Addr().Interface(), currentKey)
		} else if fieldVal.Kind() == reflect.Ptr && fieldVal.Elem().Kind() == reflect.Struct {
			bindAllEnvVars(v, fieldVal.Interface(), currentKey)
		}
	}
}

// GetDefaultConfig 获取默认配置
func GetDefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         8080,
			Mode:         "debug",
			ReadTimeout:  120, // 增加超时时间以支持AI分析 (30s -> 120s)
			WriteTimeout: 120,
			BaseURL:      "http://localhost:8080",
		},
		Database: DatabaseConfig{
			Host:              "", // 初始置空，强制要求外部提供配置
			Port:              5432,
			User:              "postgres",
			Password:          "",
			DBName:            "postgres", // 默认设为 postgres 以对齐 Supabase
			SSLMode:           "require",  // 针对 Supabase Pooler，默认要求 SSL
			Timezone:          "",         // 移除硬编码时区，避免 URL 编码导致的 unknown time zone 错误
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
		COS: COSConfig{
			BucketURL: "",
			SecretID:  "",
			SecretKey: "",
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
		Doubao: GetDefaultDoubaoConfig(),
		Analysis: AnalysisConfig{
			Timeout:    120,
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

// GetDefaultDoubaoConfig 获取默认豆包配置
func GetDefaultDoubaoConfig() DoubaoConfig {
	return DoubaoConfig{
		APIKey:  "",
		BaseURL: "https://ark.cn-beijing.volces.com/api/v3",
		Model:   "ep-20240604095209-xxxxx", // 示例 Endpoint ID
	}
}

// maskDSN 遮掩 DSN 中的密码部分
func maskDSN(dsn string) string {
	u, err := url.Parse(dsn)
	if err != nil {
		return "invalid dsn"
	}
	if u.User != nil {
		if _, set := u.User.Password(); set {
			u.User = url.UserPassword(u.User.Username(), "******")
		}
	}
	return u.String()
}

package grds

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config 数据库配置
type Config struct {
	// 基础连接配置
	Host     string `json:"host" yaml:"host"`         // 主机地址，默认 127.0.0.1
	Port     int    `json:"port" yaml:"port"`         // 端口号，默认 3306
	Username string `json:"username" yaml:"username"` // 用户名
	Password string `json:"password" yaml:"password"` // 密码
	Database string `json:"database" yaml:"database"` // 数据库名

	// 连接池配置
	MaxIdleConns    int           `json:"max_idle_conns" yaml:"max_idle_conns"`         // 最大空闲连接数，默认 10
	MaxOpenConns    int           `json:"max_open_conns" yaml:"max_open_conns"`         // 最大打开连接数，默认 100
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime" yaml:"conn_max_lifetime"`   // 连接最大生命周期，默认 1小时
	ConnMaxIdleTime time.Duration `json:"conn_max_idle_time" yaml:"conn_max_idle_time"` // 连接最大空闲时间，默认 10分钟

	// DSN 额外参数
	Charset   string            `json:"charset" yaml:"charset"`       // 字符集，默认 utf8mb4
	Collation string            `json:"collation" yaml:"collation"`   // 排序规则
	ParseTime bool              `json:"parse_time" yaml:"parse_time"` // 是否解析时间，默认 true
	Loc       string            `json:"loc" yaml:"loc"`               // 时区，默认 Local
	Timeout   time.Duration     `json:"timeout" yaml:"timeout"`       // 连接超时，默认 10秒
	Params    map[string]string `json:"params" yaml:"params"`         // 额外的 DSN 参数

	// GORM 配置
	SkipDefaultTransaction bool            `json:"skip_default_transaction" yaml:"skip_default_transaction"` // 跳过默认事务，默认 false
	PrepareStmt            bool            `json:"prepare_stmt" yaml:"prepare_stmt"`                         // 预编译语句，默认 true
	DisableAutomaticPing   bool            `json:"disable_automatic_ping" yaml:"disable_automatic_ping"`     // 禁用自动 ping，默认 false
	LogLevel               logger.LogLevel `json:"log_level" yaml:"log_level"`                               // 日志级别，默认 Silent
	SlowThreshold          time.Duration   `json:"slow_threshold" yaml:"slow_threshold"`                     // 慢查询阈值，默认 200ms

	// 自定义日志
	Logger logger.Interface `json:"-" yaml:"-"` // 自定义日志接口

	// GORM 插件和回调
	Plugins []gorm.Plugin `json:"-" yaml:"-"` // 插件列表
}

// NewDefaultConfig 创建默认配置
func NewDefaultConfig() *Config {
	return &Config{
		Host:                   "127.0.0.1",
		Port:                   3306,
		MaxIdleConns:           10,
		MaxOpenConns:           100,
		ConnMaxLifetime:        time.Hour,
		ConnMaxIdleTime:        10 * time.Minute,
		Charset:                "utf8mb4",
		ParseTime:              true,
		Loc:                    "Local",
		Timeout:                10 * time.Second,
		Params:                 make(map[string]string),
		SkipDefaultTransaction: false,
		PrepareStmt:            true,
		DisableAutomaticPing:   false,
		LogLevel:               logger.Silent,
		SlowThreshold:          200 * time.Millisecond,
		Plugins:                make([]gorm.Plugin, 0),
	}
}

// NewConfig 快速创建配置
func NewConfig(host string, port int, username, password, database string) *Config {
	cfg := NewDefaultConfig()
	cfg.Host = host
	cfg.Port = port
	cfg.Username = username
	cfg.Password = password
	cfg.Database = database
	return cfg
}

// DSN 生成数据源名称
func (c *Config) DSN() string {
	// username:password@tcp(host:port)/database?charset=utf8mb4&parseTime=True&loc=Local
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
		c.Charset,
		c.ParseTime,
		c.Loc,
	)

	if c.Collation != "" {
		dsn += fmt.Sprintf("&collation=%s", c.Collation)
	}

	if c.Timeout > 0 {
		dsn += fmt.Sprintf("&timeout=%s", c.Timeout.String())
	}

	// 添加额外参数
	for k, v := range c.Params {
		dsn += fmt.Sprintf("&%s=%s", k, v)
	}

	return dsn
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("host is required")
	}
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("invalid port: %d", c.Port)
	}
	if c.Username == "" {
		return fmt.Errorf("username is required")
	}
	if c.Database == "" {
		return fmt.Errorf("database is required")
	}
	if c.MaxOpenConns < 0 {
		return fmt.Errorf("max_open_conns must be >= 0")
	}
	if c.MaxIdleConns < 0 {
		return fmt.Errorf("max_idle_conns must be >= 0")
	}
	if c.MaxIdleConns > c.MaxOpenConns && c.MaxOpenConns > 0 {
		return fmt.Errorf("max_idle_conns cannot be greater than max_open_conns")
	}
	return nil
}

// Clone 克隆配置
func (c *Config) Clone() *Config {
	newConfig := *c
	newConfig.Params = make(map[string]string)
	for k, v := range c.Params {
		newConfig.Params[k] = v
	}
	newConfig.Plugins = append([]gorm.Plugin{}, c.Plugins...)
	return &newConfig
}

// WithHost 设置主机
func (c *Config) WithHost(host string) *Config {
	c.Host = host
	return c
}

// WithPort 设置端口
func (c *Config) WithPort(port int) *Config {
	c.Port = port
	return c
}

// WithUsername 设置用户名
func (c *Config) WithUsername(username string) *Config {
	c.Username = username
	return c
}

// WithPassword 设置密码
func (c *Config) WithPassword(password string) *Config {
	c.Password = password
	return c
}

// WithDatabase 设置数据库
func (c *Config) WithDatabase(database string) *Config {
	c.Database = database
	return c
}

// WithMaxOpenConns 设置最大打开连接数
func (c *Config) WithMaxOpenConns(n int) *Config {
	c.MaxOpenConns = n
	return c
}

// WithMaxIdleConns 设置最大空闲连接数
func (c *Config) WithMaxIdleConns(n int) *Config {
	c.MaxIdleConns = n
	return c
}

// WithConnMaxLifetime 设置连接最大生命周期
func (c *Config) WithConnMaxLifetime(d time.Duration) *Config {
	c.ConnMaxLifetime = d
	return c
}

// WithCharset 设置字符集
func (c *Config) WithCharset(charset string) *Config {
	c.Charset = charset
	return c
}

// WithLogLevel 设置日志级别
func (c *Config) WithLogLevel(level logger.LogLevel) *Config {
	c.LogLevel = level
	return c
}

// WithLogger 设置自定义日志
func (c *Config) WithLogger(l logger.Interface) *Config {
	c.Logger = l
	return c
}

// WithSlowThreshold 设置慢查询阈值
func (c *Config) WithSlowThreshold(d time.Duration) *Config {
	c.SlowThreshold = d
	return c
}

// WithPrepareStmt 设置是否使用预编译语句
func (c *Config) WithPrepareStmt(enable bool) *Config {
	c.PrepareStmt = enable
	return c
}

// WithSkipDefaultTransaction 设置是否跳过默认事务
func (c *Config) WithSkipDefaultTransaction(skip bool) *Config {
	c.SkipDefaultTransaction = skip
	return c
}

// WithParam 添加额外参数
func (c *Config) WithParam(key, value string) *Config {
	if c.Params == nil {
		c.Params = make(map[string]string)
	}
	c.Params[key] = value
	return c
}

// WithPlugin 添加插件
func (c *Config) WithPlugin(plugin gorm.Plugin) *Config {
	c.Plugins = append(c.Plugins, plugin)
	return c
}

// LogLevelInfo 设置日志级别为 Info
func (c *Config) LogLevelInfo() *Config {
	c.LogLevel = logger.Info
	return c
}

// LogLevelWarn 设置日志级别为 Warn
func (c *Config) LogLevelWarn() *Config {
	c.LogLevel = logger.Warn
	return c
}

// LogLevelError 设置日志级别为 Error
func (c *Config) LogLevelError() *Config {
	c.LogLevel = logger.Error
	return c
}

// LogLevelSilent 设置日志级别为 Silent
func (c *Config) LogLevelSilent() *Config {
	c.LogLevel = logger.Silent
	return c
}

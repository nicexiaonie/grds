package grds

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Client 数据库客户端
type Client struct {
	db     *gorm.DB
	config *Config
	mu     sync.RWMutex
	closed bool
}

// NewClient 创建客户端
func NewClient(config *Config) (*Client, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	// 创建 GORM 配置
	gormConfig := &gorm.Config{
		SkipDefaultTransaction: config.SkipDefaultTransaction,
		PrepareStmt:            config.PrepareStmt,
		DisableAutomaticPing:   config.DisableAutomaticPing,
	}

	// 配置日志
	if config.Logger != nil {
		gormConfig.Logger = config.Logger
	} else if config.LogLevel != logger.Silent {
		gormConfig.Logger = logger.Default.LogMode(config.LogLevel)
	}

	// 连接数据库
	db, err := gorm.Open(mysql.Open(config.DSN()), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// 获取底层的 *sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	client := &Client{
		db:     db,
		config: config,
	}

	// 注册插件
	for _, plugin := range config.Plugins {
		if err := db.Use(plugin); err != nil {
			return nil, fmt.Errorf("failed to register plugin: %w", err)
		}
	}

	return client, nil
}

// DB 获取 GORM DB 实例
func (c *Client) DB() *gorm.DB {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.db
}

// SqlDB 获取底层的 *sql.DB
func (c *Client) SqlDB() (*sql.DB, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.db.DB()
}

// Config 获取配置
func (c *Client) Config() *Config {
	return c.config
}

// Close 关闭连接
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}

	sqlDB, err := c.db.DB()
	if err != nil {
		return err
	}

	c.closed = true
	return sqlDB.Close()
}

// IsClosed 是否已关闭
func (c *Client) IsClosed() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.closed
}

// Ping 测试连接
func (c *Client) Ping(ctx context.Context) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return fmt.Errorf("client is closed")
	}

	sqlDB, err := c.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.PingContext(ctx)
}

// Stats 获取数据库连接池统计信息
func (c *Client) Stats() sql.DBStats {
	sqlDB, err := c.db.DB()
	if err != nil {
		return sql.DBStats{}
	}
	return sqlDB.Stats()
}

// Table 开始表查询（创建新的查询会话）
func (c *Client) Table(name string, args ...interface{}) *QueryBuilder {
	return &QueryBuilder{
		client: c,
		db:     c.db.Table(name, args...),
	}
}

// Model 使用模型进行查询
func (c *Client) Model(value interface{}) *QueryBuilder {
	return &QueryBuilder{
		client: c,
		db:     c.db.Model(value),
	}
}

// Transaction 开始事务
func (c *Client) Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	return c.db.Transaction(fc, opts...)
}

// Begin 手动开始事务
func (c *Client) Begin(opts ...*sql.TxOptions) *gorm.DB {
	return c.db.Begin(opts...)
}

// Exec 执行原生 SQL
func (c *Client) Exec(sql string, values ...interface{}) *gorm.DB {
	return c.db.Exec(sql, values...)
}

// Raw 执行原生 SQL 查询
func (c *Client) Raw(sql string, values ...interface{}) *gorm.DB {
	return c.db.Raw(sql, values...)
}

// Create 创建记录
func (c *Client) Create(value interface{}) *gorm.DB {
	return c.db.Create(value)
}

// Save 保存记录（更新所有字段）
func (c *Client) Save(value interface{}) *gorm.DB {
	return c.db.Save(value)
}

// First 查询第一条记录
func (c *Client) First(dest interface{}, conds ...interface{}) *gorm.DB {
	return c.db.First(dest, conds...)
}

// Last 查询最后一条记录
func (c *Client) Last(dest interface{}, conds ...interface{}) *gorm.DB {
	return c.db.Last(dest, conds...)
}

// Find 查询多条记录
func (c *Client) Find(dest interface{}, conds ...interface{}) *gorm.DB {
	return c.db.Find(dest, conds...)
}

// Delete 删除记录
func (c *Client) Delete(value interface{}, conds ...interface{}) *gorm.DB {
	return c.db.Delete(value, conds...)
}

// Where 添加查询条件
func (c *Client) Where(query interface{}, args ...interface{}) *gorm.DB {
	return c.db.Where(query, args...)
}

// AutoMigrate 自动迁移表结构
func (c *Client) AutoMigrate(dst ...interface{}) error {
	return c.db.AutoMigrate(dst...)
}

// Migrator 获取迁移器
func (c *Client) Migrator() gorm.Migrator {
	return c.db.Migrator()
}

// Use 使用插件
func (c *Client) Use(plugin gorm.Plugin) error {
	return c.db.Use(plugin)
}

// Session 创建新会话
func (c *Client) Session(config *gorm.Session) *gorm.DB {
	return c.db.Session(config)
}

// WithContext 设置上下文
func (c *Client) WithContext(ctx context.Context) *gorm.DB {
	return c.db.WithContext(ctx)
}

// Debug 开启调试模式
func (c *Client) Debug() *gorm.DB {
	return c.db.Debug()
}

// Scopes 应用作用域
func (c *Client) Scopes(funcs ...func(*gorm.DB) *gorm.DB) *gorm.DB {
	return c.db.Scopes(funcs...)
}

// HealthCheck 健康检查
func (c *Client) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return c.Ping(ctx)
}

// StatsInfo 获取格式化的统计信息
func (c *Client) StatsInfo() string {
	stats := c.Stats()
	return fmt.Sprintf(
		"OpenConnections: %d, InUse: %d, Idle: %d, WaitCount: %d, WaitDuration: %v, MaxIdleClosed: %d, MaxLifetimeClosed: %d",
		stats.OpenConnections,
		stats.InUse,
		stats.Idle,
		stats.WaitCount,
		stats.WaitDuration,
		stats.MaxIdleClosed,
		stats.MaxLifetimeClosed,
	)
}

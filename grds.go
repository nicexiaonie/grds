package grds

import (
	"context"

	"gorm.io/gorm"
)

// Version 版本号
const Version = "2.0.0"

// 全局默认客户端
var defaultClient *Client

// Connect 连接数据库并设置为默认客户端
func Connect(config *Config) error {
	client, err := NewClient(config)
	if err != nil {
		return err
	}
	defaultClient = client
	return nil
}

// MustConnect 连接数据库，失败则 panic
func MustConnect(config *Config) {
	if err := Connect(config); err != nil {
		panic(err)
	}
}

// SetDefaultClient 设置默认客户端
func SetDefaultClient(client *Client) {
	defaultClient = client
}

// GetDefaultClient 获取默认客户端
func GetDefaultClient() *Client {
	if defaultClient == nil {
		panic("default client not initialized, call Connect() first")
	}
	return defaultClient
}

// Close 关闭默认客户端
func Close() error {
	if defaultClient != nil {
		return defaultClient.Close()
	}
	return nil
}

// DB 获取默认客户端的 GORM DB 实例
func DB() *gorm.DB {
	return GetDefaultClient().DB()
}

// Table 使用默认客户端创建表查询
func Table(name string, args ...interface{}) *QueryBuilder {
	return GetDefaultClient().Table(name, args...)
}

// Model 使用默认客户端创建模型查询
func Model(value interface{}) *QueryBuilder {
	return GetDefaultClient().Model(value)
}

// Tx 执行事务（使用默认客户端）
func Tx(fc TxFunc) error {
	return GetDefaultClient().Transaction(fc)
}

// TxWithContext 带上下文执行事务
func TxWithContext(ctx context.Context, fc TxFunc) error {
	return TransactionWithContext(ctx, GetDefaultClient().DB(), fc)
}

// Create 创建记录（使用默认客户端）
func Create(value interface{}) error {
	return GetDefaultClient().Create(value).Error
}

// Save 保存记录（使用默认客户端）
func Save(value interface{}) error {
	return GetDefaultClient().Save(value).Error
}

// First 查询第一条记录（使用默认客户端）
func First(dest interface{}, conds ...interface{}) error {
	return GetDefaultClient().First(dest, conds...).Error
}

// Find 查询记录（使用默认客户端）
func Find(dest interface{}, conds ...interface{}) error {
	return GetDefaultClient().Find(dest, conds...).Error
}

// Delete 删除记录（使用默认客户端）
func Delete(value interface{}, conds ...interface{}) error {
	return GetDefaultClient().Delete(value, conds...).Error
}

// Where 添加查询条件（使用默认客户端）
func Where(query interface{}, args ...interface{}) *gorm.DB {
	return GetDefaultClient().Where(query, args...)
}

// Exec 执行原生 SQL（使用默认客户端）
func Exec(sql string, values ...interface{}) error {
	return GetDefaultClient().Exec(sql, values...).Error
}

// Raw 执行原生 SQL 查询（使用默认客户端）
func Raw(sql string, values ...interface{}) *gorm.DB {
	return GetDefaultClient().Raw(sql, values...)
}

// AutoMigrate 自动迁移（使用默认客户端）
func AutoMigrate(dst ...interface{}) error {
	return GetDefaultClient().AutoMigrate(dst...)
}

// Ping 测试连接（使用默认客户端）
func Ping() error {
	ctx := context.Background()
	return GetDefaultClient().Ping(ctx)
}

// Stats 获取统计信息（使用默认客户端）
func Stats() string {
	return GetDefaultClient().StatsInfo()
}

// HealthCheck 健康检查（使用默认客户端）
func HealthCheck() error {
	return GetDefaultClient().HealthCheck()
}

// Debug 开启调试模式（使用默认客户端）
func Debug() *gorm.DB {
	return GetDefaultClient().Debug()
}

// WithContext 设置上下文（使用默认客户端）
func WithContext(ctx context.Context) *gorm.DB {
	return GetDefaultClient().WithContext(ctx)
}

// Use 使用插件（使用默认客户端）
func Use(plugin gorm.Plugin) error {
	return GetDefaultClient().Use(plugin)
}

// RegisterCallbacks 获取回调注册器（使用默认客户端）
func RegisterCallbacks() *CallbackRegistry {
	return NewCallbackRegistry(GetDefaultClient().DB())
}

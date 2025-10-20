package grds

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
)

// TxFunc 事务函数类型
type TxFunc func(tx *gorm.DB) error

// Transaction 执行事务（自动提交/回滚）
func Transaction(db *gorm.DB, fc TxFunc, opts ...*sql.TxOptions) error {
	return db.Transaction(fc, opts...)
}

// TransactionWithContext 带上下文的事务
func TransactionWithContext(ctx context.Context, db *gorm.DB, fc TxFunc, opts ...*sql.TxOptions) error {
	return db.WithContext(ctx).Transaction(fc, opts...)
}

// Begin 开始事务
func Begin(db *gorm.DB, opts ...*sql.TxOptions) *gorm.DB {
	return db.Begin(opts...)
}

// Commit 提交事务
func Commit(tx *gorm.DB) error {
	return tx.Commit().Error
}

// Rollback 回滚事务
func Rollback(tx *gorm.DB) error {
	return tx.Rollback().Error
}

// SavePoint 创建保存点
func SavePoint(tx *gorm.DB, name string) error {
	return tx.SavePoint(name).Error
}

// RollbackTo 回滚到保存点
func RollbackTo(tx *gorm.DB, name string) error {
	return tx.RollbackTo(name).Error
}

// TxManager 事务管理器
type TxManager struct {
	db *gorm.DB
}

// NewTxManager 创建事务管理器
func NewTxManager(db *gorm.DB) *TxManager {
	return &TxManager{db: db}
}

// Execute 执行事务
func (tm *TxManager) Execute(fc TxFunc, opts ...*sql.TxOptions) error {
	return tm.db.Transaction(fc, opts...)
}

// ExecuteWithContext 带上下文执行事务
func (tm *TxManager) ExecuteWithContext(ctx context.Context, fc TxFunc, opts ...*sql.TxOptions) error {
	return tm.db.WithContext(ctx).Transaction(fc, opts...)
}

// ReadCommitted 读已提交事务
func (tm *TxManager) ReadCommitted(fc TxFunc) error {
	return tm.Execute(fc, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
}

// RepeatableRead 可重复读事务
func (tm *TxManager) RepeatableRead(fc TxFunc) error {
	return tm.Execute(fc, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
}

// Serializable 串行化事务
func (tm *TxManager) Serializable(fc TxFunc) error {
	return tm.Execute(fc, &sql.TxOptions{Isolation: sql.LevelSerializable})
}

// ReadOnly 只读事务
func (tm *TxManager) ReadOnly(fc TxFunc) error {
	return tm.Execute(fc, &sql.TxOptions{ReadOnly: true})
}

// WithSavepoint 使用保存点执行操作
func (tm *TxManager) WithSavepoint(tx *gorm.DB, name string, fc func(*gorm.DB) error) error {
	// 创建保存点
	if err := SavePoint(tx, name); err != nil {
		return err
	}

	// 执行操作
	if err := fc(tx); err != nil {
		// 回滚到保存点
		_ = RollbackTo(tx, name)
		return err
	}

	return nil
}

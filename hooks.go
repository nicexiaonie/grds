package grds

import (
	"gorm.io/gorm"
)

// HookFunc GORM 回调函数类型
type HookFunc func(*gorm.DB)

// CallbackRegistry 回调注册器
type CallbackRegistry struct {
	db *gorm.DB
}

// NewCallbackRegistry 创建回调注册器
func NewCallbackRegistry(db *gorm.DB) *CallbackRegistry {
	return &CallbackRegistry{db: db}
}

// BeforeCreate 注册创建前回调
func (cr *CallbackRegistry) BeforeCreate(name string, fn HookFunc) error {
	return cr.db.Callback().Create().Before("gorm:create").Register(name, fn)
}

// AfterCreate 注册创建后回调
func (cr *CallbackRegistry) AfterCreate(name string, fn HookFunc) error {
	return cr.db.Callback().Create().After("gorm:create").Register(name, fn)
}

// BeforeUpdate 注册更新前回调
func (cr *CallbackRegistry) BeforeUpdate(name string, fn HookFunc) error {
	return cr.db.Callback().Update().Before("gorm:update").Register(name, fn)
}

// AfterUpdate 注册更新后回调
func (cr *CallbackRegistry) AfterUpdate(name string, fn HookFunc) error {
	return cr.db.Callback().Update().After("gorm:update").Register(name, fn)
}

// BeforeDelete 注册删除前回调
func (cr *CallbackRegistry) BeforeDelete(name string, fn HookFunc) error {
	return cr.db.Callback().Delete().Before("gorm:delete").Register(name, fn)
}

// AfterDelete 注册删除后回调
func (cr *CallbackRegistry) AfterDelete(name string, fn HookFunc) error {
	return cr.db.Callback().Delete().After("gorm:delete").Register(name, fn)
}

// BeforeQuery 注册查询前回调
func (cr *CallbackRegistry) BeforeQuery(name string, fn HookFunc) error {
	return cr.db.Callback().Query().Before("gorm:query").Register(name, fn)
}

// AfterQuery 注册查询后回调
func (cr *CallbackRegistry) AfterQuery(name string, fn HookFunc) error {
	return cr.db.Callback().Query().After("gorm:query").Register(name, fn)
}

// BeforeRow 注册查询行前回调
func (cr *CallbackRegistry) BeforeRow(name string, fn HookFunc) error {
	return cr.db.Callback().Row().Before("gorm:row").Register(name, fn)
}

// AfterRow 注册查询行后回调
func (cr *CallbackRegistry) AfterRow(name string, fn HookFunc) error {
	return cr.db.Callback().Row().After("gorm:row").Register(name, fn)
}

// BeforeRaw 注册原生SQL前回调
func (cr *CallbackRegistry) BeforeRaw(name string, fn HookFunc) error {
	return cr.db.Callback().Raw().Before("gorm:raw").Register(name, fn)
}

// AfterRaw 注册原生SQL后回调
func (cr *CallbackRegistry) AfterRaw(name string, fn HookFunc) error {
	return cr.db.Callback().Raw().After("gorm:raw").Register(name, fn)
}

// Remove 移除回调
func (cr *CallbackRegistry) Remove(processor string, name string) error {
	switch processor {
	case "create":
		return cr.db.Callback().Create().Remove(name)
	case "update":
		return cr.db.Callback().Update().Remove(name)
	case "delete":
		return cr.db.Callback().Delete().Remove(name)
	case "query":
		return cr.db.Callback().Query().Remove(name)
	case "row":
		return cr.db.Callback().Row().Remove(name)
	case "raw":
		return cr.db.Callback().Raw().Remove(name)
	}
	return nil
}

// Replace 替换回调
func (cr *CallbackRegistry) Replace(processor string, name string, fn HookFunc) error {
	switch processor {
	case "create":
		return cr.db.Callback().Create().Replace(name, fn)
	case "update":
		return cr.db.Callback().Update().Replace(name, fn)
	case "delete":
		return cr.db.Callback().Delete().Replace(name, fn)
	case "query":
		return cr.db.Callback().Query().Replace(name, fn)
	case "row":
		return cr.db.Callback().Row().Replace(name, fn)
	case "raw":
		return cr.db.Callback().Raw().Replace(name, fn)
	}
	return nil
}

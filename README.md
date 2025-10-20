# GRDS - Go Relational Database Simplifier

基于 **GORM v2** 的 MySQL 数据库工具库，提供开箱即用、功能强大、全面、简洁的 MySQL 管理工具。

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.16-blue)](https://golang.org/)
[![GORM Version](https://img.shields.io/badge/GORM-v2-green)](https://gorm.io/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

## ✨ 特性

- 🚀 **开箱即用** - 简洁的 API 设计，快速上手
- 💪 **基于 GORM v2** - 享受 GORM 的强大功能和生态
- 🔗 **链式调用** - 流畅的查询构建器
- 📦 **灵活配置** - 支持丰富的配置选项
- 🔄 **事务支持** - 简化事务操作，支持多种隔离级别
- 🪝 **钩子系统** - 灵活的回调系统
- 🔒 **并发安全** - 线程安全的设计
- 📊 **连接池管理** - 完善的连接池配置和监控
- ⚡ **高性能** - 基于 GORM v2，性能优异
- 🎯 **全局/独立** - 支持全局默认客户端和多数据库实例

## 📦 安装

```bash
go get -u github.com/nicexiaonie/grds
```

## 🚀 快速开始

### 基础使用

```go
package main

import (
    "log"
    "github.com/nicexiaonie/grds"
)

type User struct {
    ID   uint   `gorm:"primaryKey"`
    Name string
    Age  int
}

func main() {
    // 1. 创建配置
    config := grds.NewConfig("127.0.0.1", 3306, "root", "password", "testdb")
    
    // 2. 连接数据库（设置为全局默认客户端）
    if err := grds.Connect(config); err != nil {
        log.Fatal(err)
    }
    defer grds.Close()
    
    // 3. 查询数据
    var users []User
    err := grds.Model(&User{}).
        WhereGt("age", 18).
        OrderByDesc("id").
        Limit(10).
        Find(&users)
    
    if err != nil {
        log.Fatal(err)
    }
}
```

### 使用独立客户端

```go
// 创建独立客户端（支持多数据库实例）
client, err := grds.NewClient(config)
if err != nil {
    log.Fatal(err)
}
defer client.Close()

// 使用客户端进行查询
var users []User
err = client.Model(&User{}).Find(&users)
```

## 📖 详细文档

### 配置选项

```go
// 创建默认配置
config := grds.NewDefaultConfig()

// 或者使用快速配置
config := grds.NewConfig("127.0.0.1", 3306, "root", "password", "testdb")

// 链式配置
config.WithMaxOpenConns(100).
    WithMaxIdleConns(10).
    WithConnMaxLifetime(time.Hour).
    WithLogLevelInfo().  // 开启日志
    WithPrepareStmt(true) // 使用预编译语句
```

#### 配置参数说明

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| Host | string | 127.0.0.1 | 数据库主机地址 |
| Port | int | 3306 | 数据库端口 |
| Username | string | - | 用户名 |
| Password | string | - | 密码 |
| Database | string | - | 数据库名 |
| MaxOpenConns | int | 100 | 最大打开连接数 |
| MaxIdleConns | int | 10 | 最大空闲连接数 |
| ConnMaxLifetime | duration | 1h | 连接最大生命周期 |
| ConnMaxIdleTime | duration | 10m | 连接最大空闲时间 |
| Charset | string | utf8mb4 | 字符集 |
| ParseTime | bool | true | 是否解析时间类型 |
| PrepareStmt | bool | true | 是否使用预编译语句 |
| LogLevel | logger.LogLevel | Silent | 日志级别 |
| SlowThreshold | duration | 200ms | 慢查询阈值 |

### 查询操作

#### 基础查询

```go
// 查询所有记录
var users []User
err := grds.Find(&users)

// 查询第一条记录
var user User
err := grds.First(&user)

// 使用 Model 查询
err := grds.Model(&User{}).Find(&users)

// 使用 Table 查询
err := grds.Table("users").Find(&users)

// 统计数量
count, err := grds.Model(&User{}).Count()

// 检查是否存在
exists, err := grds.Model(&User{}).WhereEq("name", "admin").Exists()
```

#### WHERE 条件

```go
// 基础条件
grds.Model(&User{}).Where("age > ?", 18)

// 便捷方法
grds.Model(&User{}).WhereEq("name", "John")    // 等于
grds.Model(&User{}).WhereNe("status", 0)       // 不等于
grds.Model(&User{}).WhereGt("age", 18)         // 大于
grds.Model(&User{}).WhereGte("age", 18)        // 大于等于
grds.Model(&User{}).WhereLt("age", 60)         // 小于
grds.Model(&User{}).WhereLte("age", 60)        // 小于等于
grds.Model(&User{}).WhereLike("name", "John%") // LIKE
grds.Model(&User{}).WhereIn("id", []int{1, 2, 3}) // IN
grds.Model(&User{}).WhereNotIn("status", []int{0, 1}) // NOT IN
grds.Model(&User{}).WhereBetween("age", 18, 60) // BETWEEN
grds.Model(&User{}).WhereNull("deleted_at")     // IS NULL
grds.Model(&User{}).WhereNotNull("email")       // IS NOT NULL

// NOT 条件
grds.Model(&User{}).Not("age > ?", 60)

// OR 条件
grds.Model(&User{}).Where("age > ?", 18).Or("is_vip = ?", true)
```

#### 排序、分组、分页

```go
// 排序
grds.Model(&User{}).Order("created_at DESC")
grds.Model(&User{}).OrderByAsc("age")
grds.Model(&User{}).OrderByDesc("created_at")

// 分组
grds.Model(&User{}).
    Select("age, COUNT(*) as count").
    GroupBy("age").
    Having("COUNT(*) > ?", 10)

// 限制和偏移
grds.Model(&User{}).Limit(10).Offset(20)

// 分页（page 从 1 开始）
grds.Model(&User{}).Page(1, 20) // 第1页，每页20条

// 去重
grds.Model(&User{}).Distinct("age")
```

#### 联表查询

```go
// JOIN
grds.Model(&User{}).Joins("JOIN profiles ON users.id = profiles.user_id")

// LEFT JOIN
grds.Model(&User{}).LeftJoin("profiles", "users.id = profiles.user_id")

// RIGHT JOIN
grds.Model(&User{}).RightJoin("orders", "users.id = orders.user_id")

// INNER JOIN
grds.Model(&User{}).InnerJoin("profiles", "users.id = profiles.user_id")

// 复杂联表
grds.Model(&User{}).
    Select("users.*, profiles.bio, COUNT(orders.id) as order_count").
    LeftJoin("profiles", "users.id = profiles.user_id").
    LeftJoin("orders", "users.id = orders.user_id").
    GroupBy("users.id")
```

#### 预加载

```go
// 预加载关联
grds.Model(&User{}).Preload("Orders").Find(&users)

// 嵌套预加载
grds.Model(&User{}).Preload("Orders.Items").Find(&users)

// 条件预加载
grds.Model(&User{}).Preload("Orders", "status = ?", "completed").Find(&users)
```

### 创建操作

```go
// 创建单条记录
user := User{Name: "John", Age: 25}
err := grds.Create(&user)

// 批量创建
users := []User{
    {Name: "John", Age: 25},
    {Name: "Jane", Age: 30},
}
err := grds.Model(&User{}).CreateInBatches(users, 100)
```

### 更新操作

```go
// 更新单个字段
err := grds.Model(&User{}).WhereEq("id", 1).Update("age", 26)

// 更新多个字段（使用 map）
err := grds.Model(&User{}).WhereEq("id", 1).Updates(map[string]interface{}{
    "name": "John Updated",
    "age":  26,
})

// 更新多个字段（使用结构体）
err := grds.Model(&User{}).WhereEq("id", 1).Updates(User{Name: "John", Age: 26})

// 保存所有字段
user := User{ID: 1, Name: "John", Age: 26}
err := grds.Save(&user)

// 更新列（不触发钩子）
err := grds.Model(&User{}).WhereEq("id", 1).UpdateColumn("age", 26)
err := grds.Model(&User{}).WhereEq("id", 1).UpdateColumns(map[string]interface{}{"age": 26})
```

### 删除操作

```go
// 删除记录
err := grds.Delete(&User{}, 1) // 根据主键删除

// 条件删除
err := grds.Model(&User{}).WhereEq("age", 0).Delete(&User{})

// 批量删除
err := grds.Model(&User{}).Where("created_at < ?", time.Now().AddDate(0, -6, 0)).Delete(&User{})
```

### 事务操作

#### 自动事务

```go
// 基础事务
err := grds.Tx(func(tx *gorm.DB) error {
    // 在事务中执行操作
    if err := tx.Create(&user).Error; err != nil {
        return err // 自动回滚
    }
    
    if err := tx.Create(&order).Error; err != nil {
        return err // 自动回滚
    }
    
    return nil // 自动提交
})

// 带上下文的事务
err := grds.TxWithContext(ctx, func(tx *gorm.DB) error {
    // 事务操作
    return nil
})
```

#### 手动事务

```go
// 开始事务
tx := grds.DB().Begin()

// 执行操作
if err := tx.Create(&user).Error; err != nil {
    tx.Rollback()
    return err
}

// 提交事务
if err := tx.Commit().Error; err != nil {
    return err
}
```

#### 事务隔离级别

```go
client := grds.GetDefaultClient()
txMgr := grds.NewTxManager(client.DB())

// 读已提交
err := txMgr.ReadCommitted(func(tx *gorm.DB) error {
    // 事务操作
    return nil
})

// 可重复读
err := txMgr.RepeatableRead(txFunc)

// 串行化
err := txMgr.Serializable(txFunc)

// 只读事务
err := txMgr.ReadOnly(txFunc)
```

#### 保存点

```go
tx := grds.DB().Begin()

// 创建保存点
grds.SavePoint(tx, "sp1")

// 执行操作
if err := tx.Create(&user).Error; err != nil {
    // 回滚到保存点
    grds.RollbackTo(tx, "sp1")
}

tx.Commit()
```

### 钩子系统

```go
// 获取回调注册器
callbacks := grds.RegisterCallbacks()

// 注册创建前回调
callbacks.BeforeCreate("log_before_create", func(db *gorm.DB) error {
    log.Println("Before create")
    return nil
})

// 注册创建后回调
callbacks.AfterCreate("log_after_create", func(db *gorm.DB) error {
    log.Println("After create")
    return nil
})

// 其他回调
callbacks.BeforeUpdate("before_update", hookFunc)
callbacks.AfterUpdate("after_update", hookFunc)
callbacks.BeforeDelete("before_delete", hookFunc)
callbacks.AfterDelete("after_delete", hookFunc)
callbacks.BeforeQuery("before_query", hookFunc)
callbacks.AfterQuery("after_query", hookFunc)
```

### 模型定义

```go
type User struct {
    ID        uint           `gorm:"primaryKey"`
    Name      string         `gorm:"size:100;not null"`
    Email     string         `gorm:"uniqueIndex;size:100"`
    Age       int            `gorm:"default:0"`
    CreatedAt time.Time      `gorm:"autoCreateTime"`
    UpdatedAt time.Time      `gorm:"autoUpdateTime"`
    DeletedAt gorm.DeletedAt `gorm:"index"` // 软删除
}

// 自定义表名
func (User) TableName() string {
    return "my_users"
}

// 使用模型
var users []User
grds.Model(&User{}).Find(&users)
```

### 原生 SQL

```go
// 执行原生 SQL
err := grds.Exec("UPDATE users SET age = age + 1 WHERE id = ?", 1)

// 原生查询
var users []User
grds.Raw("SELECT * FROM users WHERE age > ?", 18).Scan(&users)
```

### 统计信息

```go
// 获取连接池统计信息
stats := grds.Stats()
fmt.Println(stats)

// 或者获取详细的统计
client := grds.GetDefaultClient()
dbStats := client.Stats()
fmt.Printf("Open: %d, InUse: %d, Idle: %d\n", 
    dbStats.OpenConnections, 
    dbStats.InUse, 
    dbStats.Idle)
```

### 健康检查

```go
// 健康检查
if err := grds.HealthCheck(); err != nil {
    log.Printf("Database unhealthy: %v", err)
}

// 或者使用 Ping
if err := grds.Ping(); err != nil {
    log.Printf("Ping failed: %v", err)
}
```

### 自动迁移

```go
// 自动迁移表结构
err := grds.AutoMigrate(&User{}, &Order{}, &Product{})

// 获取迁移器进行更多操作
migrator := grds.GetDefaultClient().Migrator()

// 检查表是否存在
if migrator.HasTable(&User{}) {
    // ...
}

// 创建表
migrator.CreateTable(&User{})

// 删除表
migrator.DropTable(&User{})

// 重命名表
migrator.RenameTable(&User{}, &UserV2{})

// 添加列
migrator.AddColumn(&User{}, "nickname")

// 删除列
migrator.DropColumn(&User{}, "nickname")
```

## 🔧 高级特性

### 多数据库实例

```go
// 主数据库
mainConfig := grds.NewConfig("127.0.0.1", 3306, "root", "pass", "main_db")
mainClient, _ := grds.NewClient(mainConfig)

// 从数据库
slaveConfig := grds.NewConfig("127.0.0.1", 3307, "root", "pass", "slave_db")
slaveClient, _ := grds.NewClient(slaveConfig)

// 使用不同的客户端
mainClient.Model(&User{}).Find(&users)
slaveClient.Model(&User{}).Find(&users)
```

### 作用域（Scopes）

```go
// 定义作用域
func ActiveUsers(db *gorm.DB) *gorm.DB {
    return db.Where("status = ?", "active")
}

func RecentUsers(db *gorm.DB) *gorm.DB {
    return db.Where("created_at > ?", time.Now().AddDate(0, -1, 0))
}

// 使用作用域
var users []User
grds.Model(&User{}).Scopes(ActiveUsers, RecentUsers).Find(&users)
```

### 查询构建器克隆

```go
// 创建基础查询
baseQuery := grds.Model(&User{}).WhereEq("status", "active")

// 克隆查询构建器（不影响原查询）
query1 := baseQuery.Clone().WhereGt("age", 18)
query2 := baseQuery.Clone().WhereLt("age", 60)

// 两个查询互不影响
query1.Find(&users1)
query2.Find(&users2)
```

### 调试模式

```go
// 开启调试模式（打印 SQL）
grds.Debug().Model(&User{}).Find(&users)

// 或者在客户端级别开启
client := grds.GetDefaultClient()
client.Debug().Model(&User{}).Find(&users)
```

### 上下文支持

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// 使用上下文
var users []User
grds.WithContext(ctx).Model(&User{}).Find(&users)
```

## 📊 性能建议

1. **使用预编译语句**
   ```go
   config.WithPrepareStmt(true) // 默认已开启
   ```

2. **合理设置连接池参数**
   ```go
   config.WithMaxOpenConns(100)  // 根据负载调整
   config.WithMaxIdleConns(10)   // 通常为 MaxOpenConns 的 10%
   ```

3. **批量操作**
   ```go
   // 推荐：批量创建
   grds.Model(&User{}).CreateInBatches(users, 100)
   ```

4. **使用索引和限制结果集**
   ```go
   grds.Model(&User{}).
       Where("indexed_column = ?", value).
       Limit(100) // 限制结果集
   ```

5. **监控慢查询**
   ```go
   config.WithSlowThreshold(200 * time.Millisecond)
   config.WithLogLevelWarn() // 记录慢查询
   ```

## 🌟 GORM 生态兼容

由于 GRDS 基于 GORM v2，您可以直接使用 GORM 的所有插件和生态：

```go
import (
    "gorm.io/plugin/dbresolver"
    "gorm.io/plugin/prometheus"
)

// 使用读写分离插件
grds.Use(dbresolver.Register(dbresolver.Config{
    // 配置读写分离
}))

// 使用 Prometheus 监控插件
grds.Use(prometheus.New(prometheus.Config{
    // 配置监控
}))
```

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License

## 📮 联系方式

- GitHub: [github.com/nicexiaonie/grds](https://github.com/nicexiaonie/grds)
- Issues: [github.com/nicexiaonie/grds/issues](https://github.com/nicexiaonie/grds/issues)

## 🙏 致谢

本项目基于以下优秀的开源项目：
- [GORM](https://gorm.io/) - The fantastic ORM library for Golang
- [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql) - MySQL driver for Go


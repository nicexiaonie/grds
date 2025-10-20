[toc]
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

### 模型生成器

GRDS 提供了内置的模型生成器，可以从数据库表结构自动生成 GORM 模型代码。

> 📖 **详细使用指南**: 查看 [GENERATOR_USAGE.md](./GENERATOR_USAGE.md) 了解完整的功能和配置选项。

**核心特性**：
- ✅ 自动获取表注释和字段注释（参考 gormt）
- ✅ 自定义数据库类型到 Go 类型的映射
- ✅ 完整的 GORM 标签支持（包括类型、默认值、注释等）
- ✅ 灵活的 JSON 标签命名风格（snake_case、camelCase、original）
- ✅ 支持表前缀去除
- ✅ 支持选择性生成表

#### 快速开始

##### 1. 在您的项目中引入 grds

```bash
go get github.com/nicexiaonie/grds
```

##### 2. 安装命令行工具

```bash
go install github.com/nicexiaonie/grds/cmd/grds-gen@latest
```

安装成功后，`grds-gen` 命令会被添加到 `$GOPATH/bin` 目录（确保该目录在您的 PATH 中）。

##### 3. 初始化配置文件

在您的项目根目录运行：

```bash
cd your-project
grds-gen -init
```

这将创建 `.grds.yaml` 配置文件：

```yaml
# GRDS 模型生成器配置文件
database:
  host: 127.0.0.1
  port: 3306
  username: root
  password: your_password
  database: your_database

generator:
  # 输出目录
  out_dir: ./models
  # 输出文件名
  out_file: models.go
  # 包名
  package_name: models
  # 指定要生成的表（留空则生成所有表）
  tables: []
  # 表前缀（生成时会去除）
  table_prefix: ""
```

##### 4. 编辑配置文件

编辑 `.grds.yaml`，填写您的数据库连接信息。

##### 5. 生成模型

运行：

```bash
grds-gen
```

生成成功后，您会看到类似的输出：

```
📝 使用配置文件: .grds.yaml
正在生成模型...
数据库: root@127.0.0.1:3306/mydb
输出目录: ./models
输出文件: models.go
包名: models
生成所有表
--------------------------------------------------
✅ 模型生成成功！
📁 文件位置: ./models/models.go
```

##### 6. 在代码中使用生成的模型

```go
package main

import (
    "github.com/nicexiaonie/grds"
    "your-project/models"
)

func main() {
    // 连接数据库
    config := grds.NewConfig("127.0.0.1", 3306, "root", "password", "mydb")
    grds.MustConnect(config)
    defer grds.Close()
    
    // 使用生成的模型
    var users []models.Users
    grds.Model(&models.Users{}).Find(&users)
}
```

#### 高级用法

##### 指定配置文件：

```bash
grds-gen -config=./config/db.yaml
```

#### 使用命令行参数

命令行参数会覆盖配置文件：

```bash
# 生成所有表
grds-gen -database=mydb -user=root -password=secret

# 生成指定表
grds-gen -database=mydb -tables=users,orders,products

# 指定输出目录和包名
grds-gen -database=mydb -out=./internal/models -package=model

# 设置表前缀（生成时去除）
grds-gen -database=mydb -prefix=tbl_
```

#### 查看数据库信息

```bash
# 列出所有表
grds-gen -list

# 查看表结构
grds-gen -columns=users
```

#### 完整命令行参数

| 参数 | 说明 | 默认值 |
|------|------|--------|
| -config | 配置文件路径 | 自动查找 |
| -host | 数据库主机 | 127.0.0.1 |
| -port | 数据库端口 | 3306 |
| -user | 用户名 | root |
| -password | 密码 | |
| -database | 数据库名 | |
| -out | 输出目录 | ./models |
| -file | 输出文件名 | models.go |
| -package | 包名 | models |
| -tables | 指定表名(逗号分隔) | |
| -prefix | 表前缀 | |
| -list | 列出所有表 | false |
| -columns | 显示表结构 | |
| -init | 初始化配置文件 | false |
| -version | 显示版本 | false |

#### 在代码中使用生成器

如果需要在代码中使用生成器功能：

```go
package main

import (
    "log"
    "github.com/nicexiaonie/grds"
)

func main() {
    // 方式 1：快速生成
    err := grds.GenerateModels("127.0.0.1", 3306, "root", "password", "mydb", "./models")
    if err != nil {
        log.Fatal(err)
    }
    
    // 方式 2：使用配置
    config := grds.NewGeneratorConfig("127.0.0.1", 3306, "root", "password", "mydb")
    config.WithOutDir("./models").
        WithPackageName("models").
        WithTables("users", "orders").
        WithTypeMapping(map[string]string{
            "decimal": "decimal.Decimal",
        }).
        WithJSONTagStyle("camelCase")
    
    if err := config.Generate(); err != nil {
        log.Fatal(err)
    }
}
```

#### 生成的模型示例

假设数据库中有一个 `users` 表（包含注释），生成的模型如下：

```sql
CREATE TABLE `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `email` varchar(100) NOT NULL COMMENT '邮箱',
  `age` int(11) NOT NULL COMMENT '年龄',
  `balance` decimal(10,2) DEFAULT '0.00' COMMENT '余额',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';
```

生成的 Go 模型：

```go
package models

import "time"

// Users 用户表
type Users struct {
	Id        int       `gorm:"column:id;type:int(11);primaryKey;autoIncrement;not null;comment:用户ID" json:"id"` // 用户ID
	Username  string    `gorm:"column:username;type:varchar(50);not null;comment:用户名" json:"username"` // 用户名
	Email     string    `gorm:"column:email;type:varchar(100);not null;comment:邮箱" json:"email"` // 邮箱
	Age       int       `gorm:"column:age;type:int(11);not null;comment:年龄" json:"age"` // 年龄
	Balance   float64   `gorm:"column:balance;type:decimal(10,2);default:0.00;comment:余额" json:"balance"` // 余额
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;comment:创建时间" json:"created_at"` // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null;comment:更新时间" json:"updated_at"` // 更新时间
}

// TableName 指定表名
func (Users) TableName() string {
	return "users"
}
```

#### 自定义类型映射

grds-gen 支持自定义数据库类型到 Go 类型的映射，参考 gormt 的实现。在配置文件中添加：

```yaml
generator:
  # 自定义类型映射
  type_mapping:
    # 将 datetime 映射为 time.Time（默认已配置）
    datetime: time.Time
    # 将 decimal 映射为自定义类型
    decimal: decimal.Decimal
    # 将 json 映射为 json.RawMessage
    json: json.RawMessage
    # 将 text 映射为 sql.NullString
    text: sql.NullString
```

在代码中使用：

```go
config := grds.NewGeneratorConfig("127.0.0.1", 3306, "root", "password", "mydb")
config.WithTypeMapping(map[string]string{
    "decimal": "decimal.Decimal",
    "json":    "json.RawMessage",
})
config.Generate()
```

#### 默认类型映射

| 数据库类型 | Go 类型 |
|-----------|---------|
| tinyint | int8 |
| tinyint unsigned | uint8 |
| smallint | int16 |
| smallint unsigned | uint16 |
| int, integer | int |
| int unsigned | uint32 |
| bigint | int64 |
| bigint unsigned | uint64 |
| float | float32 |
| double, decimal | float64 |
| char, varchar, text | string |
| datetime, date, timestamp | time.Time |
| time | string |
| year | int |
| blob, binary | []byte |
| json | string |
| enum, set | string |

#### JSON 标签命名风格

支持三种 JSON 标签命名风格：

```yaml
generator:
  # snake_case（默认）：user_name -> "user_name"
  json_tag_style: snake_case
  
  # camelCase：user_name -> "userName"
  # json_tag_style: camelCase
  
  # original：保持原样
  # json_tag_style: original
```

#### 控制标签生成

```yaml
generator:
  # 是否生成 JSON 标签（默认: true）
  enable_json_tag: true
  # 是否生成 GORM 标签（默认: true）
  enable_gorm_tag: true
```

在代码中：

```go
config.WithEnableJSONTag(false)  // 不生成 JSON 标签
config.WithEnableGormTag(true)   // 生成 GORM 标签
config.WithJSONTagStyle("camelCase")  // 使用小驼峰命名
```

#### 表和字段注释

生成器会自动获取并生成：
- **表注释**：作为结构体注释
- **字段注释**：作为字段的行尾注释
- **GORM comment 标签**：包含在 GORM 标签中

示例数据库：
```sql
CREATE TABLE `products` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '产品ID',
  `name` varchar(100) NOT NULL COMMENT '产品名称',
  `price` decimal(10,2) NOT NULL COMMENT '价格',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB COMMENT='产品表';
```

生成的模型：
```go
// Products 产品表
type Products struct {
    Id    int64   `gorm:"column:id;type:bigint(20);primaryKey;autoIncrement;not null;comment:产品ID" json:"id"` // 产品ID
    Name  string  `gorm:"column:name;type:varchar(100);not null;comment:产品名称" json:"name"` // 产品名称
    Price float64 `gorm:"column:price;type:decimal(10,2);not null;comment:价格" json:"price"` // 价格
}
```

#### 项目集成示例

在项目的 `Makefile` 中添加：

```makefile
.PHONY: gen-models
gen-models:
	grds-gen

.PHONY: gen-models-tables
gen-models-tables:
	grds-gen -tables=users,orders,products
```

在项目的 `scripts` 目录创建 `gen.sh`：

```bash
#!/bin/bash
echo "生成数据库模型..."
grds-gen
echo "完成！"
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


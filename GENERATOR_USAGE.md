# GRDS 模型生成器使用指南

本文档详细说明 GRDS 模型生成器的使用方法，参考了 `github.com/xxjwxc/gormt` 的设计。

## 核心特性

✅ **自动获取表注释和字段注释**  
✅ **自定义类型映射（参考 gormt）**  
✅ **完整的 GORM 标签支持**  
✅ **灵活的 JSON 标签命名风格**  
✅ **支持表前缀去除**  
✅ **支持选择性生成表**  

## 快速开始

### 1. 安装

```bash
go install github.com/nicexiaonie/grds/cmd/grds-gen@latest
```

### 2. 初始化配置

```bash
cd your-project
grds-gen -init
```

### 3. 编辑配置文件

编辑 `.grds.yaml`：

```yaml
database:
  host: 127.0.0.1
  port: 3306
  username: root
  password: your_password
  database: your_database

generator:
  out_dir: ./models
  out_file: models.go
  package_name: models
  tables: []  # 留空生成所有表
  table_prefix: ""
  
  # 自定义类型映射
  type_mapping:
    datetime: time.Time
    decimal: float64
  
  enable_json_tag: true
  enable_gorm_tag: true
  json_tag_style: snake_case
```

### 4. 生成模型

```bash
grds-gen
```

## 类型映射详解

### 默认类型映射

grds-gen 提供了完善的默认类型映射：

| MySQL 类型 | Go 类型 | 说明 |
|-----------|---------|------|
| `tinyint` | `int8` | 有符号 |
| `tinyint unsigned` | `uint8` | 无符号 |
| `smallint` | `int16` | 有符号 |
| `smallint unsigned` | `uint16` | 无符号 |
| `int`, `integer` | `int` | 有符号 |
| `int unsigned` | `uint32` | 无符号 |
| `bigint` | `int64` | 有符号 |
| `bigint unsigned` | `uint64` | 无符号 |
| `float` | `float32` | 单精度 |
| `double` | `float64` | 双精度 |
| `decimal` | `float64` | 定点数 |
| `char`, `varchar`, `text` | `string` | 字符串 |
| `datetime`, `date`, `timestamp` | `time.Time` | 时间类型 |
| `time` | `string` | 时间字符串 |
| `year` | `int` | 年份 |
| `blob`, `binary` | `[]byte` | 二进制 |
| `json` | `string` | JSON 字符串 |
| `enum`, `set` | `string` | 枚举和集合 |

### 自定义类型映射

#### 在配置文件中

```yaml
generator:
  type_mapping:
    # 使用第三方库的类型
    decimal: decimal.Decimal
    json: json.RawMessage
    
    # 使用可空类型
    text: sql.NullString
    datetime: sql.NullTime
    
    # 自定义类型
    varchar: mypackage.CustomString
```

#### 在代码中

```go
config := grds.NewGeneratorConfig("127.0.0.1", 3306, "root", "pass", "db")
config.WithTypeMapping(map[string]string{
    "decimal": "decimal.Decimal",
    "json":    "json.RawMessage",
    "datetime": "time.Time",
})
config.Generate()
```

## 注释支持

### 表注释

数据库表注释会作为 Go 结构体的注释：

```sql
CREATE TABLE `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB COMMENT='用户表';
```

生成：

```go
// Users 用户表
type Users struct {
    // ...
}
```

### 字段注释

字段注释会以三种形式体现：

1. **行尾注释**
2. **GORM comment 标签**
3. **完整的字段元信息**

```sql
CREATE TABLE `products` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '产品ID',
  `name` varchar(100) NOT NULL COMMENT '产品名称',
  `price` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '价格',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态：1-上架，0-下架',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB COMMENT='产品表';
```

生成的模型：

```go
// Products 产品表
type Products struct {
    Id     int64   `gorm:"column:id;type:bigint(20);primaryKey;autoIncrement;not null;comment:产品ID" json:"id"` // 产品ID
    Name   string  `gorm:"column:name;type:varchar(100);not null;comment:产品名称" json:"name"` // 产品名称
    Price  float64 `gorm:"column:price;type:decimal(10,2);not null;default:0.00;comment:价格" json:"price"` // 价格
    Status int8    `gorm:"column:status;type:tinyint(4);not null;default:1;comment:状态：1-上架，0-下架" json:"status"` // 状态：1-上架，0-下架
}
```

## JSON 标签命名风格

### snake_case（默认）

保持数据库字段的下划线命名：

```go
type User struct {
    UserId   int    `json:"user_id"`
    UserName string `json:"user_name"`
}
```

配置：
```yaml
json_tag_style: snake_case
```

### camelCase

转换为小驼峰命名：

```go
type User struct {
    UserId   int    `json:"userId"`
    UserName string `json:"userName"`
}
```

配置：
```yaml
json_tag_style: camelCase
```

### original

保持原始字段名：

```go
type User struct {
    UserId   int    `json:"user_id"`  # 如果数据库字段是 user_id
    UserName string `json:"UserName"` # 如果数据库字段是 UserName
}
```

配置：
```yaml
json_tag_style: original
```

## 标签控制

### 只生成 GORM 标签

```yaml
generator:
  enable_json_tag: false
  enable_gorm_tag: true
```

生成：
```go
type User struct {
    Id   int    `gorm:"column:id;type:int(11);primaryKey"`
    Name string `gorm:"column:name;type:varchar(50)"`
}
```

### 只生成 JSON 标签

```yaml
generator:
  enable_json_tag: true
  enable_gorm_tag: false
```

生成：
```go
type User struct {
    Id   int    `json:"id"`
    Name string `json:"name"`
}
```

### 不生成任何标签

```yaml
generator:
  enable_json_tag: false
  enable_gorm_tag: false
```

生成：
```go
type User struct {
    Id   int
    Name string
}
```

## 表前缀处理

如果数据库表使用了统一前缀，可以在生成时去除：

```yaml
generator:
  table_prefix: "tbl_"
```

表名映射：
- `tbl_users` → `Users`
- `tbl_orders` → `Orders`
- `tbl_products` → `Products`

## 选择性生成

### 只生成指定的表

```yaml
generator:
  tables: [users, orders, products]
```

或使用命令行：

```bash
grds-gen -tables=users,orders,products
```

### 生成所有表

```yaml
generator:
  tables: []  # 空数组
```

或：

```bash
grds-gen  # 不指定 -tables 参数
```

## 高级用法

### 查看数据库信息

```bash
# 列出所有表
grds-gen -list

# 查看表结构
grds-gen -columns=users
```

### 在代码中使用

```go
package main

import (
    "log"
    "github.com/nicexiaonie/grds"
)

func main() {
    // 创建配置
    config := grds.NewGeneratorConfig(
        "127.0.0.1",
        3306,
        "root",
        "password",
        "mydb",
    )
    
    // 配置选项
    config.WithOutDir("./internal/models").
        WithPackageName("model").
        WithTables("users", "orders").
        WithTablePrefix("tbl_").
        WithTypeMapping(map[string]string{
            "decimal": "decimal.Decimal",
            "json":    "json.RawMessage",
        }).
        WithJSONTagStyle("camelCase").
        WithEnableJSONTag(true).
        WithEnableGormTag(true)
    
    // 生成模型
    if err := config.Generate(); err != nil {
        log.Fatal(err)
    }
    
    log.Println("模型生成成功！")
}
```

## 完整配置示例

```yaml
# .grds.yaml
database:
  host: 127.0.0.1
  port: 3306
  username: root
  password: your_password
  database: your_database

generator:
  # 基础配置
  out_dir: ./internal/models
  out_file: models.go
  package_name: model
  
  # 表选择
  tables: []  # 空则生成所有表
  table_prefix: "tbl_"
  
  # 类型映射（参考 gormt）
  type_mapping:
    # 时间类型
    datetime: time.Time
    date: time.Time
    timestamp: time.Time
    
    # 数值类型
    decimal: decimal.Decimal
    
    # JSON 类型
    json: json.RawMessage
    
    # 可空类型
    text: sql.NullString
    
  # 标签配置
  enable_json_tag: true
  enable_gorm_tag: true
  json_tag_style: camelCase  # snake_case, camelCase, original
```

## 对比 gormt

| 特性 | grds-gen | gormt |
|------|----------|-------|
| 表注释 | ✅ 支持 | ✅ 支持 |
| 字段注释 | ✅ 支持 | ✅ 支持 |
| 自定义类型映射 | ✅ 支持 | ✅ 支持 |
| JSON 标签风格 | ✅ 3种风格 | ✅ 支持 |
| GORM 标签 | ✅ 完整支持 | ✅ 支持 |
| 配置文件 | ✅ YAML/JSON | ✅ YAML |
| 命令行工具 | ✅ 独立工具 | ✅ 独立工具 |
| 集成方式 | ✅ go install | ✅ go install |
| 表前缀去除 | ✅ 支持 | ✅ 支持 |
| 选择性生成 | ✅ 支持 | ✅ 支持 |

## 常见问题

### Q: 如何使用自定义类型？

A: 在配置文件中指定类型映射，并确保生成的文件中导入相应的包：

```yaml
type_mapping:
  decimal: decimal.Decimal
```

生成后需要手动在文件中添加导入：
```go
import (
    "github.com/shopspring/decimal"
)
```

### Q: datetime 为什么默认映射为 time.Time？

A: 这是参考 gormt 和 GORM 的最佳实践。time.Time 可以自动处理时区，且与 GORM 完全兼容。

### Q: 如何生成可空字段？

A: 使用自定义类型映射：

```yaml
type_mapping:
  varchar: sql.NullString
  int: sql.NullInt64
  datetime: sql.NullTime
```

### Q: 生成的模型可以直接使用吗？

A: 可以！生成的模型包含完整的 GORM 标签和 TableName 方法，可以直接用于 GORM 操作。

## 总结

grds-gen 提供了完整的数据库模型生成功能，参考了 gormt 的优秀设计，并增加了更多灵活的配置选项。通过配置文件或代码方式，可以轻松定制生成的模型代码，满足不同项目的需求。


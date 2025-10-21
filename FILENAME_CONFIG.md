# 模型生成文件名配置详解

## 概述

`grds-gen` 工具支持灵活配置生成模型文件的名称、位置和组织方式。本文档详细说明文件名相关的配置选项。

## 配置项说明

### 1. 输出目录 (OutDir)

**配置键**: `out_dir` (配置文件) 或 `-out` (命令行)  
**默认值**: `./models`  
**说明**: 指定生成的模型文件存放的目录

**示例**:
```yaml
# .grds.yaml
generator:
  out_dir: ./internal/models
```

```bash
# 命令行
grds-gen -out ./internal/models
```

### 2. 输出文件名 (OutFileName)

**配置键**: `out_file` (配置文件) 或 `-file` (命令行)  
**默认值**: `models.go`  
**说明**: 指定生成的模型文件名称

**示例**:
```yaml
# .grds.yaml
generator:
  out_file: db_models.go
```

```bash
# 命令行
grds-gen -file db_models.go
```

### 3. 包名 (PackageName)

**配置键**: `package_name` (配置文件) 或 `-package` (命令行)  
**默认值**: `models`  
**说明**: 指定生成的 Go 包名，会影响文件开头的 `package` 声明

**示例**:
```yaml
# .grds.yaml
generator:
  package_name: entity
```

```bash
# 命令行
grds-gen -package entity
```

### 4. 分文件生成 (SeparateFile) ⭐ 新功能

**配置键**: `separate_file` (配置文件) 或 `-separate` (命令行)  
**默认值**: `false`  
**说明**: 是否为每个表生成单独的文件

- `false`（默认）: 所有表生成到一个文件，使用 `out_file` 指定的文件名
- `true`: 每个表生成一个独立的文件，文件名格式为 `表名_model.go`

**示例**:
```yaml
# .grds.yaml
generator:
  separate_file: true  # 启用分文件生成
  out_dir: ./models
```

```bash
# 命令行
grds-gen -separate

# 或者结合其他参数
grds-gen -separate -out ./models -tables users,orders,products
```

**生成效果对比**:

```yaml
# separate_file: false (默认)
# 生成文件: ./models/models.go
# 文件包含所有表的结构体

# separate_file: true
# 生成文件:
# - ./models/users_model.go      # User 结构体
# - ./models/orders_model.go     # Order 结构体
# - ./models/products_model.go   # Product 结构体
```

## 常见使用场景

### 场景 1: 默认配置

最简单的方式，使用所有默认值：

```yaml
generator:
  out_dir: ./models
  out_file: models.go
  package_name: models
```

生成结果：
```
./models/models.go
```

文件内容开头：
```go
package models

import "time"

// User 用户表
type User struct {
    // ...
}
```

### 场景 2: 多数据库分离

为不同的数据库生成不同的模型文件：

```yaml
# .grds.user.yaml - 用户库
database:
  database: user_db
generator:
  out_dir: ./models/user
  out_file: user_models.go
  package_name: user
```

```yaml
# .grds.order.yaml - 订单库
database:
  database: order_db
generator:
  out_dir: ./models/order
  out_file: order_models.go
  package_name: order
```

使用：
```bash
grds-gen -config .grds.user.yaml
grds-gen -config .grds.order.yaml
```

生成结果：
```
./models/
├── user/
│   └── user_models.go
└── order/
    └── order_models.go
```

### 场景 3: 按业务模块分离

为不同业务模块生成独立的模型：

```bash
# 用户模块
grds-gen -tables users,user_profiles,user_settings \
  -out ./internal/user/models \
  -file user.go \
  -package models

# 订单模块
grds-gen -tables orders,order_items,order_status \
  -out ./internal/order/models \
  -file order.go \
  -package models

# 商品模块
grds-gen -tables products,categories,inventory \
  -out ./internal/product/models \
  -file product.go \
  -package models
```

生成结果：
```
./internal/
├── user/
│   └── models/
│       └── user.go
├── order/
│   └── models/
│       └── order.go
└── product/
    └── models/
        └── product.go
```

### 场景 4: 分文件生成（推荐用于大型项目） ⭐

为每个表生成独立的文件，便于代码管理和维护：

**配置方式**:
```yaml
# .grds.yaml
generator:
  out_dir: ./models
  separate_file: true  # 启用分文件生成
```

**命令行方式**:
```bash
grds-gen -separate -out ./models
```

**生成结果**:
```
./models/
├── users_model.go           # User 结构体
├── user_profiles_model.go   # UserProfile 结构体
├── orders_model.go          # Order 结构体
├── order_items_model.go     # OrderItem 结构体
├── products_model.go        # Product 结构体
└── categories_model.go      # Category 结构体
```

**users_model.go 示例**:
```go
package models

import "time"

// User 用户表
type User struct {
    ID        uint      `gorm:"column:id;type:int unsigned;primaryKey;autoIncrement;comment:用户ID" json:"id"`
    Username  string    `gorm:"column:username;type:varchar(50);not null;comment:用户名" json:"username"`
    Email     string    `gorm:"column:email;type:varchar(100);comment:邮箱" json:"email"`
    CreatedAt time.Time `gorm:"column:created_at;type:datetime;comment:创建时间" json:"created_at"`
}

func (User) TableName() string {
    return "users"
}
```

**优势**:
- ✅ 代码组织清晰，每个模型独立管理
- ✅ Git 合并冲突大大减少
- ✅ 便于查找和定位特定模型
- ✅ 支持按需加载，提高编译速度
- ✅ 适合大型项目和团队协作

**何时使用分文件**:
- 表数量 > 10 个
- 多人协作项目
- 需要频繁更新模型
- 希望代码结构更清晰

### 场景 5: 环境分离

为不同环境使用不同的配置：

```yaml
# .grds.dev.yaml
database:
  host: localhost
generator:
  out_dir: ./models
  out_file: models_dev.go
  package_name: models
```

```yaml
# .grds.prod.yaml
database:
  host: prod-db.example.com
generator:
  out_dir: ./models
  out_file: models.go
  package_name: models
```

### 场景 6: 清晰的目录结构（推荐）

推荐的项目结构：

```yaml
generator:
  out_dir: ./internal/domain/model
  out_file: generated.go
  package_name: model
```

项目结构：
```
myproject/
├── cmd/
│   └── grds-gen/
├── internal/
│   └── domain/
│       └── model/
│           ├── generated.go      # 生成的模型
│           └── custom.go         # 自定义扩展
├── pkg/
├── .grds.yaml
└── main.go
```

## 配置优先级

当同时存在多种配置时，优先级从高到低：

1. **命令行参数** (`-file`, `-out`, `-package`)
2. **指定的配置文件** (`-config /path/to/config.yaml`)
3. **当前目录默认配置文件**（按以下顺序查找）：
   - `.grds.yaml`
   - `.grds.yml`
   - `grds.yaml`
   - `grds.yml`
   - `.grds.json`
   - `grds.json`
4. **代码默认值**

**示例**:
```bash
# 使用配置文件中的设置，但覆盖输出文件名
grds-gen -config .grds.yaml -file custom_models.go
```

## 命令行快速使用

### 查看帮助
```bash
grds-gen -h
```

### 初始化配置文件
```bash
grds-gen -init
```

这会在当前目录创建 `.grds.yaml` 模板文件，包含所有配置项的说明。

### 指定所有参数
```bash
grds-gen \
  -host 127.0.0.1 \
  -port 3306 \
  -user root \
  -password secret \
  -database mydb \
  -out ./models \
  -file models.go \
  -package models \
  -tables users,orders \
  -prefix tbl_
```

### 查看当前配置
```bash
# 列出所有表（用于确认配置是否正确）
grds-gen -list

# 查看表结构
grds-gen -columns users
```

## 文件命名最佳实践

### 1. 单文件 vs 分文件 ⭐

#### 单一文件方式 (separate_file: false)

**适用场景**: 小型项目，表数量不多（< 10 个表）

```yaml
generator:
  out_dir: ./models
  out_file: models.go
  package_name: models
  separate_file: false  # 默认值
```

**优点**:
- ✅ 简单直接，一个文件包含所有模型
- ✅ 适合小型项目
- ✅ 查看所有模型方便

**缺点**:
- ❌ 文件可能很大（> 1000 行）
- ❌ 多人协作容易产生 Git 冲突
- ❌ 查找特定模型需要搜索
- ❌ IDE 加载可能变慢

#### 分文件方式 (separate_file: true) 🌟 推荐

**适用场景**: 中大型项目，表数量 > 10 个

```yaml
generator:
  out_dir: ./models
  separate_file: true  # 每个表一个文件
  package_name: models
```

**命令行**:
```bash
grds-gen -separate -out ./models
```

**优点**:
- ✅ 每个模型独立文件，代码清晰
- ✅ 减少 Git 冲突，便于团队协作
- ✅ 文件命名规范：`表名_model.go`
- ✅ 易于查找和维护特定模型
- ✅ IDE 性能更好
- ✅ 支持按需加载

**缺点**:
- ❌ 文件数量多（但这通常是优点）

**对比示例**:

```bash
# separate_file: false
./models/
└── models.go  (3000+ 行，包含所有 50 个表)

# separate_file: true (推荐)
./models/
├── users_model.go
├── orders_model.go
├── products_model.go
├── categories_model.go
└── ... (50 个文件)
```

**最佳实践建议**:
- 📌 表数量 ≤ 5: 使用单文件（separate_file: false）
- 📌 表数量 6-10: 根据团队偏好选择
- 📌 表数量 > 10: **强烈推荐分文件**（separate_file: true）
- 📌 多人协作: **强烈推荐分文件**

### 2. 按模块拆分

**适用场景**: 中大型项目，业务模块清晰

```bash
# 分别为每个模块生成
grds-gen -tables users,roles,permissions -file auth.go
grds-gen -tables products,categories -file catalog.go
grds-gen -tables orders,payments -file commerce.go
```

**优点**:
- 代码组织清晰
- 减少合并冲突
- 职责分明

**缺点**:
- 需要多次执行命令
- 管理多个配置文件

### 3. 按数据库拆分

**适用场景**: 微服务架构，多数据库

为每个数据库创建独立的目录和包：

```
./models/
├── userdb/
│   └── models.go
├── orderdb/
│   └── models.go
└── productdb/
    └── models.go
```

### 4. 分层架构

**适用场景**: DDD 或分层架构项目

```yaml
generator:
  out_dir: ./internal/domain/entity
  out_file: db_entity.go
  package_name: entity
```

配合自定义扩展：
```
./internal/
└── domain/
    ├── entity/
    │   ├── db_entity.go      # 生成的基础实体
    │   └── user_extend.go    # 自定义扩展方法
    └── repository/
```

## 编程式使用

在代码中动态生成模型：

```go
package main

import (
    "github.com/nicexiaonie/grds"
)

func main() {
    // 方式 1: 使用默认配置
    err := grds.GenerateModels(
        "localhost", 3306,
        "root", "password",
        "mydb", "./models",
    )
    
    // 方式 2: 使用完整配置
    config := grds.NewGeneratorConfig(
        "localhost", 3306,
        "root", "password",
        "mydb",
    )
    config.OutDir = "./internal/models"
    config.OutFileName = "generated.go"
    config.PackageName = "models"
    config.Tables = []string{"users", "orders"}
    config.TablePrefix = "tbl_"
    
    err = config.Generate()
    
    // 方式 3: 链式配置
    err = grds.NewGeneratorConfig(
        "localhost", 3306,
        "root", "password",
        "mydb",
    ).WithOutFileName("my_models.go").
      WithTypeMapping(map[string]string{
          "datetime": "time.Time",
      }).
      WithJSONTagStyle("camelCase").
      Generate()
}
```

## 注意事项

1. **文件覆盖**: 每次生成会覆盖已存在的文件，请注意备份手动修改的内容
2. **包名一致性**: 同一目录下的所有 Go 文件应使用相同的包名
3. **避免冲突**: 不同模块使用不同的输出目录，避免文件名冲突
4. **版本控制**: 建议将配置文件（`.grds.yaml`）加入版本控制，但排除敏感信息（密码）

## 常见问题

### Q1: 如何生成到不存在的目录？
**A**: 工具会自动创建不存在的目录。

```bash
grds-gen -out ./path/that/does/not/exist
# 会自动创建完整目录结构
```

### Q2: 可以生成多个文件吗？
**A**: 当前版本将所有表生成到一个文件中。如需拆分，可以：
- 多次执行，每次指定不同的表和文件名
- 生成后手动拆分文件

### Q3: 文件名可以包含路径吗？
**A**: 不推荐。文件名应只包含文件名，使用 `out_dir` 指定路径：

```yaml
# 正确
generator:
  out_dir: ./models/user
  out_file: user.go

# 不推荐
generator:
  out_dir: ./models
  out_file: user/user.go  # 可能导致问题
```

### Q4: 如何为不同表生成不同文件？
**A**: 使用 `-tables` 参数多次执行：

```bash
grds-gen -tables users -file user.go
grds-gen -tables orders -file order.go
grds-gen -tables products -file product.go
```

## 相关文档

- [README.md](./README.md) - 完整使用文档
- [GENERATOR_USAGE.md](./GENERATOR_USAGE.md) - 生成器详细使用指南
- [.grds.yaml.example](./.grds.yaml.example) - 配置文件示例

## 版本历史

- v1.0.0 (2024-10): 初始版本，支持基础文件名配置


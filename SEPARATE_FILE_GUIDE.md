# 分文件生成快速指南

## 功能简介

从 v1.1.0 开始，`grds-gen` 支持为每个表生成单独的文件，文件名格式为 `表名_model.go`。

这个功能特别适合：
- 🎯 大型项目（表数量 > 10 个）
- 👥 多人团队协作
- 🔄 需要频繁更新模型的项目

## 快速开始

### 方式一：配置文件

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
  separate_file: true  # ⭐ 启用分文件生成
  package_name: models
```

然后运行：
```bash
grds-gen
```

### 方式二：命令行参数

```bash
grds-gen -separate \
  -host 127.0.0.1 \
  -port 3306 \
  -user root \
  -password your_password \
  -database your_database \
  -out ./models
```

## 生成结果对比

### 单文件模式（默认）

```bash
# separate_file: false
./models/
└── models.go  # 包含所有表的模型
```

**models.go** (可能有几千行):
```go
package models

// User 用户表
type User struct { ... }

// Order 订单表  
type Order struct { ... }

// Product 商品表
type Product struct { ... }

// ... 更多模型
```

### 分文件模式（新功能）⭐

```bash
# separate_file: true
./models/
├── users_model.go      # User 结构体
├── orders_model.go     # Order 结构体  
├── products_model.go   # Product 结构体
└── ...                 # 每个表一个文件
```

**users_model.go** (简洁清晰):
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

## 优势对比

| 特性 | 单文件模式 | 分文件模式 ⭐ |
|------|-----------|-------------|
| 文件数量 | 1 个 | N 个（每表一个） |
| 文件大小 | 可能很大 (>1000行) | 小且可控 (50-200行) |
| Git 冲突 | 容易冲突 | 大大减少 |
| 查找模型 | 需要搜索 | 直接打开文件 |
| 团队协作 | 困难 | 友好 |
| IDE 性能 | 可能变慢 | 更好 |
| 适用场景 | 小项目 (≤5表) | 中大型项目 (>10表) |

## 使用示例

### 示例 1：生成所有表（分文件）

```bash
grds-gen -separate
```

### 示例 2：生成指定表（分文件）

```bash
grds-gen -separate -tables users,orders,products
```

### 示例 3：结合其他选项

```bash
grds-gen -separate \
  -out ./internal/models \
  -package models \
  -prefix tbl_ \
  -tables tbl_users,tbl_orders
```

生成结果：
```
./internal/models/
├── users_model.go   # tbl_users -> User
└── orders_model.go  # tbl_orders -> Order
```

## 编程方式使用

```go
package main

import "github.com/nicexiaonie/grds"

func main() {
    config := grds.NewGeneratorConfig(
        "localhost", 3306,
        "root", "password",
        "mydb",
    )
    
    // 启用分文件生成
    config.WithSeparateFile(true).
        WithOutDir("./models").
        WithPackageName("models")
    
    if err := config.Generate(); err != nil {
        panic(err)
    }
}
```

## 配置优先级

当同时使用配置文件和命令行参数时：

```bash
# 配置文件中 separate_file: false
# 命令行使用 -separate
grds-gen -config .grds.yaml -separate
```

结果：使用 **命令行参数**（separate: true），因为命令行优先级更高。

## 文件命名规则

生成的文件名规则：`表名_model.go`

| 表名 | 文件名 |
|------|--------|
| users | users_model.go |
| user_profiles | user_profiles_model.go |
| order_items | order_items_model.go |
| tbl_products (prefix: tbl_) | products_model.go |

## 何时使用分文件？

### ✅ 推荐使用分文件的场景

- 表数量 > 10 个
- 多人团队协作开发
- 需要频繁修改模型
- 项目需要长期维护
- 追求代码清晰度

### 🚫 可以使用单文件的场景

- 表数量 < 5 个
- 个人小项目
- 原型开发/快速验证
- 模型很少变化

## 最佳实践

### 1. 使用版本控制

```bash
# .gitignore
models/*_model.go  # 如果模型是自动生成的
!models/custom.go  # 保留自定义扩展
```

### 2. 分离自定义代码

```
./models/
├── users_model.go        # 生成的基础模型
├── users_custom.go       # 自定义扩展方法
├── orders_model.go       # 生成的基础模型
└── orders_custom.go      # 自定义扩展方法
```

### 3. CI/CD 集成

```yaml
# .github/workflows/generate-models.yml
name: Generate Models

on:
  workflow_dispatch:

jobs:
  generate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Install grds-gen
        run: go install github.com/nicexiaonie/grds/cmd/grds-gen@latest
      - name: Generate models
        run: grds-gen -separate
        env:
          DB_HOST: ${{ secrets.DB_HOST }}
          DB_PASSWORD: ${{ secrets.DB_PASSWORD }}
      - name: Commit changes
        run: |
          git add models/
          git commit -m "Auto-generate models"
          git push
```

## 故障排除

### Q: 如何切换回单文件模式？

A: 设置 `separate_file: false` 或移除 `-separate` 参数。

### Q: 已经生成的分文件如何清理？

```bash
# 删除所有 *_model.go 文件
rm models/*_model.go

# 重新生成为单文件
grds-gen  # separate_file 默认为 false
```

### Q: 文件太多怎么办？

A: 考虑按模块拆分到不同目录：

```bash
# 用户模块
grds-gen -separate -tables users,user_profiles -out ./models/user

# 订单模块
grds-gen -separate -tables orders,order_items -out ./models/order
```

## 更多信息

- 📖 [完整配置文档](./FILENAME_CONFIG.md)
- 📖 [生成器使用指南](./GENERATOR_USAGE.md)
- 📖 [README](./README.md)

## 反馈

如有问题或建议，请提交 [Issue](https://github.com/nicexiaonie/grds/issues)。


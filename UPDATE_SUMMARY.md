# 更新总结 - 分文件生成功能

## 📅 更新日期
2024-10-21

## 🎯 新增功能

### 分文件生成模型 ⭐

现在支持为每个表生成单独的文件，文件名格式为 `表名_model.go`。

#### 核心改动

1. **generator.go**
   - 新增 `SeparateFile` 字段到 `GeneratorConfig`
   - 重构 `Generate()` 方法，拆分为 `generateSingleFile()` 和 `generateSeparateFiles()`
   - 新增 `WithSeparateFile()` 链式配置方法

2. **cmd/grds-gen/main.go**
   - 新增 `-separate` 命令行参数
   - 配置文件支持 `separate_file` 选项
   - 更新输出信息，区分单文件和分文件模式

3. **.grds.yaml.example**
   - 添加 `separate_file` 配置项及详细说明

4. **文档更新**
   - 新增 `FILENAME_CONFIG.md` - 文件名配置详解（462行）
   - 新增 `SEPARATE_FILE_GUIDE.md` - 分文件生成快速指南
   - 更新 `README.md` - 添加分文件生成使用说明

## 💡 使用方法

### 方式一：配置文件

```yaml
# .grds.yaml
generator:
  separate_file: true  # 启用分文件生成
  out_dir: ./models
```

```bash
grds-gen
```

### 方式二：命令行

```bash
grds-gen -separate -out ./models
```

### 方式三：编程方式

```go
config := grds.NewGeneratorConfig(
    "localhost", 3306,
    "root", "password",
    "mydb",
)
config.WithSeparateFile(true).
    WithOutDir("./models")
    
config.Generate()
```

## 📊 效果对比

### 单文件模式（默认）
```
./models/
└── models.go  (所有表的模型，可能 1000+ 行)
```

### 分文件模式（新功能）
```
./models/
├── users_model.go       (50-100 行)
├── orders_model.go      (50-100 行)
├── products_model.go    (50-100 行)
└── ...                  (每个表一个文件)
```

## ✨ 优势

1. **减少冲突** - Git 合并冲突大大减少
2. **代码清晰** - 每个模型独立文件，职责分明
3. **便于维护** - 快速定位特定模型文件
4. **提升性能** - IDE 加载和编译更快
5. **团队协作** - 多人可同时修改不同表的模型

## 📈 适用场景

### ✅ 推荐使用分文件

- 表数量 > 10 个
- 多人团队协作
- 需要频繁更新模型
- 大型项目

### 🔄 可以使用单文件

- 表数量 < 5 个
- 个人小项目
- 原型开发

## 📚 相关文档

| 文档 | 说明 |
|------|------|
| [FILENAME_CONFIG.md](./FILENAME_CONFIG.md) | 文件名配置详解（462行） |
| [SEPARATE_FILE_GUIDE.md](./SEPARATE_FILE_GUIDE.md) | 分文件生成快速指南 |
| [GENERATOR_USAGE.md](./GENERATOR_USAGE.md) | 生成器完整使用指南 |
| [README.md](./README.md) | 项目总体说明 |

## 🔧 配置选项

### 完整配置示例

```yaml
database:
  host: 127.0.0.1
  port: 3306
  username: root
  password: your_password
  database: your_database

generator:
  # 基础配置
  out_dir: ./models
  out_file: models.go              # 单文件模式使用
  package_name: models
  separate_file: true              # ⭐ 新增：启用分文件生成
  
  # 表选择
  tables: []                       # 留空生成所有表
  table_prefix: ""                 # 表前缀
  
  # 类型映射
  type_mapping:
    datetime: time.Time
    decimal: float64
  
  # 标签配置
  enable_json_tag: true
  enable_gorm_tag: true
  json_tag_style: snake_case       # snake_case, camelCase, original
```

### 命令行参数

```bash
grds-gen \
  -separate \                      # ⭐ 新增：启用分文件生成
  -host 127.0.0.1 \
  -port 3306 \
  -user root \
  -password secret \
  -database mydb \
  -out ./models \
  -package models \
  -tables users,orders \
  -prefix tbl_
```

## 🚀 版本兼容性

- ✅ 向后兼容：默认使用单文件模式（separate_file: false）
- ✅ 无破坏性变更：现有配置无需修改
- ✅ 可选启用：通过配置或命令行参数按需启用

## 📝 Git 提交历史

```bash
200385c ✨ 新功能：支持分文件生成模型
2d1c168 更新 README 安装说明
30cec3a 修复 .gitignore 并添加 cmd/grds-gen 目录
4b6013c 支持动态模型生成
7f02ee5 重大更新！！！
```

## 🔮 后续计划

- [ ] 支持自定义文件名模板（如 `{table}_entity.go`）
- [ ] 支持按模块自动分组生成
- [ ] 支持增量更新（只更新变化的表）
- [ ] 添加模型生成钩子（pre/post generation）

## 🙏 反馈

如有问题或建议，请：
1. 提交 [Issue](https://github.com/nicexiaonie/grds/issues)
2. 发起 [Pull Request](https://github.com/nicexiaonie/grds/pulls)

---

**版本**: v1.1.0  
**作者**: GRDS Team  
**许可**: MIT License


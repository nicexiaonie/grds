package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/nicexiaonie/grds"
	"gopkg.in/yaml.v3"
)

const version = "1.0.0"

// ConfigFile 配置文件结构
type ConfigFile struct {
	Database struct {
		Host     string `yaml:"host" json:"host"`
		Port     int    `yaml:"port" json:"port"`
		Username string `yaml:"username" json:"username"`
		Password string `yaml:"password" json:"password"`
		Database string `yaml:"database" json:"database"`
	} `yaml:"database" json:"database"`

	Generator struct {
		OutDir         string            `yaml:"out_dir" json:"out_dir"`
		OutFileName    string            `yaml:"out_file" json:"out_file"`
		PackageName    string            `yaml:"package_name" json:"package_name"`
		Tables         []string          `yaml:"tables" json:"tables"`
		TablePrefix    string            `yaml:"table_prefix" json:"table_prefix"`
		SeparateFile   bool              `yaml:"separate_file" json:"separate_file"`
		TypeMapping    map[string]string `yaml:"type_mapping" json:"type_mapping"`
		EnableJSONTag  bool              `yaml:"enable_json_tag" json:"enable_json_tag"`
		EnableGormTag  bool              `yaml:"enable_gorm_tag" json:"enable_gorm_tag"`
		JSONTagStyle   string            `yaml:"json_tag_style" json:"json_tag_style"`
		GenerateToJSON bool              `yaml:"generate_to_json" json:"generate_to_json"`
	} `yaml:"generator" json:"generator"`
}

func main() {
	// 定义命令行参数
	var (
		configFile = flag.String("config", "", "配置文件路径 (支持 .yaml, .yml, .json)")

		host     = flag.String("host", "", "数据库主机地址")
		port     = flag.Int("port", 0, "数据库端口")
		username = flag.String("user", "", "数据库用户名")
		password = flag.String("password", "", "数据库密码")
		database = flag.String("database", "", "数据库名")

		outDir      = flag.String("out", "", "输出目录")
		outFile     = flag.String("file", "", "输出文件名")
		packageName = flag.String("package", "", "包名")
		tables      = flag.String("tables", "", "指定要生成的表名，多个表用逗号分隔")
		tablePrefix = flag.String("prefix", "", "表前缀，生成时会去除")
		_           = flag.Bool("separate", false, "为每个表生成单独的文件（格式：表名_model.go）")

		listTables  = flag.Bool("list", false, "列出所有表名")
		showColumns = flag.String("columns", "", "显示指定表的列信息")
		showVersion = flag.Bool("version", false, "显示版本信息")
		initConfig  = flag.Bool("init", false, "在当前目录初始化配置文件")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "GRDS Model Generator v%s - 从数据库生成 GORM 模型\n\n", version)
		fmt.Fprintf(os.Stderr, "使用方法:\n")
		fmt.Fprintf(os.Stderr, "  grds-gen [选项]\n\n")
		fmt.Fprintf(os.Stderr, "选项:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n配置文件:\n")
		fmt.Fprintf(os.Stderr, "  默认会在当前目录查找以下配置文件（按优先级）：\n")
		fmt.Fprintf(os.Stderr, "    1. .grds.yaml\n")
		fmt.Fprintf(os.Stderr, "    2. .grds.yml\n")
		fmt.Fprintf(os.Stderr, "    3. grds.yaml\n")
		fmt.Fprintf(os.Stderr, "    4. grds.yml\n")
		fmt.Fprintf(os.Stderr, "    5. .grds.json\n")
		fmt.Fprintf(os.Stderr, "    6. grds.json\n\n")
		fmt.Fprintf(os.Stderr, "示例:\n")
		fmt.Fprintf(os.Stderr, "  # 初始化配置文件\n")
		fmt.Fprintf(os.Stderr, "  grds-gen -init\n\n")
		fmt.Fprintf(os.Stderr, "  # 使用配置文件生成\n")
		fmt.Fprintf(os.Stderr, "  grds-gen\n\n")
		fmt.Fprintf(os.Stderr, "  # 指定配置文件\n")
		fmt.Fprintf(os.Stderr, "  grds-gen -config=./config/db.yaml\n\n")
		fmt.Fprintf(os.Stderr, "  # 命令行参数（会覆盖配置文件）\n")
		fmt.Fprintf(os.Stderr, "  grds-gen -database=mydb -tables=users,orders\n\n")
		fmt.Fprintf(os.Stderr, "  # 列出所有表\n")
		fmt.Fprintf(os.Stderr, "  grds-gen -list\n\n")
	}

	flag.Parse()

	// 显示版本
	if *showVersion {
		fmt.Printf("grds-gen version %s\n", version)
		return
	}

	// 初始化配置文件
	if *initConfig {
		if err := initializeConfig(); err != nil {
			log.Fatalf("初始化配置文件失败: %v", err)
		}
		fmt.Println("✅ 配置文件已创建: .grds.yaml")
		fmt.Println("请编辑配置文件并运行 'grds-gen' 生成模型")
		return
	}

	// 加载配置
	config := loadConfig(*configFile)

	// 命令行参数覆盖配置文件
	if *host != "" {
		config.Database.Host = *host
	}
	if *port != 0 {
		config.Database.Port = *port
	}
	if *username != "" {
		config.Database.Username = *username
	}
	if *password != "" {
		config.Database.Password = *password
	}
	if *database != "" {
		config.Database.Database = *database
	}
	if *outDir != "" {
		config.Generator.OutDir = *outDir
	}
	if *outFile != "" {
		config.Generator.OutFileName = *outFile
	}
	if *packageName != "" {
		config.Generator.PackageName = *packageName
	}
	if *tables != "" {
		config.Generator.Tables = strings.Split(*tables, ",")
		for i, t := range config.Generator.Tables {
			config.Generator.Tables[i] = strings.TrimSpace(t)
		}
	}
	if *tablePrefix != "" {
		config.Generator.TablePrefix = *tablePrefix
	}
	// 注意：separateFile 是布尔值，使用 flag 包时需要特殊处理
	// 如果命令行指定了 -separate，则使用命令行值
	if flag.Lookup("separate").Value.String() == "true" {
		config.Generator.SeparateFile = true
	}

	// 验证必填参数
	if config.Database.Database == "" {
		fmt.Println("错误：数据库名不能为空")
		fmt.Println("请使用 -database 参数或在配置文件中指定")
		flag.Usage()
		os.Exit(1)
	}

	// 创建生成器配置
	genConfig := grds.NewGeneratorConfig(
		config.Database.Host,
		config.Database.Port,
		config.Database.Username,
		config.Database.Password,
		config.Database.Database,
	)
	genConfig.OutDir = config.Generator.OutDir
	genConfig.OutFileName = config.Generator.OutFileName
	genConfig.PackageName = config.Generator.PackageName
	genConfig.Tables = config.Generator.Tables
	genConfig.TablePrefix = config.Generator.TablePrefix
	genConfig.SeparateFile = config.Generator.SeparateFile
	genConfig.EnableJSONTag = config.Generator.EnableJSONTag
	genConfig.EnableGormTag = config.Generator.EnableGormTag
	genConfig.JSONTagStyle = config.Generator.JSONTagStyle
	genConfig.GenerateToJSON = config.Generator.GenerateToJSON

	// 应用自定义类型映射
	if len(config.Generator.TypeMapping) > 0 {
		genConfig.WithTypeMapping(config.Generator.TypeMapping)
	}

	// 列出所有表
	if *listTables {
		fmt.Println("正在获取表列表...")
		tableList, err := genConfig.GetTables()
		if err != nil {
			log.Fatalf("获取表列表失败: %v", err)
		}

		fmt.Printf("\n数据库 '%s' 中的表 (共 %d 个):\n", config.Database.Database, len(tableList))
		fmt.Println(strings.Repeat("-", 50))
		for i, table := range tableList {
			fmt.Printf("%3d. %s\n", i+1, table)
		}
		return
	}

	// 显示表的列信息
	if *showColumns != "" {
		fmt.Printf("正在获取表 '%s' 的列信息...\n", *showColumns)
		columns, err := genConfig.GetTableColumns(*showColumns)
		if err != nil {
			log.Fatalf("获取列信息失败: %v", err)
		}

		fmt.Printf("\n表 '%s' 的列信息 (共 %d 列):\n", *showColumns, len(columns))
		fmt.Println(strings.Repeat("-", 100))
		fmt.Printf("%-20s %-20s %-10s %-10s %-15s %-20s %s\n",
			"字段名", "类型", "可空", "键", "默认值", "额外", "注释")
		fmt.Println(strings.Repeat("-", 100))

		for _, col := range columns {
			defaultVal := "NULL"
			if col.Default != nil {
				defaultVal = fmt.Sprintf("%v", col.Default)
			}
			fmt.Printf("%-20s %-20s %-10s %-10s %-15s %-20s %s\n",
				col.Field, col.Type, col.Null, col.Key, defaultVal, col.Extra, col.Comment)
		}
		return
	}

	// 生成模型
	fmt.Println("正在生成模型...")
	fmt.Printf("数据库: %s@%s:%d/%s\n",
		config.Database.Username,
		config.Database.Host,
		config.Database.Port,
		config.Database.Database)
	fmt.Printf("输出目录: %s\n", config.Generator.OutDir)
	fmt.Printf("输出文件: %s\n", config.Generator.OutFileName)
	fmt.Printf("包名: %s\n", config.Generator.PackageName)

	if len(config.Generator.Tables) > 0 {
		fmt.Printf("生成表: %s\n", strings.Join(config.Generator.Tables, ", "))
	} else {
		fmt.Println("生成所有表")
	}

	if config.Generator.TablePrefix != "" {
		fmt.Printf("表前缀: %s (将被去除)\n", config.Generator.TablePrefix)
	}

	fmt.Println(strings.Repeat("-", 50))

	if err := genConfig.Generate(); err != nil {
		log.Fatalf("生成模型失败: %v", err)
	}

	fmt.Println("✅ 模型生成成功！")
	if config.Generator.SeparateFile {
		fmt.Printf("📁 文件位置: %s/*_model.go\n", config.Generator.OutDir)
	} else {
		fmt.Printf("📁 文件位置: %s/%s\n", config.Generator.OutDir, config.Generator.OutFileName)
	}
}

// loadConfig 加载配置文件
func loadConfig(configPath string) *ConfigFile {
	config := &ConfigFile{}

	// 设置默认值
	config.Database.Host = "127.0.0.1"
	config.Database.Port = 3306
	config.Database.Username = "root"
	config.Generator.OutDir = "./models"
	config.Generator.OutFileName = "models.go"
	config.Generator.PackageName = "models"
	config.Generator.EnableJSONTag = true
	config.Generator.EnableGormTag = true
	config.Generator.JSONTagStyle = "snake_case"

	// 如果指定了配置文件，直接加载
	if configPath != "" {
		if err := loadConfigFile(configPath, config); err != nil {
			log.Printf("警告: 加载配置文件失败: %v", err)
		}
		return config
	}

	// 尝试在当前目录查找配置文件
	configFiles := []string{
		".grds.yaml",
		".grds.yml",
		"grds.yaml",
		"grds.yml",
		".grds.json",
		"grds.json",
	}

	for _, file := range configFiles {
		if _, err := os.Stat(file); err == nil {
			if err := loadConfigFile(file, config); err != nil {
				log.Printf("警告: 加载配置文件 %s 失败: %v", file, err)
				continue
			}
			fmt.Printf("📝 使用配置文件: %s\n", file)
			break
		}
	}

	return config
}

// loadConfigFile 加载指定的配置文件
func loadConfigFile(path string, config *ConfigFile) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".yaml", ".yml":
		return yaml.Unmarshal(data, config)
	case ".json":
		return json.Unmarshal(data, config)
	default:
		// 尝试 YAML
		if err := yaml.Unmarshal(data, config); err == nil {
			return nil
		}
		// 尝试 JSON
		return json.Unmarshal(data, config)
	}
}

// initializeConfig 初始化配置文件
func initializeConfig() error {
	// 添加注释
	content := `# GRDS 模型生成器配置文件
# 运行 'grds-gen' 命令生成模型

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
  out_file: 
  # 包名
  package_name: models
  # 指定要生成的表（留空则生成所有表）
  tables: []
  # 表前缀（生成时会去除）
  table_prefix: ""
  # 是否为每个表生成单独的文件（默认: false）
  # true: 每个表生成一个文件，文件名格式为 表名_model.go
  # false: 所有表生成到一个文件（out_file 指定的文件名）
  separate_file: true
  
  # 自定义类型映射（可选）
  # 示例:
  # type_mapping:
  #   datetime: time.Time
  #   decimal: decimal.Decimal
  type_mapping: {}
  
  # 是否启用 JSON 标签（默认: true）
  enable_json_tag: true
  # 是否启用 GORM 标签（默认: true）
  enable_gorm_tag: true
  # JSON 标签命名风格: snake_case, camelCase, original（默认: snake_case）
  json_tag_style: snake_case
  # 是否生成 ToJsonString 方法（默认: false）
  # true: 为每个模型生成 ToJsonString() 方法，方便将结构体转换为 JSON 字符串
  generate_to_json: false
`

	return os.WriteFile(".grds.yaml", []byte(content), 0644)
}

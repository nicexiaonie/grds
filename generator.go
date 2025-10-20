package grds

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"
)

// GeneratorConfig 模型生成器配置
type GeneratorConfig struct {
	// 数据库配置
	Host     string
	Port     int
	Username string
	Password string
	Database string

	// 输出配置
	OutDir      string   // 输出目录，默认 "./models"
	OutFileName string   // 输出文件名，默认 "models.go"
	PackageName string   // 包名，默认 "models"
	Tables      []string // 指定要生成的表名，为空则生成所有表
	TablePrefix string   // 表前缀，生成时会去除

	// 类型映射配置（自定义数据库类型到Go类型的映射）
	// 例如: {"datetime": "time.Time", "decimal": "float64"}
	TypeMapping map[string]string

	// 是否启用JSON标签
	EnableJSONTag bool
	// 是否启用GORM标签
	EnableGormTag bool
	// JSON标签命名风格: "snake_case", "camelCase", "original"
	JSONTagStyle string
}

// NewGeneratorConfig 创建生成器配置
func NewGeneratorConfig(host string, port int, username, password, database string) *GeneratorConfig {
	return &GeneratorConfig{
		Host:          host,
		Port:          port,
		Username:      username,
		Password:      password,
		Database:      database,
		OutDir:        "./models",
		OutFileName:   "models.go",
		PackageName:   "models",
		Tables:        []string{},
		TablePrefix:   "",
		TypeMapping:   getDefaultTypeMapping(),
		EnableJSONTag: true,
		EnableGormTag: true,
		JSONTagStyle:  "snake_case",
	}
}

// getDefaultTypeMapping 获取默认的类型映射
func getDefaultTypeMapping() map[string]string {
	return map[string]string{
		// 整数类型
		"tinyint":   "int8",
		"smallint":  "int16",
		"mediumint": "int",
		"int":       "int",
		"integer":   "int",
		"bigint":    "int64",

		// 无符号整数
		"tinyint unsigned":   "uint8",
		"smallint unsigned":  "uint16",
		"mediumint unsigned": "uint32",
		"int unsigned":       "uint32",
		"integer unsigned":   "uint32",
		"bigint unsigned":    "uint64",

		// 浮点数
		"float":   "float32",
		"double":  "float64",
		"decimal": "float64",

		// 字符串
		"char":       "string",
		"varchar":    "string",
		"tinytext":   "string",
		"text":       "string",
		"mediumtext": "string",
		"longtext":   "string",

		// 时间类型
		"date":      "time.Time",
		"datetime":  "time.Time",
		"timestamp": "time.Time",
		"time":      "string", // time 类型映射为 string
		"year":      "int",

		// 二进制
		"tinyblob":   "[]byte",
		"blob":       "[]byte",
		"mediumblob": "[]byte",
		"longblob":   "[]byte",
		"binary":     "[]byte",
		"varbinary":  "[]byte",

		// JSON
		"json": "string",

		// 其他
		"enum": "string",
		"set":  "string",
	}
}

// NewGeneratorConfigFromDBConfig 从数据库配置创建生成器配置
func NewGeneratorConfigFromDBConfig(dbConfig *Config) *GeneratorConfig {
	return NewGeneratorConfig(
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Database,
	)
}

// Generate 生成模型
func (gc *GeneratorConfig) Generate() error {
	// 确保输出目录存在
	if err := os.MkdirAll(gc.OutDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// 创建临时客户端
	dbConfig := NewConfig(gc.Host, gc.Port, gc.Username, gc.Password, gc.Database)
	client, err := NewClient(dbConfig)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer client.Close()

	// 获取要生成的表列表
	tables := gc.Tables
	if len(tables) == 0 {
		// 如果未指定表，获取所有表
		allTables, err := gc.GetTables()
		if err != nil {
			return fmt.Errorf("failed to get tables: %w", err)
		}
		tables = allTables
	}

	// 生成每个表的模型
	var modelsCodes []string
	for _, tableName := range tables {
		columns, err := gc.GetTableColumns(tableName)
		if err != nil {
			return fmt.Errorf("failed to get columns for table %s: %w", tableName, err)
		}

		// 获取表注释
		tableComment, err := gc.GetTableComment(tableName)
		if err != nil {
			// 表注释获取失败不影响生成
			tableComment = ""
		}

		code, err := gc.generateTableModel(tableName, tableComment, columns)
		if err != nil {
			return fmt.Errorf("failed to generate model for table %s: %w", tableName, err)
		}

		modelsCodes = append(modelsCodes, code)
	}

	// 写入文件
	outputPath := fmt.Sprintf("%s/%s", gc.OutDir, gc.OutFileName)
	content := gc.buildFileContent(modelsCodes)

	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// generateTableModel 生成单个表的模型
func (gc *GeneratorConfig) generateTableModel(tableName, tableComment string, columns []ColumnInfo) (string, error) {
	// 去除表前缀
	structName := tableName
	if gc.TablePrefix != "" && strings.HasPrefix(tableName, gc.TablePrefix) {
		structName = strings.TrimPrefix(tableName, gc.TablePrefix)
	}

	// 转换为驼峰命名
	structName = toCamelCase(structName)

	// 生成字段
	var fields []FieldInfo
	for _, col := range columns {
		goType := gc.mapDBTypeToGoType(col.Type)
		field := FieldInfo{
			Name:    toCamelCase(col.Field),
			Type:    goType,
			Column:  col.Field,
			Comment: col.Comment,
			Tags:    gc.buildTags(col),
		}
		fields = append(fields, field)
	}

	// 使用模板生成代码
	tmpl := `// {{.StructName}}{{if .TableComment}} {{.TableComment}}{{end}}
type {{.StructName}} struct {
{{range .Fields}}	{{.Name}} {{.Type}}{{if .Tags}} {{.Tags}}{{end}}{{if .Comment}} // {{.Comment}}{{end}}
{{end}}}

// TableName 指定表名
func ({{.StructName}}) TableName() string {
	return "{{.TableName}}"
}
`

	t, err := template.New("model").Parse(tmpl)
	if err != nil {
		return "", err
	}

	data := map[string]interface{}{
		"StructName":   structName,
		"TableName":    tableName,
		"TableComment": tableComment,
		"Fields":       fields,
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// buildTags 构建字段标签
func (gc *GeneratorConfig) buildTags(col ColumnInfo) string {
	var tags []string

	// GORM 标签
	if gc.EnableGormTag {
		gormTags := []string{fmt.Sprintf("column:%s", col.Field)}

		// 数据类型
		gormTags = append(gormTags, fmt.Sprintf("type:%s", col.Type))

		// 主键
		if col.Key == "PRI" {
			gormTags = append(gormTags, "primaryKey")
		}

		// 自增
		if strings.Contains(col.Extra, "auto_increment") {
			gormTags = append(gormTags, "autoIncrement")
		}

		// NOT NULL
		if col.Null == "NO" {
			gormTags = append(gormTags, "not null")
		}

		// 默认值
		if col.Default != nil {
			defaultVal := fmt.Sprintf("%v", col.Default)
			if defaultVal != "" && defaultVal != "NULL" {
				gormTags = append(gormTags, fmt.Sprintf("default:%s", defaultVal))
			}
		}

		// 注释
		if col.Comment != "" {
			// GORM 的 comment 标签
			gormTags = append(gormTags, fmt.Sprintf("comment:%s", col.Comment))
		}

		tags = append(tags, fmt.Sprintf("gorm:\"%s\"", strings.Join(gormTags, ";")))
	}

	// JSON 标签
	if gc.EnableJSONTag {
		jsonTag := gc.formatJSONTag(col.Field)
		tags = append(tags, fmt.Sprintf("json:\"%s\"", jsonTag))
	}

	if len(tags) == 0 {
		return ""
	}

	return "`" + strings.Join(tags, " ") + "`"
}

// formatJSONTag 格式化JSON标签
func (gc *GeneratorConfig) formatJSONTag(fieldName string) string {
	switch gc.JSONTagStyle {
	case "camelCase":
		return toCamelCaseLower(fieldName)
	case "original":
		return fieldName
	default: // "snake_case"
		return fieldName
	}
}

// toCamelCaseLower 转换为小驼峰
func toCamelCaseLower(s string) string {
	parts := strings.Split(s, "_")
	if len(parts) == 0 {
		return s
	}

	result := strings.ToLower(parts[0])
	for i := 1; i < len(parts); i++ {
		if len(parts[i]) > 0 {
			result += strings.ToUpper(parts[i][:1]) + strings.ToLower(parts[i][1:])
		}
	}
	return result
}

// buildFileContent 构建文件内容
func (gc *GeneratorConfig) buildFileContent(modelsCodes []string) string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("package %s\n\n", gc.PackageName))

	// 检查是否需要导入 time 包
	needTime := false
	for _, code := range modelsCodes {
		if strings.Contains(code, "time.Time") {
			needTime = true
			break
		}
	}

	// 导入
	if needTime {
		buf.WriteString("import (\n")
		buf.WriteString("\t\"time\"\n")
		buf.WriteString(")\n\n")
	}

	// 模型代码
	for _, code := range modelsCodes {
		buf.WriteString(code)
		buf.WriteString("\n")
	}

	return buf.String()
}

// FieldInfo 字段信息
type FieldInfo struct {
	Name    string
	Type    string
	Column  string
	Comment string
	Tags    string
}

// toCamelCase 转换为驼峰命名
func toCamelCase(s string) string {
	// 移除下划线并转为驼峰
	parts := strings.Split(s, "_")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
		}
	}
	return strings.Join(parts, "")
}

// mapDBTypeToGoType 映射数据库类型到 Go 类型（使用配置的类型映射）
func (gc *GeneratorConfig) mapDBTypeToGoType(dbType string) string {
	dbTypeLower := strings.ToLower(dbType)

	// 检查是否包含 unsigned
	isUnsigned := strings.Contains(dbTypeLower, "unsigned")

	// 提取类型名（去除长度、unsigned等信息）
	re := regexp.MustCompile(`^(\w+)`)
	matches := re.FindStringSubmatch(dbTypeLower)
	baseType := dbTypeLower
	if len(matches) > 1 {
		baseType = matches[1]
	}

	// 先查找完整类型（包含 unsigned）
	if isUnsigned {
		fullType := baseType + " unsigned"
		if goType, ok := gc.TypeMapping[fullType]; ok {
			return goType
		}
	}

	// 再查找基础类型
	if goType, ok := gc.TypeMapping[baseType]; ok {
		return goType
	}

	// 特殊处理：根据 unsigned 自动转换
	if isUnsigned {
		switch baseType {
		case "tinyint":
			return "uint8"
		case "smallint":
			return "uint16"
		case "mediumint", "int", "integer":
			return "uint32"
		case "bigint":
			return "uint64"
		}
	}

	// 默认返回 interface{}
	return "interface{}"
}

// GenerateModels 快速生成模型（使用默认配置）
func GenerateModels(host string, port int, username, password, database, outDir string) error {
	config := NewGeneratorConfig(host, port, username, password, database)
	config.OutDir = outDir
	return config.Generate()
}

// GenerateModelsFromConfig 从 Config 生成模型
func GenerateModelsFromConfig(dbConfig *Config, outDir string) error {
	config := NewGeneratorConfigFromDBConfig(dbConfig)
	config.OutDir = outDir
	return config.Generate()
}

// GenerateModelsForTables 生成指定表的模型
func GenerateModelsForTables(host string, port int, username, password, database, outDir string, tables []string) error {
	config := NewGeneratorConfig(host, port, username, password, database)
	config.OutDir = outDir
	config.Tables = tables
	return config.Generate()
}

// WithOutDir 设置输出目录
func (gc *GeneratorConfig) WithOutDir(dir string) *GeneratorConfig {
	gc.OutDir = dir
	return gc
}

// WithOutFileName 设置输出文件名
func (gc *GeneratorConfig) WithOutFileName(filename string) *GeneratorConfig {
	gc.OutFileName = filename
	return gc
}

// WithPackageName 设置包名
func (gc *GeneratorConfig) WithPackageName(name string) *GeneratorConfig {
	gc.PackageName = name
	return gc
}

// WithTables 设置要生成的表
func (gc *GeneratorConfig) WithTables(tables ...string) *GeneratorConfig {
	gc.Tables = tables
	return gc
}

// WithTablePrefix 设置表前缀
func (gc *GeneratorConfig) WithTablePrefix(prefix string) *GeneratorConfig {
	gc.TablePrefix = prefix
	return gc
}

// WithTypeMapping 设置自定义类型映射
func (gc *GeneratorConfig) WithTypeMapping(typeMapping map[string]string) *GeneratorConfig {
	if gc.TypeMapping == nil {
		gc.TypeMapping = make(map[string]string)
	}
	for k, v := range typeMapping {
		gc.TypeMapping[k] = v
	}
	return gc
}

// WithJSONTagStyle 设置JSON标签风格
func (gc *GeneratorConfig) WithJSONTagStyle(style string) *GeneratorConfig {
	gc.JSONTagStyle = style
	return gc
}

// WithEnableJSONTag 设置是否启用JSON标签
func (gc *GeneratorConfig) WithEnableJSONTag(enable bool) *GeneratorConfig {
	gc.EnableJSONTag = enable
	return gc
}

// WithEnableGormTag 设置是否启用GORM标签
func (gc *GeneratorConfig) WithEnableGormTag(enable bool) *GeneratorConfig {
	gc.EnableGormTag = enable
	return gc
}

// GetTables 获取数据库中的所有表名
func (gc *GeneratorConfig) GetTables() ([]string, error) {
	// 创建临时客户端
	dbConfig := NewConfig(gc.Host, gc.Port, gc.Username, gc.Password, gc.Database)
	client, err := NewClient(dbConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}
	defer client.Close()

	var tables []string
	rows, err := client.DB().Raw("SHOW TABLES").Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to get tables: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, fmt.Errorf("failed to scan table: %w", err)
		}
		tables = append(tables, table)
	}

	return tables, nil
}

// GetTableColumns 获取表的列信息
func (gc *GeneratorConfig) GetTableColumns(tableName string) ([]ColumnInfo, error) {
	// 创建临时客户端
	dbConfig := NewConfig(gc.Host, gc.Port, gc.Username, gc.Password, gc.Database)
	client, err := NewClient(dbConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}
	defer client.Close()

	var columns []ColumnInfo
	query := fmt.Sprintf("SHOW FULL COLUMNS FROM `%s`", tableName)
	rows, err := client.DB().Raw(query).Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var col ColumnInfo
		var collation, privileges interface{}
		if err := rows.Scan(
			&col.Field,
			&col.Type,
			&collation,
			&col.Null,
			&col.Key,
			&col.Default,
			&col.Extra,
			&privileges,
			&col.Comment,
		); err != nil {
			return nil, fmt.Errorf("failed to scan column: %w", err)
		}
		columns = append(columns, col)
	}

	return columns, nil
}

// GetTableComment 获取表注释
func (gc *GeneratorConfig) GetTableComment(tableName string) (string, error) {
	// 创建临时客户端
	dbConfig := NewConfig(gc.Host, gc.Port, gc.Username, gc.Password, gc.Database)
	client, err := NewClient(dbConfig)
	if err != nil {
		return "", fmt.Errorf("failed to create client: %w", err)
	}
	defer client.Close()

	query := fmt.Sprintf(`
		SELECT TABLE_COMMENT 
		FROM information_schema.TABLES 
		WHERE TABLE_SCHEMA = '%s' AND TABLE_NAME = '%s'
	`, gc.Database, tableName)

	var comment string
	err = client.DB().Raw(query).Scan(&comment).Error
	if err != nil {
		return "", fmt.Errorf("failed to get table comment: %w", err)
	}

	return comment, nil
}

// ColumnInfo 列信息
type ColumnInfo struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default interface{}
	Extra   string
	Comment string
}

// String 格式化输出列信息
func (ci ColumnInfo) String() string {
	var parts []string
	parts = append(parts, fmt.Sprintf("Field: %s", ci.Field))
	parts = append(parts, fmt.Sprintf("Type: %s", ci.Type))
	if ci.Null == "YES" {
		parts = append(parts, "Nullable")
	}
	if ci.Key != "" {
		parts = append(parts, fmt.Sprintf("Key: %s", ci.Key))
	}
	if ci.Default != nil {
		parts = append(parts, fmt.Sprintf("Default: %v", ci.Default))
	}
	if ci.Extra != "" {
		parts = append(parts, fmt.Sprintf("Extra: %s", ci.Extra))
	}
	if ci.Comment != "" {
		parts = append(parts, fmt.Sprintf("Comment: %s", ci.Comment))
	}
	return strings.Join(parts, ", ")
}

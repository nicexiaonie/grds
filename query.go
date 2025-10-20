package grds

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// QueryBuilder 查询构建器，封装 GORM DB
type QueryBuilder struct {
	client *Client
	db     *gorm.DB
}

// DB 获取底层的 GORM DB
func (qb *QueryBuilder) DB() *gorm.DB {
	return qb.db
}

// Client 获取客户端
func (qb *QueryBuilder) Client() *Client {
	return qb.client
}

// ==================== 条件查询 ====================

// Where 添加 WHERE 条件
func (qb *QueryBuilder) Where(query interface{}, args ...interface{}) *QueryBuilder {
	qb.db = qb.db.Where(query, args...)
	return qb
}

// Not 添加 NOT 条件
func (qb *QueryBuilder) Not(query interface{}, args ...interface{}) *QueryBuilder {
	qb.db = qb.db.Not(query, args...)
	return qb
}

// Or 添加 OR 条件
func (qb *QueryBuilder) Or(query interface{}, args ...interface{}) *QueryBuilder {
	qb.db = qb.db.Or(query, args...)
	return qb
}

// ==================== 排序 ====================

// Order 排序
func (qb *QueryBuilder) Order(value interface{}) *QueryBuilder {
	qb.db = qb.db.Order(value)
	return qb
}

// OrderBy 排序（别名）
func (qb *QueryBuilder) OrderBy(column string) *QueryBuilder {
	return qb.Order(column)
}

// OrderByAsc 升序排序
func (qb *QueryBuilder) OrderByAsc(column string) *QueryBuilder {
	return qb.Order(column + " ASC")
}

// OrderByDesc 降序排序
func (qb *QueryBuilder) OrderByDesc(column string) *QueryBuilder {
	return qb.Order(column + " DESC")
}

// ==================== 分组 ====================

// GroupBy 分组
func (qb *QueryBuilder) GroupBy(name string) *QueryBuilder {
	qb.db = qb.db.Group(name)
	return qb
}

// Group 分组（GORM 原生方法）
func (qb *QueryBuilder) Group(name string) *QueryBuilder {
	qb.db = qb.db.Group(name)
	return qb
}

// Having HAVING 条件
func (qb *QueryBuilder) Having(query interface{}, args ...interface{}) *QueryBuilder {
	qb.db = qb.db.Having(query, args...)
	return qb
}

// ==================== 限制和偏移 ====================

// Limit 限制数量
func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.db = qb.db.Limit(limit)
	return qb
}

// Offset 偏移量
func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.db = qb.db.Offset(offset)
	return qb
}

// Page 分页（page 从 1 开始）
func (qb *QueryBuilder) Page(page, pageSize int) *QueryBuilder {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize
	return qb.Limit(pageSize).Offset(offset)
}

// ==================== 连接查询 ====================

// Joins 连接查询
func (qb *QueryBuilder) Joins(query string, args ...interface{}) *QueryBuilder {
	qb.db = qb.db.Joins(query, args...)
	return qb
}

// LeftJoin 左连接
func (qb *QueryBuilder) LeftJoin(tableName string, condition string) *QueryBuilder {
	return qb.Joins("LEFT JOIN " + tableName + " ON " + condition)
}

// RightJoin 右连接
func (qb *QueryBuilder) RightJoin(tableName string, condition string) *QueryBuilder {
	return qb.Joins("RIGHT JOIN " + tableName + " ON " + condition)
}

// InnerJoin 内连接
func (qb *QueryBuilder) InnerJoin(tableName string, condition string) *QueryBuilder {
	return qb.Joins("INNER JOIN " + tableName + " ON " + condition)
}

// ==================== 选择字段 ====================

// Select 选择字段
func (qb *QueryBuilder) Select(query interface{}, args ...interface{}) *QueryBuilder {
	qb.db = qb.db.Select(query, args...)
	return qb
}

// Omit 忽略字段
func (qb *QueryBuilder) Omit(columns ...string) *QueryBuilder {
	qb.db = qb.db.Omit(columns...)
	return qb
}

// Distinct 去重
func (qb *QueryBuilder) Distinct(args ...interface{}) *QueryBuilder {
	qb.db = qb.db.Distinct(args...)
	return qb
}

// ==================== 预加载 ====================

// Preload 预加载关联
func (qb *QueryBuilder) Preload(query string, args ...interface{}) *QueryBuilder {
	qb.db = qb.db.Preload(query, args...)
	return qb
}

// Clauses 添加子句
func (qb *QueryBuilder) Clauses(conds ...clause.Expression) *QueryBuilder {
	qb.db = qb.db.Clauses(conds...)
	return qb
}

// ==================== 锁 ====================

// ForUpdate 行锁（SELECT ... FOR UPDATE）
func (qb *QueryBuilder) ForUpdate() *QueryBuilder {
	qb.db = qb.db.Clauses(clause.Locking{Strength: "UPDATE"})
	return qb
}

// ForShare 共享锁（SELECT ... FOR SHARE）
func (qb *QueryBuilder) ForShare() *QueryBuilder {
	qb.db = qb.db.Clauses(clause.Locking{Strength: "SHARE"})
	return qb
}

// ==================== 查询操作 ====================

// Find 查询多条记录
func (qb *QueryBuilder) Find(dest interface{}, conds ...interface{}) error {
	return qb.db.Find(dest, conds...).Error
}

// First 查询第一条记录
func (qb *QueryBuilder) First(dest interface{}, conds ...interface{}) error {
	return qb.db.First(dest, conds...).Error
}

// Last 查询最后一条记录
func (qb *QueryBuilder) Last(dest interface{}, conds ...interface{}) error {
	return qb.db.Last(dest, conds...).Error
}

// Take 随机获取一条记录
func (qb *QueryBuilder) Take(dest interface{}, conds ...interface{}) error {
	return qb.db.Take(dest, conds...).Error
}

// Scan 扫描结果到目标
func (qb *QueryBuilder) Scan(dest interface{}) error {
	return qb.db.Scan(dest).Error
}

// Pluck 查询单列
func (qb *QueryBuilder) Pluck(column string, dest interface{}) error {
	return qb.db.Pluck(column, dest).Error
}

// Count 统计数量
func (qb *QueryBuilder) Count() (int64, error) {
	var count int64
	err := qb.db.Count(&count).Error
	return count, err
}

// Exists 检查是否存在
func (qb *QueryBuilder) Exists() (bool, error) {
	count, err := qb.Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ==================== 创建操作 ====================

// Create 创建记录
func (qb *QueryBuilder) Create(value interface{}) error {
	return qb.db.Create(value).Error
}

// CreateInBatches 批量创建
func (qb *QueryBuilder) CreateInBatches(value interface{}, batchSize int) error {
	return qb.db.CreateInBatches(value, batchSize).Error
}

// ==================== 更新操作 ====================

// Update 更新单个字段
func (qb *QueryBuilder) Update(column string, value interface{}) error {
	return qb.db.Update(column, value).Error
}

// Updates 更新多个字段
func (qb *QueryBuilder) Updates(values interface{}) error {
	return qb.db.Updates(values).Error
}

// UpdateColumn 更新单列（不触发钩子）
func (qb *QueryBuilder) UpdateColumn(column string, value interface{}) error {
	return qb.db.UpdateColumn(column, value).Error
}

// UpdateColumns 更新多列（不触发钩子）
func (qb *QueryBuilder) UpdateColumns(values interface{}) error {
	return qb.db.UpdateColumns(values).Error
}

// Save 保存所有字段
func (qb *QueryBuilder) Save(value interface{}) error {
	return qb.db.Save(value).Error
}

// ==================== 删除操作 ====================

// Delete 删除记录
func (qb *QueryBuilder) Delete(value interface{}, conds ...interface{}) error {
	return qb.db.Delete(value, conds...).Error
}

// ==================== 聚合函数 ====================

// Sum 求和
func (qb *QueryBuilder) Sum(column string) (float64, error) {
	var result float64
	err := qb.db.Select("SUM(" + column + ")").Scan(&result).Error
	return result, err
}

// Avg 平均值
func (qb *QueryBuilder) Avg(column string) (float64, error) {
	var result float64
	err := qb.db.Select("AVG(" + column + ")").Scan(&result).Error
	return result, err
}

// Max 最大值
func (qb *QueryBuilder) Max(column string) (interface{}, error) {
	var result interface{}
	err := qb.db.Select("MAX(" + column + ")").Scan(&result).Error
	return result, err
}

// Min 最小值
func (qb *QueryBuilder) Min(column string) (interface{}, error) {
	var result interface{}
	err := qb.db.Select("MIN(" + column + ")").Scan(&result).Error
	return result, err
}

// ==================== 其他操作 ====================

// Raw 原生 SQL 查询
func (qb *QueryBuilder) Raw(sql string, values ...interface{}) *QueryBuilder {
	qb.db = qb.db.Raw(sql, values...)
	return qb
}

// Exec 执行 SQL
func (qb *QueryBuilder) Exec(sql string, values ...interface{}) error {
	return qb.db.Exec(sql, values...).Error
}

// Model 指定模型
func (qb *QueryBuilder) Model(value interface{}) *QueryBuilder {
	qb.db = qb.db.Model(value)
	return qb
}

// Table 指定表名
func (qb *QueryBuilder) Table(name string, args ...interface{}) *QueryBuilder {
	qb.db = qb.db.Table(name, args...)
	return qb
}

// Session 创建新会话
func (qb *QueryBuilder) Session(config *gorm.Session) *QueryBuilder {
	qb.db = qb.db.Session(config)
	return qb
}

// Scopes 应用作用域
func (qb *QueryBuilder) Scopes(funcs ...func(*gorm.DB) *gorm.DB) *QueryBuilder {
	qb.db = qb.db.Scopes(funcs...)
	return qb
}

// Debug 开启调试模式
func (qb *QueryBuilder) Debug() *QueryBuilder {
	qb.db = qb.db.Debug()
	return qb
}

// Error 获取错误
func (qb *QueryBuilder) Error() error {
	return qb.db.Error
}

// RowsAffected 获取影响的行数
func (qb *QueryBuilder) RowsAffected() int64 {
	return qb.db.RowsAffected
}

// Clone 克隆查询构建器
func (qb *QueryBuilder) Clone() *QueryBuilder {
	return &QueryBuilder{
		client: qb.client,
		db:     qb.db.Session(&gorm.Session{NewDB: true}),
	}
}

// ==================== 便捷方法 ====================

// FindOne 查询单条（First 的便捷方法）
func (qb *QueryBuilder) FindOne(dest interface{}) error {
	return qb.First(dest)
}

// FindAll 查询所有（Find 的便捷方法）
func (qb *QueryBuilder) FindAll(dest interface{}) error {
	return qb.Find(dest)
}

// WhereEq 等于条件
func (qb *QueryBuilder) WhereEq(column string, value interface{}) *QueryBuilder {
	return qb.Where(column+" = ?", value)
}

// WhereNe 不等于条件
func (qb *QueryBuilder) WhereNe(column string, value interface{}) *QueryBuilder {
	return qb.Where(column+" != ?", value)
}

// WhereGt 大于条件
func (qb *QueryBuilder) WhereGt(column string, value interface{}) *QueryBuilder {
	return qb.Where(column+" > ?", value)
}

// WhereGte 大于等于条件
func (qb *QueryBuilder) WhereGte(column string, value interface{}) *QueryBuilder {
	return qb.Where(column+" >= ?", value)
}

// WhereLt 小于条件
func (qb *QueryBuilder) WhereLt(column string, value interface{}) *QueryBuilder {
	return qb.Where(column+" < ?", value)
}

// WhereLte 小于等于条件
func (qb *QueryBuilder) WhereLte(column string, value interface{}) *QueryBuilder {
	return qb.Where(column+" <= ?", value)
}

// WhereLike LIKE 条件
func (qb *QueryBuilder) WhereLike(column string, value string) *QueryBuilder {
	return qb.Where(column+" LIKE ?", value)
}

// WhereIn IN 条件
func (qb *QueryBuilder) WhereIn(column string, values interface{}) *QueryBuilder {
	return qb.Where(column+" IN ?", values)
}

// WhereNotIn NOT IN 条件
func (qb *QueryBuilder) WhereNotIn(column string, values interface{}) *QueryBuilder {
	return qb.Where(column+" NOT IN ?", values)
}

// WhereBetween BETWEEN 条件
func (qb *QueryBuilder) WhereBetween(column string, start, end interface{}) *QueryBuilder {
	return qb.Where(column+" BETWEEN ? AND ?", start, end)
}

// WhereNull IS NULL 条件
func (qb *QueryBuilder) WhereNull(column string) *QueryBuilder {
	return qb.Where(column + " IS NULL")
}

// WhereNotNull IS NOT NULL 条件
func (qb *QueryBuilder) WhereNotNull(column string) *QueryBuilder {
	return qb.Where(column + " IS NOT NULL")
}

package grds

import (
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type DbConfig struct {
	// 连接方式 tcp:
	Net string
	// 地址  0.0.0.0:3306
	Host string
	// 用户名
	User string
	// 密码
	Passwd string
	// 数据库
	DBName string

	MaxIdleConns int
	MaxOpenConns int
	LogMode      bool
}
type Client struct {
	Config *DbConfig
	DB     *gorm.DB
}

type tableModelInterface interface {
	TableName() string
}
type Table struct {
	Database Client
	model    tableModelInterface
	Handle   *gorm.DB
}

func NewDb(c *DbConfig) (Client, error) {
	r := Client{
		Config: c,
	}
	config := mysql.NewConfig()
	config.Net = c.Net
	config.Addr = c.Host
	config.User = c.User
	config.Passwd = c.Passwd
	config.DBName = c.DBName
	DSN := config.FormatDSN()

	db, err := gorm.Open("mysql", DSN)
	if err != nil {
		return r, errors.New("mysql初始化失败")
	}

	db.DB().SetMaxIdleConns(c.MaxIdleConns)
	db.DB().SetMaxOpenConns(c.MaxOpenConns)
	db.LogMode(c.LogMode)
	r.DB = db

	return r, nil
}

func (current *Client) NewTable(model tableModelInterface) *Table {
	r := Table{
		Database: *current,
		model:    model,
		Handle:   current.DB,
	}
	return &r
}
func (t *Table) Clone() *Table {
	var r = &Table{
		Handle:   t.Database.DB,
		model:    t.model,
		Database: t.Database,
	}
	*r.Handle = *t.Handle
	r.Handle = r.Handle.Table(r.model.TableName())
	return r
}
func (t *Table) Where(query interface{}, args ...interface{}) *Table {
	t.Handle = t.Handle.Where(query, args)
	return t
}
func (t *Table) Count() int64 {
	var r int64
	t.Handle.Count(&r)
	return r
}
func (t *Table) Find(out interface{}) {
	t.Handle.Find(out)
}
func (t *Table) Limit(limit interface{}) {
	t.Handle.Limit(limit)
}
func (t *Table) Offset(offset interface{}) {
	t.Handle.Offset(offset)
}

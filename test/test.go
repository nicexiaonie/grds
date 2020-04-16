package main

import (
	"github.com/nicexiaonie/grds"
)

type DDAlbumAuditLogModel struct {
	Id        int64 `gorm:"id"`
	ContentId int64 `gorm:"content_id"`
}

func (t *DDAlbumAuditLogModel) TableName() string {
	return "dandan_album_audit_log"
}

type LogDao struct {
	Table *grds.Table

}
func main() {

	c := &grds.DbConfig{
		Key:          "11",
		Net:          "tcp",
		Host:         "172.25.20.233:8306",
		User:         "qukan",
		Passwd:       "iDLvO3yQGv",
		DBName:       "audit",
		MaxIdleConns: 100,
		MaxOpenConns: 200,
		LogMode:      true,
	}

	db, _ := grds.NewDb(c)


	dao := LogDao{
		Table : db.NewTable(&DDAlbumAuditLogModel{}),
	}

	dao.Table.Clone().Handle.Row()




	// 测试变量污染
	//go func() {
	//
	//	t := table.Table().Where("id > 3")
	//	time.Sleep(time.Second * 10)
	//	count := t.Where("id < 100").Count()
	//	fmt.Printf("10s: %d", count)
	//
	//}()
	//
	//go func() {
	//	time.Sleep(time.Second * 2)
	//	t := table.Table().Where("id > 66")
	//	count := t.Where("id < 100").Count()
	//	fmt.Printf("step: %d", count)
	//}()

	select {}

}

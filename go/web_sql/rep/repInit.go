package rep

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	username := "root"  //账号
	password := "root"  //密码
	host := "localhost" //数据库地址
	port := "3306"      //端口
	dbname := "shop"    //数据库名

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local", username, password, host, port, dbname) //连接mysql，获得DB类型实例，用于后面数据库的读写操作
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic("连接数据库失败，error=" + err.Error())
	}
	DB = db
	//连接成功
	fmt.Println("连接数据库成功")
	err = db.AutoMigrate()
	if err != nil {
		panic("自动迁移失败，error=" + err.Error())
	}
	fmt.Println("自动迁移成功")
}

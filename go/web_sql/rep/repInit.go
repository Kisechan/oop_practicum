package rep

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 定义一个全局变量db，用于后面数据库的读写操作,通常就放在全局里面
var DB *gorm.DB

func init() {
	username := "root"  //账号
	password := "root"  //密码
	host := "localhost" //数据库地址
	port := "3306"      //端口
	Dnname := "shop"    //数据库名
	timeout := "10s"    //连接超时，10s

	//root:root@tcp(127.0.0.1:3306)/test？
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?timeout=%s", username, password, host, port, Dnname, timeout)
	//连接mysql，获得DB类型实例，用于后面数据库的读写操作
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic("连接数据库失败，error=" + err.Error())
	}
	DB = db
	//连接成功
	fmt.Println("连接数据库成功")
}

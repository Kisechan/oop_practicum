package main

import (
	"web_sql/control"
	"web_sql/route"
	"web_sql/ui"
)

func main() {
	control.RedisInit()
	go route.APIInit()
	go control.RepAPIInit()
	// go func() {
	// 	time.Sleep(3 * time.Second)
	// 	fmt.Println("CodE Dream! \nIt's My GO!!!!!")
	// }()
	ui.Show()
}

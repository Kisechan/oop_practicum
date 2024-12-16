package main

import (
	"web_sql/route"
	"web_sql/ui"
)

func main() {
	go route.APIInit()
	// go service.WSInit()
	ui.Show()
}

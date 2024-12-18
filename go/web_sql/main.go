package main

import (
	"web_sql/rep"
	"web_sql/route"
	"web_sql/ui"
)

func main() {
	go route.APIInit()
	go rep.RepAPIInit()
	ui.Show()
}

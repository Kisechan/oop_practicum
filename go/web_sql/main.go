package main

import (
	"fmt"
	"web_sql/route"

	"fyne.io/fyne/v2/app"
	"github.com/gin-gonic/gin"
)

func initWeb() {
	router := gin.Default()
	route.SetupUserRoutes(router)
	if err := router.Run(":8080"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Server started on port 8080")
	}
}
func main() {
	initWeb()
	a := app.New()
	a.Icon()
}

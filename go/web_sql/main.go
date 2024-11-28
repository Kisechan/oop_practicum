package main

import (
	"fmt"
	"web_sql/route"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	route.SetupUserRoutes(router)
	if err := router.Run(":8080"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Server started on port 8080")
	}
}

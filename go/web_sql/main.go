package main

import (
	"web_sql/route"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	route.SetupRoutes(r)
	r.Run(":8080")
}

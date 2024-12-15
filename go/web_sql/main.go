package main

import (
	"web_sql/ui"
)

func main() {
	// r := gin.Default()
	// route.SetupRoutes(r)
	// r.Run(":8080")
	ui.Show()
}

// package main

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"
// 	"time"

// 	"github.com/gin-gonic/gin"
// )

// func main() {
// 	r := gin.Default()

// 	// Go端作为HTTP服务器，处理C++端的请求
// 	r.GET("/go-ping", func(c *gin.Context) {
// 		c.JSON(http.StatusOK, gin.H{
// 			"message": "Pong from Go",
// 		})
// 	})

// 	// Go端作为HTTP客户端，向C++端发送请求
// 	go func() {
// 		for {
// 			// 发送GET请求到C++端
// 			resp, err := http.Get("http://localhost:8081/cpp-ping")
// 			if err != nil {
// 				fmt.Println("Error:", err)
// 				continue
// 			}
// 			defer resp.Body.Close()

// 			body, _ := ioutil.ReadAll(resp.Body)
// 			fmt.Println("Response from C++ server:", string(body))

// 			// 每隔5秒发送一次请求
// 			time.Sleep(5 * time.Second)
// 		}
// 	}()

// 	r.Run(":8080")
// }

package service

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"sync"

// 	"github.com/gorilla/websocket"
// )

// type Request struct {
// 	Type    string          `json:"type"`
// 	Payload json.RawMessage `json:"payload"`
// }

// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool { return true },
// }

// var (
// 	requestChan = make(chan Request, 100)
// 	mutex       sync.Mutex // 互斥锁，保护 WebSocket 连接
// )

// func WSInit() {
// 	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
// 		// 升级 HTTP 连接为 WebSocket 连接
// 		conn, err := upgrader.Upgrade(w, r, nil)
// 		if err != nil {
// 			fmt.Println("Upgrade error:", err)
// 			return
// 		}
// 		defer conn.Close()

// 		fmt.Println("C++ client connected")

// 		// 启动一个 goroutine 来处理管道中的数据
// 		go func() {
// 			for {
// 				// 从管道中读取数据
// 				request := <-requestChan
// 				mutex.Lock() // 加锁
// 				// 将请求转换为 JSON 格式
// 				requestJSON, err := json.Marshal(request)
// 				if err != nil {
// 					fmt.Println("JSON marshal error:", err)
// 					mutex.Unlock() // 解锁
// 					continue
// 				}

// 				// 通过 WebSocket 发送数据给 C++ 端
// 				err = conn.WriteMessage(websocket.TextMessage, requestJSON)
// 				if err != nil {
// 					fmt.Println("Write error:", err)
// 					mutex.Unlock() // 解锁
// 					return
// 				}
// 				fmt.Println("Sent to C++:", string(requestJSON))
// 				mutex.Unlock() // 解锁
// 			}
// 		}()

// 		// 主线程负责接收 C++ 端的消息
// 		for {
// 			_, msg, err := conn.ReadMessage()
// 			if err != nil {
// 				fmt.Println("Read error:", err)
// 				break
// 			}
// 			fmt.Println("Received from C++:", string(msg))

// 			var req Request
// 			err = json.Unmarshal(msg, &req)
// 			if err != nil {
// 				fmt.Println("JSON unmarshal error:", err)
// 				continue
// 			}

// 			// 处理请求并发送响应
// 			resp, err := ReqHandler[req.Type](req)
// 			if err != nil {
// 				fmt.Println("ReqHandle error:", err)
// 				return
// 			}

// 			if resp != nil {
// 				requestChan <- Request{
// 					Type:    "Info",
// 					Payload: resp,
// 				}
// 			}
// 		}
// 	})

// 	// 启动 WebSocket 服务端
// 	fmt.Println("WebSocket server started at :8080")
// 	if err := http.ListenAndServe(":8080", nil); err != nil {
// 		fmt.Println("WebSocket server error:", err)
// 	}
// }

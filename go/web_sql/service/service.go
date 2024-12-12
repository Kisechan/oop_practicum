package service

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
)

// 定义请求和响应结构体
type Request struct {
	Action    string `json:"action"`
	UserID    int    `json:"user_id,omitempty"`
	OrderID   int    `json:"order_id,omitempty"`
	ProductID int    `json:"product_id,omitempty"`
	Quantity  int    `json:"quantity,omitempty"`
}

type Response struct {
	Status   int         `json:"status"`
	Message  string      `json:"message,omitempty"`
	UserInfo interface{} `json:"user_info,omitempty"`
}

func SocketInit() {
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal("Socket启动失败:", err.Error())
	}
	defer listener.Close()

	log.Println("Socket服务已启动，监听端口: 8081")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("连接接受失败:", err)
			continue
		}
		go handle(conn)
	}
}

// handle 处理Socket连接
func handle(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	message, err := reader.ReadString('\n')
	if err != nil {
		log.Println("读取数据失败:", err)
		return
	}

	// 解析请求
	var req Request
	if err := json.Unmarshal([]byte(message), &req); err != nil {
		log.Println("解析请求失败:", err)
		return
	}

	// 处理请求
	var res Response
	switch req.Action {
	case "get_user_info":
		res = handleGetUserInfo(req)
	case "submit_order":
		res = handleSubmitOrder(req)
	default:
		res = Response{Status: 400, Message: "无效的操作类型"}
	}

	// 返回响应
	resBytes, _ := json.Marshal(res)
	conn.Write(append(resBytes, '\n'))
}

// handleGetUserInfo 处理查询用户信息请求
func handleGetUserInfo(req Request) Response {
	// 模拟从数据库获取用户信息
	userInfo := map[string]string{
		"name":  "Alice",
		"email": "alice@example.com",
	}
	return Response{Status: 200, UserInfo: userInfo}
}

// handleSubmitOrder 处理提交订单请求
func handleSubmitOrder(req Request) Response {
	// 模拟保存订单到数据库
	return Response{Status: 200, Message: "订单提交成功"}
}

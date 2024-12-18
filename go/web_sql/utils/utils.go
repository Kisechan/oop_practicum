package utils

// import (
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"net/http"
// 	"net/url"
// )

// // 调用 C++ 程序的接口
// func CallCppAPI(method string, params ...interface{}) (interface{}, error) {
// 	// 构建请求 URL
// 	baseURL := "http://localhost:8000/api" // C++ 程序的 API 地址
// 	u, err := url.Parse(baseURL)
// 	if err != nil {
// 		return nil, err
// 	}
// 	u.Path = method

// 	// 构建请求参数
// 	var queryParams url.Values
// 	if len(params) > 0 {
// 		queryParams = make(url.Values)
// 		for i, param := range params {
// 			queryParams.Add(fmt.Sprintf("param%d", i), fmt.Sprintf("%v", param))
// 		}
// 	}
// 	u.RawQuery = queryParams.Encode()

// 	// 发送 HTTP 请求
// 	resp, err := http.Get(u.String())
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	// 解析响应
// 	var result interface{}
// 	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
// 		return nil, err
// 	}

// 	// 检查响应状态
// 	if resp.StatusCode != http.StatusOK {
// 		return nil, errors.New(fmt.Sprintf("C++ API 返回错误: %v", result))
// 	}

// 	return result, nil
// }

#include <iostream>
#include <string>
#include <boost/beast/core.hpp>
#include <boost/beast/http.hpp>
#include <boost/beast/version.hpp>
#include <boost/asio/connect.hpp>
#include <boost/asio/ip/tcp.hpp>
#include <hiredis/hiredis.h>
#include <boost/json.hpp>

namespace beast = boost::beast;         // 简化命名空间
namespace http = beast::http;           // 简化命名空间
namespace net = boost::asio;            // 简化命名空间
using tcp = net::ip::tcp;               // 简化命名空间

// 检查库存是否充足
bool checkInventory(int productId, int quantity) {
    // 模拟库存检查（实际场景中需要从数据库或缓存中查询）
    std::cout << "Checking inventory for product " << productId << " with quantity " << quantity << std::endl;
    return false; // 假设库存不足
}

// 发送秒杀结果到 Go 端
void sendSeckillResult(int orderId, bool success, const std::string& message) {
    try {
        // 创建 I/O 上下文
        net::io_context ioc;

        // 解析域名
        tcp::resolver resolver(ioc);
        auto const results = resolver.resolve("localhost", "8080"); // Go 端的地址

        // 创建连接
        tcp::socket socket(ioc);
        net::connect(socket, results.begin(), results.end());

        // 构建 HTTP POST 请求
        http::request<http::string_body> req{ http::verb::post, "/api/seckill/result", 11 };
        req.set(http::field::host, "localhost");
        req.set(http::field::user_agent, BOOST_BEAST_VERSION_STRING);
        req.set(http::field::content_type, "application/json");
        req.body() = "{\"order_id\": " + std::to_string(orderId) + ", \"status\": \"" + (success ? "success" : "failed") + "\", \"message\": \"" + message + "\"}";
        req.prepare_payload();

        // 发送请求
        http::write(socket, req);

        // 读取响应
        beast::flat_buffer buffer;
        http::response<http::dynamic_body> res;
        http::read(socket, buffer, res);

        // 输出响应
        std::cout << "Response: " << res << std::endl;

        // 关闭连接
        socket.shutdown(tcp::socket::shutdown_both);
    }
    catch (std::exception const& e) {
        std::cerr << "Error: " << e.what() << std::endl;
    }
}

// 创建订单并调用持久层
void createOrder(const std::string& orderJson) {
    // 解析 JSON 消息
    boost::json::value json = boost::json::parse(orderJson);
    int orderId = json.at("order_id").as_int64();
    int userId = json.at("user_id").as_int64();
    int productId = json.at("product_id").as_int64();
    int quantity = json.at("quantity").as_int64();

    std::cout << "Processing seckill order: " << orderId << std::endl;

    // 检查库存
    if (!checkInventory(productId, quantity)) {
        std::cerr << "Insufficient inventory for product " << productId << std::endl;
        sendSeckillResult(orderId, false, "Insufficient inventory");
        return;
    }

    // 构建订单创建请求
    try {
        // 创建 I/O 上下文
        net::io_context ioc;

        // 解析域名
        tcp::resolver resolver(ioc);
        auto const results = resolver.resolve("localhost", "8081"); // 持久层的地址

        // 创建连接
        tcp::socket socket(ioc);
        net::connect(socket, results.begin(), results.end());

        // 构建 HTTP POST 请求
        http::request<http::string_body> req{ http::verb::post, "/api/orders/create", 11 };
        req.set(http::field::host, "localhost");
        req.set(http::field::user_agent, BOOST_BEAST_VERSION_STRING);
        req.set(http::field::content_type, "application/json");
        req.body() = orderJson;
        req.prepare_payload();

        // 发送请求
        http::write(socket, req);

        // 读取响应
        beast::flat_buffer buffer;
        http::response<http::dynamic_body> res;
        http::read(socket, buffer, res);

        // 输出响应
        std::cout << "Response: " << res << std::endl;

        // 关闭连接
        socket.shutdown(tcp::socket::shutdown_both);

        // 发送秒杀成功结果
        sendSeckillResult(orderId, true, "Order created successfully");
    }
    catch (std::exception const& e) {
        std::cerr << "Error: " << e.what() << std::endl;
        sendSeckillResult(orderId, false, "Failed to create order");
    }
}

// 从 Redis 消费消息
void consumeRedisMessages(redisContext* context) {
    redisReply* reply;
    while (true) {
        reply = (redisReply*)redisCommand(context, "BLPOP seckill_queue 0");
        if (reply && reply->type == REDIS_REPLY_ARRAY && reply->element[1]) {
            std::string message(reply->element[1]->str, reply->element[1]->len);
            std::cout << "Received message: " << message << std::endl;

            // 处理秒杀订单
            createOrder(message);
        }
        freeReplyObject(reply);
    }
}

int main() {
    std::cout << "Start Test" << std::endl;
    // 连接 Redis
    redisContext* context = redisConnect("127.0.0.1", 6379);
    if (context == NULL || context->err) {
        if (context) {
            std::cerr << "Redis connection error: " << context->errstr << std::endl;
        }
        else {
            std::cerr << "Redis connection error: cannot allocate redis context" << std::endl;
        }
        return 1;
    }

    std::cout << "Connected to Redis!" << std::endl;

    // 开始消费消息
    consumeRedisMessages(context);

    // 关闭 Redis 连接
    redisFree(context);

    return 0;
}
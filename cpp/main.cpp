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
    return true; // 假设库存充足
}

// 发送秒杀结果到 Go 端
void sendSeckillResult(const std::string& orderId, bool success, const std::string& message) {
    try {
        std::cout << "\nStart Sending Seckill Results\n" << std::endl;
        // 创建 I/O 上下文
        net::io_context ioc;

        // 解析域名
        tcp::resolver resolver(ioc);
        auto const results = resolver.resolve("localhost", "8080"); // Go 端的地址

        // 创建连接
        tcp::socket socket(ioc);
        net::connect(socket, results.begin(), results.end());

        // 构建 HTTP POST 请求
        http::request<http::string_body> req{ http::verb::post, "/orders/checkout/result", 11 };
        req.set(http::field::host, "localhost");
        req.set(http::field::user_agent, BOOST_BEAST_VERSION_STRING);
        req.set(http::field::content_type, "application/json");
        req.body() = "{\"order_id\": \"" + orderId + "\", \"status\": \"" + (success ? "success" : "failed") + "\", \"message\": \"" + message + "\"}";
        req.prepare_payload();

        // 发送请求
        http::write(socket, req);

        // 读取响应
        beast::flat_buffer buffer;
        http::response<http::dynamic_body> res;
        http::read(socket, buffer, res);

        // 输出响应
        std::cout << "Sending Results Response: " << res << std::endl;

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
    std::string orderId = json.at("order_id").as_string().c_str();
    int userId = json.at("user_id").as_int64();
    int productId = json.at("product_id").as_int64();
    int quantity = json.at("quantity").as_int64();
    double total = json.at("total").as_double();

    std::cout << "Processing seckill order: " << orderId << std::endl;

    // 检查库存
    if (!checkInventory(productId, quantity)) {
        std::cerr << "Insufficient inventory for product " << productId << std::endl;
        sendSeckillResult(orderId, false, "Insufficient inventory");
        return;
    }
    std::cout << "\nStart Making Order Requests\n" << std::endl;
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
        req.body() = "{\"id\": \"" + orderId + "\", \"user_id\": " + std::to_string(userId) + ", \"product_id\": " + std::to_string(productId) + ", \"quantity\": " + std::to_string(quantity) + ", \"total\": " + std::to_string(total) + ", \"status\": \"pending\"}";
        req.prepare_payload();

        // 发送请求
        http::write(socket, req);

        // 读取响应
        beast::flat_buffer buffer;
        http::response<http::dynamic_body> res;
        http::read(socket, buffer, res);

        // 输出响应
        std::cout << "Make Order Response: " << res << std::endl;

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
            if (message == "") 
            {
                std::cout << "Message is null, skipped" << std::endl;
                continue;
            }
            // 处理秒杀订单
            std::cout << "\nStart Resolving Seckill Requests\n" << std::endl;
            createOrder(message);
        }
        freeReplyObject(reply);
    }
}

int main() {
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

    std::cout << "Connected to Redis!Let's Start!" << std::endl;

    // 开始消费消息
    consumeRedisMessages(context);

    // 关闭 Redis 连接
    redisFree(context);

    return 0;
}
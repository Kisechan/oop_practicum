#include "inventory.h"
#include <boost/beast/core.hpp>
#include <boost/beast/http.hpp>
#include <boost/beast/version.hpp>
#include <boost/asio/connect.hpp>
#include <boost/asio/ip/tcp.hpp>
#include <boost/json.hpp>
#include <iostream>
#include <string>

namespace beast = boost::beast;
namespace http = beast::http;
namespace net = boost::asio;
namespace json = boost::json;
using tcp = net::ip::tcp;

// 从 API 获取商品库存数据
std::vector<std::pair<int, int>> fetchInventoryFromDatabase() {
    std::vector<std::pair<int, int>> inventoryData;

    try {
        // 创建 I/O 上下文
        net::io_context ioc;

        // 解析域名
        tcp::resolver resolver(ioc);
        auto const results = resolver.resolve("localhost", "8081"); // API 地址

        // 创建连接
        tcp::socket socket(ioc);
        net::connect(socket, results.begin(), results.end());

        // 构建 HTTP GET 请求
        http::request<http::string_body> req{ http::verb::get, "/api/stock", 11 };
        req.set(http::field::host, "localhost");
        req.set(http::field::user_agent, BOOST_BEAST_VERSION_STRING);

        // 发送请求
        http::write(socket, req);

        // 读取响应
        beast::flat_buffer buffer;
        http::response<http::dynamic_body> res;
        http::read(socket, buffer, res);

        // 解析 JSON 响应
        std::string responseBody = beast::buffers_to_string(res.body().data());
        json::value jsonResponse = json::parse(responseBody);

        // 提取库存数据
        for (const auto& item : jsonResponse.as_array()) {
            int productId = item.at("id").as_int64();
            int stock = item.at("stock").as_int64();
            inventoryData.push_back({ productId, stock });
        }

        // 关闭连接
        socket.shutdown(tcp::socket::shutdown_both);
    }
    catch (const std::exception& e) {
        std::cerr << "Error fetching inventory from API: " << e.what() << std::endl;
    }

    return inventoryData;
}

// 检查库存是否充足
bool checkInventory(redisContext* context, int productId, int quantity) {
    std::string key = "inventory:" + std::to_string(productId);

    // 使用 DECRBY 原子性地减少库存
    redisReply* reply = (redisReply*)redisCommand(context, "DECRBY %s %d", key.c_str(), quantity);
    if (reply == nullptr || reply->type == REDIS_REPLY_ERROR) {
        std::cerr << "Redis error: " << (reply ? reply->str : "Unknown error") << std::endl;
        freeReplyObject(reply);
        return false;
    }

    // 检查库存是否充足
    int remainingStock = reply->integer;
    freeReplyObject(reply);

    if (remainingStock >= 0) {
        std::cout << "Inventory check passed for product " << productId << ". Remaining stock: " << remainingStock << std::endl;
        return true;
    }
    else {
        // 库存不足，恢复库存
        redisReply* restoreReply = (redisReply*)redisCommand(context, "INCRBY %s %d", key.c_str(), quantity);
        if (restoreReply == nullptr || restoreReply->type == REDIS_REPLY_ERROR) {
            std::cerr << "Failed to restore inventory for product " << productId << std::endl;
        }
        freeReplyObject(restoreReply);
        std::cerr << "Insufficient inventory for product " << productId << std::endl;
        return false;
    }
}

// 初始化 Redis 库存
void initializeInventory(redisContext* context) {
    // 从数据库或 API 中获取商品库存数据
    std::vector<std::pair<int, int>> inventoryData = fetchInventoryFromDatabase();

    // 将库存数据写入 Redis
    for (const auto& item : inventoryData) {
        int productId = item.first;
        int stock = item.second;
        std::string key = "inventory:" + std::to_string(productId);

        // 设置库存值
        redisReply* reply = (redisReply*)redisCommand(context, "SET %s %d", key.c_str(), stock);
        if (reply == nullptr || reply->type == REDIS_REPLY_ERROR) {
            std::cerr << "Failed to initialize inventory for product " << productId << ": " << (reply ? reply->str : "Unknown error") << std::endl;
        }
        freeReplyObject(reply);
    }

    std::cout << "Inventory initialized successfully!" << std::endl;
}
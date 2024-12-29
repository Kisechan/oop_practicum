#define _CRT_SECURE_NO_WARNINGS
#include "order_processor.h"
#include <boost/json.hpp>
#include <hiredis/hiredis.h>
#include <boost/beast/core.hpp>
#include <boost/beast/http.hpp>
#include <boost/beast/version.hpp>
#include <boost/asio/connect.hpp>
#include <boost/asio/ip/tcp.hpp>
#include <iostream>
#include <string>

namespace beast = boost::beast;         // �������ռ�
namespace http = beast::http;           // �������ռ�
namespace net = boost::asio;            // �������ռ�
using tcp = net::ip::tcp;               // �������ռ�
namespace json = boost::json;
std::vector<std::pair<int, int>> fetchInventoryFromDatabase() {
    std::vector<std::pair<int, int>> inventoryData;

    try {
        // ���� I/O ������
        net::io_context ioc;

        // ��������
        tcp::resolver resolver(ioc);
        auto const results = resolver.resolve("localhost", "8081"); // API ��ַ

        // ��������
        tcp::socket socket(ioc);
        net::connect(socket, results.begin(), results.end());

        // ���� HTTP GET ����
        http::request<http::string_body> req{ http::verb::get, "/api/stock", 11 };
        req.set(http::field::host, "localhost");
        req.set(http::field::user_agent, BOOST_BEAST_VERSION_STRING);

        // ��������
        http::write(socket, req);

        // ��ȡ��Ӧ
        beast::flat_buffer buffer;
        http::response<http::dynamic_body> res;
        http::read(socket, buffer, res);

        // ���� JSON ��Ӧ
        std::string responseBody = beast::buffers_to_string(res.body().data());
        json::value jsonResponse = json::parse(responseBody);

        // ��ȡ�������
        for (const auto& item : jsonResponse.as_array()) {
            int productId = item.at("id").as_int64();
            int stock = item.at("stock").as_int64();
            inventoryData.push_back({ productId, stock });
        }

        // �ر�����
        socket.shutdown(tcp::socket::shutdown_both);
    }
    catch (const std::exception& e) {
        std::cerr << "Error fetching inventory from API: " << e.what() << std::endl;
    }

    return inventoryData;
}

// ������Ƿ����
bool checkInventory(redisContext* context, int productId, int quantity) {
    std::string key = "inventory:" + std::to_string(productId);

    // ʹ�� DECRBY ԭ���Եؼ��ٿ��
    redisReply* reply = (redisReply*)redisCommand(context, "DECRBY %s %d", key.c_str(), quantity);
    if (reply == nullptr || reply->type == REDIS_REPLY_ERROR) {
        std::cerr << "Redis error: " << (reply ? reply->str : "Unknown error") << std::endl;
        freeReplyObject(reply);
        return false;
    }

    // ������Ƿ����
    int remainingStock = reply->integer;
    freeReplyObject(reply);

    if (remainingStock >= 0) {
        std::cout << "Inventory check passed for product " << productId << ". Remaining stock: " << remainingStock << std::endl;
        return true;
    }
    else {
        // ��治�㣬�ָ����
        redisReply* restoreReply = (redisReply*)redisCommand(context, "INCRBY %s %d", key.c_str(), quantity);
        if (restoreReply == nullptr || restoreReply->type == REDIS_REPLY_ERROR) {
            std::cerr << "Failed to restore inventory for product " << productId << std::endl;
        }
        freeReplyObject(restoreReply);
        std::cerr << "Insufficient inventory for product " << productId << std::endl;
        return false;
    }
}

// ��ʼ�� Redis ���
void initializeInventory(redisContext* context) {
    // �����ݿ�� API �л�ȡ��Ʒ�������
    std::vector<std::pair<int, int>> inventoryData = fetchInventoryFromDatabase();

    // ���������д�� Redis
    for (const auto& item : inventoryData) {
        int productId = item.first;
        int stock = item.second;
        std::string key = "inventory:" + std::to_string(productId);

        // ���ÿ��ֵ
        redisReply* reply = (redisReply*)redisCommand(context, "SET %s %d", key.c_str(), stock);
        if (reply == nullptr || reply->type == REDIS_REPLY_ERROR) {
            std::cerr << "Failed to initialize inventory for product " << productId << ": " << (reply ? reply->str : "Unknown error") << std::endl;
        }
        freeReplyObject(reply);
    }

    std::cout << "Inventory initialized successfully!" << std::endl;
}

void sendSeckillResult(const std::string& orderId, bool success, const std::string& message) {
    try {
        std::cout << "\nStart Sending Seckill Results\n" << std::endl;
        // ���� I/O ������
        net::io_context ioc;

        // ��������
        tcp::resolver resolver(ioc);
        auto const results = resolver.resolve("localhost", "8080"); // Go �˵ĵ�ַ

        // ��������
        tcp::socket socket(ioc);
        net::connect(socket, results.begin(), results.end());

        // ���� HTTP POST ����
        http::request<http::string_body> req{ http::verb::post, "/orders/checkout/result", 11 };
        req.set(http::field::host, "localhost");
        req.set(http::field::user_agent, BOOST_BEAST_VERSION_STRING);
        req.set(http::field::content_type, "application/json");
        req.body() = "{\"order_id\": \"" + orderId + "\", \"status\": \"" + (success ? "success" : "failed") + "\", \"message\": \"" + message + "\"}";
        req.prepare_payload();

        // ��������
        http::write(socket, req);

        // ��ȡ��Ӧ
        beast::flat_buffer buffer;
        http::response<http::dynamic_body> res;
        http::read(socket, buffer, res);

        // �����Ӧ
        std::cout << "Sending Results Response: " << res << std::endl;

        // �ر�����
        socket.shutdown(tcp::socket::shutdown_both);
    }
    catch (std::exception const& e) {
        std::cerr << "Error: " << e.what() << std::endl;
    }
}

void createOrder(redisContext* context, const std::string& orderJson) {
    // ���� JSON ��Ϣ
    boost::json::value json = boost::json::parse(orderJson);
    std::string orderId = json.at("order_id").as_string().c_str();
    int userId = json.at("user_id").as_int64();
    int productId = json.at("product_id").as_int64();
    int quantity = json.at("quantity").as_int64();
    double total = json.at("total").as_double();

    std::cout << "Processing seckill order: " << orderId << std::endl;

    // �����
    if (!checkInventory(context, productId, quantity)) {
        std::cerr << "Insufficient inventory for product " << productId << std::endl;
        sendSeckillResult(orderId, false, "Insufficient inventory");
        return;
    }
    std::cout << "\nStart Making Order Requests\n" << std::endl;
    // ����������������
    try {
        // ���� I/O ������
        net::io_context ioc;

        // ��������
        tcp::resolver resolver(ioc);
        auto const results = resolver.resolve("localhost", "8081"); // �־ò�ĵ�ַ

        // ��������
        tcp::socket socket(ioc);
        net::connect(socket, results.begin(), results.end());

        // ���� HTTP POST ����
        http::request<http::string_body> req{ http::verb::post, "/api/orders/create", 11 };
        req.set(http::field::host, "localhost");
        req.set(http::field::user_agent, BOOST_BEAST_VERSION_STRING);
        req.set(http::field::content_type, "application/json");
        req.body() = "{\"id\": \"" + orderId + "\", \"user_id\": " + std::to_string(userId) + ", \"product_id\": " + std::to_string(productId) + ", \"quantity\": " + std::to_string(quantity) + ", \"total\": " + std::to_string(total) + ", \"status\": \"pending\"}";
        req.prepare_payload();

        // ��������
        http::write(socket, req);

        // ��ȡ��Ӧ
        beast::flat_buffer buffer;
        http::response<http::dynamic_body> res;
        http::read(socket, buffer, res);

        // �����Ӧ
        std::cout << "Make Order Response: " << res << std::endl;

        // �ر�����
        socket.shutdown(tcp::socket::shutdown_both);

        // ������ɱ�ɹ����
        sendSeckillResult(orderId, true, "Order created successfully");
    }
    catch (std::exception const& e) {
        std::cerr << "Error: " << e.what() << std::endl;
        sendSeckillResult(orderId, false, "Failed to create order");
    }
}

void asyncUpdateStock(const std::string& productId, int quantity) {
    boost::asio::io_context io_context;
    boost::asio::ip::tcp::resolver resolver(io_context);
    boost::asio::ip::tcp::socket socket(io_context);

    // ���������Ͷ˿�
    auto const results = resolver.resolve("your-api-server.com", "80");

    // ���ӷ�����
    boost::asio::connect(socket, results.begin(), results.end());

    // ���� HTTP ����
    http::request<http::string_body> req{ http::verb::post, "/update-stock", 11 };
    req.set(http::field::host, "your-api-server.com");
    req.set(http::field::content_type, "application/json");
    req.body() = R"({"product_id":")" + productId + R"(","quantity":)" + std::to_string(quantity) + "}";
    req.prepare_payload();

    // ��������
    http::write(socket, req);

    // ��ȡ��Ӧ
    boost::beast::flat_buffer buffer;
    http::response<http::string_body> res;
    http::read(socket, buffer, res);

    std::cout << "Update stock response: " << res.body() << std::endl;
}
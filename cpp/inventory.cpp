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

// �� API ��ȡ��Ʒ�������
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
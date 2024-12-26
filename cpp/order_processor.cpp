#define _CRT_SECURE_NO_WARNINGS
#include "order_processor.h"
#include <boost/json.hpp>
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

bool checkInventory(int productId, int quantity) {
    // ģ�����飨ʵ�ʳ�������Ҫ�����ݿ�򻺴��в�ѯ��
    std::cout << "Checking inventory for product " << productId << " with quantity " << quantity << std::endl;
    return true; // ���������
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

void createOrder(const std::string& orderJson) {
    // ���� JSON ��Ϣ
    boost::json::value json = boost::json::parse(orderJson);
    std::string orderId = json.at("order_id").as_string().c_str();
    int userId = json.at("user_id").as_int64();
    int productId = json.at("product_id").as_int64();
    int quantity = json.at("quantity").as_int64();
    double total = json.at("total").as_double();

    std::cout << "Processing seckill order: " << orderId << std::endl;

    // �����
    if (!checkInventory(productId, quantity)) {
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
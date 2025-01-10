#include <iostream>
#include <thread>
#include <vector>
#include <boost/asio.hpp>
#include <boost/beast.hpp>
#include <boost/json.hpp>
#include <chrono>
#include <random>

namespace asio = boost::asio;
namespace beast = boost::beast;
namespace http = beast::http;
namespace json = boost::json;
using tcp = asio::ip::tcp;

std::string generateOrderNumber() {
    auto now = std::chrono::system_clock::now();
    auto timestamp = std::chrono::duration_cast<std::chrono::seconds>(now.time_since_epoch()).count();

    std::random_device rd;
    std::mt19937 gen(rd());
    std::uniform_int_distribution<> dis(1000, 9999);

    return "ORDER_" + std::to_string(timestamp) + "_" + std::to_string(dis(gen));
}

void sendHttpRequest(int threadId, const std::string& orderNumber) {
    try {
        asio::io_context ioContext;

        tcp::resolver resolver(ioContext);
        auto const results = resolver.resolve("localhost", "8080");

        tcp::socket socket(ioContext);
        asio::connect(socket, results.begin(), results.end());

        json::value requestBody = {
            {"user_id", 6050},
            {"product_id", 2205},
            {"quantity", 1},
            {"coupon_code", ""},
            {"discount", 0.00},
            {"payable", 100.00},
            {"total", 100.00},
            {"order_number", orderNumber}
        };

        http::request<http::string_body> req{ http::verb::post, "/orders/checkout", 11 };
        req.set(http::field::host, "localhost");
        req.set(http::field::content_type, "application/json");
        req.body() = json::serialize(requestBody);
        req.prepare_payload();

        http::write(socket, req);

        beast::flat_buffer buffer;
        http::response<http::dynamic_body> res;
        http::read(socket, buffer, res);

        // 解析响应
        std::string responseBody = beast::buffers_to_string(res.body().data());
        json::value jsonResponse = json::parse(responseBody);

        if (res.result() == http::status::ok) {
            std::cout << "Thread " << threadId << ": Checkout successful. Order number: " << orderNumber << std::endl;

            // 轮询获取订单结果
            while (true) {
                // 构造 GET 请求
                http::request<http::string_body> getReq{ http::verb::get, "/orders/checkout/result/" + orderNumber, 11 };
                getReq.set(http::field::host, "localhost");

                http::write(socket, getReq);

                beast::flat_buffer getBuffer;
                http::response<http::dynamic_body> getRes;
                http::read(socket, getBuffer, getRes);

                std::string getResponseBody = beast::buffers_to_string(getRes.body().data());
                json::value getJsonResponse = json::parse(getResponseBody);

                if (getRes.result() == http::status::ok) {
                    std::cout << "Thread " << threadId << ": Order result for " << orderNumber << ": "
                        << getJsonResponse.at("result").as_string() << std::endl;
                    break;
                }
                else {
                    std::this_thread::sleep_for(std::chrono::seconds(1)); // 等待 1 秒后重试
                }
            }
        }
        else {
            std::cerr << "Thread " << threadId << ": Checkout failed. Response: " << responseBody << std::endl;
        }

        socket.shutdown(tcp::socket::shutdown_both);
    }
    catch (const std::exception& e) {
        std::cerr << "Thread " << threadId << ": Error: " << e.what() << std::endl;
    }
}

int main() {
    const int n = 10; // 线程数量
    std::vector<std::thread> threads;

    for (int i = 0; i < n; ++i) {
        std::string orderNumber = generateOrderNumber();
        threads.emplace_back(sendHttpRequest, i + 1, orderNumber);
    }

    for (auto& t : threads) {
        t.join();
    }

    return 0;
}
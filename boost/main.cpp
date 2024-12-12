#include "user_service.h"
#include <boost/asio/ip/tcp.hpp>
#include <boost/beast/core.hpp>
#include <boost/beast/http.hpp>
#include <boost/beast/version.hpp>
#include <boost/asio/strand.hpp>
#include <boost/config.hpp>
#include <cstdlib>
#include <iostream>
#include <memory>
#include <string>
#include <thread>

namespace beast = boost::beast;         // 从 Boost.Beast 命名空间导入
namespace http = beast::http;           // 从 Boost.Beast 命名空间导入
namespace net = boost::asio;            // 从 Boost.Asio 命名空间导入
using tcp = boost::asio::ip::tcp;       // 从 Boost.Asio 命名空间导入

// 处理 HTTP 请求
class HttpConnection : public std::enable_shared_from_this<HttpConnection> {
public:
    HttpConnection(tcp::socket socket) : socket_(std::move(socket)) {}

    void start() {
        readRequest();
        checkDeadline();
    }

private:
    tcp::socket socket_;
    beast::flat_buffer buffer_{ 8192 };
    http::request<http::dynamic_body> request_;
    http::response<http::dynamic_body> response_;
    net::steady_timer deadline_{ socket_.get_executor(), std::chrono::seconds(60) };

    void readRequest() {
        auto self = shared_from_this();

        http::async_read(socket_, buffer_, request_,
            [self](beast::error_code ec, std::size_t bytes_transferred) {
                boost::ignore_unused(bytes_transferred);
                if (!ec)
                    self->processRequest();
            });
    }

    void processRequest() {
        response_.version(request_.version());
        response_.keep_alive(false);

        switch (request_.method()) {
        case http::verb::get:
            handleGet();
            break;
        case http::verb::post:
            handlePost();
            break;
        default:
            response_.result(http::status::bad_request);
            response_.set(http::field::content_type, "text/plain");
            beast::ostream(response_.body()) << "Invalid request method '"
                << std::string(request_.method_string())
                << "'";
            break;
        }

        writeResponse();
    }

    void handleGet() {
        auto path = std::string(request_.target());
        if (path == "/api/getAllUsers") {
            auto users = UserService::getAllUsers();
            response_.result(http::status::ok);
            response_.set(http::field::content_type, "application/json");
            beast::ostream(response_.body()) << users;
        }
        else if (path == "/api/getUserByID") {
            auto id = request_.target().substr(13); // 提取 ID
            auto user = UserService::getUserByID(std::stoi(id));
            response_.result(http::status::ok);
            response_.set(http::field::content_type, "application/json");
            beast::ostream(response_.body()) << user;
        }
        else {
            response_.result(http::status::not_found);
            response_.set(http::field::content_type, "text/plain");
            beast::ostream(response_.body()) << "File not found\r\n";
        }
    }

    void handlePost() {
        auto path = std::string(request_.target());
        if (path == "/api/createUser") {
            beast::ostream(response_.body()) << "User created";
            response_.result(http::status::ok);
            response_.set(http::field::content_type, "text/plain");
        }
        else {
            response_.result(http::status::not_found);
            response_.set(http::field::content_type, "text/plain");
            beast::ostream(response_.body()) << "File not found\r\n";
        }
    }

    void writeResponse() {
        auto self = shared_from_this();

        response_.content_length(response_.body().size());

        http::async_write(socket_, response_,
            [self](beast::error_code ec, std::size_t) {
                self->socket_.shutdown(tcp::socket::shutdown_send, ec);
                self->deadline_.cancel();
            });
    }

    void checkDeadline() {
        auto self = shared_from_this();

        deadline_.async_wait(
            [self](beast::error_code ec) {
                if (!ec) {
                    self->socket_.close(ec);
                }
            });
    }
};

// 启动 HTTP 服务器
void HttpServer(net::io_context& ioc, tcp::endpoint endpoint) {
    tcp::acceptor acceptor{ ioc, endpoint };
    for (;;) {
        auto socket = std::make_shared<tcp::socket>(ioc);
        acceptor.accept(*socket);
        std::make_shared<HttpConnection>(std::move(*socket))->start();
    }
}

int main(int argc, char* argv[]) {
    try {
        if (argc != 3) {
            std::cerr << "Usage: http_server <address> <port>\n";
            return EXIT_FAILURE;
        }

        auto const address = net::ip::make_address(argv[1]);
        auto const port = static_cast<unsigned short>(std::stoi(argv[2]));

        net::io_context ioc{ 1 };

        std::thread t{ [&ioc] { ioc.run(); } };

        HttpServer(ioc, tcp::endpoint{ address, port });

        t.join();
    }
    catch (const std::exception& e) {
        std::cerr << "Error: " << e.what() << std::endl;
        return EXIT_FAILURE;
    }
}
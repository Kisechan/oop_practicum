#include <boost/beast/core.hpp>
#include <boost/beast/websocket.hpp>
#include <boost/asio/connect.hpp>
#include <boost/asio/ip/tcp.hpp>
#include <cstdlib>
#include <iostream>
#include <string>

namespace beast = boost::beast;         // from <boost/beast.hpp>
namespace http = beast::http;           // from <boost/beast/http.hpp>
namespace websocket = beast::websocket; // from <boost/beast/websocket.hpp>
namespace net = boost::asio;            // from <boost/asio.hpp>
using tcp = boost::asio::ip::tcp;       // from <boost/asio/ip/tcp.hpp>

int main(int argc, char** argv) {
    try {
        
        std::string host = "localhost";
        std::string port = "8080";

        // 创建I/O上下文
        net::io_context ioc;

        // 解析主机名和端口
        tcp::resolver resolver(ioc);
        auto const results = resolver.resolve(host, port);

        // 创建并连接WebSocket客户端
        websocket::stream<tcp::socket> ws(ioc);
        net::connect(ws.next_layer(), results.begin(), results.end());

        // 设置WebSocket选项并完成握手
        ws.set_option(websocket::stream_base::decorator(
            [](websocket::request_type& req) {
                req.set(http::field::user_agent,
                std::string(BOOST_BEAST_VERSION_STRING) + " websocket-client-coro");
            }));
        ws.handshake(host, "/ws");

        std::cout << "Connected to Go WebSocket server" << std::endl;

        // 发送消息到Go服务器
        std::string msg = R"({"type":"ReadUser", "payload": {"id":6045}})";
        ws.write(net::buffer(msg));
        std::cout << "Sent to Go: " << msg << std::endl;

        // 读取Go服务器的响应
        beast::flat_buffer buffer;
        ws.read(buffer);
        std::cout << "Received from Go: " << beast::make_printable(buffer.data()) << std::endl;

        // 关闭WebSocket连接
        ws.close(websocket::close_code::normal);
        std::cout << "WebSocket connection closed" << std::endl;
    }
    catch (std::exception const& e) {
        std::cerr << "Error: " << e.what() << std::endl;
        return EXIT_FAILURE;
    }

    return EXIT_SUCCESS;
}
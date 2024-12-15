#include <boost/asio.hpp>
#include <boost/beast.hpp>
#include <boost/beast/http.hpp>
#include <iostream>

namespace beast = boost::beast;         // 从Boost Beast库中导入
namespace http = beast::http;           // 导入HTTP模块
namespace net = boost::asio;            // 导入Boost Asio库
using tcp = net::ip::tcp;               // 导入TCP模块

int main() {
    while(1)
    {
        try {
            // 创建I/O上下文
            net::io_context ioc;

            // 解析服务器地址和端口
            tcp::resolver resolver(ioc);
            auto const results = resolver.resolve("localhost", "8080");

            // 创建连接
            tcp::socket socket(ioc);
            net::connect(socket, results);

            // 创建HTTP GET请求
            http::request<http::string_body> req{ http::verb::get, "/products/search?keyword=手机&category=手机", 11 };
            req.set(http::field::host, "localhost");
            req.set(http::field::user_agent, BOOST_BEAST_VERSION_STRING);

            // 发送请求
            http::write(socket, req);

            // 读取响应
            beast::flat_buffer buffer;
            http::response<http::dynamic_body> res;
            http::read(socket, buffer, res);

            // 输出响应
            std::cout << "Status: " << res.result_int() << std::endl;
            std::cout << "Body: " << beast::buffers_to_string(res.body().data()) << std::endl;

            // 关闭连接
            socket.shutdown(tcp::socket::shutdown_both);
        }
        catch (std::exception const& e) {
            std::cerr << "Error: " << e.what() << std::endl;
        }
    }

    return 0;
}
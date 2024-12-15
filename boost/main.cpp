#include <boost/asio.hpp>
#include <boost/beast.hpp>
#include <boost/beast/http.hpp>
#include <iostream>

namespace beast = boost::beast;         // ��Boost Beast���е���
namespace http = beast::http;           // ����HTTPģ��
namespace net = boost::asio;            // ����Boost Asio��
using tcp = net::ip::tcp;               // ����TCPģ��

int main() {
    while(1)
    {
        try {
            // ����I/O������
            net::io_context ioc;

            // ������������ַ�Ͷ˿�
            tcp::resolver resolver(ioc);
            auto const results = resolver.resolve("localhost", "8080");

            // ��������
            tcp::socket socket(ioc);
            net::connect(socket, results);

            // ����HTTP GET����
            http::request<http::string_body> req{ http::verb::get, "/products/search?keyword=�ֻ�&category=�ֻ�", 11 };
            req.set(http::field::host, "localhost");
            req.set(http::field::user_agent, BOOST_BEAST_VERSION_STRING);

            // ��������
            http::write(socket, req);

            // ��ȡ��Ӧ
            beast::flat_buffer buffer;
            http::response<http::dynamic_body> res;
            http::read(socket, buffer, res);

            // �����Ӧ
            std::cout << "Status: " << res.result_int() << std::endl;
            std::cout << "Body: " << beast::buffers_to_string(res.body().data()) << std::endl;

            // �ر�����
            socket.shutdown(tcp::socket::shutdown_both);
        }
        catch (std::exception const& e) {
            std::cerr << "Error: " << e.what() << std::endl;
        }
    }

    return 0;
}
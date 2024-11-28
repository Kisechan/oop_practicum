#include <boost/beast.hpp>
#include <boost/asio.hpp>
#include <iostream>
#include <hiredis.h>

namespace beast = boost::beast;
namespace http = beast::http;
namespace net = boost::asio;
using tcp = net::ip::tcp;

void handleRequest(http::request<http::string_body>& req,
    http::response<http::string_body>& res,
    OrderManager& orderManager) 
{
    if (req.method() != http::verb::post) 
    {
        res.result(http::status::bad_request);
        res.body() = "Invalid HTTP method";
        return;
    }
    auto body = req.body();
    int userId;
    std::vector<int> productIds;

    try 
    {
        auto pos = body.find("userId:");
        userId = std::stoi(body.substr(pos + 7, body.find(",") - (pos + 7)));

        auto products = body.substr(body.find("productIds:") + 11);
        size_t start = 0, end;
        while ((end = products.find(",", start)) != std::string::npos) 
        {
            productIds.push_back(std::stoi(products.substr(start, end - start)));
            start = end + 1;
        }
        productIds.push_back(std::stoi(products.substr(start)));
    }
    catch (...) 
    {
        res.result(http::status::bad_request);
        res.body() = "Invalid request format.";
        return;
    }
    std::string result = orderManager.createOrder(userId, productIds);
    res.result(http::status::ok);
    res.body() = result;
}

int main() 
{
    try 
    {
        net::io_context ioc;
        tcp::acceptor acceptor(ioc, tcp::endpoint(tcp::v4(), 8080));


        std::cout << "Server running on http://127.0.0.1:8080/" << std::endl;

        for (;;) 
        {
            tcp::socket socket(ioc);
            acceptor.accept(socket);
            beast::flat_buffer buffer;
            http::request<http::string_body> req;
            http::read(socket, buffer, req);

            http::response<http::string_body> res;
            handleRequest(req, res, orderManager);
            redisConnect("localhost",8080);
            http::write(socket, res);
        }
    }
    catch (std::exception& e) 
    {
        std::cerr << "Error: " << e.what() << std::endl;
    }

    return 0;
}

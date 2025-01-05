#include "order_processor.h"
#include <iostream>

// 全局变量
redisContext* context;
std::mutex redisMutex;
std::mutex httpMutex;

// 获取数据
std::vector<std::pair<int, int>> fetchInventoryFromDatabase() {
    std::vector<std::pair<int, int>> inventoryData;

    try {
        asio::io_context ioContext;
        tcp::resolver resolver(ioContext);
        auto const results = resolver.resolve("localhost", "8081");

        tcp::socket socket(ioContext);
        asio::connect(socket, results.begin(), results.end());

        // 构造 HTTP GET 请求
        http::request<http::string_body> req{ http::verb::get, "/products/stock", 11 };
        req.set(http::field::host, "localhost");
        req.set(http::field::user_agent, BOOST_BEAST_VERSION_STRING);

        http::write(socket, req);

        beast::flat_buffer buffer;
        http::response<http::dynamic_body> res;
        http::read(socket, buffer, res);

        std::string responseBody = beast::buffers_to_string(res.body().data());
        json::value jsonResponse = json::parse(responseBody);

        // 检查响应状态
        if (jsonResponse.at("status").as_string() != "success") {
            std::cerr << "API returned an error status" << std::endl;
            return inventoryData;
        }

        // 解析库存数据
        json::array inventoryArray = jsonResponse.at("data").as_array();
        for (const auto& item : inventoryArray) {
            int productId = item.at("id").as_int64();
            int stock = item.at("stock").as_int64();
            std::cout << "productId=" << productId << ", stock=" << stock << std::endl;
            inventoryData.push_back({ productId, stock });
        }

        socket.shutdown(tcp::socket::shutdown_both);
    }
    catch (const std::exception& e) {
        std::cerr << "Error fetching inventory: " << e.what() << std::endl;
    }

    return inventoryData;
}

// 初始化库存到 Redis
void initializeInventory(redisContext* context) {
    std::vector<std::pair<int, int>> inventoryData = fetchInventoryFromDatabase();
    
    if (inventoryData.empty()) {
        std::cout << "Inventory Data is Empty!" << std::endl;
        return;
    }
    
    std::lock_guard<std::mutex> lock(redisMutex);
    for (const auto& item : inventoryData) {
        int productId = item.first;
        int stock = item.second;
        std::string key = "inventory:" + std::to_string(productId);

        redisReply* reply = (redisReply*)redisCommand(context, "SET %s %d", key.c_str(), stock);
        if (reply == nullptr || reply->type == REDIS_REPLY_ERROR) {
            std::cerr << "Failed to initialize inventory for product " << productId << std::endl;
        }
        freeReplyObject(reply);
    }

    std::cout << "Inventory initialized successfully!" << std::endl;
}

// 原子化减少库存
bool decreaseInventory(redisContext* context, int productId, int quantity) {
    std::lock_guard<std::mutex> lock(redisMutex);
    std::string key = "inventory:" + std::to_string(productId);

    redisReply* reply = (redisReply*)redisCommand(context, "DECRBY %s %d", key.c_str(), quantity);
    if (reply == nullptr || reply->type == REDIS_REPLY_ERROR) {
        std::cerr << "Redis error: " << (reply ? reply->str : "Unknown error") << std::endl;
        freeReplyObject(reply);
        return false;
    }

    int remainingStock = reply->integer;
    freeReplyObject(reply);

    if (remainingStock >= 0) {
        std::cout << "Inventory decreased for product " << productId << ". Remaining stock: " << remainingStock << std::endl;
        return true;
    }
    else {
        // 回滚库存
        redisReply* restoreReply = (redisReply*)redisCommand(context, "INCRBY %s %d", key.c_str(), quantity);
        if (restoreReply == nullptr || restoreReply->type == REDIS_REPLY_ERROR) {
            std::cerr << "Failed to restore inventory for product " << productId << std::endl;
        }
        freeReplyObject(restoreReply);
        std::cerr << "Insufficient inventory for product " << productId << std::endl;
        return false;
    }
}

// 使用优惠券
bool useCoupon(const std::string& couponCode, int userId) {
    try {
        asio::io_context ioContext;
        tcp::resolver resolver(ioContext);
        auto const results = resolver.resolve("localhost", "8080");

        tcp::socket socket(ioContext);
        asio::connect(socket, results.begin(), results.end());

        // 构造 JSON 请求
        json::value couponRequest = {
            {"user_id", userId},
            {"coupon_code", couponCode}
        };

        // 构造 HTTP POST 请求
        http::request<http::string_body> req{ http::verb::post, "/coupons/use", 11 };
        req.set(http::field::host, "localhost");
        req.set(http::field::content_type, "application/json");
        req.body() = json::serialize(couponRequest);
        req.prepare_payload();

        // 发送请求
        http::write(socket, req);

        // 接收响应
        beast::flat_buffer buffer;
        http::response<http::dynamic_body> res;
        http::read(socket, buffer, res);

        // 解析响应
        std::string responseBody = beast::buffers_to_string(res.body().data());
        json::value jsonResponse = json::parse(responseBody);

        // 关闭连接
        socket.shutdown(tcp::socket::shutdown_both);

        // 检查响应状态
        if (jsonResponse.at("status").as_string() == "success") {
            std::cout << "Coupon used successfully" << std::endl;
            return true;
        }
        else {
            std::cerr << "Failed to use coupon: " << responseBody << std::endl;
            return false;
        }
    }
    catch (const std::exception& e) {
        std::cerr << "Error using coupon: " << e.what() << std::endl;
        return false;
    }
}

// 异步持久化订单信息
void persistOrderAsync(const json::value& order) {
    std::thread([order]() {
        try {
            asio::io_context ioContext;
            tcp::resolver resolver(ioContext);
            auto const results = resolver.resolve("localhost", "8081");

            tcp::socket socket(ioContext);
            asio::connect(socket, results.begin(), results.end());

            http::request<http::string_body> req{ http::verb::post, "/orders/create", 11 };
            req.set(http::field::host, "localhost");
            req.set(http::field::content_type, "application/json");
            req.body() = json::serialize(order);
            req.prepare_payload();

            http::write(socket, req);

            beast::flat_buffer buffer;
            http::response<http::dynamic_body> res;
            http::read(socket, buffer, res);

            socket.shutdown(tcp::socket::shutdown_both);
        }
        catch (const std::exception& e) {
            std::cerr << "Error persisting order: " << e.what() << std::endl;
        }
        }).detach();
}

// 异步更新库存信息
void updateStockAsync(int productId, int quantity) {
    std::thread([productId, quantity]() {
        try {
            asio::io_context ioContext;
            tcp::resolver resolver(ioContext);
            auto const results = resolver.resolve("localhost", "8081");

            tcp::socket socket(ioContext);
            asio::connect(socket, results.begin(), results.end());

            json::value stockRequest = {
                {"product_id", productId},
                {"quantity", -quantity}
            };

            http::request<http::string_body> req{ http::verb::post, "/products/stock/update", 11 };
            req.set(http::field::host, "localhost");
            req.set(http::field::content_type, "application/json");
            req.body() = json::serialize(stockRequest);
            req.prepare_payload();

            http::write(socket, req);

            beast::flat_buffer buffer;
            http::response<http::dynamic_body> res;
            http::read(socket, buffer, res);

            socket.shutdown(tcp::socket::shutdown_both);
        }
        catch (const std::exception& e) {
            std::cerr << "Error updating stock: " << e.what() << std::endl;
        }
        }).detach();
}

void processCheckoutRequest(const std::string& requestJson) {
    std::cout << "Starting Process Checkout Requests" << std::endl;

    try {
        // 解析 JSON 请求
        json::value request = json::parse(requestJson);
        std::string orderNumber = request.at("order_number").as_string().c_str();
        int userId = request.at("user_id").as_int64();
        int productId = request.at("product_id").as_int64();
        int quantity = request.at("quantity").as_int64();
        std::string couponCode = request.at("coupon_code").as_string().c_str();

        // 处理 discount 字段
        double discount;
        if (request.at("discount").is_double()) {
            discount = request.at("discount").as_double();
        }
        else if (request.at("discount").is_int64()) {
            discount = static_cast<double>(request.at("discount").as_int64());
        }
        else {
            throw std::runtime_error("Invalid type for discount");
        }

        // 处理 payable 字段
        double payable;
        if (request.at("payable").is_double()) {
            payable = request.at("payable").as_double();
        }
        else if (request.at("payable").is_int64()) {
            payable = static_cast<double>(request.at("payable").as_int64());
        }
        else {
            throw std::runtime_error("Invalid type for payable");
        }

        // 处理 total 字段
        double total;
        if (request.at("total").is_double()) {
            total = request.at("total").as_double();
        }
        else if (request.at("total").is_int64()) {
            total = static_cast<double>(request.at("total").as_int64());
        }
        else {
            throw std::runtime_error("Invalid type for total");
        }

        std::cout << "Starting Checking Inventory" << std::endl;

        // 原子化减少库存
        if (!decreaseInventory(context, productId, quantity)) {
            json::value result = {
                {"order_number", orderNumber},
                {"status", "failed"},
                {"message", "Insufficient inventory"}
            };
            redisCommand(context, "SET order_result:%s %s", orderNumber.c_str(), json::serialize(result).c_str());
            return;
        }
        std::cout << "Start Checking Coupons" << std::endl;
        // 处理优惠券
        if (!useCoupon(couponCode, userId)) {
            // 回滚库存
            redisCommand(context, "INCRBY inventory:%d %d", productId, quantity);

            json::value result = {
                {"order_number", orderNumber},
                {"status", "failed"},
                {"message", "Invalid coupon"}
            };
            redisCommand(context, "SET order_result:%s %s", orderNumber.c_str(), json::serialize(result).c_str());
            return;
        }

        // 异步持久化订单信息
        std::cout << "Start Persisting Orders Asyncly" << std::endl;

        json::value order = {
            {"order_number", orderNumber},
            {"user_id", userId},
            {"product_id", productId},
            {"quantity", quantity},
            {"discount", discount},
            {"payable", payable},
            {"total", total},
            {"status", "pending"}
        };
        persistOrderAsync(order);

        // 异步更新库存信息
        std::cout << "Start Updating Stock Asyncly" << std::endl;

        updateStockAsync(productId, quantity);

        // 推送成功结果到 Redis
        json::value result = {
            {"order_number", orderNumber},
            {"status", "success"},
            {"message", "Order completed"}
        };
        redisCommand(context, "SET order_result:%s %s", orderNumber.c_str(), json::serialize(result).c_str());

        std::cout << "All things completed" << std::endl;
    }
    catch (const std::exception& e) {
        std::cerr << "Error processing checkout request: " << e.what() << std::endl;
    }
}
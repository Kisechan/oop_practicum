#include <iostream>
#include <string>
#include <hiredis/hiredis.h>
//#include <boost/json.hpp>
#include <boost/asio.hpp>

//// 处理订单结算逻辑
//void processOrderSettlement(const std::string& orderJson) {
//    // 解析 JSON 消息
//    boost::json::value json = boost::json::parse(orderJson);
//    int orderId = json.at("order_id").as_int64();
//    int userId = json.at("user_id").as_int64();
//    double amount = json.at("amount").as_double();
//    std::string paymentMethod = json.at("payment_method").as_string().c_str();
//
//    std::cout << "Processing order settlement: " << orderId << std::endl;
//
//    // 模拟订单结算逻辑
//    // 1. 检查用户余额
//    // 2. 扣除用户余额
//    // 3. 更新订单状态
//
//    // 调用持久层更新订单状态
//    cpr::Response r = cpr::Post(
//        cpr::Url{ "http://localhost:8081/api/orders/update" },
//        cpr::Body{ orderJson },
//        cpr::Header{ {"Content-Type", "application/json"} }
//    );
//
//    if (r.status_code == 200) {
//        std::cout << "Order " << orderId << " successfully updated to paid." << std::endl;
//    }
//    else {
//        std::cerr << "Failed to update order: " << r.status_code << " - " << r.text << std::endl;
//    }
//}
//
//// 从 Redis 消费消息
//void consumeRedisMessages(redisContext* context) {
//    redisReply* reply;
//    while (true) {
//        reply = (redisReply*)redisCommand(context, "BLPOP order_queue 0");
//        if (reply && reply->type == REDIS_REPLY_ARRAY && reply->element[1]) {
//            std::string message(reply->element[1]->str, reply->element[1]->len);
//            std::cout << "Received message: " << message << std::endl;
//
//            // 处理订单结算
//            processOrderSettlement(message);
//        }
//        freeReplyObject(reply);
//    }
//}

int main() {
    // 连接 Redis
    redisContext* context = redisConnect("127.0.0.1", 6379);
    if (context == NULL || context->err) {
        if (context) {
            std::cerr << "Redis connection error: " << context->errstr << std::endl;
        }
        else {
            std::cerr << "Redis connection error: cannot allocate redis context" << std::endl;
        }
        return 1;
    }

    std::cout << "Connected to Redis." << std::endl;

    // 开始消费消息
    //consumeRedisMessages(context);

    // 关闭 Redis 连接
    redisFree(context);
    std::cout << "Redis Freed" << std::endl;

    return 0;
}
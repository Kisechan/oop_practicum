#include <iostream>
#include <string>
#include <hiredis/hiredis.h>
//#include <boost/json.hpp>
#include <boost/asio.hpp>

//// �����������߼�
//void processOrderSettlement(const std::string& orderJson) {
//    // ���� JSON ��Ϣ
//    boost::json::value json = boost::json::parse(orderJson);
//    int orderId = json.at("order_id").as_int64();
//    int userId = json.at("user_id").as_int64();
//    double amount = json.at("amount").as_double();
//    std::string paymentMethod = json.at("payment_method").as_string().c_str();
//
//    std::cout << "Processing order settlement: " << orderId << std::endl;
//
//    // ģ�ⶩ�������߼�
//    // 1. ����û����
//    // 2. �۳��û����
//    // 3. ���¶���״̬
//
//    // ���ó־ò���¶���״̬
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
//// �� Redis ������Ϣ
//void consumeRedisMessages(redisContext* context) {
//    redisReply* reply;
//    while (true) {
//        reply = (redisReply*)redisCommand(context, "BLPOP order_queue 0");
//        if (reply && reply->type == REDIS_REPLY_ARRAY && reply->element[1]) {
//            std::string message(reply->element[1]->str, reply->element[1]->len);
//            std::cout << "Received message: " << message << std::endl;
//
//            // ����������
//            processOrderSettlement(message);
//        }
//        freeReplyObject(reply);
//    }
//}

int main() {
    // ���� Redis
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

    // ��ʼ������Ϣ
    //consumeRedisMessages(context);

    // �ر� Redis ����
    redisFree(context);
    std::cout << "Redis Freed" << std::endl;

    return 0;
}
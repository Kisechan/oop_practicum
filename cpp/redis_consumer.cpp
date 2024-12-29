#include "redis_consumer.h"
#include "order_processor.h"
#include <iostream>

void consumeRedisMessages(redisContext* context) {
    redisReply* reply;
    while (true) {
        reply = (redisReply*)redisCommand(context, "BLPOP seckill_queue 0");
        if (reply && reply->type == REDIS_REPLY_ARRAY && reply->element[1]) {
            std::string message(reply->element[1]->str, reply->element[1]->len);
            std::cout << "Received message: " << message << std::endl;
            if (message == "")
            {
                std::cout << "Message is null, skipped" << std::endl;
                continue;
            }
            // ´¦ÀíÃëÉ±¶©µ¥
            std::cout << "\nStart Resolving Seckill Requests\n" << std::endl;
            createOrder(context, message);
        }
        freeReplyObject(reply);
    }
}
#include "redis_consumer.h"
#include "order_processor.h"
#include <iostream>

void consumeRedisMessages(redisContext* context) {
    std::cout << "Starting Consuming Redis Messages" << std::endl;
    while (true) {
        redisReply* reply = (redisReply*)redisCommand(context, "BLPOP checkout_queue 0");
        if (reply && reply->type == REDIS_REPLY_ARRAY && reply->element[1]) {
            std::string message(reply->element[1]->str, reply->element[1]->len);
            std::cout << "Received message: " << message << std::endl;
            if (message == "") {
                std::cout << "Message is empty! Continue!" << std::endl;
                continue;
            }
            processCheckoutRequest(message);
        }
        freeReplyObject(reply);
    }
}
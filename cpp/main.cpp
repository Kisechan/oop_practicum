#include "redis_consumer.h"
#include "order_processor.h"
#include <hiredis/hiredis.h>
#include <iostream>

int main() {
    // 初始化 Redis
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

    std::cout << "Connected to Redis!Let's Start!" << std::endl;

    initializeInventory(context);
    consumeRedisMessages(context);
    redisFree(context);

    return 0;
}
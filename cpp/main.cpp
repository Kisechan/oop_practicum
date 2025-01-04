#include "redis_consumer.h"
#include "order_processor.h"
#include <hiredis/hiredis.h>
#include <iostream>

int main() {
    // 初始化 Redis
    context = redisConnect("127.0.0.1", 6379);
    if (context == nullptr || context->err) {
        std::cerr << "Redis connection error: " << (context ? context->errstr : "Cannot allocate context") << std::endl;
        return 1;
    }

    std::cout << "Connected to Redis! Let's start!" << std::endl;

    // 初始化库存
    initializeInventory(context);

    // 消费 Redis 消息队列
    consumeRedisMessages(context);

    redisFree(context);
    return 0;
}
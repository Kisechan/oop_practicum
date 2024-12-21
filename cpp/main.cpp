#include "redis_consumer.h"
#include <hiredis/hiredis.h>
#include <iostream>

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

    std::cout << "Connected to Redis!Let's Start!" << std::endl;

    // 开始消费消息
    consumeRedisMessages(context);

    // 关闭 Redis 连接
    redisFree(context);

    return 0;
}
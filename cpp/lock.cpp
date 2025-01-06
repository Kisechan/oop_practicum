#include "lock.h"

// 获取分布式锁
bool acquireLock(redisContext* context, const std::string& lockKey, int timeout) {
    std::string command = "SET " + lockKey + " locked NX EX " + std::to_string(timeout);
    redisReply* reply = (redisReply*)redisCommand(context, command.c_str());
    if (reply == nullptr || reply->type == REDIS_REPLY_ERROR) {
        std::cerr << "Failed to acquire lock: " << (reply ? reply->str : "Unknown error") << std::endl;
        freeReplyObject(reply);
        return false;
    }
    bool lockAcquired = (reply->type == REDIS_REPLY_STATUS && std::string(reply->str) == "OK");
    freeReplyObject(reply);
    return lockAcquired;
}

// 释放分布式锁
void releaseLock(redisContext* context, const std::string& lockKey) {
    redisReply* reply = (redisReply*)redisCommand(context, "DEL %s", lockKey.c_str());
    if (reply == nullptr || reply->type == REDIS_REPLY_ERROR) {
        std::cerr << "Failed to release lock: " << (reply ? reply->str : "Unknown error") << std::endl;
    }
    freeReplyObject(reply);
}
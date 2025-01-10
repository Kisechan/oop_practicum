#pragma once
#include <hiredis/hiredis.h>
#include <iostream>
#include <string>
#include <thread>
#include <chrono>

// 获取分布式锁
bool acquireLock(redisContext* context, const std::string& lockKey, int timeout);

// 释放分布式锁
void releaseLock(redisContext* context, const std::string& lockKey);
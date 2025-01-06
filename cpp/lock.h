#pragma once
#include <hiredis/hiredis.h>
#include <iostream>
#include <string>
#include <thread>
#include <chrono>

// ��ȡ�ֲ�ʽ��
bool acquireLock(redisContext* context, const std::string& lockKey, int timeout);

// �ͷŷֲ�ʽ��
void releaseLock(redisContext* context, const std::string& lockKey);
#pragma once
#include <hiredis/hiredis.h>

void consumeRedisMessages(redisContext* context);

#pragma once
#include <hiredis/hiredis.h>
#include <string>

bool checkInventory(redisContext* context, int productId, int quantity);
void initializeInventory(redisContext* context);
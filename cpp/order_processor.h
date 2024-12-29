#pragma once
#include <string>
#include <hiredis/hiredis.h>

void createOrder(redisContext* context, const std::string& orderJson);
void sendSeckillResult(const std::string& orderId, bool success, const std::string& message);
bool checkInventory(redisContext* context, int productId, int quantity);
void initializeInventory(redisContext* context);
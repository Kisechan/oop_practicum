#pragma once
#include <string>

void createOrder(const std::string& orderJson);
bool checkInventory(int productId, int quantity);
void sendSeckillResult(const std::string& orderId, bool success, const std::string& message);
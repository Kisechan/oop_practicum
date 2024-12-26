#pragma once
#include <string>

void createOrder(const std::string& orderJson);
void sendSeckillResult(const std::string& orderId, bool success, const std::string& message);
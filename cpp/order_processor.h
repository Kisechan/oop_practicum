#pragma once
#include <string>
#include <hiredis/hiredis.h>
#include <boost/json.hpp>
#include <boost/asio.hpp>
#include <boost/beast.hpp>
#include <mutex>

namespace json = boost::json;
namespace asio = boost::asio;
namespace beast = boost::beast;
namespace http = beast::http;
using boost::asio::ip::tcp;

// 全局 Redis 上下文
extern redisContext* context;

// 同步锁
extern std::mutex redisMutex;
extern std::mutex httpMutex;

// 初始化库存
void initializeInventory(redisContext* context);

// 处理结算请求
void processCheckoutRequest(const std::string& requestJson);

// 消费 Redis 消息队列
void consumeRedisMessages(redisContext* context);
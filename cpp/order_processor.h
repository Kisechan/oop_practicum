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

// ȫ�� Redis ������
extern redisContext* context;

// ͬ����
extern std::mutex redisMutex;
extern std::mutex httpMutex;

// ��ʼ�����
void initializeInventory(redisContext* context);

// �����������
void processCheckoutRequest(const std::string& requestJson);

// ���� Redis ��Ϣ����
void consumeRedisMessages(redisContext* context);
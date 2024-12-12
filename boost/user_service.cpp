#include "user_service.h"
#include "json.hpp"
#include <iostream>

using json = nlohmann::json;

// ��ȡ�����û�
std::string UserService::getAllUsers() {
    json users = {
        {{"id", 1}, {"username", "Alice"}, {"email", "alice@example.com"}},
        {{"id", 2}, {"username", "Bob"}, {"email", "bob@example.com"}}
    };
    return users.dump();
}

// ��ȡ�����û�
std::string UserService::getUserByID(int id) {
    json user = { {"id", id}, {"username", "Alice"}, {"email", "alice@example.com"} };
    return user.dump();
}

// �����û�
std::string UserService::createUser(const std::string& userData) {
    json user = json::parse(userData);
    std::cout << "User created: " << user.dump() << std::endl;
    return "User created";
}
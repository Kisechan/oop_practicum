#pragma once

#include <string>

class UserService {
public:
    static std::string getAllUsers();
    static std::string getUserByID(int id);
    static std::string createUser(const std::string& userData);
};
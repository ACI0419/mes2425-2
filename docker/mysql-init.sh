#!/bin/bash
set -e

# 初始化MySQL
service mysql start

# 等待MySQL启动
while ! mysqladmin ping -h"localhost" --silent; do
    echo "Waiting for MySQL to start..."
    sleep 1
done

# 创建数据库和用户
mysql -u root <<-EOSQL
    CREATE DATABASE IF NOT EXISTS mes_system;
    CREATE USER IF NOT EXISTS 'mes_user'@'%' IDENTIFIED BY 'mes_password';
    GRANT ALL PRIVILEGES ON mes_system.* TO 'mes_user'@'%';
    FLUSH PRIVILEGES;
EOSQL

echo "MySQL initialization completed"
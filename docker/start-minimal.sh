#!/bin/bash
set -e

# 初始化MySQL数据目录
if [ ! -d "/var/lib/mysql/mysql" ]; then
    echo "Initializing MySQL data directory..."
    mysqld --initialize-insecure --user=mysql --datadir=/var/lib/mysql
fi

# 启动MySQL
echo "Starting MySQL..."
service mysql start

# 等待MySQL启动
while ! mysqladmin ping -h"127.0.0.1" --silent; do
    echo "Waiting for MySQL to start..."
    sleep 2
done

# 创建数据库和用户
echo "Setting up database..."
mysql -u root <<-EOSQL
    CREATE DATABASE IF NOT EXISTS mes_system;
    CREATE USER IF NOT EXISTS 'mes_user'@'localhost' IDENTIFIED BY 'mes_password';
    GRANT ALL PRIVILEGES ON mes_system.* TO 'mes_user'@'localhost';
    FLUSH PRIVILEGES;
EOSQL

# 停止MySQL（supervisor会重新启动）
service mysql stop

# 启动supervisor管理所有服务
echo "Starting all services with supervisor..."
exec /usr/bin/supervisord -c /etc/supervisor/conf.d/supervisord.conf
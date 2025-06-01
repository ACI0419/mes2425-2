-- MES系统数据库初始化脚本

-- 设置字符集
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- 创建默认管理员用户
USE mes_system;

-- 插入默认管理员账户（密码：admin123，已加密）
INSERT INTO `users` (`username`, `password`, `email`, `real_name`, `phone`, `role`, `status`, `created_at`, `updated_at`) 
VALUES 
('admin', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin@mes.com', '系统管理员', '13800138000', 'admin', 1, NOW(), NOW()),
('operator', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'operator@mes.com', '操作员', '13800138001', 'operator', 1, NOW(), NOW()),
('inspector', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'inspector@mes.com', '质检员', '13800138002', 'inspector', 1, NOW(), NOW());

-- 插入示例产品数据
INSERT INTO `products` (`code`, `name`, `description`, `unit`, `price`, `status`, `created_at`, `updated_at`) 
VALUES 
('P001', '产品A', '这是产品A的描述', '个', 100.00, 1, NOW(), NOW()),
('P002', '产品B', '这是产品B的描述', '个', 200.00, 1, NOW(), NOW()),
('P003', '产品C', '这是产品C的描述', '套', 300.00, 1, NOW(), NOW());

-- 插入示例物料数据
INSERT INTO `materials` (`code`, `name`, `category`, `unit`, `price`, `min_stock`, `max_stock`, `current_stock`, `status`, `created_at`, `updated_at`) 
VALUES 
('M001', '原料A', '原材料', 'kg', 10.00, 100, 1000, 500, 1, NOW(), NOW()),
('M002', '原料B', '原材料', 'kg', 15.00, 50, 500, 200, 1, NOW(), NOW()),
('M003', '包装材料', '包装', '个', 2.00, 1000, 10000, 5000, 1, NOW(), NOW());

-- 插入示例设备数据
INSERT INTO `equipment` (`code`, `name`, `type`, `model`, `manufacturer`, `location`, `status`, `purchase_date`, `created_at`, `updated_at`) 
VALUES 
('E001', '生产线1', '生产设备', 'PL-2023', '设备制造商A', '车间A', 'normal', '2023-01-01', NOW(), NOW()),
('E002', '检测设备1', '检测设备', 'QC-2023', '设备制造商B', '质检室', 'normal', '2023-02-01', NOW(), NOW()),
('E003', '包装机1', '包装设备', 'PK-2023', '设备制造商C', '包装车间', 'normal', '2023-03-01', NOW(), NOW());

SET FOREIGN_KEY_CHECKS = 1;
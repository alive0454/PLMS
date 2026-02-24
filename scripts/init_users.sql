-- 初始化系统用户数据
-- 默认密码：admin / admin123，user / user123

-- 如果用户已存在则删除
DELETE FROM sys_user WHERE username IN ('admin', 'user');

-- 插入默认用户（使用正确的 bcrypt 哈希）
-- admin123 的哈希
INSERT INTO sys_user (username, password, name, role, is_default_password, status) VALUES 
('admin', '$2a$10$bJPcUrvN1brNoGHFXbIqJebTljKm.pwjfBd48/Bb1uIk5wvVjJ6Ie', '系统管理员', 'admin', 1, 1);

-- user123 的哈希
INSERT INTO sys_user (username, password, name, role, is_default_password, status) VALUES 
('user', '$2a$10$8LLkeFX41.0BCdnYSdP8j.23SCbDsV2tP8LuVV5vpCpgNnCwGwgde', '普通用户', 'user', 1, 1);

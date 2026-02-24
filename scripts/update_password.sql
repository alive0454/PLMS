-- 更新用户密码为正确的 bcrypt 哈希
-- 执行前请确认数据库连接信息

-- 更新 admin 用户密码为 admin123
UPDATE sys_user 
SET password = '$2a$10$bJPcUrvN1brNoGHFXbIqJebTljKm.pwjfBd48/Bb1uIk5wvVjJ6Ie',
    is_default_password = 1,
    status = 1
WHERE username = 'admin';

-- 更新 user 用户密码为 user123
UPDATE sys_user 
SET password = '$2a$10$8LLkeFX41.0BCdnYSdP8j.23SCbDsV2tP8LuVV5vpCpgNnCwGwgde',
    is_default_password = 1,
    status = 1
WHERE username = 'user';

-- 验证更新结果
SELECT username, name, role, is_default_password, status FROM sys_user WHERE username IN ('admin', 'user');

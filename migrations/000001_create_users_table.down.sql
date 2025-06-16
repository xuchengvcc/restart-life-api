-- 删除索引
DROP INDEX IF EXISTS idx_users_username ON users;
DROP INDEX IF EXISTS idx_users_email ON users;
DROP INDEX IF EXISTS idx_users_created_at ON users;
DROP INDEX IF EXISTS idx_users_is_active ON users;

-- 删除用户表
DROP TABLE IF EXISTS users; 
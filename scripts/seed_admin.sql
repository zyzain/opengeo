-- ============================================================
-- OpenGEO 超级管理员初始化脚本
-- 默认账号: admin
-- 默认密码: Admin@123456
-- ============================================================

USE opengeo;

-- 插入超级管理员用户（tenant_id=0 表示超级管理员）
-- 密码: Admin@123456 (bcrypt hash)
INSERT INTO users (tenant_id, username, password, email, status) VALUES
(0, 'admin', '$2a$10$EKD2ZjDX11ocpxA1V/3ZgOwkMUScdJAom6MdDhwOMsnZN4bNA4Aka', 'admin@opengeo.com', 1)
ON DUPLICATE KEY UPDATE password=VALUES(password), email=VALUES(email);

-- 获取管理员用户ID
SET @admin_id = (SELECT id FROM users WHERE username = 'admin' LIMIT 1);

-- 分配管理员角色
INSERT INTO user_roles (user_id, role_id)
SELECT @admin_id, id FROM roles WHERE name = 'admin'
ON DUPLICATE KEY UPDATE user_id=VALUES(user_id);

-- 授予所有权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p WHERE r.name = 'admin'
ON DUPLICATE KEY UPDATE role_id=VALUES(role_id);

-- 验证
SELECT u.id, u.username, u.email, u.status, r.name as role
FROM users u
LEFT JOIN user_roles ur ON u.id = ur.user_id
LEFT JOIN roles r ON ur.role_id = r.id
WHERE u.username = 'admin';

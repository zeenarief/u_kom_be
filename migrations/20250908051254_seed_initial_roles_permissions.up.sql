-- Seed initial roles, permissions, and admin user

-- 1. Insert admin role
INSERT INTO roles (id, name, description, is_default, created_at, updated_at)
VALUES (UUID(), 'admin', 'Administrator role', FALSE, NOW(), NOW());

-- 2. Insert user role (default)
INSERT INTO roles (id, name, description, is_default, created_at, updated_at)
VALUES (UUID(), 'user', 'Regular user role', TRUE, NOW(), NOW());

-- 3. Insert global permissions (hanya sekali saja)
INSERT INTO permissions (id, name, description, created_at, updated_at) VALUES
                                                                            (UUID(), 'users.read', 'Read all users', NOW(), NOW()),
                                                                            (UUID(), 'users.create', 'Create new users', NOW(), NOW()),
                                                                            (UUID(), 'users.update', 'Update users', NOW(), NOW()),
                                                                            (UUID(), 'users.delete', 'Delete users', NOW(), NOW()),
                                                                            (UUID(), 'users.manage_roles', 'Manage user roles', NOW(), NOW()),
                                                                            (UUID(), 'users.manage_permissions', 'Manage user permissions', NOW(), NOW()),
                                                                            (UUID(), 'roles.manage', 'Manage roles', NOW(), NOW()),
                                                                            (UUID(), 'permissions.manage', 'Manage permissions', NOW(), NOW()),
                                                                            (UUID(), 'profile.read', 'Read own profile', NOW(), NOW()),
                                                                            (UUID(), 'profile.update', 'Update own profile', NOW(), NOW()),
                                                                            (UUID(), 'auth.logout', 'Logout from system', NOW(), NOW());

-- 4. Assign all permissions to admin role
INSERT INTO role_permission (role_id, permission_id, created_at)
SELECT
    (SELECT id FROM roles WHERE name = 'admin'),
    p.id,
    NOW()
FROM permissions p;

-- 5. Assign only limited permissions to user role
INSERT INTO role_permission (role_id, permission_id, created_at)
SELECT
    (SELECT id FROM roles WHERE name = 'user'),
    p.id,
    NOW()
FROM permissions p
WHERE p.name IN ('profile.read', 'profile.update', 'auth.logout');

-- 6. Insert default admin user
-- NOTE: ubah hash password sesuai hasil bcrypt dari utils.HashPassword("password_anda")
INSERT INTO users (id, username, name, email, password, created_at, updated_at)
VALUES (UUID(), 'admin', 'Super Admin', 'admin@example.com', '$2a$10$Y4ZQaUO.VTUMoYJJSU3VYe2UIRfDg./SqdbQ71E8gm2CHavcUMx42', NOW(), NOW());

-- 7. Assign admin role to admin user
INSERT INTO user_role (user_id, role_id, created_at)
VALUES (
           (SELECT id FROM users WHERE username = 'admin'),
           (SELECT id FROM roles WHERE name = 'admin'),
           NOW()
       );

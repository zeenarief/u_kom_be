-- Seed initial roles, permissions, and admin user

-- Insert admin role
INSERT INTO roles (id, name, description, is_default, created_at, updated_at)
VALUES (UUID(), 'admin', 'Administrator role', FALSE, NOW(), NOW());

-- Insert default permissions
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

-- Assign all permissions to admin role
INSERT INTO role_permission (role_id, permission_id, created_at)
SELECT
    (SELECT id FROM roles WHERE name = 'admin'),
    p.id,
    NOW()
FROM permissions p;

-- Insert default admin user
-- Ganti <HASHED_PASSWORD> dengan hasil bcrypt/generate hash sesuai implementasi utils.HashPassword
INSERT INTO users (id, username, name, email, password, created_at, updated_at)
# VALUES (UUID(), 'admin', 'Super Admin', 'admin@example.com', '<HASHED_PASSWORD>', NOW(), NOW());
VALUES (UUID(), 'admin', 'Super Admin', 'admin@example.com', '$2a$10$akAv3b3Pf2xxmrtdBT4mgO6/jiSngVraG01v7dAyNyKaYzXU1Q1.2', NOW(), NOW());

-- Assign admin role to admin user
INSERT INTO user_role (user_id, role_id, created_at)
VALUES (
           (SELECT id FROM users WHERE username = 'admin'),
           (SELECT id FROM roles WHERE name = 'admin'),
           NOW()
       );

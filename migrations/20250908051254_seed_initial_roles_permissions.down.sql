-- Remove admin user and its relations
DELETE FROM user_role WHERE user_id = (SELECT id FROM users WHERE username = 'admin');
DELETE FROM users WHERE username = 'admin';

-- Remove admin role and its relations
DELETE FROM role_permission WHERE role_id = (SELECT id FROM roles WHERE name = 'admin');
DELETE FROM roles WHERE name = 'admin';

-- Remove seeded permissions
DELETE FROM permissions
WHERE name LIKE 'users.%'
   OR name LIKE 'roles.%'
   OR name LIKE 'permissions.%'
   OR name LIKE 'profile.%'
   OR name LIKE 'auth.logout';

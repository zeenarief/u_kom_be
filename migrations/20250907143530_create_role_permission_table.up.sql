-- Tabel roles
CREATE TABLE roles (
                       id CHAR(36) PRIMARY KEY,
                       name VARCHAR(255) NOT NULL UNIQUE,
                       description TEXT,
                       is_default BOOLEAN DEFAULT FALSE,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Tabel permissions
CREATE TABLE permissions (
                             id CHAR(36) PRIMARY KEY,
                             name VARCHAR(255) NOT NULL UNIQUE,
                             description TEXT,
                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                             updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Tabel role_permission (many-to-many)
CREATE TABLE role_permission (
                                 role_id CHAR(36) NOT NULL,
                                 permission_id CHAR(36) NOT NULL,
                                 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                 PRIMARY KEY (role_id, permission_id),
                                 FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
                                 FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
);

-- Tabel user_role (many-to-many)
CREATE TABLE user_role (
                           user_id CHAR(36) NOT NULL,
                           role_id CHAR(36) NOT NULL,
                           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                           PRIMARY KEY (user_id, role_id),
                           FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
                           FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);

-- Tabel user_permission (many-to-many untuk permission insidental)
CREATE TABLE user_permission (
                                 user_id CHAR(36) NOT NULL,
                                 permission_id CHAR(36) NOT NULL,
                                 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                 PRIMARY KEY (user_id, permission_id),
                                 FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
                                 FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
);
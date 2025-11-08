ALTER TABLE parents
    ADD COLUMN user_id CHAR(36) NULL UNIQUE,
    ADD CONSTRAINT fk_parents_user_id
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL;
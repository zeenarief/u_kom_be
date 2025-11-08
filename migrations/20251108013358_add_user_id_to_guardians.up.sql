ALTER TABLE guardians
    ADD COLUMN user_id CHAR(36) NULL UNIQUE,
    ADD CONSTRAINT fk_guardians_user_id
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL;
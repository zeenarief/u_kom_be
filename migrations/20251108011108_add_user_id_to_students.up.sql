ALTER TABLE students
    ADD COLUMN user_id CHAR(36) NULL UNIQUE,
    ADD CONSTRAINT fk_students_user_id
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL;
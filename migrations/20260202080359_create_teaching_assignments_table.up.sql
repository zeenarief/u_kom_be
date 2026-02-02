CREATE TABLE IF NOT EXISTS teaching_assignments (
                                                    id CHAR(36) PRIMARY KEY,
                                                    classroom_id CHAR(36) NOT NULL,
                                                    subject_id CHAR(36) NOT NULL,
                                                    teacher_id CHAR(36) NOT NULL, -- FK ke employees
                                                    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                                    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

                                                    FOREIGN KEY (classroom_id) REFERENCES classrooms(id) ON DELETE CASCADE,
                                                    FOREIGN KEY (subject_id) REFERENCES subjects(id) ON DELETE CASCADE,
                                                    FOREIGN KEY (teacher_id) REFERENCES employees(id) ON DELETE CASCADE,

    -- Constraint: Satu mapel di satu kelas hanya boleh 1 record (Guru)
                                                    UNIQUE KEY unique_class_subject (classroom_id, subject_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
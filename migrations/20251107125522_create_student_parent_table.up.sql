CREATE TABLE student_parent (
                                student_id CHAR(36) NOT NULL,
                                parent_id CHAR(36) NOT NULL,

    -- Kolom tambahan untuk data relasi
                                relationship_type VARCHAR(50) NOT NULL, -- Cth: 'FATHER', 'MOTHER', 'STEP-MOTHER'

                                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- Kunci utama gabungan
                                PRIMARY KEY (student_id, parent_id),

    -- Foreign keys
                                FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE,
                                FOREIGN KEY (parent_id) REFERENCES parents(id) ON DELETE CASCADE
);
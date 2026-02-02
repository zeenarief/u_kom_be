-- Tabel Classroom (Rombel)
CREATE TABLE IF NOT EXISTS classrooms (
                                          id CHAR(36) PRIMARY KEY,
                                          academic_year_id CHAR(36) NOT NULL,
                                          homeroom_teacher_id CHAR(36), -- Bisa NULL jika belum ada wali kelas
                                          name VARCHAR(50) NOT NULL, -- Contoh: X-IPA-1
                                          level VARCHAR(10) NOT NULL, -- Contoh: 10, 11, 12
                                          major VARCHAR(50), -- Contoh: IPA, IPS, TKJ (Bisa NULL untuk SD/SMP)
                                          description TEXT,
                                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                          updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

                                          FOREIGN KEY (academic_year_id) REFERENCES academic_years(id) ON DELETE CASCADE,
                                          FOREIGN KEY (homeroom_teacher_id) REFERENCES employees(id) ON DELETE SET NULL,

                                          INDEX idx_classroom_ay (academic_year_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Tabel Pivot Student Classrooms (History Kelas Siswa)
CREATE TABLE IF NOT EXISTS student_classrooms (
                                                  id CHAR(36) PRIMARY KEY,
                                                  classroom_id CHAR(36) NOT NULL,
                                                  student_id CHAR(36) NOT NULL,
                                                  status VARCHAR(20) DEFAULT 'ACTIVE', -- ACTIVE, TRANSFERRED, PROMOTED, RETAINED
                                                  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                                  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

                                                  FOREIGN KEY (classroom_id) REFERENCES classrooms(id) ON DELETE CASCADE,
                                                  FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE,

    -- Constraint: Satu siswa hanya boleh ada 1x di kelas yang sama
                                                  UNIQUE KEY unique_student_classroom (classroom_id, student_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
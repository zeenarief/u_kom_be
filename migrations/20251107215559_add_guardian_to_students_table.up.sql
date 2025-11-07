-- Menambahkan dua kolom baru ke tabel 'students'
ALTER TABLE students
    ADD COLUMN guardian_id CHAR(36) NULL,
    ADD COLUMN guardian_type VARCHAR(20) NULL;

-- Menambahkan index untuk mempercepat pencarian
CREATE INDEX idx_student_guardian ON students (guardian_type, guardian_id);
-- Tabel Sesi Presensi (Jurnal Mengajar)
CREATE TABLE IF NOT EXISTS attendance_sessions (
                                                   id CHAR(36) PRIMARY KEY,
                                                   schedule_id CHAR(36) NOT NULL,      -- Referensi ke jadwal (Mapel & Jam)
                                                   date DATE NOT NULL,                 -- Tanggal pertemuan
                                                   topic TEXT,                         -- Materi/Topik yang diajarkan (Jurnal)
                                                   notes TEXT,                         -- Catatan Guru (misal: "Kelas ribut")
                                                   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                                   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

                                                   FOREIGN KEY (schedule_id) REFERENCES schedules(id) ON DELETE CASCADE,

    -- Constraint: Satu jadwal hanya boleh ada 1 sesi per tanggal (Mencegah dobel absen)
                                                   UNIQUE KEY unique_session_date (schedule_id, date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Tabel Detail Presensi Siswa
CREATE TABLE IF NOT EXISTS attendance_details (
                                                  id CHAR(36) PRIMARY KEY,
                                                  attendance_session_id CHAR(36) NOT NULL,
                                                  student_id CHAR(36) NOT NULL,
                                                  status VARCHAR(10) NOT NULL DEFAULT 'PRESENT', -- PRESENT, SICK, PERMISSION, ABSENT
                                                  notes VARCHAR(255),                            -- Keterangan (misal: "Sakit Demam")
                                                  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                                  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

                                                  FOREIGN KEY (attendance_session_id) REFERENCES attendance_sessions(id) ON DELETE CASCADE,
                                                  FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE,

    -- Index untuk rekapitulasi cepat
                                                  INDEX idx_attendance_student (student_id, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
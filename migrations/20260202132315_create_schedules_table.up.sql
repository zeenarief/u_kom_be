CREATE TABLE IF NOT EXISTS schedules (
                                         id CHAR(36) PRIMARY KEY,
                                         teaching_assignment_id CHAR(36) NOT NULL,
                                         day_of_week TINYINT NOT NULL, -- 1=Senin, 2=Selasa, ..., 7=Minggu
                                         start_time TIME NOT NULL,     -- Format '07:30:00'
                                         end_time TIME NOT NULL,       -- Format '09:00:00'
                                         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                         updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

                                         FOREIGN KEY (teaching_assignment_id) REFERENCES teaching_assignments(id) ON DELETE CASCADE,

    -- Index untuk optimasi query validasi bentrok
                                         INDEX idx_schedule_day (day_of_week),
                                         INDEX idx_schedule_assignment (teaching_assignment_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
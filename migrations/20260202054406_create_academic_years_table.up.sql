CREATE TABLE IF NOT EXISTS academic_years (
                                              id CHAR(36) PRIMARY KEY,
                                              name VARCHAR(50) NOT NULL,
                                              status VARCHAR(20) NOT NULL DEFAULT 'INACTIVE', -- ACTIVE / INACTIVE
                                              start_date DATE NOT NULL,
                                              end_date DATE NOT NULL,
                                              created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                              updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

                                              INDEX idx_academic_year_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
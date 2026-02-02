CREATE TABLE IF NOT EXISTS subjects (
                                        id CHAR(36) PRIMARY KEY,
                                        code VARCHAR(20) NOT NULL UNIQUE, -- Kode Mapel (cth: MTK, IND, IPA)
                                        name VARCHAR(100) NOT NULL,       -- Nama Mapel (cth: Matematika Wajib)
                                        type VARCHAR(50),                 -- Jenis (cth: Muatan Nasional, Muatan Kewilayahan, Produktif)
                                        description TEXT,
                                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
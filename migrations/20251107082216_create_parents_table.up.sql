CREATE TABLE parents (
                         id CHAR(36) PRIMARY KEY,
                         full_name VARCHAR(100) NOT NULL,

    -- Data sensitif, akan dienkripsi
                         nik TEXT,
                         nik_hash VARCHAR(64) UNIQUE,

    -- Data demografi dasar
                         gender VARCHAR(10), -- 'male' atau 'female'
                         place_of_birth VARCHAR(100),
                         date_of_birth DATE,

    -- Status yang Anda minta
                         life_status VARCHAR(10) NOT NULL DEFAULT 'alive', -- 'alive', 'deceased'
                         marital_status VARCHAR(10), -- 'married', 'divorced', 'widowed'

    -- Informasi kontak
                         phone_number VARCHAR(20) UNIQUE,
                         email VARCHAR(100) UNIQUE,

    -- Informasi tambahan (opsional tapi berguna)
                         education_level VARCHAR(50), -- 'SMA', 'S1', 'Tidak Sekolah'
                         occupation VARCHAR(100), -- 'PNS', 'Wiraswasta', 'Ibu Rumah Tangga'
                         income_range VARCHAR(50), -- '< 1jt', '1-3jt', '> 10jt'

    -- Alamat (kita samakan dengan Student)
                         address TEXT,
                         rt VARCHAR(3),
                         rw VARCHAR(3),
                         sub_district VARCHAR(100),
                         district VARCHAR(100),
                         city VARCHAR(100),
                         province VARCHAR(100),
                         postal_code VARCHAR(5),

    -- Timestamps
                         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                         updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
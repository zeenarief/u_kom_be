CREATE TABLE guardians (
                           id CHAR(36) PRIMARY KEY,
                           full_name VARCHAR(100) NOT NULL,

    -- Data sensitif, akan dienkripsi
                           nik TEXT,
                           nik_hash VARCHAR(64) UNIQUE,

    -- Info dasar
                           gender VARCHAR(10), -- 'male' atau 'female'

    -- Info kontak krusial
                           phone_number VARCHAR(20) UNIQUE NOT NULL,
                           email VARCHAR(100) UNIQUE,

    -- Info alamat
                           address TEXT,
                           rt VARCHAR(3),
                           rw VARCHAR(3),
                           sub_district VARCHAR(100),
                           district VARCHAR(100),
                           city VARCHAR(100),
                           province VARCHAR(100),
                           postal_code VARCHAR(5),

    -- Field ini penting untuk tahu relasinya
                           relationship_to_student VARCHAR(50), -- cth: 'Paman', 'Bibi', 'Kakak', 'Kakek'

    -- Timestamps
                           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
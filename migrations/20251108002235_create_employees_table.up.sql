CREATE TABLE employees (
                           id CHAR(36) PRIMARY KEY,

    -- Foreign Key ke tabel 'users'
    -- 'NULLABLE' karena pegawai bisa didata tanpa punya akun
    -- 'UNIQUE' karena satu akun user hanya bisa untuk satu pegawai
                           user_id CHAR(36) NULL UNIQUE,

    -- Data Profil Pegawai
                           full_name VARCHAR(100) NOT NULL,
                           nip VARCHAR(50) UNIQUE, -- Nomor Induk Pegawai
                           job_title VARCHAR(100) NOT NULL, -- Cth: 'Guru Matematika', 'Admin TU', 'Satpam'

    -- Data sensitif (akan dienkripsi)
                           nik TEXT,

    -- Info dasar & Kontak
                           gender VARCHAR(10), -- 'male', 'female'
                           phone_number VARCHAR(20) UNIQUE,
                           address TEXT,
                           date_of_birth DATE,
                           join_date DATE, -- Tanggal mulai bergabung
                           employment_status VARCHAR(20), -- Cth: 'Full-time', 'Part-time', 'Contract'

    -- Timestamps
                           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    -- Constraint Foreign Key
    -- ON DELETE SET NULL: Jika akun 'user' dihapus, data 'employee' tetap ada,
    -- tapi 'user_id' di-set ke NULL.
                           FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);
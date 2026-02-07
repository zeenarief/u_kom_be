CREATE TABLE students (
                          id CHAR(36) PRIMARY KEY,
                          full_name VARCHAR(100) NOT NULL,

                          no_kk TEXT,
                          nik TEXT,
                          nik_hash VARCHAR(64) UNIQUE,

                          nisn VARCHAR(20) UNIQUE,
                          nim VARCHAR(20) UNIQUE,

                          gender VARCHAR(10),

                          place_of_birth VARCHAR(100),
                          date_of_birth DATE,
                          address TEXT,

                          rt VARCHAR(3),
                          rw VARCHAR(3),

                          sub_district VARCHAR(100),
                          district VARCHAR(100),
                          city VARCHAR(100),
                          province VARCHAR(100),

                          postal_code VARCHAR(5),

                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
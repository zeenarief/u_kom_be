ALTER TABLE employees ADD COLUMN nik_hash VARCHAR(64) UNIQUE;
CREATE INDEX idx_employees_nik_hash ON employees(nik_hash);

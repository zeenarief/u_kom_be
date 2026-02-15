CREATE TABLE IF NOT EXISTS violation_categories (
    id CHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)
);

CREATE TABLE IF NOT EXISTS violation_types (
    id CHAR(36) PRIMARY KEY,
    category_id CHAR(36) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    default_points INT NOT NULL DEFAULT 0,
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    
    FOREIGN KEY (category_id) REFERENCES violation_categories(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS student_violations (
    id CHAR(36) PRIMARY KEY,
    student_id CHAR(36) NOT NULL,
    violation_type_id CHAR(36) NOT NULL,
    violation_date DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    points INT NOT NULL, -- Snapshot of points at time of violation
    action_taken TEXT,
    notes TEXT,
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    
    FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE,
    FOREIGN KEY (violation_type_id) REFERENCES violation_types(id) ON DELETE RESTRICT
);

CREATE INDEX idx_violation_types_category_id ON violation_types(category_id);
CREATE INDEX idx_student_violations_student_id ON student_violations(student_id);
CREATE INDEX idx_student_violations_violation_type_id ON student_violations(violation_type_id);

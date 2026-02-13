CREATE TABLE IF NOT EXISTS assessments (
    id CHAR(36) PRIMARY KEY,
    teaching_assignment_id CHAR(36) NOT NULL,
    title VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    max_score INT DEFAULT 100,
    date DATE NOT NULL,
    description TEXT,
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    
    FOREIGN KEY (teaching_assignment_id) REFERENCES teaching_assignments(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS student_scores (
    id CHAR(36) PRIMARY KEY,
    assessment_id CHAR(36) NOT NULL,
    student_id CHAR(36) NOT NULL,
    score DOUBLE DEFAULT 0,
    feedback TEXT,
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    
    FOREIGN KEY (assessment_id) REFERENCES assessments(id) ON DELETE CASCADE,
    FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE,
    
    UNIQUE KEY idx_assessment_student (assessment_id, student_id)
);

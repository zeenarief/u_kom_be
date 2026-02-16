CREATE TABLE IF NOT EXISTS finance_donors (
    id CHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    email VARCHAR(255),
    address TEXT,
    notes TEXT,
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    
    INDEX idx_finance_donors_phone (phone),
    INDEX idx_finance_donors_name (name)
);

CREATE TABLE IF NOT EXISTS finance_donations (
    id CHAR(36) PRIMARY KEY,
    donor_id CHAR(36) NOT NULL,
    employee_id CHAR(36) NOT NULL, -- Receiver (Staff/Admin)
    date DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    type ENUM('MONEY', 'GOODS', 'MIXED') NOT NULL,
    payment_method ENUM('CASH', 'TRANSFER', 'QRIS', 'GOODS') NOT NULL,
    total_amount DECIMAL(15, 2) DEFAULT 0, -- For money donations
    proof_file VARCHAR(255), -- Path to photo/receipt
    description TEXT,
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    
    FOREIGN KEY (donor_id) REFERENCES finance_donors(id) ON DELETE RESTRICT,
    FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE RESTRICT,
    
    INDEX idx_finance_donations_date (date),
    INDEX idx_finance_donations_donor (donor_id),
    INDEX idx_finance_donations_employee (employee_id)
);

CREATE TABLE IF NOT EXISTS finance_donation_items (
    id CHAR(36) PRIMARY KEY,
    donation_id CHAR(36) NOT NULL,
    item_name VARCHAR(255) NOT NULL,
    quantity DECIMAL(10, 2) NOT NULL DEFAULT 1,
    unit VARCHAR(50) NOT NULL, -- kg, liter, pcs, box, etc
    estimated_value DECIMAL(15, 2) DEFAULT 0, -- Optional valuation
    notes TEXT,
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    
    FOREIGN KEY (donation_id) REFERENCES finance_donations(id) ON DELETE CASCADE,
    
    INDEX idx_finance_donation_items_donation (donation_id)
);

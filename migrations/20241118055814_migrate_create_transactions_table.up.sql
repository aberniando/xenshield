CREATE TABLE IF NOT EXISTS transactions (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    ip_address VARCHAR(15) NOT NULL,
    masked_card_number VARCHAR(19) NOT NULL,
    status VARCHAR(7) NOT NULL,
    reason VARCHAR(20) NULL,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_ip_address ON transactions (ip_address);

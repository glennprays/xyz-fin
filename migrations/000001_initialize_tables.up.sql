CREATE TABLE consumers (
    nik VARCHAR(16) PRIMARY KEY,
    phone_number VARCHAR(20) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL, 
    full_name VARCHAR(255) NOT NULL,
    legal_name VARCHAR(255),
    tempat_lahir VARCHAR(100) NOT NULL, 
    tanggal_lahir DATE NOT NULL, 
    gaji NUMERIC(15, 2) NOT NULL,
    foto_ktp_path VARCHAR(512), 
    foto_selfie_path VARCHAR(512), 
    created_at TIMESTAMP NOT NULL DEFAULT NOW(), 
    updated_at TIMESTAMP NOT NULL DEFAULT NOW() 
);

CREATE INDEX idx_consumers_phone_number ON consumers (phone_number); 

CREATE TABLE consumer_limits (
    consumer_nik VARCHAR(16) REFERENCES consumers(nik) ON DELETE CASCADE, 
    tenor INT NOT NULL, 
    limit_amount NUMERIC(15, 2) NOT NULL, 
    created_at TIMESTAMP NOT NULL DEFAULT NOW(), 
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (consumer_nik, tenor) 
);

CREATE TABLE transactions (
    nomor_kontrak VARCHAR(100) PRIMARY KEY,
    consumer_nik VARCHAR(16) REFERENCES consumers(nik) ON DELETE CASCADE NOT NULL,
    otr NUMERIC(15, 2) NOT NULL, 
    admin_fee NUMERIC(15, 2) NOT NULL,
    jumlah_cicilan INT NOT NULL, 
    jumlah_bunga NUMERIC(15, 2) NOT NULL, 
    nama_asset VARCHAR(255) NOT NULL, 
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE', 
    created_at TIMESTAMP NOT NULL DEFAULT NOW(), 
    updated_at TIMESTAMP NOT NULL DEFAULT NOW() 
);

CREATE INDEX idx_transactions_consumer_nik ON transactions (consumer_nik);
CREATE INDEX idx_transactions_consumer_nik_status ON transactions (consumer_nik, status);

CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
NEW.updated_at = NOW();
RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_consumers_timestamp BEFORE UPDATE ON consumers FOR EACH ROW EXECUTE PROCEDURE update_timestamp();
CREATE TRIGGER update_consumer_limits_timestamp BEFORE UPDATE ON consumer_limits FOR EACH ROW EXECUTE PROCEDURE update_timestamp();
CREATE TRIGGER update_transactions_timestamp BEFORE UPDATE ON transactions FOR EACH ROW EXECUTE PROCEDURE update_timestamp();

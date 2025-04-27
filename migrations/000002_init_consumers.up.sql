INSERT INTO consumers (
    nik,
    phone_number,
    password_hash,
    full_name,
    legal_name,
    tempat_lahir,
    tanggal_lahir,
    gaji,
    foto_ktp_path,
    foto_selfie_path
) VALUES (
    '1111',
    '081234567890',     
    '$argon2id$v=19$m=65536,t=3,p=4$BEMOnl7KwywiS5LHLqtmJg$Cm5gGNJjzFo0navnegsRAJIFvPd/nI577+STF/t2z8E',
    'Budi',            
    'Budi Santoso',   
    'Jakarta',       
    '1990-05-15',   
    5000000.00,    
    '/dummy/ktp/budi_ktp.jpg', 
    '/dummy/selfie/budi_selfie.jpg' 
);

INSERT INTO consumers (
    nik,
    phone_number,
    password_hash,
    full_name,
    legal_name,
    tempat_lahir,
    tanggal_lahir,
    gaji,
    foto_ktp_path,
    foto_selfie_path
) VALUES (
    '2222', 
    '087654321098',     
    '$argon2id$v=19$m=65536,t=3,p=4$BEMOnl7KwywiS5LHLqtmJg$Cm5gGNJjzFo0navnegsRAJIFvPd/nI577+STF/t2z8E',
    'Annisa',           
    'Annisa Putri',    
    'Bandung',        
    '1992-08-22',    
    8000000.00,     
    '/dummy/ktp/annisa_ktp.png',
    '/dummy/selfie/annisa_selfie.png' 
);

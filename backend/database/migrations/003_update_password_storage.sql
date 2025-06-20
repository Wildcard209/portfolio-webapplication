-- Version 1: SHA-256 + bcrypt with salt and pepper (legacy)
-- Version 2: bcrypt only (new secure method)
ALTER TABLE admins ADD COLUMN IF NOT EXISTS hash_version INTEGER DEFAULT 1;

ALTER TABLE admins ALTER COLUMN password_salt DROP NOT NULL;

CREATE INDEX IF NOT EXISTS idx_admins_hash_version ON admins(hash_version);

UPDATE admins SET hash_version = 1 WHERE hash_version IS NULL;

COMMENT ON COLUMN admins.hash_version IS 'Password hashing method version: 1=legacy SHA256+bcrypt+salt+pepper, 2=bcrypt-only';
COMMENT ON COLUMN admins.password_salt IS 'Salt for legacy password hashing (version 1 only)';

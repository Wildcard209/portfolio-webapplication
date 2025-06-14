SELECT id, username, password_hash, password_salt, last_login, current_token, token_expiration, created_at, updated_at
FROM admins 
WHERE current_token = $1 AND token_expiration > CURRENT_TIMESTAMP;

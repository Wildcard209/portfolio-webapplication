SELECT id, username, password_hash, password_salt, hash_version, last_login, current_token, token_expiration, created_at, updated_at
FROM admins 
WHERE username = $1;

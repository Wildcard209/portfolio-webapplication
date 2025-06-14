INSERT INTO admins (username, password_hash, password_salt, created_at, updated_at) 
VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING id, username, password_hash, password_salt, last_login, current_token, token_expiration, created_at, updated_at;

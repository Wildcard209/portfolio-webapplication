UPDATE admins 
SET current_token = $1, token_expiration = $2, last_login = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = $3;

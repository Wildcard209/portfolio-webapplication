UPDATE admins 
SET current_token = NULL, token_expiration = NULL, updated_at = CURRENT_TIMESTAMP
WHERE token_expiration < CURRENT_TIMESTAMP;

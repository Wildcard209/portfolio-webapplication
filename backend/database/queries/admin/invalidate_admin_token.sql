UPDATE admins 
SET current_token = NULL, token_expiration = NULL, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

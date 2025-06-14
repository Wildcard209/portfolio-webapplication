SELECT COUNT(*) 
FROM login_attempts 
WHERE ip_address = $1 AND success = FALSE AND attempt_at >= $2;

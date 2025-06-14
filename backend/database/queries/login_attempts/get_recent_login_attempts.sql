SELECT id, ip_address, user_agent, success, attempt_at, details
FROM login_attempts 
WHERE ip_address = $1 AND attempt_at >= $2
ORDER BY attempt_at DESC;

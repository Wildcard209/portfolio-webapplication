INSERT INTO login_attempts (ip_address, user_agent, success, attempt_at, details)
VALUES ($1, $2, $3, CURRENT_TIMESTAMP, $4);

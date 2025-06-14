DELETE FROM login_attempts WHERE attempt_at < $1;

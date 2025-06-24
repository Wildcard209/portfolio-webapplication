'use client';

import { useState, useEffect } from 'react';
import { AuthService } from '../../../lib/auth/authService';
import { InputValidator } from '../../../lib/validation/inputValidator';

interface LoginFlowProps {
  adminToken: string;
}

const LoginFlow = ({ adminToken }: LoginFlowProps) => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [validationErrors, setValidationErrors] = useState<{
    username?: string;
    password?: string;
  }>({});

  useEffect(() => {
    const checkAuthStatus = async () => {
      const isAuth = await AuthService.checkAuthStatus();
      if (isAuth) {
        setIsLoggedIn(true);
      }
    };

    checkAuthStatus();
  }, [adminToken]);

  const validateInputs = (): boolean => {
    const errors: { username?: string; password?: string } = {};

    const usernameValidation = InputValidator.validateUsername(username);
    if (!usernameValidation.isValid) {
      errors.username = usernameValidation.error;
    }

    const passwordValidation = InputValidator.validateString(password, 'Password', {
      required: true,
      minLength: 1,
      maxLength: 128,
    });
    if (!passwordValidation.isValid) {
      errors.password = passwordValidation.error;
    }

    setValidationErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const handleInputChange = (field: 'username' | 'password', value: string) => {
    const sanitizedValue = InputValidator.sanitizeInput(value);

    if (field === 'username') {
      setUsername(sanitizedValue);
      if (validationErrors.username) {
        setValidationErrors(prev => ({ ...prev, username: undefined }));
      }
    } else {
      setPassword(value);
      if (validationErrors.password) {
        setValidationErrors(prev => ({ ...prev, password: undefined }));
      }
    }
  };

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    if (!validateInputs()) {
      return;
    }

    setIsLoading(true);

    try {
      const result = await AuthService.login(username, password);

      if (result.success) {
        setIsLoggedIn(true);
      } else {
        setError(result.error || 'Login failed');
      }      } catch {
        setError('An unexpected error occurred. Please try again.');
      }finally {
      setIsLoading(false);
    }
  };

  const handleLogout = async () => {
    setIsLoading(true);
    await AuthService.logout();
    setIsLoggedIn(false);
    setUsername('');
    setPassword('');
    setError('');
    setIsLoading(false);
  };

  if (isLoggedIn) {
    const user = AuthService.getUser();
    return (
      <div style={{ textAlign: 'center', padding: '20px' }}>
        <h2>Admin Panel</h2>
        <p>Welcome back, {user?.username}!</p>
        <p>You are now authenticated and can access admin features.</p>
        <button
          onClick={handleLogout}
          disabled={isLoading}
          style={{
            padding: '10px 20px',
            backgroundColor: '#dc3545',
            color: 'white',
            border: 'none',
            borderRadius: '4px',
            cursor: isLoading ? 'not-allowed' : 'pointer',
          }}
        >
          {isLoading ? 'Logging out...' : 'Logout'}
        </button>
      </div>
    );
  }

  return (
    <div style={{ maxWidth: '400px', margin: '0 auto', padding: '20px' }}>
      <form onSubmit={handleLogin} style={{ textAlign: 'center' }}>
        <h2>Admin Login</h2>

        {error && (
          <div
            style={{
              color: '#dc3545',
              backgroundColor: '#f8d7da',
              border: '1px solid #f5c6cb',
              borderRadius: '4px',
              padding: '10px',
              marginBottom: '15px',
            }}
          >
            {error}
          </div>
        )}

        <div style={{ marginBottom: '15px' }}>
          <label style={{ display: 'block', marginBottom: '5px', fontWeight: 'bold' }}>
            Username:
          </label>
          <input
            type="text"
            value={username}
            onChange={e => handleInputChange('username', e.target.value)}
            required
            disabled={isLoading}
            style={{
              width: '100%',
              padding: '8px',
              border: `1px solid ${validationErrors.username ? '#dc3545' : '#ccc'}`,
              borderRadius: '4px',
              fontSize: '16px',
            }}
          />
          {validationErrors.username && (
            <div
              style={{
                color: '#dc3545',
                fontSize: '12px',
                marginTop: '5px',
              }}
            >
              {validationErrors.username}
            </div>
          )}
        </div>

        <div style={{ marginBottom: '20px' }}>
          <label style={{ display: 'block', marginBottom: '5px', fontWeight: 'bold' }}>
            Password:
          </label>
          <input
            type="password"
            value={password}
            onChange={e => handleInputChange('password', e.target.value)}
            required
            disabled={isLoading}
            style={{
              width: '100%',
              padding: '8px',
              border: `1px solid ${validationErrors.password ? '#dc3545' : '#ccc'}`,
              borderRadius: '4px',
              fontSize: '16px',
            }}
          />
          {validationErrors.password && (
            <div
              style={{
                color: '#dc3545',
                fontSize: '12px',
                marginTop: '5px',
              }}
            >
              {validationErrors.password}
            </div>
          )}
        </div>

        <button
          type="submit"
          disabled={isLoading}
          style={{
            width: '100%',
            padding: '12px',
            backgroundColor: isLoading ? '#6c757d' : '#007bff',
            color: 'white',
            border: 'none',
            borderRadius: '4px',
            fontSize: '16px',
            cursor: isLoading ? 'not-allowed' : 'pointer',
          }}
        >
          {isLoading ? 'Logging in...' : 'Login'}
        </button>
      </form>
    </div>
  );
};

export default LoginFlow;

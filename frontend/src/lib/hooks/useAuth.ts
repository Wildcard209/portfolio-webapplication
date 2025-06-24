'use client';

import { useState, useEffect } from 'react';
import { AuthService } from '../auth/authService';

export interface UseAuthReturn {
  isAuthenticated: boolean;
  user: any | null;
  login: (username: string, password: string) => Promise<{ success: boolean; error?: string }>;
  logout: () => Promise<{ success: boolean; error?: string }>;
  loading: boolean;
}

export const useAuth = (): UseAuthReturn => {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const checkAuth = () => {
      const authenticated = AuthService.isAuthenticated();
      const currentUser = AuthService.getUser();

      setIsAuthenticated(authenticated);
      setUser(currentUser);
      setLoading(false);
    };

    checkAuth();

    const interval = setInterval(checkAuth, 60000);

    return () => clearInterval(interval);
  }, []);

  const login = async (username: string, password: string) => {
    setLoading(true);
    const result = await AuthService.login(username, password);

    if (result.success) {
      setIsAuthenticated(true);
      setUser(AuthService.getUser());
    }

    setLoading(false);
    return result;
  };

  const logout = async () => {
    setLoading(true);
    const result = await AuthService.logout();

    setIsAuthenticated(false);
    setUser(null);
    setLoading(false);

    return result;
  };

  return {
    isAuthenticated,
    user,
    login,
    logout,
    loading,
  };
};

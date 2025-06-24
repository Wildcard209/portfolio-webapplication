'use client';

import { useState, useEffect } from 'react';
import { useLogout, useCheckAuth } from '../api/hooks/useApi';

export interface UseAuthReturn {
  isAuthenticated: boolean;
  user: Record<string, unknown> | null;
  logout: () => Promise<{ success: boolean; error?: string }>;
  loading: boolean;
}

export const useAuth = (): UseAuthReturn => {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [user, setUser] = useState<Record<string, unknown> | null>(null);
  const [loading, setLoading] = useState(true);

  const { logout: apiLogout } = useLogout();
  const { checkAuth } = useCheckAuth();

  useEffect(() => {
    const getAuthStatus = async () => {
      const userStr = localStorage.getItem('auth_user');
      const currentUser = userStr ? JSON.parse(userStr) : null;

      if (currentUser) {
        const isValid = await checkAuth();
        setIsAuthenticated(isValid);
        setUser(isValid ? currentUser : null);
      } else {
        setIsAuthenticated(false);
        setUser(null);
      }

      setLoading(false);
    };

    getAuthStatus();

    const interval = setInterval(getAuthStatus, 60000);
    return () => clearInterval(interval);
  }, [checkAuth]);

  const logout = async (): Promise<{ success: boolean; error?: string }> => {
    setLoading(true);
    const result = await apiLogout();

    setIsAuthenticated(false);
    setUser(null);
    setLoading(false);

    return result;
  };

  return {
    isAuthenticated,
    user,
    logout,
    loading,
  };
};

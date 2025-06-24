'use client';

import { useState, useEffect } from 'react';
import { AuthService } from '../auth/authService';

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

  useEffect(() => {
    const checkAuth = (): void => {
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

  const logout = async (): Promise<{ success: boolean; error?: string }> => {
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
    logout,
    loading,
  };
};

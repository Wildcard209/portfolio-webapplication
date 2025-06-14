"use client";

import { useState, useEffect } from 'react';
import { AuthService } from '../auth/authService';

export interface UseAuthReturn {
  isAuthenticated: boolean;
  user: any | null;
  token: string | null;
  login: (username: string, password: string) => Promise<{ success: boolean; error?: string }>;
  logout: () => Promise<{ success: boolean; error?: string }>;
  loading: boolean;
}

export const useAuth = (): UseAuthReturn => {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [user, setUser] = useState(null);
  const [token, setToken] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Check authentication state on mount
    const checkAuth = () => {
      const authenticated = AuthService.isAuthenticated();
      const currentUser = AuthService.getUser();
      const currentToken = AuthService.getToken();

      setIsAuthenticated(authenticated);
      setUser(currentUser);
      setToken(currentToken);
      setLoading(false);
    };

    checkAuth();

    // Optional: Set up an interval to periodically check token validity
    const interval = setInterval(checkAuth, 60000); // Check every minute

    return () => clearInterval(interval);
  }, []);

  const login = async (username: string, password: string) => {
    setLoading(true);
    const result = await AuthService.login(username, password);
    
    if (result.success) {
      setIsAuthenticated(true);
      setUser(AuthService.getUser());
      setToken(AuthService.getToken());
    }
    
    setLoading(false);
    return result;
  };

  const logout = async () => {
    setLoading(true);
    const result = await AuthService.logout();
    
    setIsAuthenticated(false);
    setUser(null);
    setToken(null);
    setLoading(false);
    
    return result;
  };

  return {
    isAuthenticated,
    user,
    token,
    login,
    logout,
    loading,
  };
};

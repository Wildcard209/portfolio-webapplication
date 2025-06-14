interface LoginRequest {
  username: string;
  password: string;
}

interface LoginResponse {
  token: string;
  expiresAt: string;
  user: {
    id: number;
    username: string;
    lastLogin?: string;
  };
}

interface ErrorResponse {
  error: string;
  message?: string;
}

export class AuthService {
  private static readonly TOKEN_KEY = 'auth_token';
  private static readonly USER_KEY = 'auth_user';
  private static readonly ADMIN_TOKEN_KEY = 'admin_token';

  // Store the admin token (from URL parameter)
  static setAdminToken(token: string): void {
    if (typeof window !== 'undefined') {
      localStorage.setItem(this.ADMIN_TOKEN_KEY, token);
    }
  }

  // Get the admin token
  static getAdminToken(): string | null {
    if (typeof window !== 'undefined') {
      return localStorage.getItem(this.ADMIN_TOKEN_KEY);
    }
    return null;
  }

  // Store JWT token and user info
  static setAuthData(token: string, user: any): void {
    if (typeof window !== 'undefined') {
      localStorage.setItem(this.TOKEN_KEY, token);
      localStorage.setItem(this.USER_KEY, JSON.stringify(user));
    }
  }

  // Get stored JWT token
  static getToken(): string | null {
    if (typeof window !== 'undefined') {
      return localStorage.getItem(this.TOKEN_KEY);
    }
    return null;
  }

  // Get stored user info
  static getUser(): any | null {
    if (typeof window !== 'undefined') {
      const userStr = localStorage.getItem(this.USER_KEY);
      return userStr ? JSON.parse(userStr) : null;
    }
    return null;
  }

  // Check if user is authenticated
  static isAuthenticated(): boolean {
    const token = this.getToken();
    if (!token) return false;

    try {
      // Simple token expiration check (you could decode JWT for more precise check)
      const user = this.getUser();
      return !!user;
    } catch {
      return false;
    }
  }

  // Login with backend API
  static async login(username: string, password: string): Promise<{ success: boolean; error?: string }> {
    const adminToken = this.getAdminToken();
    if (!adminToken) {
      return { success: false, error: 'Admin token not found' };
    }

    try {
      const apiUrl = process.env.NEXT_PUBLIC_BASE_API_URL || 'http://localhost/api';
      const response = await fetch(`${apiUrl}/${adminToken}/admin/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username, password }),
      });

      if (response.ok) {
        const data: LoginResponse = await response.json();
        this.setAuthData(data.token, data.user);
        return { success: true };
      } else {
        const errorData: ErrorResponse = await response.json();
        return { 
          success: false, 
          error: errorData.message || errorData.error || 'Login failed' 
        };
      }
    } catch (error) {
      return { 
        success: false, 
        error: error instanceof Error ? error.message : 'Network error' 
      };
    }
  }

  // Logout with backend API
  static async logout(): Promise<{ success: boolean; error?: string }> {
    const adminToken = this.getAdminToken();
    const jwtToken = this.getToken();
    
    if (!adminToken || !jwtToken) {
      this.clearAuthData();
      return { success: true };
    }

    try {
      const apiUrl = process.env.NEXT_PUBLIC_BASE_API_URL || 'http://localhost/api';
      const response = await fetch(`${apiUrl}/${adminToken}/admin/logout`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${jwtToken}`,
        },
      });

      // Clear local data regardless of API response
      this.clearAuthData();

      if (response.ok) {
        return { success: true };
      } else {
        // Even if logout fails on backend, we still clear local data
        console.warn('Backend logout failed, but local data cleared');
        return { success: true };
      }
    } catch (error) {
      // Clear local data even on network error
      this.clearAuthData();
      return { success: true };
    }
  }

  // Clear all authentication data
  static clearAuthData(): void {
    if (typeof window !== 'undefined') {
      localStorage.removeItem(this.TOKEN_KEY);
      localStorage.removeItem(this.USER_KEY);
      localStorage.removeItem(this.ADMIN_TOKEN_KEY);
    }
  }

  // Get authorization header for API calls
  static getAuthHeader(): { Authorization: string } | {} {
    const token = this.getToken();
    return token ? { Authorization: `Bearer ${token}` } : {};
  }
}

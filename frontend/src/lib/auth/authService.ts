// Interface for login requests - currently unused but kept for future use
// interface LoginRequest {
//   username: string;
//   password: string;
// }

interface LoginResponse {
  token: string; // Will be empty now, kept for backward compatibility
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
  private static readonly USER_KEY = 'auth_user';

  // Store only user info in localStorage, no tokens
  static setUserData(user: Record<string, unknown>): void {
    if (typeof window !== 'undefined') {
      localStorage.setItem(this.USER_KEY, JSON.stringify(user));
    }
  }

  static getUser(): Record<string, unknown> | null {
    if (typeof window !== 'undefined') {
      const userStr = localStorage.getItem(this.USER_KEY);
      return userStr ? JSON.parse(userStr) : null;
    }
    return null;
  }

  static isAuthenticated(): boolean {
    // Check if user data exists and make a test request to verify session
    const user = this.getUser();
    return !!user;
  }

  static async checkAuthStatus(): Promise<boolean> {
    try {
      const apiUrl = process.env.NEXT_PUBLIC_BASE_API_URL || 'http://localhost/api';
      const response = await fetch(`${apiUrl}/admin/refresh`, {
        method: 'POST',
        credentials: 'include', // Important for cookies
      });

      if (response.ok) {
        return true;
      } else {
        // If refresh fails, clear user data
        this.clearAuthData();
        return false;
      }
    } catch {
      this.clearAuthData();
      return false;
    }
  }

  static async login(
    username: string,
    password: string
  ): Promise<{ success: boolean; error?: string }> {
    try {
      const apiUrl = process.env.NEXT_PUBLIC_BASE_API_URL || 'http://localhost/api';
      const response = await fetch(`${apiUrl}/admin/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include', // Important for cookies
        body: JSON.stringify({ username, password }),
      });

      if (response.ok) {
        const data: LoginResponse = await response.json();
        this.setUserData(data.user);
        return { success: true };
      } else {
        const errorData: ErrorResponse = await response.json();
        return {
          success: false,
          error: errorData.message || errorData.error || 'Login failed',
        };
      }
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Network error',
      };
    }
  }

  static async logout(): Promise<{ success: boolean; error?: string }> {
    try {
      const apiUrl = process.env.NEXT_PUBLIC_BASE_API_URL || 'http://localhost/api';
      const response = await fetch(`${apiUrl}/admin/logout`, {
        method: 'POST',
        credentials: 'include', // Important for cookies
      });

      this.clearAuthData();

      if (response.ok) {
        return { success: true };
      } else {
        console.warn('Backend logout failed, but local data cleared');
        return { success: true };
      }
    } catch {
      this.clearAuthData();
      return { success: true };
    }
  }

  static clearAuthData(): void {
    if (typeof window !== 'undefined') {
      localStorage.removeItem(this.USER_KEY);
    }
  }

  // For API calls that need authentication, no need to manually add headers
  // The cookies will be automatically included with credentials: 'include'
  static getRequestOptions(): RequestInit {
    return {
      credentials: 'include', // This will include HTTP-only cookies automatically
    };
  }
}

import { AuthService } from '../auth/authService';

type FetchOptions = {
  cache?: 'force-cache' | 'no-store';
  next?: {
    revalidate?: number;
    tags?: string[];
  };
  requireAuth?: boolean; // Whether this endpoint requires authentication
};

type ApiResponse<T> = {
  data: T | null;
  error: string | null;
  status: number;
};

const apiUrl = process.env.NEXT_PUBLIC_BASE_API_URL ?? 'http://localhost/api';

export class EnhancedApiHandler {
  private static getDefaultHeaders(requireAuth: boolean = false): HeadersInit {
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
    };

    // Always include auth token if available (as requested by user)
    if (typeof window !== 'undefined') {
      const token = AuthService.getToken();
      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      } else if (requireAuth) {
        // If auth is required but no token, this will fail
        console.warn('Authentication required but no token found');
      }
    }

    return headers;
  }

  private static async fetchWithErrorHandling<T>(
    url: string,
    options: RequestInit & FetchOptions = {}
  ): Promise<ApiResponse<T>> {
    const { requireAuth = false, ...fetchOptions } = options;

    try {
      const response = await fetch(url, {
        ...fetchOptions,
        headers: {
          ...this.getDefaultHeaders(requireAuth),
          ...options.headers,
        },
      });

      // Handle different response types
      let data = null;
      const contentType = response.headers.get('content-type');
      
      if (contentType && contentType.includes('application/json')) {
        data = await response.json();
      } else {
        // For non-JSON responses, get text
        const text = await response.text();
        data = text ? { message: text } : null;
      }

      if (!response.ok) {
        // If it's an auth error, might want to clear tokens
        if (response.status === 401) {
          AuthService.clearAuthData();
        }

        return {
          data: null,
          error: data?.error || data?.message || `HTTP error! status: ${response.status}`,
          status: response.status,
        };
      }

      return {
        data,
        error: null,
        status: response.status,
      };
    } catch (error) {
      return {
        data: null,
        error: error instanceof Error ? error.message : 'Network error occurred',
        status: 500,
      };
    }
  }

  static async get<T>(
    endpoint: string,
    options: FetchOptions = {}
  ): Promise<ApiResponse<T>> {
    const url = `${apiUrl}${endpoint}`;
    return this.fetchWithErrorHandling<T>(url, {
      method: 'GET',
      ...options,
    });
  }

  static async post<T>(
    endpoint: string,
    data?: unknown,
    options: FetchOptions = {}
  ): Promise<ApiResponse<T>> {
    const url = `${apiUrl}${endpoint}`;
    return this.fetchWithErrorHandling<T>(url, {
      method: 'POST',
      body: data ? JSON.stringify(data) : undefined,
      ...options,
    });
  }

  static async put<T>(
    endpoint: string,
    data: unknown,
    options: FetchOptions = {}
  ): Promise<ApiResponse<T>> {
    const url = `${apiUrl}${endpoint}`;
    return this.fetchWithErrorHandling<T>(url, {
      method: 'PUT',
      body: JSON.stringify(data),
      ...options,
    });
  }

  static async delete<T>(
    endpoint: string,
    options: FetchOptions = {}
  ): Promise<ApiResponse<T>> {
    const url = `${apiUrl}${endpoint}`;
    return this.fetchWithErrorHandling<T>(url, {
      method: 'DELETE',
      ...options,
    });
  }

  // Convenience method for authenticated requests
  static async authenticatedRequest<T>(
    method: 'GET' | 'POST' | 'PUT' | 'DELETE',
    endpoint: string,
    data?: unknown,
    options: FetchOptions = {}
  ): Promise<ApiResponse<T>> {
    const requestOptions = { ...options, requireAuth: true };

    switch (method) {
      case 'GET':
        return this.get<T>(endpoint, requestOptions);
      case 'POST':
        return this.post<T>(endpoint, data, requestOptions);
      case 'PUT':
        return this.put<T>(endpoint, data, requestOptions);
      case 'DELETE':
        return this.delete<T>(endpoint, requestOptions);
      default:
        throw new Error(`Unsupported method: ${method}`);
    }
  }
}

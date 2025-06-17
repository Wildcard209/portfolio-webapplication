type FetchOptions = {
  cache?: 'force-cache' | 'no-store';
  next?: {
    revalidate?: number;
    tags?: string[];
  };
};

type ApiResponse<T> = {
  data: T | null;
  error: string | null;
  status: number;
};

const apiUrl = process.env.NEXT_PUBLIC_BASE_API_URL ?? '';
const allowedOrigin = process.env.NEXT_PUBLIC_ALLOWED_ORIGIN ?? '';

export class ApiHandler {
  private static getDefaultHeaders(): HeadersInit {
    let authHeaders = {};
    if (typeof window !== 'undefined') {
      const token = localStorage.getItem('auth_token');
      if (token) {
        authHeaders = { Authorization: `Bearer ${token}` };
      }
    }

    return {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
      'Access-Control-Allow-Origin': allowedOrigin,
      'Access-Control-Allow-Methods': 'GET, POST, PUT, DELETE, OPTIONS',
      'Access-Control-Allow-Headers': 'Content-Type, Authorization',
      ...authHeaders,
    };
  }

  private static async fetchWithErrorHandling<T>(
    url: string,
    options: RequestInit & FetchOptions = {}
  ): Promise<ApiResponse<T>> {
    try {
      const response = await fetch(url, {
        ...options,
        credentials: 'include',
        headers: {
          ...this.getDefaultHeaders(),
          ...options.headers,
        },
      });

      if (!response.ok) {
        return {
          data: null,
          error: `HTTP error! status: ${response.status}`,
          status: response.status,
        };
      }

      const data = await response.json();
      return {
        data,
        error: null,
        status: response.status,
      };
    } catch (error) {
      return {
        data: null,
        error: error instanceof Error ? error.message : 'Unknown error occurred',
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
    data: unknown,
    options: FetchOptions = {}
  ): Promise<ApiResponse<T>> {
    const url = `${apiUrl}${endpoint}`;
    return this.fetchWithErrorHandling<T>(url, {
      method: 'POST',
      body: JSON.stringify(data),
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
  
  static async uploadFile<T>(
    endpoint: string,
    formData: FormData,
    options: RequestInit & FetchOptions = {}
  ): Promise<ApiResponse<T>> {
    const url = `${apiUrl}${endpoint}`;
    console.log('Upload URL:', url);
    console.log('FormData entries:', Array.from(formData.entries()));
    
    let authHeaders = {};
    if (typeof window !== 'undefined') {
      const token = localStorage.getItem('auth_token');
      if (token) {
        authHeaders = { Authorization: `Bearer ${token}` };
        console.log('Using auth token:', token.substring(0, 20) + '...');
      } else {
        console.log('No auth token found');
      }
    }

    try {
      const response = await fetch(url, {
        method: 'POST',
        body: formData,
        credentials: 'include',
        headers: {
          'Accept': 'application/json',
          'Access-Control-Allow-Origin': allowedOrigin,
          'Access-Control-Allow-Methods': 'GET, POST, PUT, DELETE, OPTIONS',
          'Access-Control-Allow-Headers': 'Content-Type, Authorization',
          ...authHeaders,
          ...options.headers,
        },
        ...options,
      });

      if (!response.ok) {
        const errorText = await response.text();
        console.error('Upload failed:', response.status, errorText);
        return {
          data: null,
          error: `HTTP error! status: ${response.status} - ${errorText}`,
          status: response.status,
        };
      }

      const data = await response.json();
      return {
        data,
        error: null,
        status: response.status,
      };
    } catch (error) {
      console.error('Upload error:', error);
      return {
        data: null,
        error: error instanceof Error ? error.message : 'Unknown error occurred',
        status: 500,
      };
    }
  }

  static getAssetUrl(endpoint: string): string {
    return `${apiUrl}${endpoint}?t=${Date.now()}`;
  }
}

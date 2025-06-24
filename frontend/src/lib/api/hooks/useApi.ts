import { useCallback, useEffect, useState } from 'react';
import { ApiHandler } from '../apiHandler';

export type UseApiOptions = {
  revalidateTime?: number;
  cacheTag?: string;
  lazy?: boolean;
};

function getApiEndpoint(endpoint: string): string {
  if (endpoint.startsWith('/api/')) {
    return endpoint;
  }

  if (endpoint.startsWith('/')) {
    return `/api${endpoint}`;
  }

  return `/api/${endpoint}`;
}

export function useApi<T>(endpoint: string, options: UseApiOptions = {}) {
  const [data, setData] = useState<T | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(!options.lazy);

  const fetchData = useCallback(async () => {
    setIsLoading(true);
    try {
      const response = await ApiHandler.get<T>(getApiEndpoint(endpoint), {
        next: {
          revalidate: options.revalidateTime,
          tags: options.cacheTag ? [options.cacheTag] : undefined,
        },
      });

      if (response.error) {
        setError(response.error);
        setData(null);
      } else {
        setData(response.data);
        setError(null);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred');
      setData(null);
    } finally {
      setIsLoading(false);
    }
  }, [endpoint, options.revalidateTime, options.cacheTag]);

  useEffect(() => {
    if (!options.lazy) {
      fetchData();
    }
  }, [fetchData, options.lazy]);

  return {
    data,
    error,
    isLoading,
    refetch: fetchData,
  };
}

export function useApiMutation<T, TData = unknown>(
  endpoint: string,
  method: 'POST' | 'PUT' | 'DELETE' = 'POST',
  options: UseApiOptions = {}
) {
  const [data, setData] = useState<T | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const mutate = useCallback(
    async (payload?: TData) => {
      setIsLoading(true);
      try {
        let response;

        switch (method) {
          case 'POST':
            response = await ApiHandler.post<T>(getApiEndpoint(endpoint), payload, {
              next: {
                tags: options.cacheTag ? [options.cacheTag] : undefined,
              },
            });
            break;
          case 'PUT':
            response = await ApiHandler.put<T>(getApiEndpoint(endpoint), payload, {
              next: {
                tags: options.cacheTag ? [options.cacheTag] : undefined,
              },
            });
            break;
          case 'DELETE':
            response = await ApiHandler.delete<T>(getApiEndpoint(endpoint), {
              next: {
                tags: options.cacheTag ? [options.cacheTag] : undefined,
              },
            });
            break;
        }

        if (response.error) {
          setError(response.error);
          setData(null);
          return null;
        }

        setData(response.data);
        setError(null);
        return response.data;
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : 'An error occurred';
        setError(errorMessage);
        setData(null);
        return null;
      } finally {
        setIsLoading(false);
      }
    },
    [endpoint, method, options.cacheTag]
  );

  return {
    data,
    error,
    isLoading,
    mutate,
  };
}

export function useApiFileUpload<T>(endpoint: string, options: UseApiOptions = {}) {
  const [data, setData] = useState<T | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const uploadFile = useCallback(
    async (formData: FormData) => {
      setIsLoading(true);
      try {
        const response = await ApiHandler.uploadFile<T>(getApiEndpoint(endpoint), formData, {
          next: {
            tags: options.cacheTag ? [options.cacheTag] : undefined,
          },
        });

        if (response.error) {
          setError(response.error);
          setData(null);
          return null;
        }

        setData(response.data);
        setError(null);
        return response.data;
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : 'An error occurred';
        setError(errorMessage);
        setData(null);
        return null;
      } finally {
        setIsLoading(false);
      }
    },
    [endpoint, options.cacheTag]
  );

  return {
    data,
    error,
    isLoading,
    uploadFile,
  };
}

function getAdminEndpoint(endpoint: string): string {
  return `/api/admin${endpoint.startsWith('/') ? endpoint : `/${endpoint}`}`;
}

export function useAdminApi<T>(endpoint: string, options: UseApiOptions = {}) {
  const adminEndpoint = getAdminEndpoint(endpoint);
  return useApi<T>(adminEndpoint, options);
}

export function useAdminApiMutation<T, TData = unknown>(
  endpoint: string,
  method: 'POST' | 'PUT' | 'DELETE' = 'POST',
  options: UseApiOptions = {}
) {
  const adminEndpoint = getAdminEndpoint(endpoint);
  return useApiMutation<T, TData>(adminEndpoint, method, options);
}

export function useAdminApiFileUpload<T>(endpoint: string, options: UseApiOptions = {}) {
  const adminEndpoint = getAdminEndpoint(endpoint);
  return useApiFileUpload<T>(adminEndpoint, options);
}

export function getApiAssetUrl(endpoint: string): string {
  return ApiHandler.getAssetUrl(getApiEndpoint(endpoint));
}

// Auth-specific hooks
export function useLogin() {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  const login = useCallback(async (username: string, password: string) => {
    setIsLoading(true);
    setError(null);
    
    try {
      const response = await ApiHandler.post<{
        token: string;
        expiresAt: string;
        user: {
          id: number;
          username: string;
          lastLogin?: string;
        };
      }>(getAdminEndpoint('/login'), { username, password });
      
      setIsLoading(false);
      
      if (response.error) {
        setError(response.error);
        return { success: false, error: response.error };
      }
      
      // Store user data in localStorage
      if (response.data?.user && typeof window !== 'undefined') {
        localStorage.setItem('auth_user', JSON.stringify(response.data.user));
      }
      
      return { success: true, user: response.data?.user || null };
    } catch (err) {
      setIsLoading(false);
      const errorMessage = err instanceof Error ? err.message : 'Network error';
      setError(errorMessage);
      return { success: false, error: errorMessage };
    }
  }, []);
  
  return { login, isLoading, error };
}

export function useLogout() {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  const logout = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    
    try {
      await ApiHandler.post(getAdminEndpoint('/logout'), null);
      
      // Always clear local storage, even if the API call fails
      if (typeof window !== 'undefined') {
        localStorage.removeItem('auth_user');
      }
      
      setIsLoading(false);
      return { success: true };
    } catch (err) {
      setIsLoading(false);
      const errorMessage = err instanceof Error ? err.message : 'Network error';
      setError(errorMessage);
      
      // Still clear local storage on error
      if (typeof window !== 'undefined') {
        localStorage.removeItem('auth_user');
      }
      
      return { success: true }; // Return success:true because we cleared local data
    }
  }, []);
  
  return { logout, isLoading, error };
}

export function useCheckAuth() {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  const checkAuth = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    
    try {
      const response = await ApiHandler.post<{ success: boolean }>(
        getAdminEndpoint('/refresh'), 
        null
      );
      
      setIsLoading(false);
      
      if (response.error || !response.data?.success) {
        // Clear auth data if refresh fails
        if (typeof window !== 'undefined') {
          localStorage.removeItem('auth_user');
        }
        setError(response.error || 'Authentication failed');
        return false;
      }
      
      return true;
    } catch (err) {
      setIsLoading(false);
      const errorMessage = err instanceof Error ? err.message : 'Network error';
      setError(errorMessage);
      
      // Clear auth data on error
      if (typeof window !== 'undefined') {
        localStorage.removeItem('auth_user');
      }
      
      return false;
    }
  }, []);
  
  return { checkAuth, isLoading, error };
}

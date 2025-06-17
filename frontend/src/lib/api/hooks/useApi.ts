import { useCallback, useEffect, useState } from 'react';
import { ApiHandler } from '../apiHandler';

export type UseApiOptions = {
  revalidateTime?: number;
  cacheTag?: string;
  lazy?: boolean;
};

export function useApi<T>(
  endpoint: string,
  options: UseApiOptions = {}
) {
  const [data, setData] = useState<T | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(!options.lazy);

  const fetchData = useCallback(async () => {
    setIsLoading(true);
    try {
      const response = await ApiHandler.get<T>(endpoint, {
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
            response = await ApiHandler.post<T>(endpoint, payload, {
              next: {
                tags: options.cacheTag ? [options.cacheTag] : undefined,
              },
            });
            break;
          case 'PUT':
            response = await ApiHandler.put<T>(endpoint, payload, {
              next: {
                tags: options.cacheTag ? [options.cacheTag] : undefined,
              },
            });
            break;
          case 'DELETE':
            response = await ApiHandler.delete<T>(endpoint, {
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

export function useApiFileUpload<T>(
  endpoint: string,
  options: UseApiOptions = {}
) {
  const [data, setData] = useState<T | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const uploadFile = useCallback(
    async (formData: FormData) => {
      setIsLoading(true);
      try {
        const response = await ApiHandler.uploadFile<T>(endpoint, formData, {
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
  const adminToken = typeof window !== 'undefined' ? localStorage.getItem('admin_token') : null;
  if (!adminToken) {
    return '';
  }
  return `/${adminToken}/admin${endpoint}`;
}

export function useAdminApi<T>(
  endpoint: string,
  options: UseApiOptions = {}
) {
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

export function useAdminApiFileUpload<T>(
  endpoint: string,
  options: UseApiOptions = {}
) {
  const adminEndpoint = getAdminEndpoint(endpoint);
  return useApiFileUpload<T>(adminEndpoint, options);
}

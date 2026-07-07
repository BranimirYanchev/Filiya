import React, { createContext, useState, useContext, useEffect } from 'react';
import {
  api,
  setAuthToken,
  clearAuthToken,
  getStoredAccessToken,
  getStoredRefreshToken,
  storeAuthTokens,
  clearStoredAuthTokens,
} from '../utils/api';

const AuthContext = createContext(null);

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const extractToken = (payload) => (
    payload?.token
    || payload?.access_token
    || payload?.data?.token
    || null
  );

  const extractUser = (payload) => (
    payload?.user
    || payload?.data
    || null
  );

  const extractRefreshToken = (payload) => (
    payload?.refresh_token
    || payload?.data?.refresh_token
    || null
  );

  const syncTokensFromPayload = (payload) => {
    const token = extractToken(payload);
    const refreshToken = extractRefreshToken(payload);

    if (token) {
      setAuthToken(token);
    }

    if (token || refreshToken) {
      storeAuthTokens({ token, refreshToken });
    }

    return { token, refreshToken };
  };

  const fetchAuthenticatedUser = async () => {
    const userResponse = await api.get('/users/me');
    const currentUser = extractUser(userResponse);
    if (!currentUser) {
      throw new Error('Missing user data for authenticated session');
    }

    setUser(currentUser);
    return currentUser;
  };

  const restoreSession = async () => {
    const storedRefreshToken = getStoredRefreshToken();
    const refreshResponse = await api.post(
      '/auth/refresh',
      storedRefreshToken ? { refresh_token: storedRefreshToken } : {},
    );
    const { token } = syncTokensFromPayload(refreshResponse);
    if (!token) {
      throw new Error('Missing access token after refresh');
    }

    return fetchAuthenticatedUser();
  };

  const handleAuthFailure = () => {
    setUser(null);
    clearAuthToken();
    clearStoredAuthTokens();
  };

  // Check if user is authenticated on mount and when needed
  const checkAuth = async () => {
    try {
      setLoading(true);
      setError(null);
      const storedAccessToken = getStoredAccessToken();

      if (storedAccessToken) {
        setAuthToken(storedAccessToken);
        try {
          await fetchAuthenticatedUser();
          return;
        } catch (err) {
          clearAuthToken();
        }
      }

      await restoreSession();
    } catch (err) {
      console.error('Auth check failed:', err);
      handleAuthFailure();
      // Don't set error for normal unauthenticated state
    } finally {
      setLoading(false);
    }
  };

  // Login function
  const login = async (email, password) => {
    try {
      setLoading(true);
      setError(null);
      const response = await api.post('/auth/login', {
        email,
        password,
      });
      const { token } = syncTokensFromPayload(response);
      if (token) {
        setAuthToken(token);
      }

      await fetchAuthenticatedUser();
      
      return { success: true, data: response.data };
    } catch (err) {
      const errorMessage = typeof err === 'string' ? err : (err.message || 'Login failed');
      setError(errorMessage);
      handleAuthFailure();
      throw new Error(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  // Register function
  const register = async (userData) => {
    try {
      setLoading(true);
      setError(null);
      const response = await api.post('/auth/register', {
        full_name: userData.fullName,
        email: userData.email,
        password: userData.password,
        repeated_password: userData.repeatPassword,
      });
      const { token } = syncTokensFromPayload(response);
      if (token) {
        setAuthToken(token);
      }

      await fetchAuthenticatedUser();
      
      return { success: true, data: response.data };
    } catch (err) {
      const errorMessage = typeof err === 'string' ? err : (err.message || 'Registration failed');
      setError(errorMessage);
      handleAuthFailure();
      throw new Error(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  // Logout function
  const logout = async () => {
    try {
      setLoading(true);
      await api.post('/auth/logout');
      setUser(null);
      setError(null);
      clearAuthToken();
    } catch (err) {
      console.error('Logout failed:', err);
      // Clear user state even if logout request fails
      handleAuthFailure();
    } finally {
      setLoading(false);
    }
  };

  const completeGoogleAuth = async () => {
    try {
      setLoading(true);
      setError(null);
      return await restoreSession();
    } catch (err) {
      const errorMessage = typeof err === 'string' ? err : (err.message || 'Google authentication failed');
      setError(errorMessage);
      handleAuthFailure();
      throw new Error(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  // Refresh user data
  const refreshUser = async () => {
    try {
      await restoreSession();
    } catch (err) {
      console.error('Failed to refresh user data:', err);
      handleAuthFailure();
    }
  };

  // Check authentication on mount
  useEffect(() => {
    checkAuth();
  }, []);

  const value = {
    user,
    loading,
    error,
    isAuthenticated: !!user,
    login,
    register,
    logout,
    checkAuth,
    completeGoogleAuth,
    refreshUser,
    setUser,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

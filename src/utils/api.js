const API_BASE_URL = (
  process.env.REACT_APP_API_BASE_URL
  || process.env.REACT_APP_BACKEND_API
  || 'https://filiya-backend.onrender.com/api'
).replace(/\/+$/, '');
const API_ORIGIN = new URL(API_BASE_URL).origin;

export const GOOGLE_LOGIN_URL = `${API_BASE_URL}/auth/google/login`;
const ACCESS_TOKEN_STORAGE_KEY = 'filia_access_token';
const REFRESH_TOKEN_STORAGE_KEY = 'filia_refresh_token';

const getStorage = () => {
  if (typeof window === 'undefined') {
    return null;
  }

  return window.sessionStorage;
};

// Token storage (in memory for security)
let authToken = getStorage()?.getItem(ACCESS_TOKEN_STORAGE_KEY) || null;

export const setAuthToken = (token) => {
  authToken = token;
};

export const clearAuthToken = () => {
  authToken = null;
};

export const getStoredAccessToken = () => getStorage()?.getItem(ACCESS_TOKEN_STORAGE_KEY) || null;

export const getStoredRefreshToken = () => getStorage()?.getItem(REFRESH_TOKEN_STORAGE_KEY) || null;

export const storeAuthTokens = ({ token, refreshToken }) => {
  const storage = getStorage();
  if (!storage) {
    return;
  }

  if (token) {
    storage.setItem(ACCESS_TOKEN_STORAGE_KEY, token);
  }

  if (refreshToken) {
    storage.setItem(REFRESH_TOKEN_STORAGE_KEY, refreshToken);
  }
};

export const clearStoredAuthTokens = () => {
  const storage = getStorage();
  if (!storage) {
    return;
  }

  storage.removeItem(ACCESS_TOKEN_STORAGE_KEY);
  storage.removeItem(REFRESH_TOKEN_STORAGE_KEY);
};

const createApiError = (message, status, payload) => {
  const error = new Error(message);
  error.status = status;
  error.payload = payload;
  return error;
};

const isFormDataBody = (body) => (
  typeof FormData !== 'undefined' && body instanceof FormData
);

export const resolveApiUrl = (path) => {
  if (!path) {
    return '';
  }

  if (/^https?:\/\//i.test(path)) {
    return path;
  }

  if (path.startsWith('/')) {
    return `${API_ORIGIN}${path}`;
  }

  return `${API_BASE_URL}/${path.replace(/^\/+/, '')}`;
};

export const apiRequest = async (endpoint, options = {}) => {
  const defaultHeaders = {};
  const isFormData = isFormDataBody(options.body);

  if (!isFormData) {
    defaultHeaders['Content-Type'] = 'application/json';
  }

  // Add Authorization header if token is available
  if (authToken) {
    defaultHeaders['Authorization'] = `Bearer ${authToken}`;
  }

  const defaultOptions = {
    credentials: 'include',
    headers: {
      ...defaultHeaders,
      ...options.headers,
    },
  };

  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...defaultOptions,
    ...options,
    headers: {
      ...defaultOptions.headers,
      ...options.headers,
    },
  });

  const contentType = response.headers.get('content-type') || '';
  const isJsonResponse = contentType.includes('application/json');
  const data = isJsonResponse
    ? await response.json()
    : await response.text();

  // Handle errors in response (some endpoints return 200 with error field)
  if (isJsonResponse && data?.error && data.error !== "") {
    throw createApiError(data.error, response.status, data);
  }

  // Also throw on HTTP error status
  if (!response.ok) {
    const errorMessage = (
      (isJsonResponse && data?.error)
      || (typeof data === 'string' && data)
      || `HTTP ${response.status}: ${response.statusText}`
    );
    throw createApiError(errorMessage, response.status, data);
  }

  return data;
};

// Convenience methods
export const api = {
  get: (endpoint, options = {}) => 
    apiRequest(endpoint, { ...options, method: 'GET' }),
  
  post: (endpoint, body, options = {}) => 
    apiRequest(endpoint, { 
      ...options, 
      method: 'POST', 
      body: JSON.stringify(body) 
    }),

  postForm: (endpoint, body, options = {}) =>
    apiRequest(endpoint, {
      ...options,
      method: 'POST',
      body,
    }),
  
  put: (endpoint, body, options = {}) => 
    apiRequest(endpoint, { 
      ...options, 
      method: 'PUT', 
      body: JSON.stringify(body) 
    }),
  
  delete: (endpoint, options = {}) => 
    apiRequest(endpoint, { ...options, method: 'DELETE' }),
};

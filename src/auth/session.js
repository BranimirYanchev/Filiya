const API_BASE_URL = 'https://filiya-backend.onrender.com/api';

const AUTH_SESSION_HINT_KEY = 'filia-auth-session';
const AUTH_NOTICE_KEY = 'filia-auth-notice';
const AUTH_CHANGED_EVENT = 'filia:auth-changed';
const SESSION_EXPIRED_EVENT = 'filia:session-expired';

let refreshInFlight = null;

const createRequestOptions = ({ token, headers = {}, ...options } = {}) => ({
  ...options,
  mode: 'cors',
  credentials: 'include',
  headers: {
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
    ...headers
  }
});

const hasAuthPayload = (payload) => (
  !payload
  || typeof payload !== 'object'
  || Boolean(payload.token || payload.access_token || payload.user || payload.authenticated)
);

const emitAuthChanged = () => {
  window.dispatchEvent(new CustomEvent(AUTH_CHANGED_EVENT, {
    detail: { authenticated: isAuthenticated() }
  }));
};

const emitSessionExpired = (message = 'Сесията ти изтече. Моля, влез отново.') => {
  window.sessionStorage.setItem(AUTH_NOTICE_KEY, message);
  window.dispatchEvent(new CustomEvent(SESSION_EXPIRED_EVENT, {
    detail: { message }
  }));
};

const parseJsonSafe = async (response) => {
  try {
    return await response.json();
  } catch (error) {
    return null;
  }
};

const isUnauthorizedPayload = (payload) => (
  Boolean(payload)
  && typeof payload === 'object'
  && (payload.authenticated === false || payload.authorized === false)
);

const HUMAN_ERROR_MAP = [
  { pattern: /invalid full name/i, message: 'Моля, въведи валидно име и фамилия.' },
  { pattern: /invalid name length/i, message: 'Името трябва да е с подходяща дължина.' },
  { pattern: /invalid name characters/i, message: 'Името съдържа невалидни символи.' },
  { pattern: /invalid email/i, message: 'Имейлът не е валиден.' },
  { pattern: /email.*already/i, message: 'Този имейл вече е регистриран.' },
  { pattern: /password.*match|passwords.*match/i, message: 'Паролите не съвпадат.' },
  { pattern: /invalid password length|password.*too short/i, message: 'Паролата е твърде кратка.' },
  { pattern: /invalid credentials|wrong password|unauthorized/i, message: 'Невалиден имейл или парола.' },
  { pattern: /user.*not found|unexisting user/i, message: 'Няма потребител с такъв имейл.' },
  { pattern: /invalid confirmation password|invalid password/i, message: 'Текущата парола е невалидна.' },
  { pattern: /refresh token required/i, message: 'Сесията е изтекла. Влез отново.' },
  { pattern: /authenticated false|unauthenticated/i, message: 'Трябва да влезеш в профила си.' }
];

const translateErrorMessage = (rawMessage) => {
  if (!rawMessage || typeof rawMessage !== 'string') {
    return null;
  }

  const normalized = rawMessage
    .split('|')
    .map((part) => part.trim())
    .filter(Boolean);

  if (!normalized.length) {
    return null;
  }

  const translated = normalized.map((part) => {
    const match = HUMAN_ERROR_MAP.find((entry) => entry.pattern.test(part));
    return match ? match.message : part;
  });

  return [...new Set(translated)].join(' ');
};

const extractErrorMessage = (payload, fallbackMessage) => {
  if (!payload) return fallbackMessage;

  if (isUnauthorizedPayload(payload)) {
    return 'Трябва да влезеш в профила си.';
  }

  if (typeof payload.message === 'string' && payload.message.trim()) {
    return translateErrorMessage(payload.message) || fallbackMessage;
  }

  if (typeof payload.error === 'string' && payload.error.trim()) {
    return translateErrorMessage(payload.error) || fallbackMessage;
  }

  if (typeof payload.data === 'string' && payload.data.trim()) {
    return translateErrorMessage(payload.data) || fallbackMessage;
  }

  return fallbackMessage;
};

const shouldRetryWithRefresh = (response, payload) => (
  Boolean(response)
  && (response.status === 400 || response.status === 401 || isUnauthorizedPayload(payload))
);

const postJson = async (path, body, token) => {
  try {
    const response = await fetch(`${API_BASE_URL}${path}`, createRequestOptions({
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(body),
      token
    }));

    const payload = await parseJsonSafe(response);

    if (!response.ok) {
      throw new Error(extractErrorMessage(payload, 'Заявката не беше успешна.'));
    }

    return payload;
  } catch (error) {
    if (error instanceof TypeError) {
      throw new Error('Няма връзка с бекенда. Провери дали API-то работи на https://filiya-backend.onrender.com/.');
    }

    throw error;
  }
};

export const storeSession = () => {
  window.localStorage.setItem(AUTH_SESSION_HINT_KEY, '1');
  emitAuthChanged();
};

export const clearSession = () => {
  window.localStorage.removeItem(AUTH_SESSION_HINT_KEY);
  emitAuthChanged();
};

export const consumeAuthNotice = () => {
  const message = window.sessionStorage.getItem(AUTH_NOTICE_KEY);
  if (message) {
    window.sessionStorage.removeItem(AUTH_NOTICE_KEY);
  }
  return message;
};

export const isAuthenticated = () => window.localStorage.getItem(AUTH_SESSION_HINT_KEY) === '1';

const invalidateSession = (message = 'Сесията ти изтече. Моля, влез отново.') => {
  clearSession();
  emitSessionExpired(message);
};

export const refreshSession = async () => {
  if (refreshInFlight) {
    return refreshInFlight;
  }

  if (!isAuthenticated()) {
    invalidateSession('Трябва да влезеш в профила си.');
    throw new Error('Трябва да влезеш в профила си.');
  }

  refreshInFlight = (async () => {
    try {
      const refreshed = await postJson('/auth/refresh', {});

      if (!hasAuthPayload(refreshed)) {
        invalidateSession('Сесията ти изтече. Моля, влез отново.');
        throw new Error('Сесията е изтекла. Влез отново.');
      }

      storeSession();
      return refreshed;
    } catch (error) {
      invalidateSession(error.message || 'Сесията ти изтече. Моля, влез отново.');
      throw error;
    } finally {
      refreshInFlight = null;
    }
  })();

  return refreshInFlight;
};

export const ensureValidSession = async ({ force = false } = {}) => {
  if (!isAuthenticated()) {
    return false;
  }

  if (!force) {
    return true;
  }

  await refreshSession();
  return true;
};

const requestJson = async (path, options = {}) => {
  const makeRequest = async (token) => {
    const response = await fetch(`${API_BASE_URL}${path}`, createRequestOptions({
      ...options,
      headers: {
        ...(options.body ? { 'Content-Type': 'application/json' } : {}),
        ...(options.headers || {})
      },
      token
    }));

    const payload = await parseJsonSafe(response);
    return { response, payload };
  };

  try {
    let { response, payload } = await makeRequest();

    if (!response.ok && shouldRetryWithRefresh(response, payload) && isAuthenticated()) {
      await refreshSession();
      ({ response, payload } = await makeRequest());
    }

    if (!response.ok) {
      if (isUnauthorizedPayload(payload) || response.status === 401) {
        invalidateSession('Сесията ти изтече. Моля, влез отново.');
      }
      throw new Error(extractErrorMessage(payload, 'Заявката не беше успешна.'));
    }

    return payload;
  } catch (error) {
    if (error instanceof TypeError) {
      throw new Error('Няма връзка с бекенда. Провери дали API-то работи на https://filiya-backend.onrender.com/.');
    }

    throw error;
  }
};

const requestFormData = async (path, formData, options = {}) => {
  const makeRequest = async (token) => {
    const response = await fetch(`${API_BASE_URL}${path}`, createRequestOptions({
      method: options.method || 'POST',
      ...options,
      headers: {
        ...(options.headers || {})
      },
      body: formData,
      token
    }));

    const payload = await parseJsonSafe(response);
    return { response, payload };
  };

  try {
    let { response, payload } = await makeRequest();

    if (!response.ok && shouldRetryWithRefresh(response, payload) && isAuthenticated()) {
      await refreshSession();
      ({ response, payload } = await makeRequest());
    }

    if (!response.ok) {
      if (isUnauthorizedPayload(payload) || response.status === 401) {
        invalidateSession('Сесията ти изтече. Моля, влез отново.');
      }
      throw new Error(extractErrorMessage(payload, 'Заявката не беше успешна.'));
    }

    return payload;
  } catch (error) {
    if (error instanceof TypeError) {
      throw new Error('Няма връзка с бекенда. Провери дали API-то работи на https://filiya-backend.onrender.com/.');
    }

    throw error;
  }
};

export const loginUser = async ({ email, password }) => {
  const loginPayload = await postJson('/auth/login', { email, password });

  storeSession();
  return loginPayload;
};

export const registerUser = async ({ email, fullName, password, repeatedPassword }) => {
  const registerPayload = await postJson('/auth/register', {
    email,
    full_name: fullName,
    password,
    repeated_password: repeatedPassword
  });

  storeSession();
  return registerPayload;
};

export const logoutUser = async () => {
  try {
    await postJson('/auth/logout', {});
  } catch (error) {
    // Ignore logout transport errors; local session still needs to be cleared.
  } finally {
    clearSession();
  }
};

export const fetchCurrentUser = async () => {
  const payload = await requestJson('/users/me');
  return payload?.user || null;
};

export const fetchPosts = async ({ limit = 20, offset = 0 } = {}) => {
  const payload = await requestJson(`/posts?limit=${limit}&offset=${offset}`);
  return Array.isArray(payload?.data) ? payload.data : [];
};

export const fetchCategories = async ({ limit = 100, offset = 0 } = {}) => {
  const payload = await requestJson(`/categories?limit=${limit}&offset=${offset}`);
  return Array.isArray(payload?.data) ? payload.data : [];
};

export const fetchUserFriends = async () => requestJson('/users/friends');

export const fetchCurrentUserPosts = async (userId) => {
  if (!userId) return [];
  return requestJson(`/users/${userId}/posts`);
};

export const updateProfileBasic = async ({ bio }) => requestJson('/users/profile/basic', {
  method: 'PUT',
  body: JSON.stringify({
    bio: {
      valid: false,
      value: bio
    }
  })
});

export const updateProfileSensitive = async ({ fullName, email, oldPassword }) => {
  const payload = { old_password: oldPassword };

  if (fullName) {
    payload.full_name = fullName;
  }

  if (email) {
    payload.email = email;
  }

  return requestJson('/users/profile/sensitive', {
    method: 'PUT',
    body: JSON.stringify(payload)
  });
};

export const createPost = async ({
  title,
  content,
  categoryIds,
  attachments = [],
  isPrivate = false
}) => {
  const formData = new FormData();
  formData.append('title', title);
  formData.append('content', content);
  formData.append('category_ids', JSON.stringify(categoryIds));
  formData.append('is_private', String(isPrivate));

  attachments.forEach((file) => {
    formData.append('attachments', file);
  });

  const payload = await requestFormData('/posts', formData);
  return payload?.data || null;
};

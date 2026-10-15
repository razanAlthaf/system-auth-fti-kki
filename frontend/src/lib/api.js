const API_URL = import.meta.env.VITE_API_URL ?? 'http://192.168.0.102:8080/api';

export class ApiError extends Error {
  constructor(status) {
    super(`Request failed with status ${status}`);
    this.status = status;
  }
}

async function request(path, { method = 'GET', body, token } = {}) {
  const headers = {};
  if (body) headers['Content-Type'] = 'application/json';
  if (token) headers.Authorization = `Bearer ${token}`;

  const response = await fetch(`${API_URL}${path}`, {
    method,
    headers,
    body: body ? JSON.stringify(body) : undefined
  });

  if (!response.ok) throw new ApiError(response.status);

  const text = await response.text();
  return text ? JSON.parse(text) : null;
}

export const api = {
  login: (credentials) => request('/login', { method: 'POST', body: credentials }),
  register: (payload) => request('/register', { method: 'POST', body: payload }),
  logout: (token) => request('/logout', { method: 'POST', token }),
  getProfile: (token) => request('/profile', { token }),
  getUsers: (token) => request('/users', { token }),
  forceLogout: (token, userId) =>
    request(`/force-logout?user_id=${encodeURIComponent(userId)}`, { method: 'POST', token })
};
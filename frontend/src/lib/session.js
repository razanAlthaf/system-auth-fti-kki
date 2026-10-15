const TOKEN_KEY = 'token';
const ROLE_KEY = 'role';

export function loadSession() {
  return {
    token: localStorage.getItem(TOKEN_KEY) ?? '',
    role: localStorage.getItem(ROLE_KEY) ?? ''
  };
}

export function saveSession({ token, role }) {
  localStorage.setItem(TOKEN_KEY, token);
  localStorage.setItem(ROLE_KEY, role);
}

export function clearSession() {
  localStorage.removeItem(TOKEN_KEY);
  localStorage.removeItem(ROLE_KEY);
}
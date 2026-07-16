import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface AuthState {
  token: string | null;
  username: string | null;
  setToken: (token: string, username: string) => void;
  clear: () => void;
}

/**
 * Persist admin token in localStorage. This is a known XSS risk trade-off:
 * the API still uses Bearer tokens. In a future iteration we plan to move
 * to HttpOnly cookies on the backend and drop this persistence.
 */
export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      token: null,
      username: null,
      setToken: (token, username) => set({ token, username }),
      clear: () => set({ token: null, username: null }),
    }),
    {
      name: 'rt-admin-auth',
      version: 1,
    },
  ),
);

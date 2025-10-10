import { api } from '$lib/api';

export type User = {
  id: number;
  name: string;
  profileImage?: string;
  team_id?: number;
  role: string;
  nationality?: string;
  joined_at?: string;
  // add fields you actually return
};

export async function getInfo(): Promise<User | null> {
  try {
    return await api<User>('/info');
  } catch {
    return null;
  }
}

export async function login(
  email: string,
  password: string
): Promise<User> {
  return api<User>('/login', {
    method: 'POST',
    body: JSON.stringify({ email, password })
  });
}

export async function logout(): Promise<void> {
  await api<void>('/logout', { method: 'POST' });
}

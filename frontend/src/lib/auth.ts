import { api } from '$lib/api';
import type { User } from '$lib/types';

export async function getInfo(): Promise<any | null> {
	try {
		return await api<any>('/info');
	} catch {
		return null;
	}
}

export async function login(email: string, password: string): Promise<any> {
	return api<any>('/login', {
		headers: { 'content-type': 'application/json' },
		method: 'POST',
		body: JSON.stringify({ email, password })
	});
}

export async function register(email: string, password: string, name: string): Promise<User> {
	return api<User>('/register', {
		headers: { 'content-type': 'application/json' },
		method: 'POST',
		body: JSON.stringify({ email, password, name })
	});
}

export async function logout(): Promise<void> {
	await api<any>('/logout', { method: 'POST' });
}

import { api } from '$lib/api';
import type { User, PaginatedResponse } from '$lib/types';

export async function getUsers(page = 1, limit = 20): Promise<PaginatedResponse<User>> {
	const offset = (page - 1) * limit;
	const response = await api<{ total: number; users: User[] }>(`/users?offset=${offset}&limit=${limit}`);
	return {
		success: true,
		data: response.users,
		pagination: {
			total: response.total,
			page,
			per_page: limit,
			pages: Math.ceil(response.total / limit)
		}
	};
}

export async function getUserData(id: number): Promise<User> {
	return api<User>(`/users/${id}`);
}

export async function updateUser(id: number, name: string, country: string): Promise<any> {
	return api<any>(`/users`, {
		headers: { 'content-type': 'application/json' },
		method: 'PATCH',
		body: JSON.stringify({ id, name, country })
	});
}

export async function updateUserRole(userId: number, role: string): Promise<any> {
	return api<any>('/users/role', {
		headers: { 'content-type': 'application/json' },
		method: 'PATCH',
		body: JSON.stringify({ user_id: userId, new_role: role })
	});
}
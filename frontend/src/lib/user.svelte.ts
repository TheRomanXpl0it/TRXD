import { api } from '$lib/api';
import type { User, PaginatedResponse } from '$lib/types';

import { authState } from '$lib/stores/auth';

export async function getUsers(page = 1, limit = 20): Promise<PaginatedResponse<User>> {
	const offset = (page - 1) * limit;
	const isUserMode = authState.userMode;

	if (isUserMode) {
		const response = await api<{ total: number; teams: any[] }>(
			`/teams?offset=${offset}&limit=${limit}`
		);
		// Map teams to users structure
		const teams = response?.teams || [];
		const users = teams.map((t: any) => ({ ...t, role: t.role || 'User' }));
		return {
			success: true,
			data: users,
			pagination: {
				total: response.total,
				page,
				per_page: limit,
				pages: Math.ceil(response.total / limit)
			}
		};
	} else {
		const response = await api<{ total: number; users: User[] }>(
			`/users?offset=${offset}&limit=${limit}`
		);
		return {
			success: true,
			data: response?.users || [],
			pagination: {
				total: response.total,
				page,
				per_page: limit,
				pages: Math.ceil(response.total / limit)
			}
		};
	}
}

export async function getUserData(id: number, userMode = authState.userMode): Promise<User> {
	if (userMode) {
		const team = await api<any>(`/teams/${id}`);
		return { ...team, role: team.role || 'User' };
	} else {
		return api<User>(`/users/${id}`);
	}
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

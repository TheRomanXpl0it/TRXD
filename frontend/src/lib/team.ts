import { api } from '$lib/api';
import type { Team, PaginatedResponse } from '$lib/types';

export async function getTeams(page = 1, limit = 20): Promise<PaginatedResponse<Team>> {
	const offset = (page - 1) * limit;
	const response = await api<{ total: number; teams: Team[] }>(
		`/teams?offset=${offset}&limit=${limit}`
	);
	return {
		success: true,
		data: response.teams,
		pagination: {
			total: response.total,
			page,
			per_page: limit,
			pages: Math.ceil(response.total / limit)
		}
	};
}

export async function getTeam(id: number): Promise<Team> {
	return api<Team>(`/teams/${id}`);
}

export async function joinTeam(name: string, password: string): Promise<any> {
	return api<any>(`/teams/join`, {
		headers: { 'content-type': 'application/json' },
		method: 'POST',
		body: JSON.stringify({ name, password })
	});
}

export async function createTeam(name: string, password: string): Promise<any> {
	return api<any>(`/teams/register`, {
		headers: { 'content-type': 'application/json' },
		method: 'POST',
		body: JSON.stringify({ name, password })
	});
}

export async function updateTeam(
	id: number,
	name: string,
	country: string,
	tags: string[] = []
): Promise<any> {
	return api<any>(`/teams`, {
		headers: { 'content-type': 'application/json' },
		method: 'PATCH',
		body: JSON.stringify({ id, name, country, tags })
	});
}

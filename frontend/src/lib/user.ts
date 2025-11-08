import { api } from '$lib/api';

export async function getUserData(id: number): Promise<any> {
	return api<any>(`/users/${id}`);
}

export async function updateUser(id: number, name: string, country: string, image: string): Promise<any> {
	return api<any>(`/users`, {
		headers: { 'content-type': 'application/json' },
		method: 'PATCH',
		body: JSON.stringify({ id, name, country, image })
	});
}
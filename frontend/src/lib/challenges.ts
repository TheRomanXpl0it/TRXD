import { api } from '$lib/api';
import type { Challenge, Category, Solve } from '$lib/types';

export async function getSolves(chall_id: string | number): Promise<Solve[]> {
	const challenge = await getChallenge(chall_id);
	return Array.isArray(challenge?.solves_list) ? challenge.solves_list : [];
}

export async function getChallenges(): Promise<Challenge[]> {
	return api<Challenge[]>('/challenges');
}

export async function getChallenge(chall_id: string | number): Promise<Challenge> {
	return api<Challenge>(`/challenges/${chall_id}`);
}

export async function submitFlag(
	chall_id: string,
	flag: string
): Promise<{ status: string; first_blood?: boolean }> {
	return api<{ first_blood: boolean; status: string }>(`/submissions`, {
		headers: { 'content-type': 'application/json' },
		method: 'POST',
		body: JSON.stringify({ flag, chall_id })
	});
}

export async function getCategories(): Promise<Category[]> {
	return api<Category[]>('/categories');
}

export async function createCategory(name: string, icon: string): Promise<any> {
	return api<any>('/categories', {
		headers: { 'content-type': 'application/json' },
		method: 'POST',
		body: JSON.stringify({ name, icon })
	});
}

export async function createChallenge(
	name: string,
	category: string,
	description: string,
	type: string,
	max_points: number,
	score_type: string
): Promise<any> {
	return api<any>('/challenges', {
		headers: { 'content-type': 'application/json' },
		method: 'POST',
		body: JSON.stringify({ name, category, description, type, max_points, score_type })
	});
}

export async function deleteChallenge(chall_id: string): Promise<any> {
	return api<any>(`/challenges`, {
		headers: { 'content-type': 'application/json' },
		method: 'DELETE',
		body: JSON.stringify({ chall_id })
	});
}

export async function updateChallenge(data: any): Promise<any> {
	return api<any>(`/challenges`, {
		headers: { 'content-type': 'application/json' },
		method: 'PATCH',
		body: JSON.stringify(data)
	});
}

export async function uploadAttachments(fd: FormData): Promise<any> {
	return api<any>(`/attachments`, {
		method: 'POST',
		body: fd
	});
}

export async function deleteAttachments(chall_id: number, names: string[]): Promise<any> {
	return api<any>(`/attachments`, {
		headers: { 'content-type': 'application/json' },
		method: 'DELETE',
		body: JSON.stringify({ chall_id, names })
	});
}

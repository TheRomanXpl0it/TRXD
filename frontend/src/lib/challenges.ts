import { api } from '$lib/api';

export async function getSolves(chall_id: string): Promise<any[]> {
	const challenge = await getChallenge(chall_id);
	const solves = Array.isArray(challenge?.solves_list) ? challenge.solves_list : [];
	return solves;
}

export async function getChallenges(): Promise<any[]> {
	return api<any[]>('/challenges');
}

export async function getChallenge(chall_id: string): Promise<any> {
	return api<any>(`/challenges/${chall_id}`);
}

export async function submitFlag(chall_id: string, flag: string): Promise<{ status: string }> {
	return api<{ first_blood: boolean; status: string }>(`/submissions`, {
		headers: { 'content-type': 'application/json' },
		method: 'POST',
		body: JSON.stringify({ flag, chall_id })
	});
}

export async function getCategories(): Promise<any[]> {
	return api<any[]>('/categories');
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

/** Sends multipart/form-data to PATCH /challenges (server reads chall_id) */
export async function updateChallengeMultipart(fd: FormData): Promise<any> {
	return api<any>(`/challenges`, {
		method: 'PATCH',
		body: fd
	});
}

import { api } from '$lib/api';

export async function getSolves(
  chall_id: string
): Promise<any[]>{
  return api<any[]>(`/challenges/${chall_id}/solves`);
}

export async function getChallenges(): Promise<any[]> {
  return api<any[]>('/challenges');
}

export async function getChallenge(chall_id: string): Promise<any> {
  return api<any>(`/challenges/${chall_id}`);
}

export async function submitFlag(
  chall_id: string,
  flag: string
): Promise<{status: string }> {
  return await  api<{ first_blood: boolean; status: string }>(`/submissions`, {
    method: 'POST',
    body: JSON.stringify({ flag,chall_id })
  });
}

export async function getCategories(): Promise<any[]> {
  return api<any[]>('/categories');
}

export async function createCategory(name: string,icon:string): Promise<any> {
  return api<any>('/categories', {
    method: 'POST',
    body: JSON.stringify({ name, icon })
  });
}

export async function createChallenge(
    name:        string,
		category:    string,
		description: string,
		type:        string,
		max_points:   number,
		score_type:   string,
): Promise<any> {
  return api<any>('/challenges', {
    method: 'POST',
    body: JSON.stringify({ name, category, description, type, max_points, score_type })
  });
}

export async function updateChallenge(
    id:          number,
    name:        string,
    category:    string,
    description: string,
    type:        string,
    max_points:   number,
    score_type:   string,
): Promise<any> {
  return api<any>(`/challenges/${id}`, {
    method: 'PATCH',
    body: JSON.stringify({ name, category, description, type, max_points, score_type })
  });
}
import { api } from '$lib/api';

export async function getSolves(
  chall_id: string
): Promise<any[]>{
  return api<any[]>(`/challenges/${chall_id}/solves`);
}

export async function getChallenges(): Promise<any[]> {
  return api<any[]>('/challenges');
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
import { api} from '$lib/api';

export async function getTeam(id: number): Promise<any> {
  return api<any>(`/teams/${id}`);
}

export async function joinTeam(name: string, password: string): Promise<any> {
  return api<any>(`/teams/join`, {
    method: 'POST',
    body: JSON.stringify({ name, password })
  });
}

export async function createTeam(name: string, password:string): Promise<any> {
  return api<any>(`/teams/register`, {
    method: 'POST',
    body: JSON.stringify({ name,password })
  });
}
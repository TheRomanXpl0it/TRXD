import { api} from '$lib/api';

export async function getTeam(id: number): Promise<any> {
  return api<any>(`/teams/${id}`);
}

export async function joinTeam(name: string, password: string): Promise<any> {
  return api<any>(`/teams/join`, {
    headers: { "content-type": "application/json" },
    method: 'POST',
    body: JSON.stringify({ name, password })
  });
}

export async function createTeam(name: string, password:string): Promise<any> {
  return api<any>(`/teams/register`, {
    headers: { "content-type": "application/json" },
    method: 'POST',
    body: JSON.stringify({ name,password })
  });
}

export async function updateTeam(id: number, name:string, bio: string, image: string, country: string): Promise<any>{
  return api<any>(`/teams`, {
    headers: { "content-type": "application/json" },
    method: 'PATCH',
    body: JSON.stringify({ id, name, bio, image, country })
  });
}
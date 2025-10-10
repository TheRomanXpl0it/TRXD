import { api} from '$lib/api';

export async function getTeam(id: number): Promise<any> {
  return api<any>(`/teams/${id}`);
}
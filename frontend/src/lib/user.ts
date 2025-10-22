import { api } from '$lib/api';

export async function getUserData(id:number): Promise<any> {
  return api<any>(`/users/${id}`);
}
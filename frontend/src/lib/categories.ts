import { api } from "./api";

export function createCategory(name: string, icon:string): Promise<any>{
  return api<any>(`/categories`,{
		headers: { 'content-type': 'application/json' },
		method: 'POST',
		body: JSON.stringify({name,icon})
	});
}
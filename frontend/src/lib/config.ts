import { api } from '$lib/api';

export async function getConfigs(): Promise<any> {
	return api<any>(`/configs`);
}

export async function updateConfigs(configs: any): Promise<any> {
	return api<any>(`/configs`, {
		headers: { 'content-type': 'application/json' },
		method: 'PATCH',
		body: JSON.stringify(configs)
	});
}

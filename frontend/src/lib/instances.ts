import { api } from './api';

export async function startInstance(
	chall_id: string
): Promise<{ host: string; port: number; timeout: number }> {
	return api<{ host: string; port: number; timeout: number }>(`/instances`, {
		headers: { 'content-type': 'application/json' },
		method: 'POST',
		body: JSON.stringify({ chall_id })
	});
}

export async function stopInstance(chall_id: string): Promise<void> {
	return api<void>(`/instances`, {
		headers: { 'content-type': 'application/json' },
		method: 'DELETE',
		body: JSON.stringify({ chall_id })
	});
}

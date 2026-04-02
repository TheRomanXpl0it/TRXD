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

export async function getInstances(): Promise<any[]> {
	return api<any[]>('/instances', {
		method: 'GET'
	});
}

export async function adminStopInstance(teamId: number, challId: number): Promise<void> {
	return api<void>('/instances', {
		headers: { 'content-type': 'application/json' },
		method: 'DELETE',
		body: JSON.stringify({ team_id: teamId, chall_id: challId })
	});
}

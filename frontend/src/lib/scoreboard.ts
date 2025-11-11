import { api } from '$lib/api';

export async function getScoreboard(): Promise<any[]> {
	return api<any[]>('/scoreboard');
}

export async function getGraphData(): Promise<any[]> {
	return api<any[]>('/scoreboard/graph');
}

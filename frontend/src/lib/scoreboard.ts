import { api } from '$lib/api';

export interface PaginatedResponse<T> {
	success: boolean;
	data: T[];
	pagination?: {
		next?: number;
		prev?: number;
		page?: number;
		per_page?: number;
		total?: number;
		pages?: number;
	};
}

export async function getScoreboard(page = 1, limit = 20): Promise<PaginatedResponse<any>> {
	const offset = (page - 1) * limit;
	const response = await api<{ total: number; teams: any[] }>(`/scoreboard?offset=${offset}&limit=${limit}`);
	return {
		success: true,
		data: response.teams,
		pagination: {
			total: response.total,
			page,
			per_page: limit,
			pages: Math.ceil(response.total / limit)
		}
	};
}

export async function getGraphData(): Promise<any> {
	return api<any>('/scoreboard/graph');
}

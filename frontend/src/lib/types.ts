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

export interface Badge {
	name: string;
	description: string;
}

export interface User {
	id: number;
	name: string;
	email?: string;
	image?: string;
	profileImage?: string;
	team_id?: number | null;
	role: 'User' | 'Admin' | 'Author';
	country?: string;
	joined_at?: string;
	solves?: any[]; // Populated by GetUser / GetTeam
	badges?: Badge[];
	score?: number;
	total_category_challenges?: { category: string; count: number }[];
}

export interface Team {
	id: number;
	name: string;
	email?: string;
	country?: string;
	tags?: string[];
	captain_id?: number;
	members?: User[];
	solves?: any[]; // Populated by GetTeam
	badges?: Badge[];
	score?: number;
	total_category_challenges?: { category: string; count: number }[];
}

export interface Category {
	id?: number;
	name: string;
	color?: string;
}

export interface Challenge {
	id: number;
	name: string;
	description: string;
	category: string;
	type: 'Container' | 'Compose' | 'Normal';
	score_type: 'static' | 'dynamic';
	points: number;
	max_points?: number;
	min_points?: number;
	initial_points?: number;
	decay?: number;

	// Metadata
	tags: string[];
	authors: string[];
	attachments?: string[];
	difficulty?: string;
	colors?: string[];

	// Host and connection info
	host?: string;
	port?: number;
	conn_type?: 'NONE' | 'TCP' | 'TCP_TLS' | 'HTTP' | 'HTTPS';
	connection_info?: string;

	// Instance info
	instance: boolean;
	instance_host?: string | null;
	instance_port?: number | null;
	hash_domain?: boolean;
	timeout?: number | null;

	// Solve state
	solved?: boolean;
	solves?: number;
	solves_list?: Solve[];

	// Admin fields
	hidden?: boolean;
}

export interface Solve {
	id: number;
	name: string;
	timestamp: string;
	// Optional context-specific fields
	user_id?: number;
	category?: string;
	points?: number;
	first_blood?: boolean;
}

export interface ApiError {
	success: boolean;
	errors?: string[];
	message?: string;
}

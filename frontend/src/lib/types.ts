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

export interface User {
    id: number;
    name: string;
    email?: string;
    profileImage?: string;
    team_id?: number | null;
    role: 'User' | 'Admin';
    country?: string;
    joined_at?: string;
    solves?: any[]; // Populated by GetUser / GetTeam
}

export interface Team {
    id: number;
    name: string;
    country?: string;
    tags?: string[];
    captain_id?: number;
    members?: User[];
    solves?: any[]; // Populated by GetTeam
}

export interface Category {
    id?: number;
    name: string;
    icon?: string;
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

    // Connection info
    host?: string;
    port?: number;
    connection_info?: string;

    // Instance info
    instance: boolean;
    instance_host?: string | null;
    instance_port?: number | null;
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

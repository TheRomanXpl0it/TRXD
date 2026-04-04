import { render, screen } from '@testing-library/svelte';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import Scoreboard from '../../../../src/routes/scoreboard/+page.svelte';
import { authState } from '$lib/stores/auth';

// Mock the scoreboard API functions
vi.mock('@/scoreboard', () => ({
	getScoreboard: vi.fn(() => Promise.resolve({ data: [], pagination: { total: 0 } })),
	getGraphData: vi.fn(() => Promise.resolve([]))
}));

// Mock svelte-query
vi.mock('@tanstack/svelte-query', () => ({
	createQuery: vi.fn(() => ({
		data: { data: [], pagination: { total: 0 } },
		isLoading: false,
		error: null
	}))
}));

describe('Scoreboard Page', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		// Reset authState defaults
		authState.ready = true;
		authState.startTime = null;
		authState.user = { id: 1, name: 'Test User', role: 'Player' } as any;
	});

	it('shows Scoreboard even if competition is upcoming', async () => {
		authState.startTime = new Date(Date.now() + 1000000).toISOString();

		render(Scoreboard);

		// Should NOT show WaitingPage, but the Scoreboard header
		expect(screen.getByText('Scoreboard')).toBeInTheDocument();
		expect(screen.queryByText('Starting soon')).not.toBeInTheDocument();
	});

	it('shows Scoreboard when competition has started', async () => {
		authState.startTime = new Date(Date.now() - 1000000).toISOString();

		render(Scoreboard);

		expect(screen.getByText('Scoreboard')).toBeInTheDocument();
	});
});

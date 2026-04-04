import { render, screen } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import Page from '../../../../../src/routes/account/[[id]]/+page.svelte';
import { authState } from '$lib/stores/auth';
import { createQuery } from '@tanstack/svelte-query';
import { page } from '$app/stores';

// Mock auth store
vi.mock('$lib/stores/auth', () => ({
	authState: {
		user: { id: 1, role: 'User' },
		ready: true,
		userMode: true,
		startTime: null
	},
	loadUser: vi.fn()
}));

vi.mock('@tanstack/svelte-query', () => ({
	createQuery: vi.fn(),
	useQueryClient: vi.fn(() => ({
		setQueryData: vi.fn(),
		refetchQueries: vi.fn()
	}))
}));

vi.mock('$app/stores', () => ({
	page: {
		subscribe: (fn: (v: any) => void) => {
			fn({
				params: { id: '1' },
				url: new URL('http://localhost/account/1')
			});
			return () => {};
		}
	}
}));

vi.mock('$lib/user', () => ({
	getUserData: vi.fn()
}));

vi.mock('$lib/team', () => ({
	getTeam: vi.fn()
}));

vi.mock('$lib/components/RadarChart.svelte', () => ({
	default: vi.fn(() => ({ render: () => '' }))
}));

describe('Account Profile Page', () => {
	beforeEach(() => {
		vi.clearAllMocks();

		// Default mock for user query
		(createQuery as any).mockReturnValue({
			data: { id: 1, name: 'Test User', solves: [], badges: [] },
			isLoading: false,
			error: null
		});
	});

	it('hides statistics cards when competition is upcoming', () => {
		// Set authState to upcoming
		authState.ready = true;
		authState.startTime = new Date(Date.now() + 100000).toISOString();

		render(Page);

		// Cards that should be HIDDEN
		expect(screen.queryByText('Category Breakdown')).not.toBeInTheDocument();
		expect(screen.queryByText('Challenges Solved')).not.toBeInTheDocument();

		// Profile header should still be there
		expect(screen.getByText('Profile')).toBeInTheDocument();
	});

	it('shows statistics cards when competition has started', () => {
		// Set authState to NOT upcoming
		authState.ready = true;
		authState.startTime = new Date(Date.now() - 100000).toISOString();

		render(Page);

		// Cards that should be VISIBLE
		expect(screen.getByText('Category Breakdown')).toBeInTheDocument();
		expect(screen.getByText('Challenges Solved')).toBeInTheDocument();
	});
});

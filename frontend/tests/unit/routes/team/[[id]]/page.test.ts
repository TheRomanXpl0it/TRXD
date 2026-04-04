import { render, screen } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import Page from '$routes/team/[[id]]/+page.svelte';
import { authState } from '$lib/stores/auth';
import { createQuery } from '@tanstack/svelte-query';
import { page } from '$app/stores';

// Mock auth store
vi.mock('$lib/stores/auth', () => ({
    authState: {
        user: { id: 1, team_id: 10 },
        ready: true,
        userMode: true,
        startTime: null,
    },
    loadUser: vi.fn(),
}));

vi.mock('@tanstack/svelte-query', () => ({
    createQuery: vi.fn(),
    useQueryClient: vi.fn(() => ({
        setQueryData: vi.fn(),
        refetchQueries: vi.fn(),
    })),
}));

vi.mock('$app/stores', () => ({
    page: {
        subscribe: (fn: (v: any) => void) => {
            fn({
                params: { id: '10' },
                url: new URL('http://localhost/team/10'),
            });
            return () => {};
        }
    }
}));

vi.mock('$lib/team', () => ({
    getTeam: vi.fn(),
}));

vi.mock('$lib/components/RadarChart.svelte', () => ({
    default: vi.fn(() => ({ render: () => '' }))
}));

describe('Team Profile Page', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        
        // Default mock for team query
        (createQuery as any).mockReturnValue({
            data: { id: 10, name: 'Test Team', solves: [], totalChallenges: [] },
            isLoading: false,
            error: null,
        });
    });

    it('hides statistics cards when competition is upcoming', () => {
        // Set authState to upcoming
        authState.ready = true;
        authState.startTime = new Date(Date.now() + 100000).toISOString();

        render(Page);

        // Cards that should be HIDDEN
        expect(screen.queryByText('Category Breakdown')).not.toBeInTheDocument();
        expect(screen.queryByText('Team Status')).not.toBeInTheDocument();
        
        // Basic team info should still be there
        expect(screen.getByText('Test Team')).toBeInTheDocument();
    });

    it('shows statistics cards when competition has started', () => {
        // Set authState to NOT upcoming
        authState.ready = true;
        authState.startTime = new Date(Date.now() - 100000).toISOString();

        render(Page);

        // Cards that should be VISIBLE
        expect(screen.getByText('Category Breakdown')).toBeInTheDocument();
        expect(screen.getByText('Team Status')).toBeInTheDocument();
    });
});

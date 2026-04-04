import { render, screen } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import Page from '$routes/challenges/+page.svelte';
import { authState } from '$lib/stores/auth';
import { createQuery } from '@tanstack/svelte-query';

// Mock dependencies
vi.mock('$lib/stores/auth', () => ({
    authState: {
        user: null,
        ready: true,
        userMode: true,
        startTime: null,
    }
}));

vi.mock('@tanstack/svelte-query', () => ({
    createQuery: vi.fn(),
    useQueryClient: vi.fn(() => ({
        setQueryData: vi.fn(),
        refetchQueries: vi.fn(),
    })),
}));

vi.mock('$lib/challenges', () => ({
    getChallenges: vi.fn(),
    deleteChallenge: vi.fn(),
}));

vi.mock('$lib/categories', () => ({
    getCategories: vi.fn(),
}));

vi.mock('$lib/env', () => ({
    config: {
        startTime: null
    }
}));

// Mock the admin-only modal to prevent async RPC error from dynamic import
vi.mock('$lib/components/challenges/CreateChallengeModal.svelte', () => ({
    default: {
        name: 'CreateChallengeModal',
        render: () => ({ html: '', css: { code: '', map: null }, head: '' })
    }
}));

describe('Challenges Page', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        
        // Default mock for queries
        (createQuery as any).mockReturnValue({
            data: [],
            isLoading: false,
            error: null,
        });
    });

    it('hides the "Challenges" header when competition is upcoming', () => {
        // Set authState to upcoming
        authState.ready = true;
        authState.startTime = new Date(Date.now() + 100000).toISOString();
        authState.user = { role: 'User' } as any;

        render(Page);

        // Header should not be in the document
        const header = screen.queryByText('Challenges');
        expect(header).not.toBeInTheDocument();
    });

    it('shows the "Challenges" header when competition has started', () => {
        // Set authState to NOT upcoming
        authState.ready = true;
        authState.startTime = new Date(Date.now() - 100000).toISOString();
        authState.user = { role: 'User' } as any;

        render(Page);

        // Header should be in the document
        const header = screen.getByText('Challenges');
        expect(header).toBeInTheDocument();
    });

    it('shows the WaitingPage when isNotStarted is true', () => {
        // Set authState to upcoming (which makes isNotStarted true for non-admins)
        authState.ready = true;
        authState.startTime = new Date(Date.now() + 100000).toISOString();
        authState.user = { role: 'User' } as any;

        render(Page);

        // WaitingPage content should be visible
        expect(screen.getByText('Get your horses ready.')).toBeInTheDocument();
    });

    it('shows challenges for Admins even if competition is upcoming', async () => {
        // Set authState to upcoming but user is Admin
        authState.ready = true;
        authState.startTime = new Date(Date.now() + 100000).toISOString();
        authState.user = { role: 'Admin' } as any;

        // Mock queries to return some data for admin
        (createQuery as any).mockReturnValue({
            data: [{ id: 1, name: 'Admin Chall', category: 'Web', points: 100 }],
            isLoading: false,
            error: null,
        });

        render(Page);
        
        // Wait for dynamic imports of admin-only components
        await new Promise(resolve => setTimeout(resolve, 0));

        // Header should NOT be in the document (the user requested it hidden for EVERYONE if upcoming)
        const header = screen.queryByText('Challenges');
        expect(header).not.toBeInTheDocument();
        
        // Admin controls should be visible (check for a button inside AdminControls)
        expect(screen.getByText('Create Challenge')).toBeInTheDocument();
    });
});

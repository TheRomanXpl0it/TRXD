import { render, screen, waitFor } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import TeamMemberlist from '../TeamMemberlist.svelte';

const mockGetUserData = vi.fn();
vi.mock('$lib/user', () => ({
	getUserData: (...args: any[]) => mockGetUserData(...args)
}));

const team = {
	id: 1,
	name: 'CyberCats',
	score: 4242,
	members: [
		{ id: 10, name: 'Alice', role: 'Captain', score: 1200 },
		{ id: 11, name: 'Bob Sea', role: 'Member', score: 800 },
		{ id: 12, name: 'Seann', role: 'Member', score: 600 },
		{ id: 13, name: 'Zed', role: 'Member', score: 200 }
	]
};

describe('TeamMemberlist', () => {
	beforeEach(() => {
		vi.clearAllMocks();
	});

	it('renders members sorted by score', async () => {
		mockGetUserData.mockResolvedValue({ image: null });
		render(TeamMemberlist, { props: { team } });

		expect(screen.getByText('Alice')).toBeInTheDocument();
		expect(screen.getByText('Bob Sea')).toBeInTheDocument();
		expect(screen.getByText('Seann')).toBeInTheDocument();
		expect(screen.getByText('Zed')).toBeInTheDocument();

		expect(screen.getByText(/1.*200.*pts/)).toBeInTheDocument();
		expect(screen.getByText(/800.*pts/)).toBeInTheDocument();
		expect(screen.getByText('Captain')).toBeInTheDocument();
	});

	it('fetches user images for members', async () => {
		mockGetUserData.mockResolvedValue({ image: 'http://img.png' });
		render(TeamMemberlist, { props: { team } });

		await waitFor(() => {
			expect(mockGetUserData).toHaveBeenCalledWith(10);
			expect(mockGetUserData).toHaveBeenCalledWith(11);
			expect(mockGetUserData).toHaveBeenCalledWith(12);
			expect(mockGetUserData).toHaveBeenCalledWith(13);
		});
	});

	it('renders links to account pages', () => {
		render(TeamMemberlist, { props: { team } });

		const aliceLink = screen.getByRole('link', { name: /Alice/i });
		expect(aliceLink).toHaveAttribute('href', '/account/10');
	});
});

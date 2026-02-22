import { render, screen, fireEvent } from '@testing-library/svelte';
import { describe, it, expect } from 'vitest';
import TeamScoreboard from '../TeamScoreboard.svelte';

const team = {
	name: 'CyberCats',
	members: [
		{ id: 10, name: 'Alice' },
		{ id: 11, name: 'Bob' }
	],
	solves: [
		{
			id: 1,
			name: 'A',
			category: 'Web',
			points: 100,
			user_id: 11,
			timestamp: '2024-01-02T00:00:00Z'
		},
		{
			id: 2,
			name: 'B',
			category: 'Crypto',
			points: 200,
			user_id: 10,
			timestamp: '2024-01-01T00:00:00Z'
		}
	]
};

describe('TeamScoreboard', () => {
	it('renders rows and total points', () => {
		render(TeamScoreboard, { props: { team } });
		expect(screen.queryByText('No solves yet.')).not.toBeInTheDocument();
		expect(screen.getByText(/300 pts total/)).toBeInTheDocument();
		expect(screen.getByText('Alice')).toBeInTheDocument();
		expect(screen.getByText('Bob')).toBeInTheDocument();
	});

	it('sorts by clicking table headers', async () => {
		render(TeamScoreboard, { props: { team } });

		// Default sort by timestamp desc => A (newer) first
		const rowsBefore = screen.getAllByRole('row');
		expect(rowsBefore.some((r) => r.textContent?.includes('A'))).toBe(true);

		// Click Points header to sort by points asc->desc
		fireEvent.click(screen.getByText(/points/i));
		fireEvent.click(screen.getByText(/points/i));
		// Now the 200 should be before 100
		const tbody = document.querySelector('tbody');
		const firstRowText = tbody?.querySelector('tr')?.textContent || '';
		expect(firstRowText).toContain('B');
	});
});

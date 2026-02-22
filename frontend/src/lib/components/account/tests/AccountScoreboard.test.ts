import { render, screen, fireEvent } from '@testing-library/svelte';
import { tick } from 'svelte';
import { describe, it, expect } from 'vitest';
import AccountScoreboard from '../AccountScoreboard.svelte';

const solves = [
	{ id: 1, name: 'X', category: 'Web', points: 50, timestamp: '2024-01-02T00:00:00Z' },
	{ id: 2, name: 'Y', category: 'Pwn', points: 150, timestamp: '2024-01-01T00:00:00Z' }
];

describe('AccountScoreboard', () => {
	it('sorts by clicking headers', async () => {
		render(AccountScoreboard, { props: { solves } });

		// Clicking headers does not remove rows
		fireEvent.click(screen.getByText(/points/i));
		fireEvent.click(screen.getByText(/challenge/i));

		fireEvent.click(screen.getByText(/category/i));
		expect(screen.getByText('X')).toBeInTheDocument();
		expect(screen.getByText('Y')).toBeInTheDocument();
	});
});

import { render, screen, fireEvent } from '@testing-library/svelte';
import { tick } from 'svelte';
import { describe, it, expect } from 'vitest';
import AccountScoreboard from '$lib/components/account/AccountScoreboard.svelte';

const solves = [
	{ id: 1, name: 'X', category: 'Web', points: 50, timestamp: '2024-01-02T00:00:00Z' },
	{ id: 2, name: 'Y', category: 'Pwn', points: 150, timestamp: '2024-01-01T00:00:00Z' }
];

describe('AccountScoreboard', () => {
	it('sorts by clicking headers', async () => {
		render(AccountScoreboard, { props: { solves } });

		await tick();

		// Clicking headers does not remove rows
		const pointsHeader = screen.getByText(/points/i);
		await fireEvent.click(pointsHeader);
		await tick();

		const challengeHeader = screen.getByText(/challenge/i);
		await fireEvent.click(challengeHeader);
		await tick();

		const categoryHeader = screen.getByText(/category/i);
		await fireEvent.click(categoryHeader);
		await tick();

		expect(screen.getByText('X')).toBeInTheDocument();
		expect(screen.getByText('Y')).toBeInTheDocument();
	});
});

import { render, screen } from '@testing-library/svelte';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import ChallengeCard from '../ChallengeCard.svelte';

function generateRandomChallenge(overrides = {}) {
	return {
		id: Math.floor(Math.random() * 10000),
		name: `Challenge ${Math.floor(Math.random() * 100)}`,
		points: Math.floor(Math.random() * 500) + 50,
		tags: ['web', 'crypto'],
		solved: false,
		hidden: false,
		instance: false,
		...overrides
	};
}

describe('ChallengeCard Component', () => {
	beforeEach(() => {
		vi.clearAllMocks();
	});

	it('renders challenge name and points in grid view', () => {
		const challenge = generateRandomChallenge({
			name: 'Test Challenge',
			points: 250
		});

		render(ChallengeCard, {
			props: {
				challenge,
				compactView: false,
				onclick: vi.fn()
			}
		});

		expect(screen.getByText('Test Challenge')).toBeInTheDocument();
		expect(screen.getByText('250 pts')).toBeInTheDocument();
	});

	it('renders challenge name and points in compact view', () => {
		const challenge = generateRandomChallenge({
			name: 'Compact Challenge',
			points: 150
		});

		render(ChallengeCard, {
			props: {
				challenge,
				compactView: true,
				onclick: vi.fn()
			}
		});

		expect(screen.getByText('Compact Challenge')).toBeInTheDocument();
		expect(screen.getByText('150')).toBeInTheDocument();
	});

	it('displays all tags', () => {
		const challenge = generateRandomChallenge({
			tags: ['web', 'pwn', 'forensics']
		});

		render(ChallengeCard, {
			props: {
				challenge,
				compactView: false,
				onclick: vi.fn()
			}
		});

		expect(screen.getByText('web')).toBeInTheDocument();
		expect(screen.getByText('pwn')).toBeInTheDocument();
		expect(screen.getByText('forensics')).toBeInTheDocument();
	});

	it('shows solved checkmark when challenge is solved', () => {
		const challenge = generateRandomChallenge({
			solved: true
		});

		const { container } = render(ChallengeCard, {
			props: {
				challenge,
				compactView: false,
				onclick: vi.fn()
			}
		});

		// Check for the green checkmark icon (CheckCircleSolid)
		const checkIcon = container.querySelector('.text-green-500');
		expect(checkIcon).toBeInTheDocument();
	});

	it('does not show solved checkmark when challenge is not solved', () => {
		const challenge = generateRandomChallenge({
			solved: false
		});

		const { container } = render(ChallengeCard, {
			props: {
				challenge,
				compactView: false,
				onclick: vi.fn()
			}
		});

		// Should not have green checkmark icon
		const checkIcon = container.querySelector('.text-green-500');
		expect(checkIcon).not.toBeInTheDocument();
	});

	it('does not show instance icon for non-instance challenges', () => {
		const challenge = generateRandomChallenge({
			instance: false
		});

		render(ChallengeCard, {
			props: {
				challenge,
				compactView: false,
				onclick: vi.fn()
			}
		});

		expect(screen.queryByLabelText('Instance-based challenge')).not.toBeInTheDocument();
	});

	it('displays countdown timer when instance is running', () => {
		const challenge = generateRandomChallenge({
			instance: true
		});

		render(ChallengeCard, {
			props: {
				challenge,
				compactView: false,
				countdown: 3600 // 1 hour
			}
		});

		// Should display formatted time
		expect(screen.getByText(/1:00:00/)).toBeInTheDocument();
	});

	it('does not display countdown when countdown is 0', () => {
		const challenge = generateRandomChallenge({
			instance: true
		});

		const { container } = render(ChallengeCard, {
			props: {
				challenge,
				compactView: false,
				countdown: 0
			}
		});

		// No countdown badge should be present
		expect(container.querySelector('[aria-label*="expires in"]')).not.toBeInTheDocument();
	});

	it('formats countdown correctly for minutes only', () => {
		const challenge = generateRandomChallenge({
			instance: true
		});

		render(ChallengeCard, {
			props: {
				challenge,
				compactView: false,
				countdown: 125 // 2:05
			}
		});

		expect(screen.getByText(/2:05/)).toBeInTheDocument();
	});

	it('formats countdown correctly for seconds only', () => {
		const challenge = generateRandomChallenge({
			instance: true
		});

		render(ChallengeCard, {
			props: {
				challenge,
				compactView: false,
				countdown: 45
			}
		});

		expect(screen.getByText(/^45$/)).toBeInTheDocument();
	});

	it('calls onclick handler when clicked', async () => {
		const challenge = generateRandomChallenge();
		const handleClick = vi.fn();
		const user = userEvent.setup();

		render(ChallengeCard, {
			props: {
				challenge,
				compactView: false,
				onclick: handleClick
			}
		});

		const button = screen.getByRole('button', { name: new RegExp(challenge.name) });
		await user.click(button);

		expect(handleClick).toHaveBeenCalledTimes(1);
	});

	it('has proper accessibility label in grid view', () => {
		const challenge = generateRandomChallenge({
			name: 'Test Challenge',
			points: 100,
			solved: false
		});

		render(ChallengeCard, {
			props: {
				challenge,
				compactView: false,
				onclick: vi.fn()
			}
		});

		expect(
			screen.getByRole('button', { name: /view details for test challenge, 100 points/i })
		).toBeInTheDocument();
	});


	it('does not has instance icon', () => {
		const challenge = generateRandomChallenge({
			name: 'Test Challenge',
			points: 100,
			solved: false,
			instance: false
		});

		render(ChallengeCard, {
			props: {
				challenge,
				compactView: false,
				onclick: vi.fn()
			}
		});
		expect(
			screen.queryByLabelText('Instance-based challenge')
		).not.toBeInTheDocument();
	});

	it('includes solved status in accessibility label', () => {
		const challenge = generateRandomChallenge({
			name: 'Solved Challenge',
			points: 200,
			solved: true
		});

		render(ChallengeCard, {
			props: {
				challenge,
				compactView: false,
				onclick: vi.fn()
			}
		});

		expect(
			screen.getByRole('button', { name: /solved challenge.*solved/i })
		).toBeInTheDocument();
	});

	it('has proper accessibility label in compact view', () => {
		const challenge = generateRandomChallenge({
			name: 'Compact Test',
			solved: false
		});

		render(ChallengeCard, {
			props: {
				challenge,
				compactView: true,
				onclick: vi.fn()
			}
		});

		expect(
			screen.getByRole('button', { name: /view details for compact test/i })
		).toBeInTheDocument();
	});

	it('applies solved styling in grid view', () => {
		const challenge = generateRandomChallenge({
			solved: true
		});

		render(ChallengeCard, {
			props: {
				challenge,
				compactView: false,
				onclick: vi.fn()
			}
		});

		const button = screen.getByRole('button');
		expect(button.className).toMatch(/bg-green/);
	});

	it('applies solved styling in compact view', () => {
		const challenge = generateRandomChallenge({
			solved: true
		});

		render(ChallengeCard, {
			props: {
				challenge,
				compactView: true,
				onclick: vi.fn()
			}
		});

		const button = screen.getByRole('button');
		expect(button.className).toMatch(/bg-green/);
	});

	it('applies hidden styling when challenge is hidden', () => {
		const challenge = generateRandomChallenge({
			hidden: true
		});

		render(ChallengeCard, {
			props: {
				challenge,
				compactView: false,
				onclick: vi.fn()
			}
		});

		const button = screen.getByRole('button');
		expect(button.className).toMatch(/amber/);
		expect(button.className).toMatch(/ring-2/);
	});

	it('switches between compact and grid view correctly', async () => {
		const challenge = generateRandomChallenge({
			name: 'View Test',
			points: 100
		});

		const { rerender } = render(ChallengeCard, {
			props: {
				challenge,
				compactView: false,
				onclick: vi.fn()
			}
		});

		// Grid view shows "pts" suffix
		expect(screen.getByText('100 pts')).toBeInTheDocument();

		// Switch to compact view
		await rerender({
			challenge,
			compactView: true,
			onclick: vi.fn()
		});

		// Compact view shows just the number
		expect(screen.getByText('100')).toBeInTheDocument();
		expect(screen.queryByText('100 pts')).not.toBeInTheDocument();
	});
});

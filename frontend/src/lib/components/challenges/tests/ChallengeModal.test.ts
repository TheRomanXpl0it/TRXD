import { render, screen, waitFor } from '@testing-library/svelte';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import ChallengeModal from '../ChallengeModal.svelte';
import { toast } from 'svelte-sonner';

// Mock toast
vi.mock('svelte-sonner', () => ({
	toast: {
		success: vi.fn(),
		error: vi.fn()
	}
}));

// Mock API functions used by child components
vi.mock('$lib/challenges', () => ({
	submitFlag: vi.fn()
}));

vi.mock('$lib/instances', () => ({
	startInstance: vi.fn(),
	stopInstance: vi.fn()
}));

function generateRandomChallenge(overrides = {}) {
	return {
		id: Math.floor(Math.random() * 10000),
		name: `Challenge ${Math.floor(Math.random() * 100)}`,
		description: 'This is a test challenge description',
		points: Math.floor(Math.random() * 500) + 50,
		tags: ['web', 'crypto'],
		difficulty: 'medium',
		solves: 0,
		authors: ['Author1', 'Author2'],
		attachments: [],
		instance: false,
		host: null,
		port: null,
		...overrides
	};
}

describe('ChallengeModal Component', () => {
	beforeEach(() => {
		vi.clearAllMocks();
		// Mock clipboard
		Object.defineProperty(navigator, 'clipboard', {
			value: { writeText: vi.fn() },
			writable: true,
			configurable: true
		});
	});

	it('renders challenge name and description', () => {
		const challenge = generateRandomChallenge({
			name: 'Test Challenge',
			description: 'Test description here'
		});

		render(ChallengeModal, {
			props: {
				open: true,
				challenge
			}
		});

		expect(screen.getByText('Test Challenge')).toBeInTheDocument();
		expect(screen.getByText('Test description here')).toBeInTheDocument();
	});

	it('displays all tags', () => {
		const challenge = generateRandomChallenge({
			tags: ['web', 'pwn', 'forensics']
		});

		render(ChallengeModal, {
			props: {
				open: true,
				challenge
			}
		});

		expect(screen.getByText('web')).toBeInTheDocument();
		expect(screen.getByText('pwn')).toBeInTheDocument();
		expect(screen.getByText('forensics')).toBeInTheDocument();
	});



	it('shows blood icon for unsolved challenges', () => {
		const challenge = generateRandomChallenge({
			solves: 0
		});

		const { container } = render(ChallengeModal, {
			props: {
				open: true,
				challenge
			}
		});

		expect(screen.getByText('0 solves')).toBeInTheDocument();
	});

	it('shows solves count as clickable button when solves > 0', () => {
		const challenge = generateRandomChallenge({
			solves: 5
		});

		render(ChallengeModal, {
			props: {
				open: true,
				challenge
			}
		});

		const solvesButton = screen.getByRole('button', { name: /view 5 solves/i });
		expect(solvesButton).toBeInTheDocument();
	});

	it('calls onOpenSolves when solves button is clicked', async () => {
		const challenge = generateRandomChallenge({
			solves: 3
		});
		const onOpenSolves = vi.fn();
		const user = userEvent.setup();

		render(ChallengeModal, {
			props: {
				open: true,
				challenge,
				onOpenSolves
			}
		});

		const solvesButton = screen.getByRole('button', { name: /view 3 solves/i });
		await user.click(solvesButton);

		expect(onOpenSolves).toHaveBeenCalledTimes(1);
	});

	it('displays challenge authors', () => {
		const challenge = generateRandomChallenge({
			authors: ['Alice', 'Bob', 'Charlie']
		});

		render(ChallengeModal, {
			props: {
				open: true,
				challenge
			}
		});

		expect(screen.getByText(/by alice, bob, charlie/i)).toBeInTheDocument();
	});

	it('shows admin controls when isAdmin is true', () => {
		const challenge = generateRandomChallenge();

		render(ChallengeModal, {
			props: {
				open: true,
				challenge,
				isAdmin: true
			}
		});

		expect(screen.getByRole('button', { name: /edit challenge/i })).toBeInTheDocument();
		expect(screen.getByRole('button', { name: /delete challenge/i })).toBeInTheDocument();
	});

	it('hides admin controls when isAdmin is false', () => {
		const challenge = generateRandomChallenge();

		render(ChallengeModal, {
			props: {
				open: true,
				challenge,
				isAdmin: false
			}
		});

		expect(screen.queryByRole('button', { name: /edit challenge/i })).not.toBeInTheDocument();
		expect(screen.queryByRole('button', { name: /delete challenge/i })).not.toBeInTheDocument();
	});

	it('calls onEdit when edit button is clicked', async () => {
		const challenge = generateRandomChallenge();
		const onEdit = vi.fn();
		const user = userEvent.setup();

		render(ChallengeModal, {
			props: {
				open: true,
				challenge,
				isAdmin: true,
				onEdit
			}
		});

		const editButton = screen.getByRole('button', { name: /edit challenge/i });
		await user.click(editButton);

		expect(onEdit).toHaveBeenCalledWith(challenge);
	});

	it('calls onDelete when delete button is clicked', async () => {
		const challenge = generateRandomChallenge();
		const onDelete = vi.fn();
		const user = userEvent.setup();

		render(ChallengeModal, {
			props: {
				open: true,
				challenge,
				isAdmin: true,
				onDelete
			}
		});

		const deleteButton = screen.getByRole('button', { name: /delete challenge/i });
		await user.click(deleteButton);

		expect(onDelete).toHaveBeenCalledWith(challenge);
	});

	it('displays attachments section when attachments exist', () => {
		const challenge = generateRandomChallenge({
			attachments: ['/files/challenge1.zip', '/files/challenge2.txt']
		});

		render(ChallengeModal, {
			props: {
				open: true,
				challenge
			}
		});

		expect(screen.getByText(/attachments/i)).toBeInTheDocument();
		expect(screen.getByText('challenge1.zip')).toBeInTheDocument();
		expect(screen.getByText('challenge2.txt')).toBeInTheDocument();
	});

	it('does not display attachments section when no attachments', () => {
		const challenge = generateRandomChallenge({
			attachments: []
		});

		render(ChallengeModal, {
			props: {
				open: true,
				challenge
			}
		});

		expect(screen.queryByText(/attachments/i)).not.toBeInTheDocument();
	});

	it('attachment links have correct attributes', () => {
		const challenge = generateRandomChallenge({
			attachments: ['/files/test.zip']
		});

		render(ChallengeModal, {
			props: {
				open: true,
				challenge
			}
		});

		const link = screen.getByRole('link', { name: /download test.zip/i });
		expect(link).toHaveAttribute('href', `/attachments/${challenge.id}/files/test.zip`);
		expect(link).toHaveAttribute('target', '_blank');
		expect(link).toHaveAttribute('rel', 'noopener noreferrer');
	});

	it('shows connection string for non-instance challenges with host', () => {
		const challenge = generateRandomChallenge({
			instance: false,
			host: 'ctf.example.com',
			port: 1337
		});

		render(ChallengeModal, {
			props: {
				open: true,
				challenge
			}
		});

		expect(screen.getByText(/connection/i)).toBeInTheDocument();
		expect(screen.getByText('ctf.example.com:1337')).toBeInTheDocument();
	});

	it('copies connection string to clipboard when clicked', async () => {
		const challenge = generateRandomChallenge({
			instance: false,
			host: 'ctf.example.com',
			port: 8080
		});
		const user = userEvent.setup();
		const mockToast = vi.mocked(toast);
		const mockClipboard = vi.fn().mockResolvedValueOnce(undefined);

		Object.defineProperty(navigator, 'clipboard', {
			value: { writeText: mockClipboard },
			writable: true,
			configurable: true
		});

		render(ChallengeModal, {
			props: {
				open: true,
				challenge
			}
		});

		const connectionButton = screen.getByRole('button', { name: /copy connection string/i });
		await user.click(connectionButton);

		expect(mockClipboard).toHaveBeenCalledWith('ctf.example.com:8080');
		await waitFor(() => {
			expect(mockToast.success).toHaveBeenCalledWith('Copied to clipboard!');
		});
	});

	it('does not show connection section for instance challenges', () => {
		const challenge = generateRandomChallenge({
			instance: true,
			host: 'ctf.example.com',
			port: 1337
		});

		render(ChallengeModal, {
			props: {
				open: true,
				challenge
			}
		});

		// Should not show static connection info for instance challenges
		const connectionHeadings = screen.queryAllByText(/connection/i);
		// InstanceControls might have "connection" but not the static connection section
		expect(connectionHeadings.length).toBe(0);
	});

	it('handles connection string without port', () => {
		const challenge = generateRandomChallenge({
			instance: false,
			host: 'ctf.example.com',
			port: null
		});

		render(ChallengeModal, {
			props: {
				open: true,
				challenge
			}
		});

		expect(screen.getByText('ctf.example.com')).toBeInTheDocument();
	});

	it('handles clipboard copy failure gracefully', async () => {
		const challenge = generateRandomChallenge({
			instance: false,
			host: 'ctf.example.com',
			port: 1337
		});
		const user = userEvent.setup();
		const mockToast = vi.mocked(toast);
		const mockClipboard = vi.fn().mockRejectedValueOnce(new Error('Clipboard denied'));

		Object.defineProperty(navigator, 'clipboard', {
			value: { writeText: mockClipboard },
			writable: true,
			configurable: true
		});

		render(ChallengeModal, {
			props: {
				open: true,
				challenge
			}
		});

		const connectionButton = screen.getByRole('button', { name: /copy connection string/i });
		await user.click(connectionButton);

		await waitFor(() => {
			expect(mockToast.error).toHaveBeenCalledWith('Failed to copy to clipboard.');
		});
	});

	it('uses singular "solve" for solves count of 1', () => {
		const challenge = generateRandomChallenge({
			solves: 1
		});

		render(ChallengeModal, {
			props: {
				open: true,
				challenge
			}
		});

		expect(screen.getByRole('button', { name: /view 1 solve$/i })).toBeInTheDocument();
	});

	it('uses plural "solves" for solves count > 1', () => {
		const challenge = generateRandomChallenge({
			solves: 10
		});

		render(ChallengeModal, {
			props: {
				open: true,
				challenge
			}
		});

		expect(screen.getByRole('button', { name: /view 10 solves/i })).toBeInTheDocument();
	});
});

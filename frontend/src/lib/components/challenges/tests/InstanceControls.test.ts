import { render, screen, waitFor } from '@testing-library/svelte';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import InstanceControls from '../InstanceControls.svelte';
import { startInstance, stopInstance } from '$lib/instances';
import { toast } from 'svelte-sonner';

// Mock the instances API
vi.mock('$lib/instances', () => ({
	startInstance: vi.fn(),
	stopInstance: vi.fn()
}));

// Mock svelte-sonner toast
vi.mock('svelte-sonner', () => ({
	toast: {
		success: vi.fn(),
		error: vi.fn()
	}
}));

describe('InstanceControls Component', () => {
	const originalError = console.error;

	beforeEach(() => {
		vi.clearAllMocks();
		// Reset clipboard mock
		Object.defineProperty(navigator, 'clipboard', {
			value: {
				writeText: vi.fn()
			},
			writable: true,
			configurable: true
		});
		// Suppress console.error for cleaner test output
		console.error = vi.fn();
	});

	afterEach(() => {
		// Restore console.error
		console.error = originalError;
	});

	it('starts instance when clicking Start Instance button', async () => {
		const challengeId = Math.floor(Math.random() * 10000);
		const user = userEvent.setup();
		const mockStart = vi.mocked(startInstance);
		const mockToast = vi.mocked(toast);
		const onCountdownUpdate = vi.fn();

		mockStart.mockResolvedValueOnce({
			host: 'ctf.example.com',
			port: 1337,
			timeout: 3600
		});

		const challenge = { id: challengeId, host: null, port: null, timeout: null };
		render(InstanceControls, {
			props: {
				challenge,
				countdown: 0,
				onCountdownUpdate
			}
		});

		const startButton = screen.getByRole('button', { name: /start.*instance/i });
		await user.click(startButton);

		// Verify API was called
		await waitFor(() => {
			expect(mockStart).toHaveBeenCalledWith(challengeId);
		});

		// Verify success toast
		await waitFor(() => {
			expect(mockToast.success).toHaveBeenCalledWith('Instance created!');
		});

		// Verify countdown callback was called
		expect(onCountdownUpdate).toHaveBeenCalledWith(challengeId, 3600);
	});

	it('stops instance when clicking stop button', async () => {
		const challengeId = Math.floor(Math.random() * 10000);
		const user = userEvent.setup();
		const mockStop = vi.mocked(stopInstance);
		const mockToast = vi.mocked(toast);
		const onCountdownUpdate = vi.fn();

		mockStop.mockResolvedValueOnce(undefined);

		const challenge = {
			id: challengeId,
			host: 'ctf.example.com',
			port: 1337,
			timeout: 3600
		};

		render(InstanceControls, {
			props: {
				challenge,
				countdown: 1800, // 30 minutes left
				onCountdownUpdate
			}
		});

		// Should show running state with stop button
		const stopButton = screen.getByRole('button', { name: /stop instance/i });
		await user.click(stopButton);

		// Verify API was called
		await waitFor(() => {
			expect(mockStop).toHaveBeenCalledWith(challengeId);
		});

		// Verify success toast
		await waitFor(() => {
			expect(mockToast.success).toHaveBeenCalledWith('Instance stopped!');
		});

		// Verify countdown callback was called with 0
		expect(onCountdownUpdate).toHaveBeenCalledWith(challengeId, 0);
	});

	it('copies connection string to clipboard when clicking running instance', async () => {
		const challengeId = Math.floor(Math.random() * 10000);
		const user = userEvent.setup();
		const mockToast = vi.mocked(toast);

		// Create a fresh mock for this test
		const mockWriteText = vi.fn().mockResolvedValueOnce(undefined);
		Object.defineProperty(navigator, 'clipboard', {
			value: { writeText: mockWriteText },
			writable: true,
			configurable: true
		});

		const challenge = {
			id: challengeId,
			host: 'ctf.example.com',
			port: 1337,
			timeout: 3600
		};

		render(InstanceControls, {
			props: {
				challenge,
				countdown: 3600
			}
		});

		// Click the running instance button (green background with connection string)
		const instanceButton = screen.getByRole('button', { name: /copy instance connection/i });
		await user.click(instanceButton);

		// Verify clipboard API was called with connection string
		expect(mockWriteText).toHaveBeenCalledWith('ctf.example.com:1337');

		// Verify success toast
		await waitFor(() => {
			expect(mockToast.success).toHaveBeenCalledWith('Copied to clipboard!');
		});
	});

	it('shows error toast when instance creation fails', async () => {
		const challengeId = Math.floor(Math.random() * 10000);
		const user = userEvent.setup();
		const mockStart = vi.mocked(startInstance);
		const mockToast = vi.mocked(toast);

		mockStart.mockRejectedValueOnce(new Error('No available instances'));

		const challenge = { id: challengeId, host: null, port: null, timeout: null };
		render(InstanceControls, {
			props: {
				challenge,
				countdown: 0
			}
		});

		const startButton = screen.getByRole('button', { name: /start.*instance/i });
		await user.click(startButton);

		// Verify error toast
		await waitFor(() => {
			expect(mockToast.error).toHaveBeenCalledWith(
				'Failed to create instance: No available instances'
			);
		});

		// API should have been called
		expect(mockStart).toHaveBeenCalledWith(challengeId);
	});

	it('shows error toast when instance stop fails', async () => {
		const challengeId = Math.floor(Math.random() * 10000);
		const user = userEvent.setup();
		const mockStop = vi.mocked(stopInstance);
		const mockToast = vi.mocked(toast);

		mockStop.mockRejectedValueOnce(new Error('Instance not found'));

		const challenge = {
			id: challengeId,
			host: 'ctf.example.com',
			port: 1337,
			timeout: 3600
		};

		render(InstanceControls, {
			props: {
				challenge,
				countdown: 1800
			}
		});

		const stopButton = screen.getByRole('button', { name: /stop instance/i });
		await user.click(stopButton);

		// Verify error toast
		await waitFor(() => {
			expect(mockToast.error).toHaveBeenCalledWith('Failed to stop instance: Instance not found');
		});
	});

	it('disables start button during instance creation', async () => {
		const challengeId = Math.floor(Math.random() * 10000);
		const user = userEvent.setup();
		const mockStart = vi.mocked(startInstance);

		// Make the API call slow
		mockStart.mockImplementationOnce(
			() =>
				new Promise((resolve) =>
					setTimeout(() => resolve({ host: 'test', port: 1337, timeout: 3600 }), 100)
				)
		);

		const challenge = { id: challengeId, host: null, port: null, timeout: null };
		render(InstanceControls, {
			props: {
				challenge,
				countdown: 0
			}
		});

		const startButton = screen.getByRole('button', { name: /start.*instance/i });

		// Click button
		await user.click(startButton);

		// Button should show "Starting..." and be disabled
		await waitFor(() => {
			expect(screen.getByText(/starting/i)).toBeInTheDocument();
			expect(startButton).toBeDisabled();
		});
	});

	it('disables stop button during instance destruction', async () => {
		const challengeId = Math.floor(Math.random() * 10000);
		const user = userEvent.setup();
		const mockStop = vi.mocked(stopInstance);

		// Make the API call slow
		mockStop.mockImplementationOnce(
			() => new Promise((resolve) => setTimeout(() => resolve(undefined), 100))
		);

		const challenge = {
			id: challengeId,
			host: 'ctf.example.com',
			port: 1337,
			timeout: 3600
		};

		render(InstanceControls, {
			props: {
				challenge,
				countdown: 1800
			}
		});

		const stopButton = screen.getByRole('button', { name: /stop instance/i });

		// Click button
		await user.click(stopButton);

		// Button should be disabled
		await waitFor(() => {
			expect(stopButton).toBeDisabled();
		});
	});

	it('prevents multiple simultaneous instance creation requests', async () => {
		const challengeId = Math.floor(Math.random() * 10000);
		const user = userEvent.setup();
		const mockStart = vi.mocked(startInstance);

		mockStart.mockImplementationOnce(
			() =>
				new Promise((resolve) =>
					setTimeout(() => resolve({ host: 'test', port: 1337, timeout: 3600 }), 100)
				)
		);

		const challenge = { id: challengeId, host: null, port: null, timeout: null };
		render(InstanceControls, {
			props: {
				challenge,
				countdown: 0
			}
		});

		const startButton = screen.getByRole('button', { name: /start.*instance/i });

		// Try to click multiple times rapidly
		await Promise.all([user.click(startButton), user.click(startButton), user.click(startButton)]);

		// Wait for request to complete
		await waitFor(() => {
			expect(mockStart).toHaveBeenCalledTimes(1);
		});
	});

	it('formats countdown timer correctly for hours', async () => {
		const challengeId = Math.floor(Math.random() * 10000);
		const challenge = { id: challengeId, host: 'test.com', port: 1337, timeout: 3600 };

		render(InstanceControls, {
			props: {
				challenge,
				countdown: 7265 // 2:01:05
			}
		});

		// Should show formatted time with hours
		expect(screen.getByText(/2:01:05/)).toBeInTheDocument();
	});

	it('formats countdown timer correctly for minutes', async () => {
		const challengeId = Math.floor(Math.random() * 10000);
		const challenge = { id: challengeId, host: 'test.com', port: 1337, timeout: 3600 };

		render(InstanceControls, {
			props: {
				challenge,
				countdown: 125 // 2:05
			}
		});

		// Should show formatted time without hours
		expect(screen.getByText(/2:05/)).toBeInTheDocument();
	});

	it('formats countdown timer correctly for seconds only', async () => {
		const challengeId = Math.floor(Math.random() * 10000);
		const challenge = { id: challengeId, host: 'test.com', port: 1337, timeout: 3600 };

		render(InstanceControls, {
			props: {
				challenge,
				countdown: 45
			}
		});

		// Should show just seconds
		expect(screen.getByText(/^45$/)).toBeInTheDocument();
	});

	it('displays connection string with port', async () => {
		const challengeId = Math.floor(Math.random() * 10000);
		const challenge = {
			id: challengeId,
			host: 'ctf.example.com',
			port: 1337,
			timeout: 3600
		};

		render(InstanceControls, {
			props: {
				challenge,
				countdown: 3600
			}
		});

		// Should show host:port
		expect(screen.getByText('ctf.example.com:1337')).toBeInTheDocument();
	});

	it('displays connection string without port when port is null', async () => {
		const challengeId = Math.floor(Math.random() * 10000);
		const challenge = {
			id: challengeId,
			host: 'ctf.example.com',
			port: null,
			timeout: 3600
		};

		render(InstanceControls, {
			props: {
				challenge,
				countdown: 3600
			}
		});

		// Should show just host (no port)
		expect(screen.getByText('ctf.example.com')).toBeInTheDocument();
	});

	it('handles clipboard copy failure gracefully', async () => {
		const challengeId = Math.floor(Math.random() * 10000);
		const user = userEvent.setup();
		const mockToast = vi.mocked(toast);

		// Mock clipboard to reject
		navigator.clipboard.writeText = vi
			.fn()
			.mockRejectedValueOnce(new Error('Clipboard access denied'));

		const challenge = {
			id: challengeId,
			host: 'ctf.example.com',
			port: 1337,
			timeout: 3600
		};

		render(InstanceControls, {
			props: {
				challenge,
				countdown: 3600
			}
		});

		const instanceButton = screen.getByRole('button', { name: /copy instance connection/i });
		await user.click(instanceButton);

		// Verify error toast
		await waitFor(() => {
			expect(mockToast.error).toHaveBeenCalledWith('Failed to copy to clipboard.');
		});
	});

	it('shows start button when countdown is 0', () => {
		const challengeId = Math.floor(Math.random() * 10000);
		const challenge = { id: challengeId, host: null, port: null, timeout: null };

		render(InstanceControls, {
			props: {
				challenge,
				countdown: 0
			}
		});

		// Should show Start Instance button
		expect(screen.getByRole('button', { name: /start.*instance/i })).toBeInTheDocument();
		// Should NOT show running instance button
		expect(
			screen.queryByRole('button', { name: /copy instance connection/i })
		).not.toBeInTheDocument();
	});

	it('shows running instance button when countdown is greater than 0', () => {
		const challengeId = Math.floor(Math.random() * 10000);
		const challenge = {
			id: challengeId,
			host: 'ctf.example.com',
			port: 1337,
			timeout: 3600
		};

		render(InstanceControls, {
			props: {
				challenge,
				countdown: 1800
			}
		});

		// Should show running instance with countdown
		expect(screen.getByRole('button', { name: /copy instance connection/i })).toBeInTheDocument();
		// Should NOT show Start Instance button
		expect(screen.queryByRole('button', { name: /start.*instance/i })).not.toBeInTheDocument();
	});

	it('updates displayed countdown when prop changes', async () => {
		const challengeId = Math.floor(Math.random() * 10000);
		const challenge = {
			id: challengeId,
			host: 'ctf.example.com',
			port: 1337,
			timeout: 3600
		};

		const { rerender } = render(InstanceControls, {
			props: {
				challenge,
				countdown: 3600 // 1:00:00
			}
		});

		// Initial countdown shows 1:00:00
		expect(screen.getByText(/1:00:00/)).toBeInTheDocument();

		// Update countdown to 30 seconds
		await rerender({
			challenge,
			countdown: 30
		});

		// Should now show 30
		expect(screen.getByText(/^30$/)).toBeInTheDocument();
		expect(screen.queryByText(/1:00:00/)).not.toBeInTheDocument();
	});

	it('switches to start button when countdown reaches 0', async () => {
		const challengeId = Math.floor(Math.random() * 10000);
		const challenge = {
			id: challengeId,
			host: 'ctf.example.com',
			port: 1337,
			timeout: 3600
		};

		const { rerender } = render(InstanceControls, {
			props: {
				challenge,
				countdown: 10 // About to expire
			}
		});

		// Should show running instance
		expect(screen.getByRole('button', { name: /copy instance connection/i })).toBeInTheDocument();

		// Update countdown to 0 (expired)
		await rerender({
			challenge: { ...challenge, host: null, port: null, timeout: null },
			countdown: 0
		});

		// Should now show start button
		expect(screen.getByRole('button', { name: /start.*instance/i })).toBeInTheDocument();
		expect(
			screen.queryByRole('button', { name: /copy instance connection/i })
		).not.toBeInTheDocument();
	});

	it('handles negative countdown values as 0', () => {
		const challengeId = Math.floor(Math.random() * 10000);
		const challenge = {
			id: challengeId,
			host: 'ctf.example.com',
			port: 1337,
			timeout: 3600
		};

		render(InstanceControls, {
			props: {
				challenge,
				countdown: -50 // Negative value
			}
		});

		// With countdown <= 0, should show start button instead of running state
		expect(screen.getByRole('button', { name: /start.*instance/i })).toBeInTheDocument();
		expect(
			screen.queryByRole('button', { name: /copy instance connection/i })
		).not.toBeInTheDocument();
	});
});

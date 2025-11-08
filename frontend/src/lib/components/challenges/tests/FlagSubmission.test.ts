import { submitFlag } from '$lib/challenges';
import { toast } from 'svelte-sonner';


import { render, screen, waitFor } from '@testing-library/svelte';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import FlagSubmission from '../FlagSubmission.svelte';

// Mock the submitFlag API
vi.mock('$lib/challenges', () => ({
	submitFlag: vi.fn()
}));

// Mock svelte-sonner toast
vi.mock('svelte-sonner', () => ({
	toast: {
		success: vi.fn(),
		error: vi.fn()
	}
}));

// Helper function to generate random flags
function generateRandomFlag(): string {
	const chars = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_';
	const length = Math.floor(Math.random() * 20) + 10; // 10-30 characters
	let flag = 'TRX{';
	for (let i = 0; i < length; i++) {
		flag += chars[Math.floor(Math.random() * chars.length)];
	}
	flag += '}';
	return flag;
}

describe('FlagSubmission Component - User Workflow', () => {
	beforeEach(() => {
		vi.clearAllMocks();
	});

	it('allows user to submit a correct flag and shows success', async () => {
		const challengeId = Math.floor(Math.random() * 10000);
		const testFlag = generateRandomFlag();
		const user = userEvent.setup();
		const mockSubmit = vi.mocked(submitFlag);
		mockSubmit.mockResolvedValueOnce({ status: 'Correct', first_blood: false });

		const onSolved = vi.fn();
		render(FlagSubmission, { 
			props: { 
				challenge: { id: challengeId, solved: false },
				onSolved
			}
		});

		// User types a flag
		const input = screen.getByPlaceholderText(/TRX\{/i) as HTMLInputElement;
		await user.click(input);
		await user.paste(testFlag);

		// User clicks submit
		const submitButton = screen.getByRole('button', { name: /submit/i });
		await user.click(submitButton);

		// Verify API was called with correct data
		await waitFor(() => {
			expect(mockSubmit).toHaveBeenCalledWith(challengeId, testFlag);
		});

		// Check that challenge is marked as solved
		await waitFor(() => {
			expect(screen.getByText(/challenge solved/i)).toBeInTheDocument();
		});
		
		// Verify onSolved callback was called
		expect(onSolved).toHaveBeenCalled();
	});

	it('clears error state when user starts typing after wrong flag', async () => {
		const challengeId = Math.floor(Math.random() * 10000);
		const testFlag = generateRandomFlag();
		const user = userEvent.setup();
		const mockSubmit = vi.mocked(submitFlag);
		mockSubmit.mockResolvedValueOnce({ status: 'Wrong' });

		render(FlagSubmission, { 
			props: { 
				challenge: { id: challengeId, solved: false },
				onSolved: vi.fn()
			}
		});

		const input = screen.getByPlaceholderText(/TRX\{/i) as HTMLInputElement;
		
		// Submit wrong flag
		await user.click(input);
		await user.paste(testFlag);
		const submitButton = screen.getByRole('button', { name: /submit/i });
		await user.click(submitButton);

		// Verify error state appears
		await waitFor(() => {
			expect(input).toHaveAttribute('aria-invalid', 'true');
		});
		const errorIcon = document.querySelector('.text-red-500');
		expect(errorIcon).toBeInTheDocument();

		// User starts typing - error should clear immediately
		await user.clear(input);
		await user.type(input, 'T');
		
		// Error state should be gone
		expect(input).toHaveAttribute('aria-invalid', 'false');
		const flagIcon = document.querySelector('.text-red-500');
		expect(flagIcon).not.toBeInTheDocument();
	});

	it('prevents submission when input is empty', async () => {
		const challengeId = Math.floor(Math.random() * 10000);
		const user = userEvent.setup();
		const mockSubmit = vi.mocked(submitFlag);

		render(FlagSubmission, { 
			props: { 
				challenge: { id: challengeId, solved: false },
				onSolved: vi.fn()
			}
		});

		const submitButton = screen.getByRole('button', { name: /submit/i });
		
		// Button should be disabled when input is empty
		expect(submitButton).toBeDisabled();
		
		// Clicking disabled button should not call API
		await user.click(submitButton);
		expect(mockSubmit).not.toHaveBeenCalled();
	});

	it('handles API errors gracefully', async () => {
		const challengeId = Math.floor(Math.random() * 10000);
		const testFlag = generateRandomFlag();
		const user = userEvent.setup();
		const mockSubmit = vi.mocked(submitFlag);
		mockSubmit.mockRejectedValueOnce(new Error('Network error'));

		render(FlagSubmission, { 
			props: { 
				challenge: { id: challengeId, solved: false },
				onSolved: vi.fn()
			}
		});

		const input = screen.getByPlaceholderText(/TRX\{/i) as HTMLInputElement;
		await user.click(input);
		await user.paste(testFlag);

		const submitButton = screen.getByRole('button', { name: /submit/i });
		await user.click(submitButton);

		// Should show error state (aria-invalid on input)
		await waitFor(() => {
			expect(input).toHaveAttribute('aria-invalid', 'true');
		});
		
		// Should show error icon instead of flag icon
		const errorIcon = document.querySelector('.text-red-500');
		expect(errorIcon).toBeInTheDocument();
	});


	it('Hides flag input and button if submission succeded', async () => {
		const challengeId = Math.floor(Math.random() * 10000);
        const testFlag = generateRandomFlag();
		const user = userEvent.setup();
		const mockSubmit = vi.mocked(submitFlag);
		
		render(FlagSubmission, { 
			props: { 
				challenge: { id: challengeId, solved: false },
				onSolved: vi.fn()
			}
		});

		const input = screen.getByPlaceholderText(/TRX\{/i) as HTMLInputElement;
		const submitButton = screen.getByRole('button', { name: /submit/i });

		await user.clear(input);
		await user.paste(testFlag);
		expect(input.value).toBe(testFlag);

		mockSubmit.mockResolvedValueOnce({ status: 'Correct' });
		await user.click(submitButton);

		await waitFor(() => {
			expect(mockSubmit).toHaveBeenCalledWith(challengeId, testFlag);
		});

        // try to click again, this should fail
        await user.click(submitButton);

		// The success span (also the success div) should be in the document
		const successSpan = screen.getByText(/challenge solved/i);
		expect(successSpan).toBeInTheDocument();

		// Also, the input element should not be present anymore in the document
		await waitFor(() => {
			expect(input).not.toBeInTheDocument();
		});

		// Same thing for the submit button
		await waitFor(() => {
			expect(submitButton).not.toBeInTheDocument();
		});



		expect(mockSubmit).toHaveBeenCalledTimes(1);
	});


	it('allows editing and resubmitting a different flag. Also test wrong flag submission behaviour', async () => {
		const challengeId = Math.floor(Math.random() * 10000);
		const testFlag1 = generateRandomFlag();
        const testFlag2 = generateRandomFlag();
        
        const user = userEvent.setup();
		const mockSubmit = vi.mocked(submitFlag);
		
		render(FlagSubmission, { 
			props: { 
				challenge: { id: challengeId, solved: false },
				onSolved: vi.fn()
			}
		});

		const input = screen.getByPlaceholderText(/TRX\{/i) as HTMLInputElement;
		const submitButton = screen.getByRole('button', { name: /submit/i });

        await user.click(input);
		await user.paste(testFlag1);
		expect(input.value).toBe(testFlag1);

		await user.clear(input);
		await user.paste(testFlag2);
		expect(input.value).toBe(testFlag2);

		mockSubmit.mockResolvedValueOnce({ status: 'Wrong' });
		await user.click(submitButton);

		await waitFor(() => {
			expect(mockSubmit).toHaveBeenCalledWith(challengeId, testFlag2);
		});
		expect(mockSubmit).toHaveBeenCalledTimes(1);
	});

	it('Verify that the first blood toast shows up correctly', async () => {
		const challengeId = Math.floor(Math.random() * 10000);
        const testFlag = generateRandomFlag();
		const user = userEvent.setup();
		const mockSubmit = vi.mocked(submitFlag);
		const mockToast = vi.mocked(toast);
		
		render(FlagSubmission, { 
			props: { 
				challenge: { id: challengeId, solved: false },
				onSolved: vi.fn()
			}
		});

		const input = screen.getByPlaceholderText(/TRX\{/i) as HTMLInputElement;
		const submitButton = screen.getByRole('button', { name: /submit/i });

		await user.click(input);
		await user.paste(testFlag);

		mockSubmit.mockResolvedValueOnce({ status: 'Correct', first_blood: true });
		await user.click(submitButton);

		await waitFor(() => {
			expect(mockSubmit).toHaveBeenCalledWith(challengeId, testFlag);
		});

		// Verify first blood toast was shown
		expect(mockToast.success).toHaveBeenCalledWith('First blood!');
		expect(mockSubmit).toHaveBeenCalledTimes(1);
		
		// Challenge should be marked as solved
		await waitFor(() => {
			expect(screen.getByText(/challenge solved/i)).toBeInTheDocument();
		});
	});


	it('Verify that the first blood toast does not show up', async () => {
		const challengeId = Math.floor(Math.random() * 10000);
        const testFlag = generateRandomFlag();
		const user = userEvent.setup();
		const mockSubmit = vi.mocked(submitFlag);
		const mockToast = vi.mocked(toast);
		
		render(FlagSubmission, { 
			props: { 
				challenge: { id: challengeId, solved: false },
				onSolved: vi.fn()
			}
		});

		const input = screen.getByPlaceholderText(/TRX\{/i) as HTMLInputElement;
		const submitButton = screen.getByRole('button', { name: /submit/i });

		await user.click(input);
		await user.paste(testFlag);

		mockSubmit.mockResolvedValueOnce({ status: 'Correct', first_blood: false });
		await user.click(submitButton);

		await waitFor(() => {
			expect(mockSubmit).toHaveBeenCalledWith(challengeId, testFlag);
		});

		// Verify first blood toast was not shown
		expect(mockToast.success).not.toHaveBeenCalledWith('First blood!');
		expect(mockSubmit).toHaveBeenCalledTimes(1);
		
		// Challenge should be marked as solved
		await waitFor(() => {
			expect(screen.getByText(/challenge solved/i)).toBeInTheDocument();
		});
	});


	it('Verify that the success toast shows up', async () => {
		const challengeId = Math.floor(Math.random() * 10000);
        const testFlag = generateRandomFlag();
		const user = userEvent.setup();
		const mockSubmit = vi.mocked(submitFlag);
		const mockToast = vi.mocked(toast);
		
		render(FlagSubmission, { 
			props: { 
				challenge: { id: challengeId, solved: false },
				onSolved: vi.fn()
			}
		});

		const input = screen.getByPlaceholderText(/TRX\{/i) as HTMLInputElement;
		const submitButton = screen.getByRole('button', { name: /submit/i });

		await user.click(input);
		await user.paste(testFlag);

		mockSubmit.mockResolvedValueOnce({ status: 'Correct', first_blood: false });
		await user.click(submitButton);

		await waitFor(() => {
			expect(mockSubmit).toHaveBeenCalledWith(challengeId, testFlag);
		});

		// Verify success toast was shown
		expect(mockToast.success).toHaveBeenCalledWith('Correct flag!');
		expect(mockSubmit).toHaveBeenCalledTimes(1);
		
		// Challenge should be marked as solved
		await waitFor(() => {
			expect(screen.getByText(/challenge solved/i)).toBeInTheDocument();
		});
	});

    it('Verify that the success toast does not show up', async () => {
		const challengeId = Math.floor(Math.random() * 10000);
        const testFlag = generateRandomFlag();
		const user = userEvent.setup();
		const mockSubmit = vi.mocked(submitFlag);
		const mockToast = vi.mocked(toast);
		
		render(FlagSubmission, { 
			props: { 
				challenge: { id: challengeId, solved: false },
				onSolved: vi.fn()
			}
		});

		const input = screen.getByPlaceholderText(/TRX\{/i) as HTMLInputElement;
		const submitButton = screen.getByRole('button', { name: /submit/i });

		await user.click(input);
		await user.paste(testFlag);

		mockSubmit.mockResolvedValueOnce({ status: 'Wrong' });
		await user.click(submitButton);

		await waitFor(() => {
			expect(mockSubmit).toHaveBeenCalledWith(challengeId, testFlag);
		});

		// Verify success toast was not shown
		expect(mockToast.error).toHaveBeenCalledWith('Incorrect flag');
		expect(mockSubmit).toHaveBeenCalledTimes(1);
		
		// Challenge should not be marked as solved
		await waitFor(() => {
			expect(screen.queryByText(/challenge solved/i)).not.toBeInTheDocument();
		});
	});

	it('Flag submission with whitespaces', async () => {
		const challengeId = Math.floor(Math.random() * 10000);
        const testFlag = generateRandomFlag();
		const user = userEvent.setup();
		const mockSubmit = vi.mocked(submitFlag);
		
		render(FlagSubmission, { 
			props: { 
				challenge: { id: challengeId, solved: false },
				onSolved: vi.fn()
			}
		});

		const input = screen.getByPlaceholderText(/TRX\{/i) as HTMLInputElement;
		const submitButton = screen.getByRole('button', { name: /submit/i });

		await user.click(input);

        const spaces = Math.floor(Math.random() * 10);
        await user.paste(' '.repeat(spaces) + testFlag + ' '.repeat(spaces));

		mockSubmit.mockResolvedValueOnce({ status: 'Correct', first_blood: false });
		await user.click(submitButton);

		await waitFor(() => {
			expect(mockSubmit).toHaveBeenCalledWith(challengeId, testFlag);
		});

		expect(mockSubmit).toHaveBeenCalledTimes(1);
		
		await waitFor(() => {
			expect(screen.getByText(/challenge solved/i)).toBeInTheDocument();
		});
	});


	it('Rapid flag submissions', async () => {
		const challengeId = Math.floor(Math.random() * 10000);
        const testFlag = generateRandomFlag();
		const user = userEvent.setup();
		const mockSubmit = vi.mocked(submitFlag);
		
		render(FlagSubmission, { 
			props: { 
				challenge: { id: challengeId, solved: false },
				onSolved: vi.fn()
			}
		});

		const input = screen.getByPlaceholderText(/TRX\{/i) as HTMLInputElement;
		const submitButton = screen.getByRole('button', { name: /submit/i });

		await user.click(input);

        await user.paste(testFlag);

		mockSubmit.mockResolvedValueOnce({ status: 'Correct', first_blood: false });
		
		await Promise.all([
			user.click(submitButton),
			user.click(submitButton),
            user.click(submitButton),
            user.click(submitButton)
		]);

		await waitFor(() => {
			expect(mockSubmit).toHaveBeenCalledWith(challengeId, testFlag);
		});

		expect(mockSubmit).toHaveBeenCalledTimes(1);
		
		await waitFor(() => {
			expect(screen.getByText(/challenge solved/i)).toBeInTheDocument();
		});
	});
});

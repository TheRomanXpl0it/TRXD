
import userEvent from '@testing-library/user-event';
import { render, screen, waitFor, fireEvent } from '@testing-library/svelte';
import { tick } from 'svelte';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import DeleteChallengeDialog from '../DeleteChallengeDialog.svelte';


describe('DeleteChallengeDialog Component', () => {
	beforeEach(() => {
		vi.clearAllMocks();
	});

	afterEach(async () => {
		await new Promise(resolve => setTimeout(resolve, 150));
	});

	it('renders dialog title', () => {
		render(DeleteChallengeDialog, {
			props: {
				open: true,
				toDelete: { name: 'Test Challenge' },
				deleting: false
			}
		});

		expect(screen.getByText('Delete challenge?')).toBeInTheDocument();
	});

	it('displays challenge name in description', () => {
		render(DeleteChallengeDialog, {
			props: {
				open: true,
				toDelete: { name: 'My Test Challenge' },
				deleting: false
			}
		});

		// The name appears in a <b> tag within the description
		const descriptionText = screen.getByText(/you're about to permanently delete/i);
		expect(descriptionText).toBeInTheDocument();
		
		// Check that the challenge name is in the document
		const challengeName = screen.getAllByText(/my test challenge/i)[0];
		expect(challengeName).toBeInTheDocument();
	});

	it('shows warning text about permanent deletion', () => {
		render(DeleteChallengeDialog, {
			props: {
				open: true,
				toDelete: { name: 'Test' },
				deleting: false
			}
		});

		expect(screen.getByText(/this action cannot be undone/i)).toBeInTheDocument();
	});

	it('shows warning about related data removal', () => {
		render(DeleteChallengeDialog, {
			props: {
				open: true,
				toDelete: { name: 'Test' },
				deleting: false
			}
		});

		expect(screen.getByText(/all related data.*may be removed/i)).toBeInTheDocument();
	});

	it('displays confirmation prompt text', () => {
		render(DeleteChallengeDialog, {
			props: {
				open: true,
				toDelete: { name: 'Test' },
				deleting: false
			}
		});

		expect(screen.getByText(/to confirm, type the following text/i)).toBeInTheDocument();
	});

	it('shows expected confirmation text', () => {
		render(DeleteChallengeDialog, {
			props: {
				open: true,
				toDelete: { name: 'Web Challenge' },
				deleting: false
			}
		});

		expect(screen.getByText("Yes, I want to delete 'Web Challenge'")).toBeInTheDocument();
	});

	it('renders confirmation input field', () => {
		render(DeleteChallengeDialog, {
			props: {
				open: true,
				toDelete: { name: 'Test' },
				deleting: false
			}
		});

		expect(screen.getByLabelText(/confirmation/i)).toBeInTheDocument();
	});

	it('renders cancel button', () => {
		render(DeleteChallengeDialog, {
			props: {
				open: true,
				toDelete: { name: 'Test' },
				deleting: false
			}
		});

		const cancelButtons = screen.getAllByRole('button', { name: /cancel/i });
		// There might be multiple due to Dialog.Close wrapper
		expect(cancelButtons.length).toBeGreaterThan(0);
	});

	it('renders delete button', () => {
		render(DeleteChallengeDialog, {
			props: {
				open: true,
				toDelete: { name: 'Test' },
				deleting: false
			}
		});

		expect(screen.getByRole('button', { name: /^delete$/i })).toBeInTheDocument();
	});

	it('delete button is disabled initially', () => {
		render(DeleteChallengeDialog, {
			props: {
				open: true,
				toDelete: { name: 'Test Challenge' },
				deleting: false
			}
		});

		const deleteButton = screen.getByRole('button', { name: /^delete$/i });
		expect(deleteButton).toBeDisabled();
	});

	it('delete button is disabled when confirmation text is incorrect', async () => {
		const user = userEvent.setup();

		render(DeleteChallengeDialog, {
			props: {
				open: true,
				toDelete: { name: 'Test' },
				deleting: false
			}
		});

		const input = screen.getByLabelText(/confirmation/i);
		await user.type(input, 'wrong text');

		const deleteButton = screen.getByRole('button', { name: /^delete$/i });
		expect(deleteButton).toBeDisabled();
	});

	it('delete button is enabled when confirmation text matches exactly', async () => {
		const user = userEvent.setup();

		render(DeleteChallengeDialog, {
			props: {
				open: true,
				toDelete: { name: 'Test Challenge' },
				deleting: false
			}
		});

		const input = await screen.findByLabelText(/confirmation/i);
		const expectedText = "Yes, I want to delete 'Test Challenge'";
		await user.clear(input);
		await user.type(input, expectedText);

		// Wait for the button to become enabled
		const deleteButton = screen.getByRole('button', { name: /^delete$/i });
		await waitFor(() => {
			expect(deleteButton).not.toBeDisabled();
		});
	});

	it('calls onconfirm when delete button clicked with correct text', async () => {
		const user = userEvent.setup();
		const handleConfirm = vi.fn();

		render(DeleteChallengeDialog, {
			props: {
			open: true,
			toDelete: { name: 'Test' },
			deleting: false,
			onconfirm: handleConfirm,
			},
		});

		const input = screen.getByLabelText(/confirmation/i);
		await user.clear(input);
		await user.type(input, "Yes, I want to delete 'Test'");

		// Let Svelte recompute canDelete
		await tick();

		const deleteButton = screen.getByRole('button', { name: /^delete$/i });
		await waitFor(() => {
			expect(deleteButton).not.toBeDisabled();
		});

		// Use fireEvent to sidestep pointer-events/body-scroll-lock weirdness
		await fireEvent.click(deleteButton);

		await waitFor(() => {
			expect(handleConfirm).toHaveBeenCalledTimes(1);
		});
	});


	it('shows spinner and "Deleting..." text when deleting is true', () => {
		render(DeleteChallengeDialog, {
			props: {
				open: true,
				toDelete: { name: 'Test' },
				deleting: true
			}
		});

		expect(screen.getByText(/deleting/i)).toBeInTheDocument();
	});

	it('disables input when deleting is true', () => {
		render(DeleteChallengeDialog, {
			props: {
				open: true,
				toDelete: { name: 'Test' },
				deleting: true
			}
		});

		const input = screen.getByLabelText(/confirmation/i);
		expect(input).toBeDisabled();
	});

	it('disables cancel button when deleting is true', () => {
		render(DeleteChallengeDialog, {
			props: {
				open: true,
				toDelete: { name: 'Test' },
				deleting: true
			}
		});

		const cancelButtons = screen.getAllByRole('button', { name: /cancel/i });
		// At least one cancel button should be disabled
		const disabledButton = cancelButtons.find(btn => btn.hasAttribute('disabled'));
		expect(disabledButton).toBeDefined();
	});

	it('disables delete button when deleting is true', () => {
		render(DeleteChallengeDialog, {
			props: {
				open: true,
				toDelete: { name: 'Test' },
				deleting: true
			}
		});

		const deleteButton = screen.getByRole('button', { name: /deleting/i });
		expect(deleteButton).toBeDisabled();
	});

	it('confirmation text is case-sensitive', async () => {
		const user = userEvent.setup();

		render(DeleteChallengeDialog, {
			props: {
				open: true,
				toDelete: { name: 'Test' },
				deleting: false
			}
		});

		const input = screen.getByLabelText(/confirmation/i);
		// Wrong case
		await user.type(input, "yes, i want to delete 'Test'");

		const deleteButton = screen.getByRole('button', { name: /^delete$/i });
		expect(deleteButton).toBeDisabled();
	});

	it('confirmation requires exact punctuation', async () => {
		const user = userEvent.setup();

		render(DeleteChallengeDialog, {
			props: {
				open: true,
				toDelete: { name: 'Test' },
				deleting: false
			}
		});

		const input = screen.getByLabelText(/confirmation/i);
		// Missing comma
		await user.type(input, "Yes I want to delete 'Test'");

		const deleteButton = screen.getByRole('button', { name: /^delete$/i });
		expect(deleteButton).toBeDisabled();
	});

	it('handles challenge name with special characters in confirmation', () => {
		render(DeleteChallengeDialog, {
			props: {
				open: true,
				toDelete: { name: "Test's \"Challenge\" & More" },
				deleting: false
			}
		});

		expect(screen.getByText("Yes, I want to delete 'Test's \"Challenge\" & More'")).toBeInTheDocument();
	});

	it('shows placeholder text in input field', () => {
		render(DeleteChallengeDialog, {
			props: {
				open: true,
				toDelete: { name: 'Test' },
				deleting: false
			}
		});

		expect(screen.getByPlaceholderText(/type here to confirm/i)).toBeInTheDocument();
	});

	it('clears confirmation text when dialog is closed and reopened', async () => {
		const user = userEvent.setup();

		const { rerender } = render(DeleteChallengeDialog, {
			props: {
				open: true,
				toDelete: { name: 'Test' },
				deleting: false
			}
		});

		const input = screen.getByLabelText(/confirmation/i);
		await user.type(input, 'some text');
		
		expect(input).toHaveValue('some text');

		// Close dialog
		await rerender({
			open: false,
			toDelete: { name: 'Test' },
			deleting: false
		});

		// Reopen dialog
		await rerender({
			open: true,
			toDelete: { name: 'Test' },
			deleting: false
		});

		const newInput = screen.getByLabelText(/confirmation/i);
		expect(newInput).toHaveValue('');
	});
});

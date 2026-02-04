import { render, screen, waitFor, fireEvent } from '@testing-library/svelte';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import AccountEdit from '../AccountEdit.svelte';
import { updateUser } from '$lib/user';
import { showError, showSuccess } from '$lib/utils/toast';
import { tick } from 'svelte';

async function flush() {
	await tick();
	await Promise.resolve();
}

vi.mock('$lib/utils/toast', () => ({
	showError: vi.fn(),
	showSuccess: vi.fn()
}));

vi.mock('$lib/user', () => ({
	updateUser: vi.fn()
}));

describe('AccountEdit Component', () => {
	const baseUser = {
		id: 5,
		name: 'John Doe',
		country: 'ITA'
	};

	beforeEach(() => {
		vi.clearAllMocks();
	});

	it('renders account edit dialog', () => {
		render(AccountEdit, {
			props: {
				open: true,
				user: baseUser
			}
		});

		expect(screen.getByLabelText(/display name/i)).toBeInTheDocument();
		expect(screen.getByLabelText(/nationality/i)).toBeInTheDocument();
		expect(screen.getByRole('button', { name: /^save$/i })).toBeInTheDocument();
	});

	it('validates empty form submission', async () => {
		const user = userEvent.setup();

		render(AccountEdit, {
			props: {
				open: true,
				user: { id: 5, name: '', country: '' }
			}
		});

		await flush();

		await user.click(screen.getByRole('button', { name: /^save$/i }));

		expect(showError).toHaveBeenCalledWith(null, 'Please fill at least one field.');
		expect(updateUser).not.toHaveBeenCalled();
	});

	it('trims whitespace from input fields', async () => {
		const user = userEvent.setup();
		vi.mocked(updateUser).mockResolvedValueOnce({ ok: true } as any);

		render(AccountEdit, {
			props: {
				open: true,
				user: baseUser
			}
		});

		await flush();

		const nameInput = screen.getByLabelText(/display name/i) as HTMLInputElement;
		await fireEvent.input(nameInput, { target: { value: ' Bob ' } });
		await flush();

		await user.click(screen.getByRole('button', { name: /^save$/i }));

		await waitFor(() => {
			expect(updateUser).toHaveBeenCalledWith(5, 'Bob', 'ITA');
		});
	});

	it('updates user profile successfully', async () => {
		const user = userEvent.setup();
		vi.mocked(updateUser).mockResolvedValueOnce({ ok: true } as any);

		render(AccountEdit, {
			props: {
				open: true,
				user: baseUser
			}
		});
		await flush();

		const nameInput = screen.getByLabelText(/display name/i) as HTMLInputElement;
		await fireEvent.input(nameInput, { target: { value: 'Bob' } });
		await flush();

		await user.click(screen.getByRole('button', { name: /^save$/i }));

		await waitFor(() => {
			expect(updateUser).toHaveBeenCalledWith(5, 'Bob', 'ITA');
		});

		await waitFor(() => {
			expect(showSuccess).toHaveBeenCalledWith('Profile updated.');
		});
	});

	it('handles update error', async () => {
		const user = userEvent.setup();
		const error = new Error('Update failed');
		vi.mocked(updateUser).mockRejectedValueOnce(error);

		render(AccountEdit, {
			props: {
				open: true,
				user: baseUser
			}
		});

		await flush();

		const nameInput = screen.getByLabelText(/display name/i) as HTMLInputElement;
		await fireEvent.input(nameInput, { target: { value: 'Bob' } });
		await flush();

		await user.click(screen.getByRole('button', { name: /^save$/i }));

		await waitFor(() => {
			expect(showError).toHaveBeenCalledWith(error, 'Failed to update profile.');
		});
	});

	it('shows loading state during update', async () => {
		const user = userEvent.setup();
		vi.mocked(updateUser).mockImplementation(
			() => new Promise((resolve) => setTimeout(() => resolve({ ok: true } as any), 200))
		);

		render(AccountEdit, {
			props: {
				open: true,
				user: baseUser
			}
		});

		await flush();

		const nameInput = screen.getByLabelText(/display name/i) as HTMLInputElement;
		await fireEvent.input(nameInput, { target: { value: 'Bob' } });
		await flush();

		await user.click(screen.getByRole('button', { name: /^save$/i }));

		// The button text changes to "Saving..."
		const savingButton = await screen.findByText(/saving\.\.\./i);
		expect(savingButton).toBeInTheDocument();
		// Button should be disabled
		const button = screen.getByRole('button', { name: /saving\.\.\./i });
		expect(button).toBeDisabled();
	});
});

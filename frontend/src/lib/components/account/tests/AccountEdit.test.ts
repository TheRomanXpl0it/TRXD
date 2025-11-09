import { render, screen, waitFor } from '@testing-library/svelte';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { toast } from 'svelte-sonner';
import AccountEdit from '../AccountEdit.svelte';
import { updateUser } from '$lib/user';

vi.mock('svelte-sonner', () => ({
	toast: {
		success: vi.fn(),
		error: vi.fn()
	}
}));

vi.mock('$lib/user', () => ({
	updateUser: vi.fn()
}));

describe('AccountEdit Component', () => {
	const baseUser = {
		id: 5,
		name: '',
		image: '',
		country: ''
	} as any;

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
		expect(screen.getByLabelText(/image url/i)).toBeInTheDocument();
		expect(screen.getByRole('button', { name: /^save$/i })).toBeInTheDocument();
	});

	it('validates empty form submission', async () => {
		const user = userEvent.setup();

		render(AccountEdit, {
			props: {
				open: true,
				user: baseUser
			}
		});

		await user.click(screen.getByRole('button', { name: /^save$/i }));

		expect(toast.error).toHaveBeenCalledWith('Please fill at least one field.');
		expect(updateUser).not.toHaveBeenCalled();
	});

	it('validates image URL format', async () => {
		const user = userEvent.setup();

		render(AccountEdit, {
			props: {
				open: true,
				user: baseUser
			}
		});

		await user.type(screen.getByLabelText(/image url/i), 'bad');
		await user.click(screen.getByRole('button', { name: /^save$/i }));

		expect(toast.error).toHaveBeenCalledWith('Image must be a valid URL.');
		expect(updateUser).not.toHaveBeenCalled();
	});

	it('trims whitespace from input fields', async () => {
		const user = userEvent.setup();
		vi.mocked(updateUser).mockResolvedValueOnce({ ok: true });

		render(AccountEdit, {
			props: {
				open: true,
				user: baseUser
			}
		});

		await user.type(screen.getByLabelText(/display name/i), ' Bob ');
		await user.type(screen.getByLabelText(/image url/i), ' http://img.png ');
		await user.click(screen.getByRole('button', { name: /^save$/i }));

		await waitFor(() => {
			expect(updateUser).toHaveBeenCalledWith(5, 'Bob', '', 'http://img.png');
		});
	});

  it('updates user profile successfully', async () => {
    const user = userEvent.setup();
    vi.mocked(updateUser).mockResolvedValueOnce({ ok: true });

    render(AccountEdit, {
      props: {
        open: true,
        user: baseUser
      }
    });

    await user.type(screen.getByLabelText(/display name/i), 'Bob');
    await user.type(screen.getByLabelText(/image url/i), 'http://img.png');
    await user.click(screen.getByRole('button', { name: /^save$/i }));

    await waitFor(() => {
      expect(updateUser).toHaveBeenCalledWith(5, 'Bob', '', 'http://img.png');
    });

    expect(toast.success).toHaveBeenCalledWith('Profile updated.');
  });

	it('handles update error', async () => {
		const user = userEvent.setup();
		vi.mocked(updateUser).mockRejectedValueOnce(new Error('Update failed'));

		render(AccountEdit, {
			props: {
				open: true,
				user: baseUser
			}
		});

		await user.type(screen.getByLabelText(/display name/i), 'Bob');
		await user.click(screen.getByRole('button', { name: /^save$/i }));

		await waitFor(() => {
			expect(toast.error).toHaveBeenCalledWith('Update failed');
		});
	});

	it('shows loading state during update', async () => {
		const user = userEvent.setup();
		vi.mocked(updateUser).mockImplementation(
			() => new Promise((resolve) => setTimeout(() => resolve({ ok: true }), 100))
		);

		render(AccountEdit, {
			props: {
				open: true,
				user: baseUser
			}
		});

		await user.type(screen.getByLabelText(/display name/i), 'Bob');
		await user.click(screen.getByRole('button', { name: /^save$/i }));

		const saveButton = screen.getByRole('button', { name: /saving/i });
		expect(saveButton).toBeDisabled();
	});

	it('allows updating individual fields', async () => {
		const user = userEvent.setup();
		vi.mocked(updateUser).mockResolvedValueOnce({ ok: true });

		render(AccountEdit, {
			props: {
				open: true,
				user: baseUser
			}
		});

		await user.type(screen.getByLabelText(/image url/i), 'http://newimage.png');
		await user.click(screen.getByRole('button', { name: /^save$/i }));

		await waitFor(() => {
			expect(updateUser).toHaveBeenCalledWith(5, '', '', 'http://newimage.png');
		});
	});

	it('validates URL must start with http or https', async () => {
		const user = userEvent.setup();

		render(AccountEdit, {
			props: {
				open: true,
				user: baseUser
			}
		});

		await user.type(screen.getByLabelText(/image url/i), 'ftp://image.png');
		await user.click(screen.getByRole('button', { name: /^save$/i }));

		expect(toast.error).toHaveBeenCalledWith('Image must be a valid URL.');
		expect(updateUser).not.toHaveBeenCalled();
	});
});


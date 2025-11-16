import { render, screen, waitFor, fireEvent } from '@testing-library/svelte';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { toast } from 'svelte-sonner';
import AccountEdit from '../AccountEdit.svelte';
import { updateUser } from '$lib/user';
import { tick } from 'svelte';

async function flush() {
  await tick();
  await Promise.resolve();
}

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

	afterEach(async () => {
		await new Promise(resolve => setTimeout(resolve, 150));
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

		await flush();

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

		await flush();

    const imgInput = screen.getByLabelText(/image url/i) as HTMLInputElement;
    await fireEvent.input(imgInput, { target: { value: 'bad' } });
    await flush();
    await user.click(screen.getByRole('button', { name: /^save$/i }));

    await waitFor(() => {
      expect(toast.error).toHaveBeenCalledWith('Image must be a valid URL.');
    });
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

    const nameInput = screen.getByLabelText(/display name/i) as HTMLInputElement;
    const imageInput = screen.getByLabelText(/image url/i) as HTMLInputElement;
    await fireEvent.input(nameInput, { target: { value: ' Bob ' } });
    await fireEvent.input(imageInput, { target: { value: ' http://img.png ' } });
    await flush();
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
    await flush();

    const name2 = screen.getByLabelText(/display name/i) as HTMLInputElement;
    const image2 = screen.getByLabelText(/image url/i) as HTMLInputElement;
    await fireEvent.input(name2, { target: { value: 'Bob' } });
    await fireEvent.input(image2, { target: { value: 'http://img.png' } });
    await flush();
    await user.click(screen.getByRole('button', { name: /^save$/i }));

    await waitFor(() => {
      expect(updateUser).toHaveBeenCalledWith(5, 'Bob', '', 'http://img.png');
    });

    await waitFor(() => {
      expect(toast.success).toHaveBeenCalledWith('Profile updated.');
    });
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

		await flush();

		const nameErr = screen.getByLabelText(/display name/i) as HTMLInputElement;
		await fireEvent.input(nameErr, { target: { value: 'Bob' } });
		await user.click(screen.getByRole('button', { name: /^save$/i }));

		await waitFor(() => {
			expect(toast.error).toHaveBeenCalledWith('Update failed');
		});
	});

	it('shows loading state during update', async () => {
		const user = userEvent.setup();
    vi.mocked(updateUser).mockImplementation(
        () => new Promise((resolve) => setTimeout(() => resolve({ ok: true }), 200))
    );

		render(AccountEdit, {
			props: {
				open: true,
				user: baseUser
			}
		});

		await flush();

    const name3 = screen.getByLabelText(/display name/i) as HTMLInputElement;
    await fireEvent.input(name3, { target: { value: 'Bob' } });
    await flush();
    await user.click(screen.getByRole('button', { name: /^save$/i }));

    const saveButton = await screen.findByRole('button', { name: /saving/i });
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

		await flush();

    const image3 = screen.getByLabelText(/image url/i) as HTMLInputElement;
    await fireEvent.input(image3, { target: { value: 'http://newimage.png' } });
    await flush();
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

    const image4 = screen.getByLabelText(/image url/i) as HTMLInputElement;
    await fireEvent.input(image4, { target: { value: 'ftp://image.png' } });
    await flush();
    await user.click(screen.getByRole('button', { name: /^save$/i }));

    await waitFor(() => {
      expect(toast.error).toHaveBeenCalledWith('Image must be a valid URL.');
    });
		expect(updateUser).not.toHaveBeenCalled();
	});
});

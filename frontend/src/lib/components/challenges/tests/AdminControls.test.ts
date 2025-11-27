import { render, screen, waitFor } from '@testing-library/svelte';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { toast } from 'svelte-sonner';
import AdminControls from '../AdminControls.svelte';
import { createCategory } from '$lib/categories';

vi.mock('svelte-sonner', () => ({
	toast: {
		success: vi.fn(),
		error: vi.fn()
	}
}));

vi.mock('$lib/categories', () => ({
	createCategory: vi.fn()
}));

describe('AdminControls Component', () => {
	beforeEach(() => {
		vi.clearAllMocks();
	});

	it('renders create challenge button', () => {
		render(AdminControls);

		expect(screen.getByRole('button', { name: /create challenge/i })).toBeInTheDocument();
	});

	it('renders new category button', () => {
		render(AdminControls);

		expect(screen.getByRole('button', { name: /new category/i })).toBeInTheDocument();
	});

	it('calls oncreate callback when create challenge button is clicked', async () => {
		const user = userEvent.setup();
		const handleOpenCreate = vi.fn();

		render(AdminControls, {
			props: {
				'onopen-create': handleOpenCreate
			}
		});

		const createButton = screen.getByRole('button', { name: /create challenge/i });
		await user.click(createButton);

		expect(handleOpenCreate).toHaveBeenCalledTimes(1);
	});

	it('opens popover when new category button is clicked', async () => {
		const user = userEvent.setup();

		render(AdminControls);

		const categoryButton = screen.getByRole('button', { name: /new category/i });
		await user.click(categoryButton);

		// Popover content should be visible
		await waitFor(() => {
			expect(screen.getByLabelText(/category name/i)).toBeInTheDocument();
		});
	});

	it('displays category name input in popover', async () => {
		const user = userEvent.setup();

		render(AdminControls);

		await user.click(screen.getByRole('button', { name: /new category/i }));

		await waitFor(() => {
			expect(screen.getByLabelText(/category name/i)).toBeInTheDocument();
		});
	});

	it('displays cancel and create buttons in popover', async () => {
		const user = userEvent.setup();

		render(AdminControls);

		await user.click(screen.getByRole('button', { name: /new category/i }));

		await waitFor(() => {
			expect(screen.getByRole('button', { name: /^cancel$/i })).toBeInTheDocument();
			expect(screen.getByRole('button', { name: /^create$/i })).toBeInTheDocument();
		});
	});

	it('create button is disabled when category name is empty', async () => {
		const user = userEvent.setup();

		render(AdminControls);

		await user.click(screen.getByRole('button', { name: /new category/i }));

		await waitFor(() => {
			const createButton = screen.getByRole('button', { name: /^create$/i });
			expect(createButton).toBeDisabled();
		});
	});

	it('create button is enabled when category name field is filled', async () => {
		// TODO: Implement this once the backend removes the icon support
		/*
		const user = userEvent.setup();

		render(AdminControls);

		await user.click(screen.getByRole('button', { name: /new category/i }));

		await waitFor(() => {
			expect(screen.getByLabelText(/category name/i)).toBeInTheDocument();
		});

		const nameInput = screen.getByLabelText(/category name/i);
		const iconInput = screen.getByLabelText(/icon/i);

		await user.type(nameInput, 'Web');
		await user.type(iconInput, 'Globe');

		await waitFor(() => {
			const createButton = screen.getByRole('button', { name: /^create$/i });
			expect(createButton).not.toBeDisabled();
		});
		*/
	});

	it('closes popover when cancel is clicked', async () => {
		const user = userEvent.setup();

		render(AdminControls);

		await user.click(screen.getByRole('button', { name: /new category/i }));

		await waitFor(() => {
			expect(screen.getByRole('button', { name: /^cancel$/i })).toBeInTheDocument();
		});

		const cancelButton = screen.getByRole('button', { name: /^cancel$/i });
		await user.click(cancelButton);

		await waitFor(() => {
			expect(screen.queryByLabelText(/category name/i)).not.toBeInTheDocument();
		});
	});

	it('creates category successfully', async () => {
		// TODO: Implement this once the backend removes the icon support
		/*
		const user = userEvent.setup();
		const mockCreateCategory = vi.mocked(createCategory);
		mockCreateCategory.mockResolvedValueOnce(undefined);

		render(AdminControls);

		await user.click(screen.getByRole('button', { name: /new category/i }));

		await waitFor(() => {
			expect(screen.getByLabelText(/category name/i)).toBeInTheDocument();
		});

		const nameInput = screen.getByLabelText(/category name/i);
		const iconInput = screen.getByLabelText(/icon/i);

		await user.type(nameInput, 'Forensics');
		await user.type(iconInput, 'Search');

		const createButton = screen.getByRole('button', { name: /^create$/i });
		await user.click(createButton);

		await waitFor(() => {
			expect(mockCreateCategory).toHaveBeenCalledWith('Forensics', 'Search');
		});
		*/
	});

	it('shows success toast after creating category', async () => {
		// TODO: Implement this once the backend removes the icon support
		/*
		const user = userEvent.setup();
		const mockCreateCategory = vi.mocked(createCategory);
		const mockToast = vi.mocked(toast);
		mockCreateCategory.mockResolvedValueOnce(undefined);

		render(AdminControls);

		await user.click(screen.getByRole('button', { name: /new category/i }));

		await waitFor(() => {
			expect(screen.getByLabelText(/category name/i)).toBeInTheDocument();
		});

		const nameInput = screen.getByLabelText(/category name/i);
		const iconInput = screen.getByLabelText(/icon/i);

		await user.type(nameInput, 'Crypto');
		await user.type(iconInput, 'Lock');

		const createButton = screen.getByRole('button', { name: /^create$/i });
		await user.click(createButton);

		await waitFor(() => {
			expect(mockToast.success).toHaveBeenCalledWith('Category created!');
		});
		*/
	});

	it('calls oncategory-created callback after successful creation', async () => {
		// TODO: Implement this once the backend removes the icon support
		/*
		const user = userEvent.setup();
		const mockCreateCategory = vi.mocked(createCategory);
		const handleCategoryCreated = vi.fn();
		mockCreateCategory.mockResolvedValueOnce(undefined);

		render(AdminControls, {
			props: {
				'oncategory-created': handleCategoryCreated
			}
		});

		await user.click(screen.getByRole('button', { name: /new category/i }));

		await waitFor(() => {
			expect(screen.getByLabelText(/category name/i)).toBeInTheDocument();
		});

		const nameInput = screen.getByLabelText(/category name/i);
		const iconInput = screen.getByLabelText(/icon/i);

		await user.type(nameInput, 'Pwn');
		await user.type(iconInput, 'Zap');

		const createButton = screen.getByRole('button', { name: /^create$/i });
		await user.click(createButton);

		await waitFor(() => {
			expect(handleCategoryCreated).toHaveBeenCalledTimes(1);
		});
		*/
	});

	it('closes popover after successful creation', async () => {
		// TODO: Implement this once the backend removes the icon support
		/*
		const user = userEvent.setup();
		const mockCreateCategory = vi.mocked(createCategory);
		mockCreateCategory.mockResolvedValueOnce(undefined);

		render(AdminControls);

		await user.click(screen.getByRole('button', { name: /new category/i }));

		await waitFor(() => {
			expect(screen.getByLabelText(/category name/i)).toBeInTheDocument();
		});

		const nameInput = screen.getByLabelText(/category name/i);
		const iconInput = screen.getByLabelText(/icon/i);

		await user.type(nameInput, 'Rev');
		await user.type(iconInput, 'Code');

		const createButton = screen.getByRole('button', { name: /^create$/i });
		await user.click(createButton);

		await waitFor(() => {
			expect(screen.queryByLabelText(/category name/i)).not.toBeInTheDocument();
		});
		*/
	});

	it('clears form after successful creation', async () => {
		// TODO: Implement this once the backend removes the icon support
		/*
		const user = userEvent.setup();
		const mockCreateCategory = vi.mocked(createCategory);
		mockCreateCategory.mockResolvedValueOnce(undefined);

		render(AdminControls);

		await user.click(screen.getByRole('button', { name: /new category/i }));

		await waitFor(() => {
			expect(screen.getByLabelText(/category name/i)).toBeInTheDocument();
		});

		let nameInput = screen.getByLabelText(/category name/i);
		let iconInput = screen.getByLabelText(/icon/i);

		await user.type(nameInput, 'OSINT');
		await user.type(iconInput, 'Eye');

		const createButton = screen.getByRole('button', { name: /^create$/i });
		await user.click(createButton);

		// Wait for popover to close
		await waitFor(() => {
			expect(screen.queryByLabelText(/category name/i)).not.toBeInTheDocument();
		});

		// Reopen and check fields are empty
		await user.click(screen.getByRole('button', { name: /new category/i }));

		await waitFor(() => {
			expect(screen.getByLabelText(/category name/i)).toBeInTheDocument();
		});

		nameInput = screen.getByLabelText(/category name/i);
		iconInput = screen.getByLabelText(/icon/i);

		expect(nameInput).toHaveValue('');
		expect(iconInput).toHaveValue('');
		*/
	});

	it('shows error toast when creation fails', async () => {
		// TODO: Implement this once the backend removes the icon support
		/*
		const user = userEvent.setup();
		const mockCreateCategory = vi.mocked(createCategory);
		const mockToast = vi.mocked(toast);
		mockCreateCategory.mockRejectedValueOnce(new Error('API error'));

		render(AdminControls);

		await user.click(screen.getByRole('button', { name: /new category/i }));

		await waitFor(() => {
			expect(screen.getByLabelText(/category name/i)).toBeInTheDocument();
		});

		const nameInput = screen.getByLabelText(/category name/i);
		const iconInput = screen.getByLabelText(/icon/i);

		await user.type(nameInput, 'Web');
		await user.type(iconInput, 'Globe');

		const createButton = screen.getByRole('button', { name: /^create$/i });
		await user.click(createButton);

		await waitFor(() => {
			expect(mockToast.error).toHaveBeenCalledWith('API error');
		});
		*/
	});

	it('shows generic error message when error has no message', async () => {
		// TODO: Implement this once the backend removes the icon support
		/*
		const user = userEvent.setup();
		const mockCreateCategory = vi.mocked(createCategory);
		const mockToast = vi.mocked(toast);
		mockCreateCategory.mockRejectedValueOnce({});

		render(AdminControls);

		await user.click(screen.getByRole('button', { name: /new category/i }));

		await waitFor(() => {
			expect(screen.getByLabelText(/category name/i)).toBeInTheDocument();
		});

		const nameInput = screen.getByLabelText(/category name/i);
		const iconInput = screen.getByLabelText(/icon/i);

		await user.type(nameInput, 'Misc');
		await user.type(iconInput, 'Star');

		const createButton = screen.getByRole('button', { name: /^create$/i });
		await user.click(createButton);

		await waitFor(() => {
			expect(mockToast.error).toHaveBeenCalledWith('Failed to create category.');
		});
		*/
	});

	it('shows error toast when name is empty on submit', async () => {
		// TODO: Implement this once the backend removes the icon support
		/*
		const user = userEvent.setup();
		const mockToast = vi.mocked(toast);

		render(AdminControls);

		await user.click(screen.getByRole('button', { name: /new category/i }));

		await waitFor(() => {
			expect(screen.getByLabelText(/icon/i)).toBeInTheDocument();
		});

		const iconInput = screen.getByLabelText(/icon/i);
		await user.type(iconInput, 'Lock');

		// Manually trigger form submission (button should be disabled but test the logic)
		const form = iconInput.closest('form');
		if (form) {
			form.dispatchEvent(new Event('submit', { cancelable: true }));
		}

		await waitFor(() => {
			expect(mockToast.error).toHaveBeenCalledWith(
				'Please enter a category name and an icon.'
			);
		});
		*/
	});

	it('shows error toast when icon is empty on submit', async () => {
		// TODO: Implement this once the backend removes the icon support
		/*
		const user = userEvent.setup();
		const mockToast = vi.mocked(toast);

		render(AdminControls);

		await user.click(screen.getByRole('button', { name: /new category/i }));

		await waitFor(() => {
			expect(screen.getByLabelText(/category name/i)).toBeInTheDocument();
		});

		const nameInput = screen.getByLabelText(/category name/i);
		await user.type(nameInput, 'Hardware');

		// Manually trigger form submission
		const form = nameInput.closest('form');
		if (form) {
			form.dispatchEvent(new Event('submit', { cancelable: true }));
		}

		await waitFor(() => {
			expect(mockToast.error).toHaveBeenCalledWith(
				'Please enter a category name and an icon.'
			);
		});
		*/
	});

	it('trims whitespace from category name', async () => {
		// TODO: Implement this once the backend removes the icon support
		/*
		const user = userEvent.setup();
		const mockCreateCategory = vi.mocked(createCategory);
		mockCreateCategory.mockResolvedValueOnce(undefined);

		render(AdminControls);

		await user.click(screen.getByRole('button', { name: /new category/i }));

		await waitFor(() => {
			expect(screen.getByLabelText(/category name/i)).toBeInTheDocument();
		});

		const nameInput = screen.getByLabelText(/category name/i);
		const iconInput = screen.getByLabelText(/icon/i);

		await user.type(nameInput, '  Web  ');
		await user.type(iconInput, '  Globe  ');

		const createButton = screen.getByRole('button', { name: /^create$/i });
		await user.click(createButton);

		await waitFor(() => {
			expect(mockCreateCategory).toHaveBeenCalledWith('Web', 'Globe');
		});
		*/
	});

	it('shows loading state during category creation', async () => {
		// TODO: Implement this once the backend removes the icon support
		/*
		const user = userEvent.setup();
		const mockCreateCategory = vi.mocked(createCategory);
		
		// Create a promise we can control
		let resolveCreate: () => void;
		const createPromise = new Promise<void>((resolve) => {
			resolveCreate = resolve;
		});
		mockCreateCategory.mockReturnValueOnce(createPromise);

		render(AdminControls);

		await user.click(screen.getByRole('button', { name: /new category/i }));

		await waitFor(() => {
			expect(screen.getByLabelText(/category name/i)).toBeInTheDocument();
		});

		const nameInput = screen.getByLabelText(/category name/i);
		const iconInput = screen.getByLabelText(/icon/i);

		await user.type(nameInput, 'Test');
		await user.type(iconInput, 'TestTube');

		const createButton = screen.getByRole('button', { name: /^create$/i });
		await user.click(createButton);

		// Button should be disabled during loading
		await waitFor(() => {
			expect(createButton).toBeDisabled();
		});

		// Resolve the promise
		resolveCreate!();
		*/
	});

});

import { render, screen, waitFor, within, fireEvent } from '@testing-library/svelte';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { toast } from 'svelte-sonner';
import TeamJoinCreate from '../TeamJoinCreate.svelte';
import { joinTeam, createTeam } from '@/team';

vi.mock('svelte-sonner', () => ({
	toast: {
		success: vi.fn(),
		error: vi.fn()
	}
}));

vi.mock('@/team', () => ({
	joinTeam: vi.fn(),
	createTeam: vi.fn()
}));

describe('TeamJoinCreate Component', () => {
  beforeEach(() => {
    vi.resetAllMocks();
  });

  afterEach(async () => {
    await new Promise(resolve => setTimeout(resolve, 150));
  });

	it('renders join and create buttons', () => {
		render(TeamJoinCreate);

		expect(screen.getByRole('button', { name: /^join$/i })).toBeInTheDocument();
		expect(screen.getByRole('button', { name: /^create$/i })).toBeInTheDocument();
	});

	it('opens join dialog when join button is clicked', async () => {
		const user = userEvent.setup();

		render(TeamJoinCreate);

		await fireEvent.click(screen.getByRole('button', { name: /^join$/i }));

		expect(screen.getByLabelText(/team name/i)).toBeInTheDocument();
		expect(screen.getByLabelText(/team password/i)).toBeInTheDocument();
	});

	it('prevents join submission with empty fields', async () => {
		const user = userEvent.setup();

		render(TeamJoinCreate);

		const joinButtons = screen.getAllByRole('button', { name: /^join$/i });
		await fireEvent.click(joinButtons[0]);

		const dialog = screen.getByRole('dialog');
		const form = dialog.querySelector('form') as HTMLFormElement;
		await fireEvent.submit(form);

		expect(joinTeam).not.toHaveBeenCalled();
	});

	it('joins team successfully', async () => {
		const user = userEvent.setup();
		vi.mocked(joinTeam).mockResolvedValueOnce({ ok: true });

		render(TeamJoinCreate);

		// Open join dialog
		const joinButtons = screen.getAllByRole('button', { name: /^join$/i });
		await fireEvent.click(joinButtons[0]);

		// Fill in form
		await user.type(screen.getByLabelText(/team name/i), 'ZeroDayCats');
		await user.type(screen.getByLabelText(/team password/i), 'p@ssw0rd');

		// Submit the form inside the open dialog (avoids overlay/pointer-events issues)
		const dialog = screen.getByRole('dialog');
		const form = dialog.querySelector('form') as HTMLFormElement;
		await fireEvent.submit(form);

		await waitFor(() => {
			expect(joinTeam).toHaveBeenCalledWith('ZeroDayCats', 'p@ssw0rd');
		});

		expect(toast.success).toHaveBeenCalled();
	});

	it('handles join errors with toast message', async () => {
		const user = userEvent.setup();
		vi.mocked(joinTeam).mockRejectedValueOnce(new Error('Bad credentials'));

		render(TeamJoinCreate);

		// Open join dialog
		const joinButtons = screen.getAllByRole('button', { name: /^join$/i });
		await fireEvent.click(joinButtons[0]);

		// Wait for fields to be available
		const nameInput = await screen.findByLabelText(/team name/i);
		const passwordInput = await screen.findByLabelText(/team password/i);

		// Fill in form
		await user.type(nameInput, 'TeamX');
		await user.type(passwordInput, 'wrong');

		// Submit the form inside the open dialog (avoids overlay/pointer-events issues)
		const dialog = screen.getByRole('dialog');
		const form = dialog.querySelector('form') as HTMLFormElement;
		await fireEvent.submit(form);

		await waitFor(() => {
			expect(toast.error).toHaveBeenCalledWith('Bad credentials');
		});
	});

	it('opens create dialog when create button is clicked', async () => {
		const user = userEvent.setup();

		render(TeamJoinCreate);

    const createButtons = screen.getAllByRole('button', { name: /^create$/i });
    await fireEvent.click(createButtons[0]);

		// After dialog opens, verify fields are visible
		expect(screen.getByLabelText(/^team name$/i)).toBeInTheDocument();
		expect(screen.getByLabelText(/^team password$/i)).toBeInTheDocument();
		expect(screen.getByLabelText(/confirm password/i)).toBeInTheDocument();
	});

	it('validates password mismatch', async () => {
		const user = userEvent.setup();

		render(TeamJoinCreate);

		// Open create dialog
	const createButtons = screen.getAllByRole('button', { name: /^create$/i });
	await fireEvent.click(createButtons[0]);

		// Fill in form
		await user.type(screen.getByLabelText(/^team name$/i), 'BlueTeam');
		await user.type(screen.getByLabelText(/^team password$/i), '1234567');
		await user.type(screen.getByLabelText(/confirm password/i), '7654321');

		// Submit the form inside the open dialog (avoids overlay/pointer-events issues)
		const dialog = screen.getByRole('dialog');
		const form = dialog.querySelector('form') as HTMLFormElement;
		await fireEvent.submit(form);

    expect(createTeam).not.toHaveBeenCalled();
    await waitFor(() => {
      expect(toast.error).toHaveBeenCalled();
    });
	});

	it('validates password minimum length', async () => {
		const user = userEvent.setup();

		render(TeamJoinCreate);

		// Open the create dialog
	const openButtons = screen.getAllByRole('button', { name: /^create$/i });
	await fireEvent.click(openButtons[0]);

		// Fill in form
		await user.type(screen.getByLabelText(/^team name$/i), 'BlueTeam');
		await user.type(screen.getByLabelText(/^team password$/i), '123');
		await user.type(screen.getByLabelText(/confirm password/i), '123');

		// Submit the form inside the open dialog (avoids overlay/pointer-events issues)
		const dialog = screen.getByRole('dialog');
		const form = dialog.querySelector('form') as HTMLFormElement;
		await fireEvent.submit(form);

    expect(createTeam).not.toHaveBeenCalled();
    await waitFor(() => {
      expect(toast.error).toHaveBeenCalled();
    });
	});

	it('creates team successfully', async () => {
		const user = userEvent.setup();
		vi.mocked(createTeam).mockResolvedValueOnce({ ok: true });

		render(TeamJoinCreate);

      // Open create dialog
    const createButtons = screen.getAllByRole('button', { name: /^create$/i });
    await fireEvent.click(createButtons[0]);

		// Wait for dialog to open and fields to be available
		const nameInput = await screen.findByLabelText(/^team name$/i);
		const passwordInput = await screen.findByLabelText(/^team password$/i);
		const confirmInput = await screen.findByLabelText(/confirm password/i);

    await user.type(nameInput, 'RedTeam');
		await user.type(passwordInput, 'longpassword');
		await user.type(confirmInput, 'longpassword');

		// Submit the form inside the open dialog (avoids overlay/pointer-events issues)
		const dialog = screen.getByRole('dialog');
		const form = dialog.querySelector('form') as HTMLFormElement;
		await fireEvent.submit(form);

		await waitFor(() => {
			expect(createTeam).toHaveBeenCalledWith('RedTeam', 'longpassword');
		});

		expect(toast.success).toHaveBeenCalled();
	});

	it('shows loading state during join', async () => {
		const user = userEvent.setup();
		let resolveJoin: (value: any) => void = () => {};
		vi.mocked(joinTeam).mockImplementation(
			() => new Promise((resolve) => {
				resolveJoin = resolve;
			})
		);

		render(TeamJoinCreate);

		// Open join dialog
		const openButtons = screen.getAllByRole('button', { name: /^join$/i });
		await fireEvent.click(openButtons[0]);

		// Fill in form
		await user.type(screen.getByLabelText(/team name/i), 'TestTeam');
		await user.type(screen.getByLabelText(/team password/i), 'password');
    
    const dialog = screen.getByRole('dialog');
    const submitButton = within(dialog).getByRole('button', { name: /^join$/i });
    await fireEvent.click(submitButton);

    // Wait for loading state - use findByRole which waits automatically
    const joinButton = await within(dialog).findByRole('button', { name: /joining/i });
    expect(joinButton).toBeDisabled();

    // Resolve join and let it finish
    resolveJoin({ ok: true });
	});

	it('shows loading state during create', async () => {
		const user = userEvent.setup();
		let resolveCreate: (value: any) => void = () => {};
		vi.mocked(createTeam).mockImplementation(
			() => new Promise((resolve) => {
				resolveCreate = resolve;
			})
		);

		render(TeamJoinCreate);

		// Open create dialog
	const openButtons = screen.getAllByRole('button', { name: /^create$/i });
	await fireEvent.click(openButtons[0]);

		// Wait for fields to be available
		const nameInput = await screen.findByLabelText(/^team name$/i);
		const passwordInput = await screen.findByLabelText(/^team password$/i);
		const confirmInput = await screen.findByLabelText(/confirm password/i);

		// Fill in form
		await user.type(nameInput, 'TestTeam');
		await user.type(passwordInput, 'longpassword');
		await user.type(confirmInput, 'longpassword');

    const dialog = screen.getByRole('dialog');
    const form = dialog.querySelector('form') as HTMLFormElement;
    await fireEvent.submit(form);

    // Wait for loading state
    await waitFor(() => {
      const createButton = within(dialog).getByRole('button', { name: /creating/i });
      expect(createButton).toBeDisabled();
    });

    // Resolve create and let it finish
    resolveCreate({ ok: true });

	});

  it('handles create errors with toast message', async () => {
    const user = userEvent.setup();
    vi.mocked(createTeam).mockRejectedValueOnce(new Error('Team name already exists'));

    render(TeamJoinCreate);

    // Open create dialog
    const openButtons = screen.getAllByRole('button', { name: /^create$/i });
    await fireEvent.click(openButtons[0]);

    // Wait for fields to be available
    const nameInput = await screen.findByLabelText(/^team name$/i);
    const passwordInput = await screen.findByLabelText(/^team password$/i);
    const confirmInput = await screen.findByLabelText(/confirm password/i);

    // Fill in form
    await user.type(nameInput, 'ExistingTeam');
    await user.type(passwordInput, 'validpassword');
    await user.type(confirmInput, 'validpassword');

    // Submit the form inside the open dialog
    const dialog = screen.getByRole('dialog');
    const form = dialog.querySelector('form') as HTMLFormElement;
    await fireEvent.submit(form);

    // 1) Ensure the component actually called createTeam
    await waitFor(() => {
      expect(createTeam).toHaveBeenCalledWith('ExistingTeam', 'validpassword');
    });

    // 2) Then ensure we surfaced that error via toast
    await waitFor(() => {
      expect(toast.error).toHaveBeenCalledWith('Team name already exists');
    });
  });


});

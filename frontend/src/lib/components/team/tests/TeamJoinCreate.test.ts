import { render, screen, waitFor, within, fireEvent } from '@testing-library/svelte';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { toast } from 'svelte-sonner';
import TeamJoinCreate from '../TeamJoinCreate.svelte';
import { joinTeam, createTeam } from '@/team';
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


		expect(screen.getByRole('button', { name: /^join team$/i })).toBeInTheDocument();
		expect(screen.getByRole('button', { name: /^create team$/i })).toBeInTheDocument();
	});

	it('opens join dialog when join button is clicked', async () => {
		const user = userEvent.setup();

		render(TeamJoinCreate);
		await flush();

		await fireEvent.click(screen.getByRole('button', { name: /^join team$/i }));

		expect(screen.getByLabelText(/team name/i)).toBeInTheDocument();
		expect(screen.getByLabelText(/team password/i)).toBeInTheDocument();
	});

	it('prevents join submission with empty fields', async () => {
		const user = userEvent.setup();

		render(TeamJoinCreate);
		await flush();

		const joinButtons = screen.getAllByRole('button', { name: /^join team$/i });
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
		await flush();

		// Open join dialog
		const joinButtons = screen.getAllByRole('button', { name: /^join team$/i });
		await fireEvent.click(joinButtons[0]);

		// Scope all queries to the open join dialog to avoid matching hidden inputs
		const dialog = await screen.findByRole('dialog');
		const nameInput = within(dialog).getByLabelText(/team name/i);
		const passInput = within(dialog).getByLabelText(/team password/i);
		await user.type(nameInput, 'ZeroDayCats');
		await user.type(passInput, 'p@ssw0rd');
		await flush();

		// Submit the form inside the open dialog (avoids overlay/pointer-events issues)
		const form = dialog.querySelector('form') as HTMLFormElement;
		await fireEvent.submit(form);

		await waitFor(() => {
			expect(joinTeam).toHaveBeenCalledWith('ZeroDayCats', 'p@ssw0rd');
		});

		await waitFor(() => {
			expect(toast.success).toHaveBeenCalled();
		});
	});

	it('handles join errors with toast message', async () => {
		const user = userEvent.setup();
		vi.mocked(joinTeam).mockRejectedValueOnce(new Error('Bad credentials'));

		render(TeamJoinCreate);
		await flush();

		// Open join dialog
		const joinButtons = screen.getAllByRole('button', { name: /^join team$/i });
		await fireEvent.click(joinButtons[0]);

		// Scope to the open dialog to avoid picking fields from other modals
		const dialog = await screen.findByRole('dialog');
		const nameInput = within(dialog).getByLabelText(/team name/i);
		const passwordInput = within(dialog).getByLabelText(/team password/i);

		// Fill in form
    await user.type(nameInput, 'TeamX');
    await user.type(passwordInput, 'wrong');
    await flush();

		// Submit the form inside the open dialog (avoids overlay/pointer-events issues)
		const form = dialog.querySelector('form') as HTMLFormElement;
		await fireEvent.submit(form);

		await waitFor(() => {
			expect(toast.error).toHaveBeenCalledWith('Bad credentials');
		});
	});

	it('opens create dialog when create button is clicked', async () => {
		const user = userEvent.setup();

		render(TeamJoinCreate);

    const createButtons = screen.getAllByRole('button', { name: /^create team$/i });
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
	const createButtons = screen.getAllByRole('button', { name: /^create team$/i });
	await fireEvent.click(createButtons[0]);

		// Fill in form
		await user.type(screen.getByLabelText(/^team name$/i), 'BlueTeam');
		await user.type(screen.getByLabelText(/^team password$/i), '1234567');
		await user.type(screen.getByLabelText(/confirm password/i), '7654321');
		await flush();

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
	const openButtons = screen.getAllByRole('button', { name: /^create team$/i });
	await fireEvent.click(openButtons[0]);

		await waitFor(() => {
			expect(screen.getByRole('dialog')).toBeInTheDocument();
		});

		// Fill in form using fireEvent instead of userEvent
		const nameInput = screen.getByLabelText(/^team name$/i);
		const passwordInput = screen.getByLabelText(/^team password$/i);
		const confirmInput = screen.getByLabelText(/confirm password/i);

		await fireEvent.input(nameInput, { target: { value: 'BlueTeam' } });
		await fireEvent.input(passwordInput, { target: { value: '123' } });
		await fireEvent.input(confirmInput, { target: { value: '123' } });

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
    const createButtons = screen.getAllByRole('button', { name: /^create team$/i });
    await fireEvent.click(createButtons[0]);

		// Wait for dialog to open and fields to be available
		const nameInput = await screen.findByLabelText(/^team name$/i);
		const passwordInput = await screen.findByLabelText(/^team password$/i);
		const confirmInput = await screen.findByLabelText(/confirm password/i);

    await fireEvent.input(nameInput, { target: { value: 'RedTeam' } });
    await fireEvent.input(passwordInput, { target: { value: 'longpassword' } });
    await fireEvent.input(confirmInput, { target: { value: 'longpassword' } });
    await flush();

		// Submit the form inside the open dialog (avoids overlay/pointer-events issues)
		const dialog = screen.getByRole('dialog');
		const form = dialog.querySelector('form') as HTMLFormElement;
		await fireEvent.submit(form);

		await waitFor(() => {
			expect(createTeam).toHaveBeenCalledWith('RedTeam', 'longpassword');
		});

		await waitFor(() => {
			expect(toast.success).toHaveBeenCalled();
		});
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
		const openButtons = screen.getAllByRole('button', { name: /^join team$/i });
		await fireEvent.click(openButtons[0]);

		// Scope inputs to the open dialog and fill the form
		const dialog = await screen.findByRole('dialog');
		const nameInput = within(dialog).getByLabelText(/team name/i);
		const passInput = within(dialog).getByLabelText(/team password/i);
		await user.type(nameInput, 'TestTeam');
		await user.type(passInput, 'password');

		// Submit the form directly to ensure onsubmit runs
		const form = dialog.querySelector('form') as HTMLFormElement;
		await fireEvent.submit(form);

		// Wait for loading state by checking the submit button becomes disabled
		const submitButton = within(dialog).getByRole('button', { name: /join/i });
		await waitFor(() => {
			expect(submitButton).toBeDisabled();
		});

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
	const openButtons = screen.getAllByRole('button', { name: /^create team$/i });
	await fireEvent.click(openButtons[0]);

		// Wait for fields to be available
		const nameInput = await screen.findByLabelText(/^team name$/i);
		const passwordInput = await screen.findByLabelText(/^team password$/i);
		const confirmInput = await screen.findByLabelText(/confirm password/i);

    // Fill in form (force sync)
    await fireEvent.input(nameInput, { target: { value: 'TestTeam' } });
    await fireEvent.input(passwordInput, { target: { value: 'longpassword' } });
    await fireEvent.input(confirmInput, { target: { value: 'longpassword' } });

    const dialog = screen.getByRole('dialog');
    const createBtn = within(dialog).getByRole('button', { name: /^create$/i });
    const form = dialog.querySelector('form') as HTMLFormElement;
    await fireEvent.submit(form);

    // Wait for loading state (submit button becomes disabled)
    await waitFor(() => {
      expect(createBtn).toBeDisabled();
    });

    // Resolve create and let it finish
    resolveCreate({ ok: true });

	});

  it('handles create errors with toast message', async () => {
    const user = userEvent.setup();
    vi.mocked(createTeam).mockRejectedValueOnce(new Error('Team name already exists'));

    render(TeamJoinCreate);

    // Open create dialog
    const openButtons = screen.getAllByRole('button', { name: /^create team$/i });
    await fireEvent.click(openButtons[0]);

    // Wait for fields to be available
    const nameInput = await screen.findByLabelText(/^team name$/i);
    const passwordInput = await screen.findByLabelText(/^team password$/i);
    const confirmInput = await screen.findByLabelText(/confirm password/i);

    // Fill in form
    await fireEvent.input(nameInput, { target: { value: 'ExistingTeam' } });
    await fireEvent.input(passwordInput, { target: { value: 'validpassword' } });
    await fireEvent.input(confirmInput, { target: { value: 'validpassword' } });
    await flush();

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

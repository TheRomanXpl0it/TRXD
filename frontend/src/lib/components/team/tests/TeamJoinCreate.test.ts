import { render, screen, waitFor, within, fireEvent } from '@testing-library/svelte';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi, beforeEach } from 'vitest';
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
		vi.clearAllMocks();
	});

	it('renders join and create buttons', () => {
		render(TeamJoinCreate);

		expect(screen.getByRole('button', { name: /^join team$/i })).toBeInTheDocument();
		expect(screen.getByRole('button', { name: /^create team$/i })).toBeInTheDocument();
	});

	it('opens join dialog and joins successfully', async () => {
		const user = userEvent.setup();
		vi.mocked(joinTeam).mockResolvedValueOnce({ ok: true } as any);

		render(TeamJoinCreate);
		await flush();

		await fireEvent.click(screen.getByRole('button', { name: /^join team$/i }));
		await flush();

		const dialog = await screen.findByRole('dialog');
		const nameInput = within(dialog).getByLabelText(/team name/i);
		const passInput = within(dialog).getByLabelText(/team password/i);

		await user.type(nameInput, 'ZeroDayCats');
		await user.type(passInput, 'p@ssw0rd');
		await flush();

		const submitBtn = within(dialog).getByRole('button', { name: /^join$/i });
		await user.click(submitBtn);

		await waitFor(() => {
			expect(joinTeam).toHaveBeenCalledWith('ZeroDayCats', 'p@ssw0rd');
			expect(toast.success).toHaveBeenCalledWith('Team Joined, welcome aboard!');
		});
	});

	it('handles join errors', async () => {
		const user = userEvent.setup();
		vi.mocked(joinTeam).mockRejectedValueOnce(new Error('Bad credentials'));

		render(TeamJoinCreate);
		await flush();

		await fireEvent.click(screen.getByRole('button', { name: /^join team$/i }));
		await flush();

		const dialog = await screen.findByRole('dialog');
		const nameInput = within(dialog).getByLabelText(/team name/i);
		const passInput = within(dialog).getByLabelText(/team password/i);

		await user.type(nameInput, 'TeamX');
		await user.type(passInput, 'wrong');
		await flush();

		const submitBtn = within(dialog).getByRole('button', { name: /^join$/i });
		await user.click(submitBtn);

		await waitFor(() => {
			expect(toast.error).toHaveBeenCalledWith('Bad credentials');
		});
	});

	it('opens create dialog and validates password mismatch', async () => {
		const user = userEvent.setup();

		render(TeamJoinCreate);
		await flush();

		await fireEvent.click(screen.getByRole('button', { name: /^create team$/i }));
		await flush();

		const dialog = await screen.findByRole('dialog');
		const nameInput = within(dialog).getByLabelText(/team name/i);
		const passInput = within(dialog).getByLabelText(/^team password$/i);
		const confirmInput = within(dialog).getByLabelText(/confirm password/i);

		await user.type(nameInput, 'BlueTeam');
		await user.type(passInput, 'password123');
		await user.type(confirmInput, 'mismatch');
		await flush();

		const submitBtn = within(dialog).getByRole('button', { name: /^create$/i });
		await user.click(submitBtn);

		expect(createTeam).not.toHaveBeenCalled();
		await waitFor(() => {
			expect(toast.error).toHaveBeenCalledWith('Passwords do not match.');
		});
	});

	it('creates team successfully', async () => {
		const user = userEvent.setup();
		vi.mocked(createTeam).mockResolvedValueOnce({ ok: true } as any);

		render(TeamJoinCreate);
		await flush();

		await fireEvent.click(screen.getByRole('button', { name: /^create team$/i }));
		await flush();

		const dialog = await screen.findByRole('dialog');
		const nameInput = within(dialog).getByLabelText(/team name/i);
		const passInput = within(dialog).getByLabelText(/^team password$/i);
		const confirmInput = within(dialog).getByLabelText(/confirm password/i);

		await user.type(nameInput, 'RedTeam');
		await user.type(passInput, 'longpassword');
		await user.type(confirmInput, 'longpassword');
		await flush();

		const submitBtn = within(dialog).getByRole('button', { name: /^create$/i });
		await user.click(submitBtn);

		await waitFor(() => {
			expect(createTeam).toHaveBeenCalledWith('RedTeam', 'longpassword');
			expect(toast.success).toHaveBeenCalledWith('Team Created!');
		});
	});

	it('shows loading states', async () => {
		const user = userEvent.setup();
		vi.mocked(joinTeam).mockImplementation(() => new Promise((resolve) => setTimeout(() => resolve({ ok: true } as any), 200)));

		render(TeamJoinCreate);
		await flush();

		await fireEvent.click(screen.getByRole('button', { name: /^join team$/i }));
		await flush();

		const dialog = await screen.findByRole('dialog');
		const nameInput = within(dialog).getByLabelText(/team name/i);
		const passInput = within(dialog).getByLabelText(/team password/i);

		await user.type(nameInput, 'TestTeam');
		await user.type(passInput, 'password');
		await flush();

		const submitBtn = within(dialog).getByRole('button', { name: /^join$/i });
		await user.click(submitBtn);

		const joiningBtn = await within(dialog).findByText(/joining\.\.\./i);
		expect(joiningBtn).toBeInTheDocument();
		expect(submitBtn).toBeDisabled();
	});
});

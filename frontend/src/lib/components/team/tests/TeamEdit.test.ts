import { render, screen, waitFor, fireEvent } from '@testing-library/svelte';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import TeamEdit from '../TeamEdit.svelte';
import { updateTeam } from '$lib/team';
import { useQueryClient } from '@tanstack/svelte-query';
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

vi.mock('$lib/team', () => ({
	updateTeam: vi.fn()
}));

vi.mock('@tanstack/svelte-query', () => ({
	useQueryClient: vi.fn(() => ({
		invalidateQueries: vi.fn()
	}))
}));

describe('TeamEdit Component', () => {
	const baseTeam = {
		id: 7,
		name: 'Team A',
		country: 'ITA',
		tags: ['tag1']
	};

	beforeEach(() => {
		vi.clearAllMocks();
	});

	it('renders team edit dialog', () => {
		render(TeamEdit, { props: { open: true, team: baseTeam } });

		expect(screen.getByLabelText(/team name/i)).toBeInTheDocument();
		expect(screen.getByRole('button', { name: /^save$/i })).toBeInTheDocument();
	});

	it('validates empty form submission', async () => {
		const user = userEvent.setup();

		render(TeamEdit, { props: { open: true, team: { id: 7, name: '', country: '' } } });
		await flush();

		await user.click(screen.getByRole('button', { name: /^save$/i }));

		expect(showError).toHaveBeenCalledWith(null, 'Please fill at least one field.');
		expect(updateTeam).not.toHaveBeenCalled();
	});

	it('trims whitespace from input fields', async () => {
		const user = userEvent.setup();
		vi.mocked(updateTeam).mockResolvedValueOnce({ ok: true } as any);

		render(TeamEdit, { props: { open: true, team: baseTeam } });
		await flush();

		const tName = screen.getByLabelText(/team name/i) as HTMLInputElement;
		await fireEvent.input(tName, { target: { value: ' New Team ' } });
		await flush();

		await user.click(screen.getByRole('button', { name: /^save$/i }));

		await waitFor(() => {
			expect(updateTeam).toHaveBeenCalledWith(7, 'New Team', 'ITA', ['tag1']);
		});
	});

	it('updates team successfully and invalidates cache', async () => {
		const user = userEvent.setup();
		const mockInvalidateQueries = vi.fn();
		vi.mocked(updateTeam).mockResolvedValueOnce({ ok: true } as any);
		vi.mocked(useQueryClient).mockReturnValue({ invalidateQueries: mockInvalidateQueries } as any);

		render(TeamEdit, { props: { open: true, team: baseTeam } });
		await flush();

		const tName2 = screen.getByLabelText(/team name/i) as HTMLInputElement;
		await fireEvent.input(tName2, { target: { value: 'New Team' } });
		await flush();

		await user.click(screen.getByRole('button', { name: /^save$/i }));

		await waitFor(() => {
			expect(updateTeam).toHaveBeenCalledWith(7, 'New Team', 'ITA', ['tag1']);
			expect(mockInvalidateQueries).toHaveBeenCalledWith({ queryKey: ['teams'] });
			expect(showSuccess).toHaveBeenCalledWith('Team updated.');
		});
	});

	it('handles update error', async () => {
		const user = userEvent.setup();
		const error = new Error('Update failed');
		vi.mocked(updateTeam).mockRejectedValueOnce(error);

		render(TeamEdit, { props: { open: true, team: baseTeam } });
		await flush();

		const nm = screen.getByLabelText(/team name/i) as HTMLInputElement;
		await fireEvent.input(nm, { target: { value: 'New Team' } });
		await flush();

		await user.click(screen.getByRole('button', { name: /^save$/i }));

		await waitFor(() => {
			expect(showError).toHaveBeenCalledWith(error, 'Failed to update team.');
		});
	});

	it('shows loading state during update', async () => {
		const user = userEvent.setup();
		vi.mocked(updateTeam).mockImplementation(
			() => new Promise((resolve) => setTimeout(() => resolve({ ok: true } as any), 200))
		);

		render(TeamEdit, { props: { open: true, team: baseTeam } });
		await flush();

		const nameInput = screen.getByLabelText(/team name/i) as HTMLInputElement;
		await fireEvent.input(nameInput, { target: { value: 'New Team' } });
		await flush();

		await user.click(screen.getByRole('button', { name: /^save$/i }));

		const savingButton = await screen.findByText(/saving\.\.\./i);
		expect(savingButton).toBeInTheDocument();
		expect(savingButton).toBeDisabled();
	});
});

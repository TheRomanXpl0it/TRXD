import { render, screen, waitFor, fireEvent } from '@testing-library/svelte';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { toast } from 'svelte-sonner';
import TeamEdit from '../TeamEdit.svelte';
import { updateTeam } from '$lib/team';
import { useQueryClient } from '@tanstack/svelte-query';
import { tick } from 'svelte';

async function flush() {
  await tick();
  await Promise.resolve();
}

vi.mock('svelte-sonner', () => ({
  toast: {
    success: vi.fn(),
    error: vi.fn(),
  },
}));

vi.mock('$lib/team', () => ({
  updateTeam: vi.fn(),
}));

vi.mock('@tanstack/svelte-query', () => ({
  useQueryClient: vi.fn(() => ({
    invalidateQueries: vi.fn(),
  })),
}));

describe('TeamEdit Component', () => {
  const baseTeam = {
    id: 7,
    name: '',
    bio: '',
    image: '',
    country: '',
  } as any;

  beforeEach(() => {
    vi.clearAllMocks();
  });

  afterEach(async () => {
    await new Promise(resolve => setTimeout(resolve, 150));
  });

  it('renders team edit dialog', () => {
		render(TeamEdit, { props: { open: true, team: baseTeam } });


    expect(screen.getByLabelText(/team name/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/^bio$/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/image url/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /^save$/i })).toBeInTheDocument();
  });

  it('validates empty form submission', async () => {
    const user = userEvent.setup();

		render(TeamEdit, { props: { open: true, team: baseTeam } });
		await flush();

    await user.click(screen.getByRole('button', { name: /^save$/i }));

    expect(toast.error).toHaveBeenCalledWith('Please fill at least one field.');
    expect(updateTeam).not.toHaveBeenCalled();
  });

  it('validates image URL format', async () => {
    const user = userEvent.setup();

		render(TeamEdit, { props: { open: true, team: baseTeam } });
		await flush();

    const imageBad = screen.getByLabelText(/image url/i) as HTMLInputElement;
    await fireEvent.input(imageBad, { target: { value: 'not-a-url' } });
    await user.click(screen.getByRole('button', { name: /^save$/i }));

    expect(toast.error).toHaveBeenCalledWith('Image must be a valid URL.');
    expect(updateTeam).not.toHaveBeenCalled();
  });

  it('trims whitespace from input fields', async () => {
    const user = userEvent.setup();
    vi.mocked(updateTeam).mockResolvedValueOnce({ ok: true } as any);

		render(TeamEdit, { props: { open: true, team: baseTeam } });
		await flush();

    const tName = screen.getByLabelText(/team name/i) as HTMLInputElement;
    const tBio = screen.getByLabelText(/^bio$/i) as HTMLTextAreaElement;
    const tImg = screen.getByLabelText(/image url/i) as HTMLInputElement;
    await fireEvent.input(tName, { target: { value: ' New Team ' } });
    await fireEvent.input(tBio, { target: { value: ' Cool bio ' } });
    await fireEvent.input(tImg, { target: { value: ' http://image.png ' } });
    await flush();
    await user.click(screen.getByRole('button', { name: /^save$/i }));

    await waitFor(() => {
      expect(updateTeam).toHaveBeenCalledWith(7, 'New Team', 'Cool bio', 'http://image.png', '');
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
    const tBio2 = screen.getByLabelText(/^bio$/i) as HTMLTextAreaElement;
    const tImg2 = screen.getByLabelText(/image url/i) as HTMLInputElement;
    await fireEvent.input(tName2, { target: { value: 'New Team' } });
    await fireEvent.input(tBio2, { target: { value: 'Cool bio' } });
    await fireEvent.input(tImg2, { target: { value: 'http://image.png' } });
    await flush();
    await user.click(screen.getByRole('button', { name: /^save$/i }));

    await waitFor(() => {
      expect(updateTeam).toHaveBeenCalledWith(7, 'New Team', 'Cool bio', 'http://image.png', '');
      expect(mockInvalidateQueries).toHaveBeenCalled();
    });
  });

  it('handles update error', async () => {
    const user = userEvent.setup();
    vi.mocked(updateTeam).mockRejectedValueOnce(new Error('Update failed'));

    render(TeamEdit, { props: { open: true, team: baseTeam } });

    const nm = screen.getByLabelText(/team name/i) as HTMLInputElement;
    await fireEvent.input(nm, { target: { value: 'New Team' } });
    await user.click(screen.getByRole('button', { name: /^save$/i }));

    await waitFor(() => {
      expect(toast.error).toHaveBeenCalledWith('Update failed');
    });
  });

  it('shows loading state during update', async () => {
    const user = userEvent.setup();
    let resolveUpdate: (value: any) => void = () => {};
    const updatePromise = new Promise((resolve) => {
      resolveUpdate = resolve;
    });
    vi.mocked(updateTeam).mockReturnValueOnce(updatePromise as any);

    render(TeamEdit, { props: { open: true, team: baseTeam } });

    await user.type(screen.getByLabelText(/team name/i), 'New Team');
    await flush();
    const submitButton = screen.getByRole('button', { name: /^save$/i });
    await user.click(submitButton);

    // While pending, the original submit button becomes disabled
    await waitFor(() => {
      expect(submitButton).toBeDisabled();
    });

    // Resolve and ensure it completes without throwing
    resolveUpdate({ ok: true });
  });

  it('allows updating individual fields', async () => {
    const user = userEvent.setup();
    vi.mocked(updateTeam).mockResolvedValueOnce({ ok: true } as any);

    render(TeamEdit, { props: { open: true, team: baseTeam } });

    const tBio3 = screen.getByLabelText(/^bio$/i) as HTMLTextAreaElement;
    await fireEvent.input(tBio3, { target: { value: 'Only bio updated' } });
    await flush();
    await user.click(screen.getByRole('button', { name: /^save$/i }));

    await waitFor(() => {
      expect(updateTeam).toHaveBeenCalledWith(7, '', 'Only bio updated', '', '');
    });
  });
});

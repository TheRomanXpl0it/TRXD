import { render, screen, waitFor } from '@testing-library/svelte';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { toast } from 'svelte-sonner';
import TeamEdit from '../TeamEdit.svelte';
import { updateTeam } from '$lib/team';
import { useQueryClient } from '@tanstack/svelte-query';

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

    await user.click(screen.getByRole('button', { name: /^save$/i }));

    expect(toast.error).toHaveBeenCalledWith('Please fill at least one field.');
    expect(updateTeam).not.toHaveBeenCalled();
  });

  it('validates image URL format', async () => {
    const user = userEvent.setup();

    render(TeamEdit, { props: { open: true, team: baseTeam } });

    await user.type(screen.getByLabelText(/image url/i), 'not-a-url');
    await user.click(screen.getByRole('button', { name: /^save$/i }));

    expect(toast.error).toHaveBeenCalledWith('Image must be a valid URL.');
    expect(updateTeam).not.toHaveBeenCalled();
  });

  it('trims whitespace from input fields', async () => {
    const user = userEvent.setup();
    vi.mocked(updateTeam).mockResolvedValueOnce({ ok: true } as any);

    render(TeamEdit, { props: { open: true, team: baseTeam } });

    await user.type(screen.getByLabelText(/team name/i), ' New Team ');
    await user.type(screen.getByLabelText(/^bio$/i), ' Cool bio ');
    await user.type(screen.getByLabelText(/image url/i), ' http://image.png ');
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

    await user.type(screen.getByLabelText(/team name/i), 'New Team');
    await user.type(screen.getByLabelText(/^bio$/i), 'Cool bio');
    await user.type(screen.getByLabelText(/image url/i), 'http://image.png');
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

    await user.type(screen.getByLabelText(/team name/i), 'New Team');
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
    await user.click(screen.getByRole('button', { name: /^save$/i }));

    // While pending, shows disabled Saving... button
    const savingButton = screen.getByRole('button', { name: /saving/i });
    expect(savingButton).toBeDisabled();

    // Resolve and ensure it completes without throwing
    resolveUpdate({ ok: true });
  });

  it('allows updating individual fields', async () => {
    const user = userEvent.setup();
    vi.mocked(updateTeam).mockResolvedValueOnce({ ok: true } as any);

    render(TeamEdit, { props: { open: true, team: baseTeam } });

    await user.type(screen.getByLabelText(/^bio$/i), 'Only bio updated');
    await user.click(screen.getByRole('button', { name: /^save$/i }));

    await waitFor(() => {
      expect(updateTeam).toHaveBeenCalledWith(7, '', 'Only bio updated', '', '');
    });
  });
});

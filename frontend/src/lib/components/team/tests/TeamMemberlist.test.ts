import { render, screen, waitFor } from '@testing-library/svelte';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import TeamMemberlist from '../TeamMemberlist.svelte';

const mockPush = vi.fn();
vi.mock('svelte-spa-router', () => ({
  push: (...args: any[]) => mockPush(...args)
}));

const mockGetUserData = vi.fn();
vi.mock('$lib/user', () => ({ getUserData: (...args: any[]) => mockGetUserData(...args) }));

const team = {
  id: 1,
  name: 'CyberCats',
  score: 4242,
  members: [
    { id: 10, name: 'Alice', role: 'Captain', score: 1200 },
    { id: 11, name: 'Bob Sea', role: 'Member', score: 800 },
    { id: 12, name: 'Seann', role: 'Member', score: 600 },
    { id: 13, name: 'Zed', role: 'Member', score: 200 }
  ]
};

describe('TeamMemberlist', () => {
  beforeEach(() => vi.clearAllMocks());

  it('renders members and supports fuzzy search', async () => {
    mockGetUserData.mockResolvedValue({ image: null });
    render(TeamMemberlist, { props: { team } });

    // Shows header and counts
    expect(screen.getByText(/members/i)).toBeInTheDocument();
    expect(screen.getByText(/showing/i)).toBeInTheDocument();

    // Filter to alice only
    const search = screen.getByPlaceholderText(/search members/i);
    await userEvent.type(search, 'alice');
    await waitFor(() => {
      expect(screen.getByText('Alice')).toBeInTheDocument();
    });
    // Non-matching disappear
    expect(screen.queryByText('Zed')).not.toBeInTheDocument();
  });

  it('navigates to account on name click', async () => {
    mockGetUserData.mockResolvedValue({ image: null });
    render(TeamMemberlist, { props: { team } });

    const nameLink = await screen.findByText('Alice');
    await userEvent.click(nameLink);
    expect(mockPush).toHaveBeenCalledWith('/account/10');
  });

  it('fetches user images for members', async () => {
    mockGetUserData.mockResolvedValue({ image: 'http://img.png' });
    render(TeamMemberlist, { props: { team } });

    await waitFor(() => {
      expect(mockGetUserData).toHaveBeenCalledWith(10);
      expect(mockGetUserData).toHaveBeenCalledWith(11);
      expect(mockGetUserData).toHaveBeenCalledWith(12);
      expect(mockGetUserData).toHaveBeenCalledWith(13);
    });
  });
});


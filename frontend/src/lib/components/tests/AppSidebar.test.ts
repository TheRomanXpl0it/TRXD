import { render, screen } from '@testing-library/svelte';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import AppSidebar from '../AppSidebar.svelte';
import { useSidebar } from '$lib/components/ui/sidebar/context.svelte.js';
import { getUserData } from '$lib/user';
import { push } from 'svelte-spa-router';

vi.mock('$lib/components/ui/sidebar/context.svelte.js', () => ({
	useSidebar: vi.fn()
}));

vi.mock('$lib/user', () => ({
	getUserData: vi.fn()
}));

vi.mock('svelte-spa-router', () => ({
	link: () => {},
	push: vi.fn()
}));

describe('AppSidebar Component', () => {
	const mockSetOpenMobile = vi.fn();

	beforeEach(() => {
		vi.clearAllMocks();
		vi.mocked(useSidebar).mockReturnValue({
			isMobile: false,
			openMobile: false,
			state: 'expanded',
			setOpenMobile: mockSetOpenMobile,
			open: true,
			setOpen: vi.fn()
		} as any);
	});

	it('renders base navigation items', () => {
		render(AppSidebar, {
			props: {
				user: null,
				userMode: true
			}
		});

		expect(screen.getByText('TRXD')).toBeInTheDocument();
		expect(screen.getByText('Home')).toBeInTheDocument();
		expect(screen.getByText('Scoreboard')).toBeInTheDocument();
		expect(screen.getByText('Challenges')).toBeInTheDocument();
		// Accounts, Teams, and Team should not show when user is null
		expect(screen.queryByText('Accounts')).not.toBeInTheDocument();
		expect(screen.queryByText('Teams')).not.toBeInTheDocument();
		expect(screen.queryByText('Team')).not.toBeInTheDocument();
	});

	it('shows sign in link when user is not logged in', () => {
		render(AppSidebar, {
			props: {
				user: null,
				userMode: true
			}
		});

		expect(screen.getByText('Sign in')).toBeInTheDocument();
		expect(screen.queryByText('Team')).not.toBeInTheDocument();
		expect(screen.queryByText('Configs')).not.toBeInTheDocument();
	});

	it('shows team item when not in user mode', () => {
		render(AppSidebar, {
			props: {
				user: {
					id: 1,
					name: 'Alice'
				},
				userMode: false
			} as any
		});

		expect(screen.getByText('Team')).toBeInTheDocument();
		expect(screen.getByText('Teams')).toBeInTheDocument();
		expect(screen.getByText('Accounts')).toBeInTheDocument();
	});

	it('shows configs item for admin users', () => {
		render(AppSidebar, {
			props: {
				user: {
					id: 1,
					name: 'Ada',
					role: 'Admin'
				},
				userMode: true
			} as any
		});

		expect(screen.getByText('Configs')).toBeInTheDocument();
	});

	it('hides configs item for non-admin users', () => {
		render(AppSidebar, {
			props: {
				user: {
					id: 2,
					name: 'Bob',
					role: 'User'
				},
				userMode: true
			} as any
		});

		expect(screen.queryByText('Configs')).not.toBeInTheDocument();
	});

	it('closes mobile sidebar when navigation item is clicked', async () => {
		const user = userEvent.setup();
		vi.mocked(useSidebar).mockReturnValue({
			isMobile: true,
			openMobile: true,
			state: 'expanded',
			setOpenMobile: mockSetOpenMobile,
			open: true,
			setOpen: vi.fn()
		} as any);

		render(AppSidebar, {
			props: {
				user: null,
				userMode: true
			}
		});

		const link = screen.getByText('Home').closest('a')!;
		await user.click(link);

		expect(mockSetOpenMobile).toHaveBeenCalledWith(false);
	});

	it('does not close sidebar on desktop', async () => {
		const user = userEvent.setup();
		vi.mocked(useSidebar).mockReturnValue({
			isMobile: false,
			openMobile: false,
			state: 'expanded',
			setOpenMobile: mockSetOpenMobile,
			open: true,
			setOpen: vi.fn()
		} as any);

		render(AppSidebar, {
			props: {
				user: null,
				userMode: true
			}
		});

		const link = screen.getByText('Home').closest('a')!;
		await user.click(link);

		expect(mockSetOpenMobile).not.toHaveBeenCalled();
	});

	it('renders user section for logged in users', async () => {
		vi.mocked(getUserData).mockResolvedValueOnce({
			image: 'http://img/avatar.png'
		});

		render(AppSidebar, {
			props: {
				user: {
					id: 42,
					name: 'Alice'
				},
				userMode: true
			} as any
		});

		expect(await screen.findByText('Alice')).toBeInTheDocument();
	});

	it('triggers sign out when logout button is clicked', async () => {
		const user = userEvent.setup();
		vi.mocked(getUserData).mockResolvedValueOnce({
			image: 'http://img/avatar.png'
		});
		vi.mocked(useSidebar).mockReturnValue({
			isMobile: true,
			openMobile: true,
			state: 'expanded',
			setOpenMobile: mockSetOpenMobile,
			open: true,
			setOpen: vi.fn()
		} as any);

		render(AppSidebar, {
			props: {
				user: {
					id: 42,
					name: 'Alice'
				},
				userMode: true
			} as any
		});

		const logoutBtn = await screen.findByRole('button', { name: /log out/i });
		await user.click(logoutBtn);

		expect(push).toHaveBeenCalledWith('/signOut');
	});

	it('displays user avatar when image is available', async () => {
		vi.mocked(getUserData).mockResolvedValueOnce({
			image: 'http://img/avatar.png'
		});

		render(AppSidebar, {
			props: {
				user: {
					id: 42,
					name: 'Alice'
				},
				userMode: true
			} as any
		});

		const avatar = await screen.findByRole('img');
		expect(avatar).toHaveAttribute('src', 'http://img/avatar.png');
	});

	it('handles missing user avatar gracefully', async () => {
		vi.mocked(getUserData).mockResolvedValueOnce({
			image: ''
		});

		render(AppSidebar, {
			props: {
				user: {
					id: 42,
					name: 'Alice'
				},
				userMode: true
			} as any
		});

		expect(await screen.findByText('Alice')).toBeInTheDocument();
	});
});

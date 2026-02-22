import { render, screen } from '@testing-library/svelte';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import ChallengeFilters from '../ChallengeFilters.svelte';

describe('ChallengeFilters Component', () => {
	const defaultProps = {
		search: '',
		filterCategories: [],
		filterTags: [],
		categories: [
			{ value: 'web', label: 'Web' },
			{ value: 'crypto', label: 'Crypto' },
			{ value: 'pwn', label: 'Pwn' }
		],
		allTags: ['easy', 'medium', 'hard'],
		compactView: false,
		activeFiltersCount: 0
	};

	beforeEach(() => {
		vi.clearAllMocks();
	});

	it('renders search input with correct placeholder', () => {
		render(ChallengeFilters, { props: defaultProps });

		const searchInput = screen.getByPlaceholderText(/search challenges by name, category, or tag/i);
		expect(searchInput).toBeInTheDocument();
	});

	it('updates search value when typing', async () => {
		const user = userEvent.setup();
		let search = '';

		const { component } = render(ChallengeFilters, {
			props: {
				...defaultProps,
				search
			}
		});

		const searchInput = screen.getByRole('textbox', { name: /search challenges/i });
		await user.type(searchInput, 'test challenge');

		// Check that input has the value
		expect(searchInput).toHaveValue('test challenge');
	});

	it('shows clear search button when search has value', async () => {
		const user = userEvent.setup();

		render(ChallengeFilters, {
			props: {
				...defaultProps,
				search: 'test'
			}
		});

		const clearButton = screen.getByRole('button', { name: /clear search/i });
		expect(clearButton).toBeInTheDocument();
	});

	it('does not show clear search button when search is empty', () => {
		render(ChallengeFilters, {
			props: {
				...defaultProps,
				search: ''
			}
		});

		const clearButton = screen.queryByRole('button', { name: /clear search/i });
		expect(clearButton).not.toBeInTheDocument();
	});

	it('clears search when clear button is clicked', async () => {
		const user = userEvent.setup();

		const { component } = render(ChallengeFilters, {
			props: {
				...defaultProps,
				search: 'test search'
			}
		});

		const clearButton = screen.getByRole('button', { name: /clear search/i });
		await user.click(clearButton);

		const searchInput = screen.getByRole('textbox');
		expect(searchInput).toHaveValue('');
	});

	it('renders categories filter button', () => {
		render(ChallengeFilters, { props: defaultProps });

		expect(screen.getByRole('button', { name: /filter by categories/i })).toBeInTheDocument();
	});

	it('renders tags filter button', () => {
		render(ChallengeFilters, { props: defaultProps });

		expect(screen.getByRole('button', { name: /filter by tags/i })).toBeInTheDocument();
	});

	it('shows clear filters button when activeFiltersCount > 0', () => {
		render(ChallengeFilters, {
			props: {
				...defaultProps,
				activeFiltersCount: 3
			}
		});

		expect(
			screen.getByRole('button', { name: /clear all filters \(3 active\)/i })
		).toBeInTheDocument();
	});

	it('does not show clear filters button when activeFiltersCount is 0', () => {
		render(ChallengeFilters, {
			props: {
				...defaultProps,
				activeFiltersCount: 0
			}
		});

		expect(screen.queryByRole('button', { name: /clear all filters/i })).not.toBeInTheDocument();
	});

	it('renders grid view button', () => {
		render(ChallengeFilters, { props: defaultProps });

		expect(screen.getByRole('button', { name: /grid view/i })).toBeInTheDocument();
	});

	it('renders compact view button', () => {
		render(ChallengeFilters, { props: defaultProps });

		expect(screen.getByRole('button', { name: /compact view/i })).toBeInTheDocument();
	});

	it('grid view button is pressed when compactView is false', () => {
		render(ChallengeFilters, {
			props: {
				...defaultProps,
				compactView: false
			}
		});

		const gridButton = screen.getByRole('button', { name: /grid view/i });
		expect(gridButton).toHaveAttribute('aria-pressed', 'true');
	});

	it('compact view button is pressed when compactView is true', () => {
		render(ChallengeFilters, {
			props: {
				...defaultProps,
				compactView: true
			}
		});

		const compactButton = screen.getByRole('button', { name: /compact view/i });
		expect(compactButton).toHaveAttribute('aria-pressed', 'true');
	});

	it('switches to compact view when compact button is clicked', async () => {
		const user = userEvent.setup();

		const { component } = render(ChallengeFilters, {
			props: {
				...defaultProps,
				compactView: false
			}
		});

		const compactButton = screen.getByRole('button', { name: /compact view/i });
		await user.click(compactButton);

		expect(compactButton).toHaveAttribute('aria-pressed', 'true');
	});

	it('switches to grid view when grid button is clicked', async () => {
		const user = userEvent.setup();

		const { component } = render(ChallengeFilters, {
			props: {
				...defaultProps,
				compactView: true
			}
		});

		const gridButton = screen.getByRole('button', { name: /grid view/i });
		await user.click(gridButton);

		expect(gridButton).toHaveAttribute('aria-pressed', 'true');
	});

	it('displays correct view mode group label', () => {
		render(ChallengeFilters, { props: defaultProps });

		expect(screen.getByRole('group', { name: /view mode/i })).toBeInTheDocument();
	});

	it('renders all available categories in categories list', async () => {
		const user = userEvent.setup();

		render(ChallengeFilters, { props: defaultProps });

		// Open categories popover
		const categoriesButton = screen.getByRole('button', { name: /filter by categories/i });
		await user.click(categoriesButton);

		// Wait for popover to open and check for categories
		expect(await screen.findByText('Web')).toBeInTheDocument();
		expect(screen.getByText('Crypto')).toBeInTheDocument();
		expect(screen.getByText('Pwn')).toBeInTheDocument();
	});

	it('renders all available tags in tags list', async () => {
		const user = userEvent.setup();

		render(ChallengeFilters, { props: defaultProps });

		// Open tags popover
		const tagsButton = screen.getByRole('button', { name: /filter by tags/i });
		await user.click(tagsButton);

		// Wait for popover to open and check for tags
		expect(await screen.findByText('easy')).toBeInTheDocument();
		expect(screen.getByText('medium')).toBeInTheDocument();
		expect(screen.getByText('hard')).toBeInTheDocument();
	});

	it('has search input in categories popover', async () => {
		const user = userEvent.setup();

		render(ChallengeFilters, { props: defaultProps });

		const categoriesButton = screen.getByRole('button', { name: /filter by categories/i });
		await user.click(categoriesButton);

		expect(await screen.findByPlaceholderText(/search categories/i)).toBeInTheDocument();
	});

	it('has search input in tags popover', async () => {
		const user = userEvent.setup();

		render(ChallengeFilters, { props: defaultProps });

		const tagsButton = screen.getByRole('button', { name: /filter by tags/i });
		await user.click(tagsButton);

		expect(await screen.findByPlaceholderText(/search tags/i)).toBeInTheDocument();
	});
});

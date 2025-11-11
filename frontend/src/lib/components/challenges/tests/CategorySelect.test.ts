import { render, screen, waitFor } from '@testing-library/svelte';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, beforeEach } from 'vitest';
import CategorySelect from '../CategorySelect.svelte';

describe('CategorySelect Component', () => {
	const categories = [
		{ value: 'web', label: 'Web' },
		{ value: 'crypto', label: 'Cryptography' },
		{ value: 'pwn', label: 'Binary Exploitation' },
		{ value: 'forensics', label: 'Forensics' },
		{ value: 'rev', label: 'Reverse Engineering' }
	];

	beforeEach(() => {
		// Clear any existing popovers
		document.body.innerHTML = '';
		
		Element.prototype.scrollIntoView = function() {};
	});

	it('renders with placeholder text', () => {
		render(CategorySelect, {
			props: {
				items: categories,
				placeholder: 'Select a category...'
			}
		});

		expect(screen.getByRole('combobox')).toBeInTheDocument();
		expect(screen.getByText('Select a category...')).toBeInTheDocument();
	});

	it('displays selected category label', () => {
		render(CategorySelect, {
			props: {
				items: categories,
				value: 'crypto'
			}
		});

		expect(screen.getByText('Cryptography')).toBeInTheDocument();
	});

	it('opens popover when button is clicked', async () => {
		const user = userEvent.setup();

		render(CategorySelect, {
			props: {
				items: categories
			}
		});

		const button = screen.getByRole('combobox');
		await user.click(button);

		await waitFor(() => {
			expect(button).toHaveAttribute('aria-expanded', 'true');
		});
	});

	it('displays all category items in popover', async () => {
		const user = userEvent.setup();

		render(CategorySelect, {
			props: {
				items: categories
			}
		});

		await user.click(screen.getByRole('combobox'));

		await waitFor(() => {
			expect(screen.getByText('Web')).toBeInTheDocument();
			expect(screen.getByText('Cryptography')).toBeInTheDocument();
			expect(screen.getByText('Binary Exploitation')).toBeInTheDocument();
			expect(screen.getByText('Forensics')).toBeInTheDocument();
			expect(screen.getByText('Reverse Engineering')).toBeInTheDocument();
		});
	});

	it('shows search input in popover', async () => {
		const user = userEvent.setup();

		render(CategorySelect, {
			props: {
				items: categories,
				searchPlaceholder: 'Search category...'
			}
		});

		await user.click(screen.getByRole('combobox'));

		await waitFor(() => {
			expect(screen.getByPlaceholderText('Search category...')).toBeInTheDocument();
		});
	});

	it('filters items based on search input', async () => {
		const user = userEvent.setup();

		render(CategorySelect, {
			props: {
				items: categories,
				searchPlaceholder: 'Search...'
			}
		});

		await user.click(screen.getByRole('combobox'));

		await waitFor(() => {
			expect(screen.getByPlaceholderText('Search...')).toBeInTheDocument();
		});

		const searchInput = screen.getByPlaceholderText('Search...');
		await user.type(searchInput, 'Crypto');

		await waitFor(() => {
			expect(screen.getByText('Cryptography')).toBeInTheDocument();
			expect(screen.queryByText('Web')).not.toBeInTheDocument();
		});
	});

	it('selects category when item is clicked', async () => {
		const user = userEvent.setup();

		let selectedValue = '';
		
		render(CategorySelect, {
			props: {
				items: categories,
				value: selectedValue
			}
		});

		await user.click(screen.getByRole('combobox'));

		await waitFor(() => {
			expect(screen.getByText('Web')).toBeInTheDocument();
		});

		await user.click(screen.getByText('Web'));

		await waitFor(() => {
			expect(screen.getByRole('combobox')).toHaveTextContent('Web');
		});
	});

	it('closes popover after selecting item', async () => {
		const user = userEvent.setup();

		render(CategorySelect, {
			props: {
				items: categories
			}
		});

		const button = screen.getByRole('combobox');
		await user.click(button);

		await waitFor(() => {
			expect(button).toHaveAttribute('aria-expanded', 'true');
		});

		await user.click(screen.getByText('Forensics'));

		await waitFor(() => {
			expect(button).toHaveAttribute('aria-expanded', 'false');
		});
	});

	it('displays selected category after selection', async () => {
		const user = userEvent.setup();

		render(CategorySelect, {
			props: {
				items: categories,
				placeholder: 'Choose...'
			}
		});

		await user.click(screen.getByRole('combobox'));

		await waitFor(() => {
			expect(screen.getByText('Binary Exploitation')).toBeInTheDocument();
		});

		await user.click(screen.getByText('Binary Exploitation'));

		await waitFor(() => {
			expect(screen.getByText('Binary Exploitation')).toBeInTheDocument();
			expect(screen.queryByText('Choose...')).not.toBeInTheDocument();
		});
	});

	it('shows checkmark for selected item', async () => {
		const user = userEvent.setup();

		render(CategorySelect, {
			props: {
				items: categories,
				value: 'rev'
			}
		});

		await user.click(screen.getByRole('combobox'));

		await waitFor(() => {
			const options = screen.getAllByRole('option');
			const revOption = options.find((opt) => opt.textContent?.includes('Reverse Engineering'));
			expect(revOption).toBeTruthy();

			const revIcon = revOption!.querySelector('svg');
			expect(revIcon).toBeTruthy();
			expect(revIcon).not.toHaveClass('text-transparent');

			const otherOptions = options.filter((opt) => opt !== revOption);
			for (const opt of otherOptions) {
				const icon = opt.querySelector('svg');
				if (icon) {
					expect(icon).toHaveClass('text-transparent');
				}
			}
		});
	});

	it('opens popover with keyboard (Enter key)', async () => {
		const user = userEvent.setup();

		render(CategorySelect, {
			props: {
				items: categories
			}
		});

		const button = screen.getByRole('combobox');
		button.focus();

		await user.keyboard('{Enter}');

		await waitFor(() => {
			expect(button).toHaveAttribute('aria-expanded', 'true');
		});
	});

	it('closes popover with keyboard (Escape key)', async () => {
		const user = userEvent.setup();

		render(CategorySelect, {
			props: {
				items: categories
			}
		});

		const button = screen.getByRole('combobox');
		await user.click(button);

		await waitFor(() => {
			expect(button).toHaveAttribute('aria-expanded', 'true');
		});

		await user.keyboard('{Escape}');

		await waitFor(() => {
			expect(button).toHaveAttribute('aria-expanded', 'false');
		});
	});

	it('returns focus to trigger button after selection', async () => {
		const user = userEvent.setup();

		render(CategorySelect, {
			props: {
				items: categories
			}
		});

		const button = screen.getByRole('combobox');
		await user.click(button);

		await waitFor(() => {
			expect(screen.getByText('Web')).toBeInTheDocument();
		});

		await user.click(screen.getByText('Web'));

		await waitFor(() => {
			expect(document.activeElement).toBe(button);
		});
	});


	it('shows "No results" when search has no matches', async () => {
		const user = userEvent.setup();

		render(CategorySelect, {
			props: {
				items: categories,
				searchPlaceholder: 'Search...'
			}
		});

		await user.click(screen.getByRole('combobox'));

		await waitFor(() => {
			expect(screen.getByPlaceholderText('Search...')).toBeInTheDocument();
		});

		const searchInput = screen.getByPlaceholderText('Search...');
		await user.type(searchInput, 'xyz123');

		await waitFor(() => {
			expect(screen.getByText('No results.')).toBeInTheDocument();
		});
	});

	it('uses custom placeholder text', () => {
		render(CategorySelect, {
			props: {
				items: categories,
				placeholder: 'Pick a category'
			}
		});

		expect(screen.getByText('Pick a category')).toBeInTheDocument();
	});

	it('uses custom search placeholder text', async () => {
		const user = userEvent.setup();

		render(CategorySelect, {
			props: {
				items: categories,
				searchPlaceholder: 'Type to search...'
			}
		});

		await user.click(screen.getByRole('combobox'));

		await waitFor(() => {
			expect(screen.getByPlaceholderText('Type to search...')).toBeInTheDocument();
		});
	});

	it('applies custom width class', () => {
		render(CategorySelect, {
			props: {
				items: categories,
				widthClass: 'w-[300px]'
			}
		});

		const button = screen.getByRole('combobox');
		expect(button).toHaveClass('w-[300px]');
	});

	it('applies custom className', () => {
		render(CategorySelect, {
			props: {
				items: categories,
				className: 'my-custom-class'
			}
		});

		const button = screen.getByRole('combobox');
		expect(button).toHaveClass('my-custom-class');
	});

	it('handles empty items array', () => {
		render(CategorySelect, {
			props: {
				items: [],
				placeholder: 'No categories'
			}
		});

		expect(screen.getByText('No categories')).toBeInTheDocument();
	});

	it('shows "No results" when items array is empty and popover is opened', async () => {
		const user = userEvent.setup();

		render(CategorySelect, {
			props: {
				items: []
			}
		});

		await user.click(screen.getByRole('combobox'));

		await waitFor(() => {
			expect(screen.getByText('No results.')).toBeInTheDocument();
		});
	});

	it('allows changing selection multiple times', async () => {
		const user = userEvent.setup();

		render(CategorySelect, {
			props: {
				items: categories,
				value: 'web'
			}
		});

		expect(screen.getByText('Web')).toBeInTheDocument();

		// Change to Crypto
		await user.click(screen.getByRole('combobox'));
		await waitFor(() => {
			expect(screen.getByText('Cryptography')).toBeInTheDocument();
		});
		await user.click(screen.getByText('Cryptography'));

		await waitFor(() => {
			expect(screen.getByRole('combobox')).toHaveTextContent('Cryptography');
		});

		// Change to Forensics
		await user.click(screen.getByRole('combobox'));
		await waitFor(() => {
			expect(screen.getByText('Forensics')).toBeInTheDocument();
		});
		await user.click(screen.getByText('Forensics'));

		await waitFor(() => {
			expect(screen.getByRole('combobox')).toHaveTextContent('Forensics');
		});
	});

	it('filters are case-insensitive', async () => {
		const user = userEvent.setup();

		render(CategorySelect, {
			props: {
				items: categories,
				searchPlaceholder: 'Search...'
			}
		});

		await user.click(screen.getByRole('combobox'));

		await waitFor(() => {
			expect(screen.getByPlaceholderText('Search...')).toBeInTheDocument();
		});

		const searchInput = screen.getByPlaceholderText('Search...');
		await user.type(searchInput, 'CRYPTO');

		await waitFor(() => {
			expect(screen.getByText('Cryptography')).toBeInTheDocument();
		});
	});

	it('clears search when popover is reopened', async () => {
		const user = userEvent.setup();

		render(CategorySelect, {
			props: {
				items: categories,
				searchPlaceholder: 'Search...'
			}
		});

		// Open and search
		await user.click(screen.getByRole('combobox'));
		await waitFor(() => {
			expect(screen.getByPlaceholderText('Search...')).toBeInTheDocument();
		});

		const searchInput = screen.getByPlaceholderText('Search...');
		await user.type(searchInput, 'Web');

		// Select an item to close
		await user.click(screen.getByText('Web'));

		// Reopen
		await user.click(screen.getByRole('combobox'));

		await waitFor(() => {
			const newSearchInput = screen.getByPlaceholderText('Search...');
			expect(newSearchInput).toHaveValue('');
		});
	});

	it('handles items with special characters in labels', async () => {
		const user = userEvent.setup();

		const specialItems = [
			{ value: 'test1', label: "Category's Name" },
			{ value: 'test2', label: 'Category "Special"' },
			{ value: 'test3', label: 'Category & More' }
		];

		render(CategorySelect, {
			props: {
				items: specialItems
			}
		});

		await user.click(screen.getByRole('combobox'));

		await waitFor(() => {
			expect(screen.getByText("Category's Name")).toBeInTheDocument();
			expect(screen.getByText('Category "Special"')).toBeInTheDocument();
			expect(screen.getByText('Category & More')).toBeInTheDocument();
		});
	});

	it('displays correct number of items', async () => {
		const user = userEvent.setup();

		render(CategorySelect, {
			props: {
				items: categories
			}
		});

		await user.click(screen.getByRole('combobox'));

		await waitFor(() => {
			const options = screen.getAllByRole('option');
			expect(options).toHaveLength(categories.length);
		});
	});

	it('maintains selection when reopening popover', async () => {
		const user = userEvent.setup();

		render(CategorySelect, {
			props: {
				items: categories,
				value: 'crypto'
			}
		});

		// Open popover
		await user.click(screen.getByRole('combobox'));

		await waitFor(() => {
			const options = screen.getAllByRole('option');
			const selectedOption = options.find((opt) => opt.textContent?.includes('Cryptography'));
			const icon = selectedOption?.querySelector('svg');
			expect(icon).not.toHaveClass('text-transparent');
		});

		// Close popover
		await user.keyboard('{Escape}');

		// Reopen popover
		await user.click(screen.getByRole('combobox'));

		await waitFor(() => {
			const options = screen.getAllByRole('option');
			const selectedOption = options.find((opt) => opt.textContent?.includes('Cryptography'));
			const icon = selectedOption?.querySelector('svg');
			expect(icon).not.toHaveClass('text-transparent');
		});
	});

	it('filters out all items when search does not match', async () => {
		const user = userEvent.setup();

		render(CategorySelect, {
			props: {
				items: categories,
				searchPlaceholder: 'Search...'
			}
		});

		await user.click(screen.getByRole('combobox'));

		await waitFor(() => {
			expect(screen.getByPlaceholderText('Search...')).toBeInTheDocument();
		});

		const searchInput = screen.getByPlaceholderText('Search...');
		await user.type(searchInput, 'xyz123');

		await waitFor(() => {
			const options = screen.queryAllByRole('option');
			expect(options).toHaveLength(0);
			expect(screen.getByText('No results.')).toBeInTheDocument();
		});
	});

	it('applies custom group label', async () => {
		const user = userEvent.setup();

		render(CategorySelect, {
			props: {
				items: categories,
				groupLabel: 'challenge-categories'
			}
		});

		await user.click(screen.getByRole('combobox'));

		await waitFor(() => {
			// The group label is used as the data-value attribute
			const group = document.querySelector('[data-command-group][data-value="challenge-categories"]');
			expect(group).toBeInTheDocument();
		});
	});

	it('handles single item in list', async () => {
		const user = userEvent.setup();

		const singleItem = [{ value: 'web', label: 'Web' }];

		render(CategorySelect, {
			props: {
				items: singleItem
			}
		});

		await user.click(screen.getByRole('combobox'));

		await waitFor(() => {
			const options = screen.getAllByRole('option');
			expect(options).toHaveLength(1);
			expect(screen.getByText('Web')).toBeInTheDocument();
		});
	});

	it('handles very long category labels', async () => {
		const user = userEvent.setup();

		const longLabelItems = [
			{
				value: 'long',
				label: 'This is a very long category name that should still be displayed correctly'
			}
		];

		render(CategorySelect, {
			props: {
				items: longLabelItems
			}
		});

		await user.click(screen.getByRole('combobox'));

		await waitFor(() => {
			expect(
				screen.getByText(
					'This is a very long category name that should still be displayed correctly'
				)
			).toBeInTheDocument();
		});
	});

	it('handles rapid selection changes', async () => {
		const user = userEvent.setup();

		render(CategorySelect, {
			props: {
				items: categories,
				value: 'web'
			}
		});

		// Rapidly change selections
		await user.click(screen.getByRole('combobox'));
		await waitFor(() => {
			expect(screen.getByText('Cryptography')).toBeInTheDocument();
		});
		await user.click(screen.getByText('Cryptography'));

		await waitFor(() => {
			expect(screen.getByRole('combobox')).toHaveTextContent('Cryptography');
		});

		// Change again immediately
		await user.click(screen.getByRole('combobox'));
		await waitFor(() => {
			expect(screen.getByText('Forensics')).toBeInTheDocument();
		});
		await user.click(screen.getByText('Forensics'));

		await waitFor(() => {
			expect(screen.getByRole('combobox')).toHaveTextContent('Forensics');
		});
	});

	it('clears selection state visually', async () => {
		const user = userEvent.setup();

		render(CategorySelect, {
			props: {
				items: categories,
				value: 'web',
				placeholder: 'Select...'
			}
		});

		expect(screen.getByText('Web')).toBeInTheDocument();

		// Select another item
		await user.click(screen.getByRole('combobox'));

		await waitFor(() => {
			expect(screen.getByText('Cryptography')).toBeInTheDocument();
		});

		await user.click(screen.getByText('Cryptography'));

		await waitFor(() => {
			expect(screen.getByRole('combobox')).toHaveTextContent('Cryptography');
			expect(screen.queryByText('Web')).not.toBeInTheDocument();
		});
	});
});

import { render, screen, fireEvent, waitFor } from '@testing-library/svelte';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi } from 'vitest';
import CountrySelect from '../country-select.svelte';
import * as countriesUtil from '$lib/utils/countries';

// Mock getCountryItems to have a small controlled set
vi.mock('$lib/utils/countries', async (importOriginal) => {
    const original = await importOriginal<any>();
    return {
        ...original,
        getCountryItems: vi.fn(() => [
            { value: 'ITA', label: 'Italy', iso2: 'IT' },
            { value: 'USA', label: 'United States', iso2: 'US' },
            { value: 'FRA', label: 'France', iso2: 'FR' }
        ])
    };
});

describe('CountrySelect Component', () => {
    it('renders with placeholder', () => {
        render(CountrySelect, {
            props: {
                placeholder: 'Choose a country'
            }
        });
        expect(screen.getByText('Choose a country')).toBeInTheDocument();
    });

    it('renders selected country label', () => {
        render(CountrySelect, {
            props: {
                value: 'ITA'
            }
        });
        expect(screen.getByText('ITA')).toBeInTheDocument();
    });

    it('opens popover on click and shows items', async () => {
        const user = userEvent.setup();
        render(CountrySelect);

        const trigger = screen.getByRole('combobox');
        await user.click(trigger);

        // Popover might be async
        await waitFor(() => {
            expect(screen.getByPlaceholderText('Search countries...')).toBeInTheDocument();
            expect(screen.getByText('Italy')).toBeInTheDocument();
            expect(screen.getByText('United States')).toBeInTheDocument();
        });
    });

    it('filters countries when searching', async () => {
        const user = userEvent.setup();
        render(CountrySelect);

        await user.click(screen.getByRole('combobox'));

        let searchInput: HTMLElement;
        await waitFor(() => {
            searchInput = screen.getByPlaceholderText('Search countries...');
            expect(searchInput).toBeInTheDocument();
        });

        await user.type(searchInput!, 'France');

        await waitFor(() => {
            expect(screen.getByText('France')).toBeInTheDocument();
            expect(screen.queryByText('Italy')).not.toBeInTheDocument();
        });
    });

    it('selects a country and updates value', async () => {
        const user = userEvent.setup();

        render(CountrySelect, {
            props: {
                value: ''
            }
        });

        await user.click(screen.getByRole('combobox'));

        let italyItem: HTMLElement;
        await waitFor(() => {
            italyItem = screen.getByText('Italy');
            expect(italyItem).toBeInTheDocument();
        });

        await user.click(italyItem!);

        await waitFor(() => {
            expect(screen.queryByPlaceholderText('Search countries...')).not.toBeInTheDocument();
        });
    });
});

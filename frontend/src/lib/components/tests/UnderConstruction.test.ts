import { render, screen } from '@testing-library/svelte';
import { describe, it, expect } from 'vitest';
import UnderConstruction from '../UnderConstruction.svelte';

describe('UnderConstruction Component', () => {
	it('renders construction image', () => {
		render(UnderConstruction);

		expect(screen.getByAltText(/under construction/i)).toBeInTheDocument();
	});

	it('renders heading text', () => {
		render(UnderConstruction);

		expect(screen.getByText(/under/i)).toBeInTheDocument();
		expect(screen.getByText(/construction/i)).toBeInTheDocument();
	});

	it('displays image with correct alt text', () => {
		render(UnderConstruction);

		const image = screen.getByAltText(/under construction/i);
		expect(image).toBeInTheDocument();
		expect(image.tagName).toBe('IMG');
	});
});


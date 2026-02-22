import '@testing-library/jest-dom';
import { vi, afterEach } from 'vitest';
import { cleanup } from '@testing-library/svelte';

// Mock window.matchMedia for components that use media queries
Object.defineProperty(window, 'matchMedia', {
	writable: true,
	value: vi.fn().mockImplementation((query) => ({
		matches: false,
		media: query,
		onchange: null,
		addListener: vi.fn(), // deprecated
		removeListener: vi.fn(), // deprecated
		addEventListener: vi.fn(),
		removeEventListener: vi.fn(),
		dispatchEvent: vi.fn()
	}))
});

// Ensure each test starts with a clean DOM
afterEach(async () => {
	cleanup();
	// Wait for any pending timers to complete before tearing down
	// This prevents bits-ui cleanup timers from firing after test teardown
	await new Promise((resolve) => setTimeout(resolve, 100));
});

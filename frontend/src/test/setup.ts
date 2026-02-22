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

// Mock ResizeObserver for components that use it (like virtual lists)
if (typeof global !== 'undefined') {
	global.ResizeObserver = class ResizeObserver {
		observe = vi.fn();
		unobserve = vi.fn();
		disconnect = vi.fn();
	};
}

// Ensure each test starts with a clean DOM
afterEach(async () => {
	cleanup();
	// Wait for any pending timers to complete before tearing down
	// This prevents bits-ui cleanup timers from firing after test teardown
	await new Promise((resolve) => setTimeout(resolve, 100));
});

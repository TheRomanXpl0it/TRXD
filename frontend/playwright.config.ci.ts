import { defineConfig, devices } from '@playwright/test';

/**
 * CI-specific Playwright configuration.
 * Infrastructure (Docker Compose) is managed by GitHub Actions, not by Playwright's webServer.
 */
export default defineConfig({
	testDir: './e2e',
	/* Global setup to build Docker images */
	globalSetup: './e2e/setup.ts',
	/* Run tests in files in parallel */
	fullyParallel: true,
	/* Fail the build on CI if you accidentally left test.only in the source code. */
	forbidOnly: !!process.env.CI,
	/* Retry on CI only */
	retries: process.env.CI ? 2 : 0,
	/* Opt out of parallel tests on CI. */
	workers: process.env.CI ? 1 : undefined,
	/* Reporter to use. See https://playwright.dev/docs/test-reporters */
	reporter: 'html',
	/* Update snapshots automatically in local dev, but fail in CI if they don't match */
	updateSnapshots: process.env.CI ? 'none' : 'missing',
	/* Shared settings for all the projects below. See https://playwright.dev/docs/api/class-testoptions. */
	use: {
		/* Base URL to use in actions like `await page.goto('/')`. */
		baseURL: 'http://localhost:80',
		/* Collect trace when retrying the failed test. See https://playwright.dev/docs/trace-viewer */
		trace: 'on-first-retry',
		screenshot: 'only-on-failure'
	},

	/* Configure projects for major browsers */
	projects: [
		{
			name: 'chromium',
			use: { ...devices['Desktop Chrome'] }
		},

		{
			name: 'firefox',
			use: { ...devices['Firefox'] }
		}
	],

	/* CI: Infrastructure is started by GitHub Actions, so no webServer config */
});

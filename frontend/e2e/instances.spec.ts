import { test, expect } from "@playwright/test";
import { loginUser } from "./helpers/auth";
import {
  makeUserAdmin,
  generateRandomCategory,
  createTeam,
  ECHO_SERVER_COMPOSE,
} from "./helpers/admin";

test.describe("Challenge Instances", () => {
  test.setTimeout(120000);

  /*
	test('user can start and stop Container instance with echo-server', async ({ page }) => {
		// Create an admin user and set up a Container challenge
		const credentials = await loginUser(page);
		await makeUserAdmin(credentials.email);
		
		// Create a team (required for instance controls to show)
		await createTeam(page);
		
		await page.goto('/#/challenges');
		await page.waitForLoadState('networkidle');
		await page.reload();
		await page.waitForLoadState('networkidle');
		
		// Create category
		const category = generateRandomCategory();
		await page.getByRole('button', { name: /New Category/i }).click();
		await page.waitForTimeout(300);
		await page.locator('#cat-name').fill(category.name);
		await page.locator('#cat-icon').fill(category.icon);
		await page.getByRole('button', { name: /^Create$/i }).last().click();
		await expect(page.getByText('Category created!').first()).toBeVisible({ timeout: 5000 });
		await page.keyboard.press('Escape');
		await page.waitForTimeout(500);
		
		// Create Container challenge
		const challengeName = `Container Test ${Date.now()}`;
		await page.getByRole('button', { name: /Create Challenge/i }).click();
		await expect(page.getByRole('dialog')).toBeVisible();
		await page.waitForTimeout(300);
		
		await page.locator('#name').fill(challengeName);
		await page.locator('#description').fill('Test container instance management with echo-server');
		await page.getByRole('combobox').first().click();
		await page.waitForTimeout(200);
		await page.getByRole('option', { name: category.name }).click();
		await page.getByRole('combobox').nth(1).click();
		await page.waitForTimeout(200);
		await page.getByRole('option', { name: 'Container' }).click();
		await page.locator('#points').fill('100');
		
		await page.getByRole('button', { name: /^Create$/i }).last().click();
		await expect(page.getByText('Challenge created!').first()).toBeVisible({ timeout: 5000 });
		
		console.log(`Created Container challenge: ${challengeName}`);
		
		// Edit challenge to add Docker image and configuration
		await page.waitForTimeout(2000);
		await page.getByText(challengeName).first().click();
		
		const dialog = page.getByRole('dialog');
		await expect(dialog).toBeVisible({ timeout: 5000 });
		await page.waitForTimeout(1000);
		
		const editButton = dialog.getByRole('button', { name: 'Edit challenge' });
		await expect(editButton).toBeVisible({ timeout: 5000 });
		await editButton.click({ force: true });
		
		const deploymentButton = page.getByRole('button', { name: 'Deployment' });
		await expect(deploymentButton).toBeVisible({ timeout: 5000 });
		await page.waitForTimeout(500);
		await deploymentButton.click();
		await page.waitForTimeout(1000);
		
		// Fill in Docker configuration
		const imageInput = page.locator('#ho-host'); // Container Image Name
		await expect(imageInput).toBeVisible({ timeout: 5000 });
		await imageInput.fill('echo-server:latest');
		
		const portInput = page.locator('#ho-port'); // Port
		if (await portInput.isVisible({ timeout: 2000 })) {
			await portInput.fill('1337');
		}
		
		const lifetimeInput = page.locator('#perf-lifetime'); // Lifetime (seconds)
		await expect(lifetimeInput).toBeVisible({ timeout: 5000 });
		await lifetimeInput.fill('3600'); // 1 hour
		
		// Save the changes
		const saveButton = page.getByRole('button', { name: /Save|Update/i }).first();
		await expect(saveButton).toBeVisible({ timeout: 3000 });
		await saveButton.click();
		
		// Wait for save success
		await expect(page.getByText(/Challenge updated|saved/i).first()).toBeVisible({ timeout: 5000 });
		await page.waitForTimeout(1000);
		
		// Close edit sheet and dialog
		await page.keyboard.press('Escape');
		await page.waitForTimeout(500);
		await page.keyboard.press('Escape');
		await page.waitForTimeout(500);
		
		console.log('Added Docker configuration: echo-server:latest on port 1337, lifetime 3600s');
		
		// Refresh the page to ensure the challenge data is updated
		await page.goto('/#/challenges');
		await page.waitForLoadState('networkidle');
		await page.waitForTimeout(1000); // Extra wait for data to load
		
		// Click on the challenge to open it
		await page.getByText(challengeName).first().click();
		await page.waitForTimeout(1500); // Wait for modal to fully load
		
		// Wait for modal to be visible
		const containerDialog = page.getByRole('dialog');
		await expect(containerDialog).toBeVisible({ timeout: 10000 });
		
		// Check if instance is already running (from previous test)
		const runningButton = page.getByRole('button', { name: /Copy instance connection/i });
		const stopButton = page.getByRole('button', { name: /Stop instance/i });
		
		if (await runningButton.isVisible({ timeout: 2000 })) {
			console.log('Instance already running from previous test, stopping it first...');
			await stopButton.click();
			// Wait for instance to stop
			await page.waitForTimeout(3000);
		}
		
		// Look for "Start Instance" button
		console.log('Looking for Start Instance button...');
		const startButton = page.getByRole('button', { name: /Start Instance|Start challenge instance/i });
		await expect(startButton).toBeVisible({ timeout: 15000 });
		console.log('Found Start Instance button!');
		
		console.log('Starting Container instance...');
		await startButton.click();
		
		// Wait for instance to start - look for the running state button
		const runningButtonAfterStart = page.getByRole('button', { name: /Copy instance connection/i });
		await expect(runningButtonAfterStart).toBeVisible({ timeout: 30000 }); // Give 30 seconds for instance to start
		
		console.log('Container instance started successfully');
		
		// Check for connection info (host/port should be visible)
		const connectionInfo = page.locator('text=/Host|Port|Connection/i').first();
		if (await connectionInfo.isVisible({ timeout: 3000 })) {
			console.log('Connection information is visible');
		}
		
		// Look for port number in the UI
		const portText = await page.locator('text=/:\\d{4,5}/').first().textContent({ timeout: 5000 }).catch(() => null);
		
		if (portText) {
			const port = portText.match(/:(\d+)/)?.[1];
			console.log(`Echo server running on port: ${port}`);
			
			// Test connectivity to echo server
			try {
				const response = await page.evaluate(async (testPort) => {
					try {
						const res = await fetch(`http://localhost:${testPort}/test`, {
							method: 'GET',
						});
						const text = await res.text();
						return { success: true, status: res.status, body: text.substring(0, 200) };
					} catch (err: any) {
						return { success: false, error: err.message };
					}
				}, port);
				
				if (response.success) {
					console.log(`Echo server responded with status ${response.status}`);
					console.log(`Response preview: ${response.body}`);
				} else {
					console.log(`Could not connect to echo server: ${response.error}`);
				}
			} catch (err) {
				console.log(`Connectivity test error:`, err);
			}
		}
		
		// Wait a bit to ensure instance is fully running
		await page.waitForTimeout(2000);
		
		// Stop the instance
		console.log('Stopping Container instance...');
		const stopButtonFinal = page.getByRole('button', { name: /Stop instance/i });
		await stopButtonFinal.click();
		
		// Wait for instance to stop - Start button should reappear
		const startButtonAgain = page.getByRole('button', { name: /Start Instance|Start challenge instance/i });
		await expect(startButtonAgain).toBeVisible({ timeout: 30000 });
		
		console.log('Container instance stopped successfully');
	});

	
	test('user can start and stop Compose instance and verify connectivity', async ({ page }) => {
		// Create an admin user and set up a Compose challenge with echo server
		const credentials = await loginUser(page);
		await makeUserAdmin(credentials.email);
		
		// Create a team (required for instance controls to show)
		await createTeam(page);
		
		await page.goto('/#/challenges');
		await page.waitForLoadState('networkidle');
		await page.reload();
		await page.waitForLoadState('networkidle');
		
		// Create category
		const category = generateRandomCategory();
		await page.getByRole('button', { name: /New Category/i }).click();
		await page.waitForTimeout(300);
		await page.locator('#cat-name').fill(category.name);
		await page.locator('#cat-icon').fill(category.icon);
		await page.getByRole('button', { name: /^Create$/i }).last().click();
		await expect(page.getByText('Category created!').first()).toBeVisible({ timeout: 5000 });
		await page.keyboard.press('Escape');
		await page.waitForTimeout(500);
		
		// Create Compose challenge
		const challengeName = `Compose Test ${Date.now()}`;
		await page.getByRole('button', { name: /Create Challenge/i }).click();
		await expect(page.getByRole('dialog')).toBeVisible();
		await page.waitForTimeout(300);
		
		await page.locator('#name').fill(challengeName);
		await page.locator('#description').fill('Test compose instance with connectivity check');
		await page.getByRole('combobox').first().click();
		await page.waitForTimeout(200);
		await page.getByRole('option', { name: category.name }).click();
		await page.getByRole('combobox').nth(1).click();
		await page.waitForTimeout(200);
		await page.getByRole('option', { name: 'Compose' }).click();
		await page.locator('#points').fill('200');
		
		await page.getByRole('button', { name: /^Create$/i }).last().click();
		await expect(page.getByText('Challenge created!').first()).toBeVisible({ timeout: 5000 });
		
		console.log(`Created Compose challenge: ${challengeName}`);
		
		// Edit challenge to add compose file
		await page.waitForTimeout(2000);
		await page.getByText(challengeName).first().click();
		
		const dialog = page.getByRole('dialog');
		await expect(dialog).toBeVisible({ timeout: 5000 });
		await page.waitForTimeout(1000);
		
		const editButton = dialog.getByRole('button', { name: 'Edit challenge' });
		await expect(editButton).toBeVisible({ timeout: 5000 });
		await editButton.click({ force: true });
		
		const deploymentButton = page.getByRole('button', { name: 'Deployment' });
		await expect(deploymentButton).toBeVisible({ timeout: 5000 });
		await page.waitForTimeout(500);
		await deploymentButton.click();
		await page.waitForTimeout(1000); // Give Monaco Editor time to initialize
		
		// Monaco Editor - target the actual contenteditable div
		// Wait for Monaco editor to be visible
		await expect(page.locator('.monaco-editor')).toBeVisible({ timeout: 5000 });
		await page.waitForTimeout(500); // Extra time for Monaco to fully initialize
		
		// Use clipboard to paste content (more reliable than typing)
		await page.evaluate((text) => {
			navigator.clipboard.writeText(text);
		}, ECHO_SERVER_COMPOSE);
		
		// Click on the Monaco editor to focus it
		await page.locator('.monaco-editor .view-lines').click();
		await page.waitForTimeout(300);
		
		// Select all and paste
		await page.keyboard.press('Control+A');
		await page.keyboard.press('Control+V');
		await page.waitForTimeout(500);
		
		// Add port configuration (required for compose instances to expose the service)
		const portInput = page.locator('#ho-port');
		if (await portInput.isVisible({ timeout: 2000 })) {
			await portInput.fill('1337'); // Match the port in the compose file
		}
		
		// Add lifetime configuration (required for instances)
		const lifetimeInput = page.locator('#perf-lifetime');
		await expect(lifetimeInput).toBeVisible({ timeout: 5000 });
		await lifetimeInput.fill('3600'); // 1 hour
		
		await page.waitForTimeout(500);
		
		// Save the changes
		const saveButton = page.getByRole('button', { name: /Save|Update/i }).first();
		await expect(saveButton).toBeVisible({ timeout: 3000 });
		await saveButton.click();
		
		// Wait for save success toast
		await expect(page.getByText(/Challenge updated|saved/i).first()).toBeVisible({ timeout: 5000 });
		await page.waitForTimeout(1000);
		
		// Close edit sheet and dialog
		await page.keyboard.press('Escape');
	await page.waitForTimeout(500);
	await page.keyboard.press('Escape');
	await page.waitForTimeout(500);
	
	console.log('Added Docker Compose file to challenge');
	
	// Refresh the page to ensure the challenge data is updated
	await page.goto('/#/challenges');
	await page.waitForLoadState('networkidle');
	await page.waitForTimeout(1000); // Extra wait for data to load
	
	// Click on the challenge
	await page.getByText(challengeName).first().click();
	await page.waitForTimeout(1500); // Wait for modal to fully load
	
	// Debug: Check what's actually in the modal
	console.log(' Debugging modal content...');
	const modalContent = await page.locator('[role="dialog"]').textContent();
	console.log('Modal text content:', modalContent?.substring(0, 500));
	
	// Wait for the modal dialog to be fully visible
	const challengeDialog = page.getByRole('dialog');
	await expect(challengeDialog).toBeVisible({ timeout: 10000 });
	
	// Check if instance is already running (from previous test)
	const runningButton = page.getByRole('button', { name: /Copy instance connection/i });
	const stopButtonCheck = page.getByRole('button', { name: /Stop instance/i });
	
	if (await runningButton.isVisible({ timeout: 2000 })) {
		console.log('Instance already running from previous test, stopping it first...');
		await stopButtonCheck.click();
		// Wait for instance to stop
		await page.waitForTimeout(3000);
	}
	
	// Start instance
	console.log('Looking for Start Instance button...');
	const startButton = page.getByRole('button', { name: /Start Instance|Start challenge instance/i });
	
	// Try multiple times with longer timeout
	await expect(startButton).toBeVisible({ timeout: 15000 });
	console.log('Found Start Instance button!');
	console.log('Starting Compose instance...');
		await startButton.click();
		
		// Wait for instance to start - look for the running state button
		const runningButtonAfterStart = page.getByRole('button', { name: /Copy instance connection/i });
		await expect(runningButtonAfterStart).toBeVisible({ timeout: 45000 }); // Compose takes longer to start
		
		console.log('Compose instance started successfully');
		
		// Get connection info (host and port)
		await page.waitForTimeout(3000); // Give time for services to fully start
		
		// Look for port number in the UI
		const portText = await page.locator('text=/:\\d{4,5}/').first().textContent({ timeout: 5000 }).catch(() => null);
		
		if (portText) {
			const port = portText.match(/:(\d+)/)?.[1];
			console.log(`Instance running on port: ${port}`);
			
			// Test connectivity using fetch (basic HTTP check)
			try {
				// The echo server should respond to HTTP requests
				const response = await page.evaluate(async (testPort) => {
					try {
						const res = await fetch(`http://localhost:${testPort}/`, {
							method: 'GET',
							mode: 'no-cors' // Since it's a simple echo server
						});
						return { success: true, status: res.status };
					} catch (err: any) {
						return { success: false, error: err.message };
					}
				}, port);
				
				console.log(`Connectivity test result:`, response);
			} catch (err) {
				console.log(`Could not test connectivity (expected for echo server):`, err);
			}
		}
		
		// Stop the instance
		console.log('Stopping instance...');
		const stopButton = page.getByRole('button', { name: /Stop instance/i });
		await stopButton.click();
		
		const startButtonAgain = page.getByRole('button', { name: /Start Instance|Start challenge instance/i });
		await expect(startButtonAgain).toBeVisible({ timeout: 30000 });
		
		console.log('Instance stopped successfully');
	});

	*/

  test("Normal challenge has no instance controls", async ({ page }) => {
    const credentials = await loginUser(page);
    await makeUserAdmin(credentials.email);

    await page.goto("/#/challenges");
    await page.waitForLoadState("networkidle");
    await page.reload();
    await page.waitForLoadState("networkidle");

    const category = generateRandomCategory();
    await page.getByRole("button", { name: /New Category/i }).click();
    await page.waitForTimeout(300);
    await page.locator("#cat-name").fill(category.name);
    await page.locator("#cat-icon").fill(category.icon);
    
    const categoryResponsePromise = page.waitForResponse(
      (response) => response.url().includes('/api/categories') && response.status() === 200,
      { timeout: 10000 }
    );
    
    await page
      .getByRole("button", { name: /^Create$/i })
      .last()
      .click();
    
    await categoryResponsePromise;
    await page.keyboard.press("Escape");
    await page.waitForTimeout(500);

    const challengeName = `Normal Test ${Date.now()}`;
    await page.getByRole("button", { name: /Create Challenge/i }).click();
    await expect(page.getByRole("dialog")).toBeVisible();
    await page.waitForTimeout(300);

    await page.locator("#name").fill(challengeName);
    await page
      .locator("#description")
      .fill("Test that normal challenges have no instances");
    await page.getByRole("combobox").first().click();
    await page.waitForTimeout(200);
    await page.getByRole("option", { name: category.name }).click();
    await page.getByRole("combobox").nth(1).click();
    await page.waitForTimeout(200);
    await page.getByRole("option", { name: "Normal" }).click();
    await page.locator("#points").fill("50");

    await page
      .getByRole("button", { name: /^Create$/i })
      .last()
      .click();
    await expect(page.getByText("Challenge created!").first()).toBeVisible({
      timeout: 5000,
    });

    await page.goto("/#/challenges");
    await page.waitForLoadState("networkidle");

    await page.getByText(challengeName).first().click();
    await page.waitForTimeout(1000);

    const startButton = page.getByRole("button", { name: /Start Instance/i });
    await expect(startButton).not.toBeVisible({ timeout: 2000 });

    const stopButton = page.getByRole("button", {
      name: /Stop Instance|Delete Instance/i,
    });
    await expect(stopButton).not.toBeVisible({ timeout: 2000 });
  });
});

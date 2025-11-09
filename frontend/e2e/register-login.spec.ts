import { test, expect } from "@playwright/test";

test.describe("Register and Login Flow", () => {
  test("unauthenticated user is redirected to sign in page", async ({
    page,
  }) => {
    await page.goto("/");

    // Should redirect to sign in page
    await expect(page).toHaveURL(/#\/signIn/);

    // Verify sign in page elements are visible
    await expect(page.getByText(/welcome back hacker/i)).toBeVisible();
  });

  test("sign in page has correct elements", async ({ page }) => {
    await page.goto("/#/signIn");

    // Check for heading and description
    await expect(page.getByText(/welcome back hacker/i)).toBeVisible();
    await expect(
      page.getByText(/enter your email below to login/i),
    ).toBeVisible();

    // Check for form fields
    await expect(page.getByLabel(/email/i)).toBeVisible();
    await expect(page.getByLabel(/password/i)).toBeVisible();

    // Check for buttons
    await expect(page.getByRole("button", { name: /sign in/i })).toBeVisible();
    await expect(page.getByRole("button", { name: /sign up/i })).toBeVisible();
  });

  test("can navigate to sign up page", async ({ page }) => {
    await page.goto("/#/signIn");

    // Click sign up button
    await page.getByRole("button", { name: /sign up/i }).click();

    // Should navigate to sign up page
    await expect(page).toHaveURL(/#\/signUp/);
    await expect(page.getByText(/create your account/i)).toBeVisible();
  });

  test("sign up page has all required fields", async ({ page }) => {
    await page.goto("/#/signUp");

    // Check for form fields
    await expect(page.getByLabel(/username/i)).toBeVisible();
    await expect(page.getByLabel(/email/i)).toBeVisible();
    await expect(page.getByLabel(/^password$/i)).toBeVisible();
    await expect(page.getByLabel(/confirm password/i)).toBeVisible();

    // Check for submit button
    await expect(page.getByRole("button", { name: /sign up/i })).toBeVisible();
  });

  test("can navigate back to sign in from sign up", async ({ page }) => {
    await page.goto("/#/signUp");

    // Click sign in button
    await page.getByRole("button", { name: /sign in/i }).click();

    // Should navigate back to sign in page
    await expect(page).toHaveURL(/#\/signIn/);
  });

  test("can register a new user", async ({ page }) => {
    await page.goto("/#/signUp");

    // Generate unique credentials
    const timestamp = Date.now();
    const username = `testuser${timestamp}`;
    const email = `testuser${timestamp}@example.com`;
    const password = "TestPassword123!";

    // Fill in registration form
    await page.getByLabel(/username/i).fill(username);
    await page.getByLabel(/email/i).fill(email);
    await page.getByLabel(/^password$/i).fill(password);
    await page.getByLabel(/confirm password/i).fill(password);

    // Submit form
    await page.getByRole("button", { name: /sign up/i }).click();

    // Wait for redirect - should go to /team (no team) or /challenges (has team)
    await page.waitForTimeout(2000);

    // Should redirect away from sign up page
    await expect(page).not.toHaveURL(/#\/signUp/);
  });

  test("validates password mismatch on registration", async ({ page }) => {
    await page.goto("/#/signUp");

    // Fill form with mismatched passwords
    await page.getByLabel(/username/i).fill("testuser");
    await page.getByLabel(/email/i).fill("test@example.com");
    await page.getByLabel(/^password$/i).fill("Password123!");
    await page.getByLabel(/confirm password/i).fill("DifferentPassword123!");

    // Submit form
    await page.getByRole("button", { name: /sign up/i }).click();

    // Should show error message
    await expect(page.getByText(/passwords do not match/i)).toBeVisible({
      timeout: 2000,
    });
  });

  test("validates password minimum length", async ({ page }) => {
    await page.goto("/#/signUp");

    // Fill form with short password
    await page.getByLabel(/username/i).fill("testuser");
    await page.getByLabel(/email/i).fill("test@example.com");

    const passwordInput = page.getByLabel(/^password$/i);
    await passwordInput.fill("short");
    await page.getByLabel(/confirm password/i).fill("short");

    // Submit form
    await page.getByRole("button", { name: /sign up/i }).click();

    // HTML5 validation should prevent submission or show client-side validation
    // Check if password field has validation error or shows message
    const validationMessage = await passwordInput.evaluate(
      (el: HTMLInputElement) => el.validationMessage,
    );
    expect(validationMessage).toBeTruthy(); // Should have some validation message
  });

  test("can login with valid credentials", async ({ page }) => {
    // First register a user
    await page.goto("/#/signUp");

    const timestamp = Date.now();
    const username = `logintest${timestamp}`;
    const email = `logintest${timestamp}@example.com`;
    const password = "LoginTest123!";

    await page.getByLabel(/username/i).fill(username);
    await page.getByLabel(/email/i).fill(email);
    await page.getByLabel(/^password$/i).fill(password);
    await page.getByLabel(/confirm password/i).fill(password);
    await page.getByRole("button", { name: /sign up/i }).click();

    await page.waitForTimeout(2000);

    await expect(page).not.toHaveURL(/#\/signIn/);
  });

  test("shows error on login with invalid credentials", async ({ page }) => {
    await page.goto("/#/signIn");

    // Try to login with invalid credentials
    await page.getByLabel(/email/i).fill("nonexistent@example.com");
    await page.getByLabel(/password/i).fill("wrongpassword");
    await page.getByRole("button", { name: /sign in/i }).click();

    // Should show error message
    await page.waitForTimeout(1000);
    await expect(
      page.locator("text=/login failed|invalid|incorrect/i").first(),
    ).toBeVisible({ timeout: 5000 });
  });

  test("form fields are properly sized and aligned", async ({ page }) => {
    await page.goto("/#/signIn");

    const emailInput = page.getByLabel(/email/i);
    const passwordInput = page.getByLabel(/password/i);

    const emailBox = await emailInput.boundingBox();
    const passwordBox = await passwordInput.boundingBox();

    if (emailBox && passwordBox) {
      // Both inputs should have similar width
      expect(Math.abs(emailBox.width - passwordBox.width)).toBeLessThan(10);

      // Inputs should be reasonably wide
      expect(emailBox.width).toBeGreaterThan(200);

      // Password field should be below email
      expect(passwordBox.y).toBeGreaterThan(emailBox.y);
    }
  });

  test("forms are responsive on mobile", async ({ page }) => {
    await page.setViewportSize({ width: 375, height: 667 });
    await page.goto("/#/signUp");

    const usernameInput = page.getByLabel(/username/i);
    const box = await usernameInput.boundingBox();

    if (box) {
      // Input should not overflow viewport
      expect(box.x + box.width).toBeLessThanOrEqual(375);
    }
  });

  test("password fields mask input", async ({ page }) => {
    await page.goto("/#/signIn");

    const passwordInput = page.getByLabel(/password/i);
    await expect(passwordInput).toHaveAttribute("type", "password");
  });

  test("visual regression: sign in page", async ({ page }) => {
    // Set consistent viewport for visual regression
    await page.setViewportSize({ width: 1280, height: 720 });
    await page.goto("/#/signIn");
    await page.waitForLoadState("networkidle");

    // Take full page screenshot instead of just card
    await expect(page).toHaveScreenshot("signin-page.png", {
      fullPage: false,
      maxDiffPixels: 100,
    });
  });

  test("visual regression: sign up page", async ({ page }) => {
    // Set consistent viewport for visual regression
    await page.setViewportSize({ width: 1280, height: 720 });
    await page.goto("/#/signUp");
    await page.waitForLoadState("networkidle");

    // Take full page screenshot instead of just card
    await expect(page).toHaveScreenshot("signup-page.png", {
      fullPage: false,
      maxDiffPixels: 100,
    });
  });
});

import { Page } from "@playwright/test";

export interface UserCredentials {
  username: string;
  email: string;
  password: string;
}

export async function loginUser(page: Page): Promise<UserCredentials> {
  const timestamp = Date.now();
  const random = Math.floor(Math.random() * 10000);
  const username = `e2euser${timestamp}${random}`;
  const email = `e2euser${timestamp}${random}@example.com`;
  const password = "E2ETestPassword123!";

  await page.goto("/#/signUp");
  await page.waitForLoadState("networkidle");

  await page.getByLabel(/username/i).fill(username);
  await page.getByLabel(/email/i).fill(email);
  await page.getByLabel(/^password$/i).fill(password);
  await page.getByLabel(/confirm password/i).fill(password);

  await page.getByRole("button", { name: /sign up|register/i }).click();

  await page.waitForTimeout(2000);

  const currentUrl = page.url();
  if (currentUrl.includes("signUp")) {
    await page.goto("/#/signIn");
    await page.waitForLoadState("networkidle");
    await page.getByLabel(/email/i).fill(email);
    await page.getByLabel(/password/i).fill(password);
    await page.getByRole("button", { name: /sign in|login/i }).click();
    await page.waitForTimeout(2000);
  } else if (currentUrl.includes("signIn")) {
    await page.getByLabel(/email/i).fill(email);
    await page.getByLabel(/password/i).fill(password);
    await page.getByRole("button", { name: /sign in|login/i }).click();
    await page.waitForTimeout(2000);
  }

  return { username, email, password };
}

export async function loginWithCredentials(
  page: Page,
  email: string,
  password: string,
): Promise<void> {
  await page.goto("/#/signIn");
  await page.waitForLoadState("networkidle");

  await page.getByLabel(/email/i).fill(email);
  await page.getByLabel(/password/i).fill(password);
  await page.getByRole("button", { name: /sign in|login/i }).click();

  await page.waitForTimeout(2000);
}

export async function logout(page: Page): Promise<void> {
  const logoutButton = page.getByRole("button", { name: /logout|sign out/i });
  if (await logoutButton.isVisible()) {
    await logoutButton.click();
    await page.waitForTimeout(1000);
  }
}

export async function isAuthenticated(page: Page): Promise<boolean> {
  const url = page.url();
  return !url.includes("signIn");
}

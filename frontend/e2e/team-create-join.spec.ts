import { test, expect } from "@playwright/test";
import { loginUser } from "./helpers/auth";

test.describe("Team Create and Join Flow", () => {
  test("user1 creates a team and user2 joins it", async ({ browser }) => {
    // Create two separate browser contexts for two different users
    const context1 = await browser.newContext();
    const context2 = await browser.newContext();

    const page1 = await context1.newPage();
    const page2 = await context2.newPage();

    try {
      // Step 1: User1 registers and logs in
      const user1Credentials = await loginUser(page1);
      console.log("User1 logged in:", user1Credentials.email);

      // Step 2: User1 should be on /team page with no team (TeamJoinCreate component)
      await expect(page1).toHaveURL(/\/#\/team/);

      // Verify the TeamJoinCreate cards are visible by checking for the action buttons
      await expect(
        page1.getByRole("button", { name: /^Join Team$/i }),
      ).toBeVisible();
      await expect(
        page1.getByRole("button", { name: /^Create Team$/i }),
      ).toBeVisible(); // Step 3: User1 creates a team
      const teamName = `TestTeam_${Date.now()}`;
      const teamPassword = "TeamPass123";

      // Click the "Create" button to open the dialog
      await page1.getByRole("button", { name: "Create" }).click();

      // Wait for the dialog to appear
      await expect(page1.getByRole("dialog")).toBeVisible();
      await page1.waitForTimeout(200);
      // Check dialog is open by looking for the confirm password field (unique to create dialog)
      await expect(page1.locator("#confirm-pass")).toBeVisible(); // Fill in the team creation form
      const teamNameInput = page1.locator("#team-name").last(); // Use .last() in case there are multiple dialogs
      const teamPassInput = page1.locator("#team-pass").last();
      const confirmPassInput = page1.locator("#confirm-pass");

      await teamNameInput.fill(teamName);
      await teamPassInput.fill(teamPassword);
      await confirmPassInput.fill(teamPassword);

      // Submit the form
      await page1
        .getByRole("button", { name: /Create/i })
        .last()
        .click();

      // Wait for success toast
      await expect(page1.getByText("Team Created!")).toBeVisible({
        timeout: 10000,
      });

      // Verify user1 is now viewing their team page with team details
      await expect(page1.getByText(teamName)).toBeVisible({ timeout: 10000 });
      await expect(
        page1.getByRole("button", { name: "Overview" }),
      ).toBeVisible();
      await expect(
        page1.getByRole("button", { name: "Members" }),
      ).toBeVisible();
      await expect(page1.getByRole("button", { name: "Solves" })).toBeVisible();

      // Step 4: User2 registers and logs in
      const user2Credentials = await loginUser(page2);
      console.log("User2 logged in:", user2Credentials.email);

      // Step 5: User2 should be on /team page with no team
      await expect(page2).toHaveURL(/\/#\/team/);
      await expect(
        page2.getByRole("button", { name: /^Join Team$/i }),
      ).toBeVisible(); // Step 6: User2 joins the team created by user1
      // Click the "Join" button to open the dialog
      await page2.getByRole("button", { name: "Join" }).click();

      // Wait for the join dialog
      await expect(page2.getByRole("dialog")).toBeVisible();
      // Verify it's the join dialog by checking for team name input
      await expect(page2.locator("#team-name").last()).toBeVisible(); // Fill in the join form
      const joinTeamNameInput = page2.locator("#team-name").last();
      const joinTeamPassInput = page2.locator("#team-pass").last();

      await joinTeamNameInput.fill(teamName);
      await joinTeamPassInput.fill(teamPassword);

      // Submit the form
      await page2.getByRole("button", { name: /Join/i }).last().click();

      // Wait for success toast
      await expect(page2.getByText("Team Joined, welcome aboard!")).toBeVisible(
        { timeout: 10000 },
      );

      // Verify user2 is now viewing the team page
      await expect(page2.getByText(teamName)).toBeVisible({ timeout: 10000 });

      // Step 7: Verify both users can see the team has 2 members
      // Refresh user1's page to see updated member count
      await page1.reload();
      await page1.waitForLoadState("networkidle");

      // The member count should be visible in the header area
      await expect(page1.getByText("2 members")).toBeVisible({
        timeout: 10000,
      });

      // Also verify on page2
      await expect(page2.getByText("2 members")).toBeVisible({
        timeout: 10000,
      });

      console.log("Test completed successfully: Team created and joined");
    } finally {
      // Cleanup
      await context1.close();
      await context2.close();
    }
  });

  test("cannot create team with mismatched passwords", async ({ page }) => {
    await loginUser(page);

    // Should be on team page
    await expect(page).toHaveURL(/\/#\/team/);

    // Click Create button
    await page.getByRole("button", { name: "Create" }).click();
    await expect(page.getByRole("dialog")).toBeVisible();

    // Fill form with mismatched passwords
    const teamNameInput = page.locator("#team-name").last();
    const teamPassInput = page.locator("#team-pass").last();
    const confirmPassInput = page.locator("#confirm-pass");

    await teamNameInput.fill(`Team_${Date.now()}`);
    await teamPassInput.fill("Password123");
    await confirmPassInput.fill("DifferentPass123");

    // Submit
    await page
      .getByRole("button", { name: /Create/i })
      .last()
      .click();

    // Should show error toast
    await expect(page.getByText("Passwords do not match.")).toBeVisible();
  });

  test("cannot create team with short password", async ({ page }) => {
    await loginUser(page);

    await expect(page).toHaveURL(/\/#\/team/);

    // Click Create button
    await page.getByRole("button", { name: "Create" }).click();
    await expect(page.getByRole("dialog")).toBeVisible();

    // Fill form with short password
    const teamNameInput = page.locator("#team-name").last();
    const teamPassInput = page.locator("#team-pass").last();
    const confirmPassInput = page.locator("#confirm-pass");

    await teamNameInput.fill(`Team_${Date.now()}`);
    await teamPassInput.fill("short");
    await confirmPassInput.fill("short");

    // Submit
    await page
      .getByRole("button", { name: /Create/i })
      .last()
      .click();

    // Should show error toast
    await expect(
      page.getByText("Password must be at least 8 characters."),
    ).toBeVisible();
  });

  test("cannot join team with wrong password", async ({ browser }) => {
    const context1 = await browser.newContext();
    const context2 = await browser.newContext();

    const page1 = await context1.newPage();
    const page2 = await context2.newPage();

    try {
      // User1 creates a team
      await loginUser(page1);
      await expect(page1).toHaveURL(/\/#\/team/);

      const teamName = `SecureTeam_${Date.now()}`;
      const correctPassword = "CorrectPass123";

      await page1.getByRole("button", { name: "Create" }).click();
      await expect(page1.getByRole("dialog")).toBeVisible();

      const teamNameInput = page1.locator("#team-name").last();
      const teamPassInput = page1.locator("#team-pass").last();
      const confirmPassInput = page1.locator("#confirm-pass");

      await teamNameInput.fill(teamName);
      await teamPassInput.fill(correctPassword);
      await confirmPassInput.fill(correctPassword);

      await page1
        .getByRole("button", { name: /Create/i })
        .last()
        .click();
      await expect(page1.getByText("Team Created!")).toBeVisible({
        timeout: 10000,
      });

      // User2 tries to join with wrong password
      await loginUser(page2);
      await expect(page2).toHaveURL(/\/#\/team/);

      await page2.getByRole("button", { name: "Join" }).click();
      await expect(page2.getByRole("dialog")).toBeVisible();

      const joinTeamNameInput = page2.locator("#team-name").last();
      const joinTeamPassInput = page2.locator("#team-pass").last();

      await joinTeamNameInput.fill(teamName);
      await joinTeamPassInput.fill("WrongPassword123");

      await page2.getByRole("button", { name: /Join/i }).last().click();

      await page2.waitForTimeout(2000); // Give time for error to appear

      await expect(
        page2.getByRole("button", { name: /^Create Team$/i }).first(),
      ).toBeVisible();
    } finally {
      await context1.close();
      await context2.close();
    }
  });

  test("cannot join non-existent team", async ({ page }) => {
    await loginUser(page);

    await expect(page).toHaveURL(/\/#\/team/);

    await page.getByRole("button", { name: "Join" }).click();
    await expect(page.getByRole("dialog")).toBeVisible();

    const joinTeamNameInput = page.locator("#team-name").last();
    const joinTeamPassInput = page.locator("#team-pass").last();

    await joinTeamNameInput.fill(`NonExistentTeam_${Date.now()}`);
    await joinTeamPassInput.fill("SomePassword123");

    await page.getByRole("button", { name: /Join/i }).last().click();

    // Wait for error
    await page.waitForTimeout(2000);

    // Should still see join/create cards - check for Create button which is unique
    await expect(
      page.getByRole("button", { name: /^Create Team$/i }).first(),
    ).toBeVisible();
  });

  test("team page shows correct statistics after creation", async ({
    page,
  }) => {
    await loginUser(page);

    await expect(page).toHaveURL(/\/#\/team/);

    // Create team
    const teamName = `StatsTeam_${Date.now()}`;
    await page.getByRole("button", { name: "Create" }).click();
    await expect(page.getByRole("dialog")).toBeVisible();

    const teamNameInput = page.locator("#team-name").last();
    const teamPassInput = page.locator("#team-pass").last();
    const confirmPassInput = page.locator("#confirm-pass");

    await teamNameInput.fill(teamName);
    await teamPassInput.fill("TeamPass123");
    await confirmPassInput.fill("TeamPass123");

    await page
      .getByRole("button", { name: /Create/i })
      .last()
      .click();
    await expect(page.getByText("Team Created!")).toBeVisible({
      timeout: 10000,
    });

    // Verify team page loads with statistics
    await expect(page.getByText(teamName)).toBeVisible();

    // Check statistics section
    await expect(page.getByText("Statistics")).toBeVisible();
    await expect(page.getByText("Total Points")).toBeVisible();

    // New team should have 1 member, 0 points, 0 solves
    await expect(page.getByText("1 member")).toBeVisible();
  });

  test("can cancel team creation dialog", async ({ page }) => {
    await loginUser(page);

    await expect(page).toHaveURL(/\/#\/team/);

    // Open create dialog
    await page.getByRole("button", { name: "Create" }).click();
    await expect(page.getByRole("dialog")).toBeVisible();

    // Click Cancel
    await page.getByRole("button", { name: "Cancel" }).first().click();

    // Dialog should close, cards should still be visible
    await expect(page.getByRole("dialog")).not.toBeVisible();
    await expect(page.getByRole("button", { name: /^Join Team$/i })).toBeVisible();
    await expect(page.getByRole("button", { name: /^Create Team$/i })).toBeVisible();
  });

  test("can cancel team join dialog", async ({ page }) => {
    await loginUser(page);

    await expect(page).toHaveURL(/\/#\/team/);

    // Open join dialog
    await page.getByRole("button", { name: "Join" }).click();
    await expect(page.getByRole("dialog")).toBeVisible();

    // Click Cancel
    await page.getByRole("button", { name: "Cancel" }).first().click();

    // Dialog should close and cards should be visible
    await expect(page.getByRole("dialog")).not.toBeVisible();
    await expect(page.getByRole("button", { name: /^Create Team$/i })).toBeVisible();
  });

  test("visual regression: team join/create page", async ({ page }) => {
    await page.setViewportSize({ width: 1280, height: 720 });
    await loginUser(page);

    await expect(page).toHaveURL(/\/#\/team/);
    await page.waitForLoadState("networkidle");

    // Take screenshot of the join/create cards
    await expect(page).toHaveScreenshot("team-join-create-page.png", {
      fullPage: false,
      maxDiffPixels: 100,
    });
  });
});

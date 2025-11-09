import { test, expect } from "@playwright/test";
import { loginUser } from "./helpers/auth";
import {
  makeUserAdmin,
  generateRandomCategory,
  generateRandomChallenge,
  ECHO_SERVER_COMPOSE,
} from "./helpers/admin";

test.describe("Admin Features", () => {
  test("setup: create admin user and verify admin controls", async ({
    page,
  }) => {
    // Create a new user and make them admin
    const credentials = await loginUser(page);

    console.log("Created user:", credentials.email);

    // Make the user an admin via database
    await makeUserAdmin(credentials.email);

    // Reload the page to get updated permissions
    await page.reload();
    await page.waitForLoadState("networkidle");

    // Go to challenges page
    await page.goto("/#/challenges");
    await page.waitForLoadState("networkidle");

    // Should see admin controls
    await expect(
      page.getByRole("button", { name: /Create Challenge/i }),
    ).toBeVisible();
    await expect(
      page.getByRole("button", { name: /New Category/i }),
    ).toBeVisible();

    console.log("Admin controls visible");
  });

  test("admin can create categories", async ({ page }) => {
    // Login as admin
    const credentials = await loginUser(page);
    await makeUserAdmin(credentials.email);

    await page.goto("/#/challenges");
    await page.waitForLoadState("networkidle");
    await page.reload(); // Reload to get admin permissions
    await page.waitForLoadState("networkidle");

    // Create 3 random categories
    for (let i = 0; i < 3; i++) {
      const category = generateRandomCategory();

      // Click "New Category" button
      await page.getByRole("button", { name: /New Category/i }).click();

      // Wait for popover
      await page.waitForTimeout(300);

      // Fill in category details
      await page.locator("#cat-name").fill(category.name);
      await page.locator("#cat-icon").fill(category.icon);

      // Submit
      await page
        .getByRole("button", { name: /^Create$/i })
        .last()
        .click();

      // Wait for success toast
      await expect(page.getByText("Category created!").first()).toBeVisible({
        timeout: 5000,
      });

      console.log(`Created category: ${category.name}`); // Wait a bit before next category
      await page.waitForTimeout(500);
    }
  });

  test("admin can create random challenges", async ({ page }) => {
    test.setTimeout(60000); // Increase timeout to 60 seconds for this test

    // Login as admin
    const credentials = await loginUser(page);
    await makeUserAdmin(credentials.email);

    await page.goto("/#/challenges");
    await page.waitForLoadState("networkidle");
    await page.reload(); // Reload to get admin permissions
    await page.waitForLoadState("networkidle");

    // First, create a test category to use
    const testCategory = generateRandomCategory();
    await page.getByRole("button", { name: /New Category/i }).click();
    await page.waitForTimeout(300);
    await page.locator("#cat-name").fill(testCategory.name);
    await page.locator("#cat-icon").fill(testCategory.icon);
    await page
      .getByRole("button", { name: /^Create$/i })
      .last()
      .click();
    await expect(page.getByText("Category created!").first()).toBeVisible({
      timeout: 5000,
    });

    // Close the popover by pressing Escape
    await page.keyboard.press("Escape");
    await page.waitForTimeout(500);

    // Create between 3-8 random challenges (reduced from 5-15 to avoid timeout)
    const numChallenges = Math.floor(Math.random() * 6) + 3; // Random between 3-8
    console.log(`Creating ${numChallenges} random challenges...`);

    for (let i = 0; i < numChallenges; i++) {
      const challenge = generateRandomChallenge(testCategory.name);

      // Click "Create Challenge" button
      await page.getByRole("button", { name: /Create Challenge/i }).click();

      // Wait for dialog
      await expect(page.getByRole("dialog")).toBeVisible();
      await page.waitForTimeout(300);

      // Fill in challenge details
      await page.locator("#name").fill(challenge.name);
      await page.locator("#description").fill(challenge.description);

      // Select category - click the combobox button labeled "Category*"
      await page.getByRole("combobox").first().click(); // This opens the category dropdown
      await page.waitForTimeout(200);
      await page.getByRole("option", { name: testCategory.name }).click();

      // Select type - click the second combobox for "Type*"
      await page.getByRole("combobox").nth(1).click(); // This opens the type dropdown
      await page.waitForTimeout(200);
      await page.getByRole("option", { name: challenge.type }).click(); // Set points
      await page.locator("#points").fill(challenge.points.toString());

      // Set dynamic score checkbox
      const checkbox = page.locator("#scoretype");
      const isChecked = await checkbox.isChecked();
      if (challenge.dynamicScore && !isChecked) {
        await checkbox.click();
      } else if (!challenge.dynamicScore && isChecked) {
        await checkbox.click();
      }

      // Submit
      await page
        .getByRole("button", { name: /^Create$/i })
        .last()
        .click();

      // Wait for success
      await expect(page.getByText("Challenge created!").first()).toBeVisible({
        timeout: 5000,
      });

      console.log(
        `Created challenge ${i + 1}/${numChallenges}: ${challenge.name} (${challenge.type})`,
      );
      await page.waitForTimeout(500);
    }
  });

  test("admin can create instanced challenges with compose", async ({
    page,
  }) => {
    test.setTimeout(90000); // Increase timeout to 90 seconds (creates 5 challenges + edits each)

    // Login as admin
    const credentials = await loginUser(page);
    await makeUserAdmin(credentials.email);

    await page.goto("/#/challenges");
    await page.waitForLoadState("networkidle");
    await page.reload(); // Reload to get admin permissions
    await page.waitForLoadState("networkidle");

    // Create a category for instanced challenges
    const instanceCategory = generateRandomCategory();
    await page.getByRole("button", { name: /New Category/i }).click();
    await page.waitForTimeout(300);
    await page.locator("#cat-name").fill(instanceCategory.name);
    await page.locator("#cat-icon").fill(instanceCategory.icon);
    await page
      .getByRole("button", { name: /^Create$/i })
      .last()
      .click();
    await expect(page.getByText("Category created!").first()).toBeVisible({
      timeout: 5000,
    });

    // Close the popover by pressing Escape or clicking outside
    await page.keyboard.press("Escape");
    await page.waitForTimeout(500);

    // Create at least 5 instanced challenges with compose
    for (let i = 0; i < 5; i++) {
      const challengeName = `Instance Test ${Date.now()}_${i}`;

      // Create challenge
      await page.getByRole("button", { name: /Create Challenge/i }).click();
      await expect(page.getByRole("dialog")).toBeVisible();
      await page.waitForTimeout(300);

      await page.locator("#name").fill(challengeName);
      await page
        .locator("#description")
        .fill("Test instanced challenge with Docker Compose");

      // Select category - click the combobox button
      await page.getByRole("combobox").first().click();
      await page.waitForTimeout(200);
      await page.getByRole("option", { name: instanceCategory.name }).click();

      // Select Compose type - click the second combobox
      await page.getByRole("combobox").nth(1).click();
      await page.waitForTimeout(200);
      await page.getByRole("option", { name: "Compose" }).click();
      await page.locator("#points").fill("500");

      // Submit
      await page
        .getByRole("button", { name: /^Create$/i })
        .last()
        .click();
      await expect(page.getByText("Challenge created!").first()).toBeVisible({
        timeout: 5000,
      });

      console.log(`Created instanced challenge ${i + 1}/5: ${challengeName}`);

      // Wait for toast to disappear before interacting with challenges
      await page.waitForTimeout(2000);

      // Now we need to edit the challenge to add the compose file
      // Click on the challenge card to open the detail modal
      await page.getByText(challengeName).first().click();

      // Wait for the challenge modal dialog to open
      const dialog = page.getByRole("dialog");
      await expect(dialog).toBeVisible({ timeout: 5000 });
      await page.waitForTimeout(1000); // Give more time for dialog to fully render

      // Find and click the Edit button within the dialog using aria-label
      const editButton = dialog.getByRole("button", { name: "Edit challenge" });
      await expect(editButton).toBeVisible({ timeout: 5000 });
      await editButton.click({ force: true }); // Force click in case something is overlaying

      // Wait for the edit sheet to actually open by checking for the Deployment button
      const deploymentButton = page.getByRole("button", { name: "Deployment" });
      await expect(deploymentButton).toBeVisible({ timeout: 5000 });
      await page.waitForTimeout(500);

      // Click on the "Deployment" button to access the compose file section
      await deploymentButton.click();
      await page.waitForTimeout(1000); // Give Monaco Editor time to initialize

      // Monaco Editor - use clipboard to paste content (more reliable than typing)
      await expect(page.locator(".monaco-editor")).toBeVisible({
        timeout: 5000,
      });
      await page.waitForTimeout(500); // Extra time for Monaco to fully initialize

      // Use clipboard to paste content
      await page.evaluate((text) => {
        navigator.clipboard.writeText(text);
      }, ECHO_SERVER_COMPOSE);

      // Click on the Monaco editor to focus it
      await page.locator(".monaco-editor .view-lines").click();
      await page.waitForTimeout(300);

      // Select all and paste
      await page.keyboard.press("Control+A");
      await page.keyboard.press("Control+V");
      await page.waitForTimeout(500);

      // Save changes
      const saveButton = page
        .getByRole("button", { name: /Save|Update/i })
        .first();
      await expect(saveButton).toBeVisible({ timeout: 3000 });
      await saveButton.click();

      // Wait for save success
      await expect(
        page.getByText(/Challenge updated|saved/i).first(),
      ).toBeVisible({ timeout: 5000 });
      await page.waitForTimeout(1000);

      // Close any open dialogs and go back to challenges list
      await page.keyboard.press("Escape");
      await page.waitForTimeout(500);
      await page.keyboard.press("Escape");
      await page.waitForTimeout(500);

      // Verify we're back on the challenges list (no dialogs open)
      await expect(page.getByRole("dialog")).not.toBeVisible();
    }
    console.log("Created 5 instanced challenges with compose files");
  });

  test("admin can edit challenge and update docker compose", async ({
    page,
  }) => {
    // Login as admin
    const credentials = await loginUser(page);
    await makeUserAdmin(credentials.email);

    await page.goto("/#/challenges");
    await page.waitForLoadState("networkidle");
    await page.reload(); // Reload to get admin permissions
    await page.waitForLoadState("networkidle");

    // Create a test category and challenge first
    const category = generateRandomCategory();
    await page.getByRole("button", { name: /New Category/i }).click();
    await page.waitForTimeout(300);
    await page.locator("#cat-name").fill(category.name);
    await page.locator("#cat-icon").fill(category.icon);
    await page
      .getByRole("button", { name: /^Create$/i })
      .last()
      .click();
    await expect(page.getByText("Category created!").first()).toBeVisible({
      timeout: 5000,
    });

    // Close the popover by pressing Escape
    await page.keyboard.press("Escape");
    await page.waitForTimeout(500);

    const challengeName = `Editable Challenge ${Date.now()}`;
    await page.getByRole("button", { name: /Create Challenge/i }).click();
    await expect(page.getByRole("dialog")).toBeVisible();
    await page.waitForTimeout(300);

    await page.locator("#name").fill(challengeName);
    await page.locator("#description").fill("Challenge to be edited");

    // Select category - click the combobox button
    await page.getByRole("combobox").first().click();
    await page.waitForTimeout(200);
    await page.getByRole("option", { name: category.name }).click();

    // Select type - click the second combobox
    await page.getByRole("combobox").nth(1).click();
    await page.waitForTimeout(200);
    await page.getByRole("option", { name: "Compose" }).click();

    await page.locator("#points").fill("1000");
    await page
      .getByRole("button", { name: /^Create$/i })
      .last()
      .click();
    await expect(page.getByText("Challenge created!").first()).toBeVisible({
      timeout: 5000,
    });

    console.log(`Created challenge: ${challengeName}`);

    // Now edit it - click on challenge card to open the detail modal
    await page.waitForTimeout(1000);
    await page.getByText(challengeName).first().click();

    // Wait for the challenge modal dialog to open
    const dialog = page.getByRole("dialog");
    await expect(dialog).toBeVisible({ timeout: 5000 });
    await page.waitForTimeout(1000); // Give more time for dialog to fully render

    // Find and click the Edit button within the dialog using aria-label
    const editButton = dialog.getByRole("button", { name: "Edit challenge" });
    await expect(editButton).toBeVisible({ timeout: 5000 });
    await editButton.click();

    // Wait for the edit sheet to open by checking for Deployment button
    const deploymentButton = page.getByRole("button", { name: "Deployment" });
    await expect(deploymentButton).toBeVisible({ timeout: 5000 });
    await page.waitForTimeout(500);

    // Update description (should be on the first/default tab)
    const descTextarea = page
      .locator('#description, textarea[name="description"]')
      .first();
    if (await descTextarea.isVisible({ timeout: 3000 })) {
      await descTextarea.fill(
        "Updated description: This challenge has been edited!",
      );
    }

    // Click on the "Deployment" button to access the compose file section
    await deploymentButton.click();
    await page.waitForTimeout(1000); // Give Monaco Editor time to initialize

    // Monaco Editor - use clipboard to paste content (more reliable than typing)
    await expect(page.locator(".monaco-editor")).toBeVisible({ timeout: 5000 });
    await page.waitForTimeout(500); // Extra time for Monaco to fully initialize

    // Use clipboard to paste content
    await page.evaluate((text) => {
      navigator.clipboard.writeText(text);
    }, ECHO_SERVER_COMPOSE);

    // Click on the Monaco editor to focus it
    await page.locator(".monaco-editor .view-lines").click();
    await page.waitForTimeout(300);

    // Select all and paste
    await page.keyboard.press("Control+A");
    await page.keyboard.press("Control+V");
    await page.waitForTimeout(500);

    // Update points
    const pointsInput = page.locator('#points, input[name="points"]').first();
    if (await pointsInput.isVisible({ timeout: 3000 })) {
      await pointsInput.fill("1500");
    }

    // Save
    const saveButton = page
      .getByRole("button", { name: /Save|Update/i })
      .first();
    await expect(saveButton).toBeVisible({ timeout: 3000 });
    await saveButton.click();

    // Wait for save success
    await expect(
      page.getByText(/Challenge updated|saved/i).first(),
    ).toBeVisible({ timeout: 5000 });
    await page.waitForTimeout(1000);

    console.log("Successfully edited challenge with compose file");
  });
});

import { exec } from "child_process";
import { promisify } from "util";

const execAsync = promisify(exec);

export async function makeUserAdmin(email: string): Promise<void> {
  try {
    const command = `docker exec trxd-postgres-1 psql -U user -d postgres -c "UPDATE users SET role = 'Admin' WHERE email = '${email}';"`;

    const { stdout, stderr } = await execAsync(command);

    if (stderr && !stderr.includes("UPDATE")) {
      console.error("Error making user admin:", stderr);
      throw new Error(`Failed to make user admin: ${stderr}`);
    }

    console.log(`Made user ${email} an admin`);
  } catch (error) {
    console.error("Failed to make user admin:", error);
    throw error;
  }
}

export function generateRandomChallenge(categoryName: string) {
  const adjectives = [
    "Easy",
    "Hard",
    "Tricky",
    "Simple",
    "Complex",
    "Advanced",
    "Basic",
  ];
  const nouns = ["Crypto", "Reverse", "Pwn", "Web", "Forensics", "Misc"];
  const randomAdj = adjectives[Math.floor(Math.random() * adjectives.length)];
  const randomNoun = nouns[Math.floor(Math.random() * nouns.length)];
  const randomNum = Math.floor(Math.random() * 1000);

  const types = ["Container", "Compose", "Normal"];
  const randomType = types[Math.floor(Math.random() * types.length)];

  const points = [100, 200, 300, 400, 500, 750, 1000][
    Math.floor(Math.random() * 7)
  ];

  return {
    name: `${randomAdj} ${randomNoun} ${randomNum}`,
    description: `This is a test challenge for ${categoryName}. Solve it to get points!`,
    category: categoryName,
    type: randomType,
    points: points,
    dynamicScore: Math.random() > 0.5,
  };
}

export function generateRandomCategory() {
  const icons = [
    "Bug",
    "Shield",
    "Lock",
    "Key",
    "Binary",
    "Code",
    "Terminal",
    "Server",
    "Database",
  ];
  const prefixes = ["Test", "Demo", "Sample", "Example"];
  const subjects = [
    "Crypto",
    "Web",
    "Pwn",
    "Reverse",
    "Forensics",
    "Misc",
    "OSINT",
  ];

  const randomPrefix = prefixes[Math.floor(Math.random() * prefixes.length)];
  const randomSubject = subjects[Math.floor(Math.random() * subjects.length)];
  const randomIcon = icons[Math.floor(Math.random() * icons.length)];
  const randomNum = Math.floor(Math.random() * 1000);

  return {
    name: `${randomPrefix} ${randomSubject} ${randomNum}`,
    icon: randomIcon,
  };
}

export async function createTeam(
  page: any,
): Promise<{ teamName: string; teamPassword: string }> {
  const teamName = `E2ETeam_${Date.now()}`;
  const teamPassword = "TeamPass123!";

  await page.goto("/#/team");
  await page.waitForLoadState("networkidle");

  await page.getByRole("button", { name: "Create" }).click();
  await page.waitForTimeout(300);

  const teamNameInput = page.locator("#team-name").last();
  const teamPassInput = page.locator("#team-pass").last();
  const confirmPassInput = page.locator("#confirm-pass");

  await teamNameInput.fill(teamName);
  await teamPassInput.fill(teamPassword);
  await confirmPassInput.fill(teamPassword);

  await page
    .getByRole("button", { name: /Create/i })
    .last()
    .click();

  await page.waitForTimeout(2000);

  console.log(`Created team: ${teamName}`);

  return { teamName, teamPassword };
}

export const ECHO_SERVER_COMPOSE = `services:
  echo:
    image: echo-server:latest
    container_name: \${CONTAINER_NAME}
    network_mode: bridge
    ports:
      - "\${INSTANCE_PORT}:1337"
`;

export const ECHO_SERVER_DOCKERFILE = `FROM hashicorp/http-echo
CMD ["-text=TRX{test_flag_12345}"]
`;
